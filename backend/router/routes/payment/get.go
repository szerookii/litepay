package payment

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/ent"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
	"github.com/szerookii/litepay/backend/utils"
)

type GetPaymentResponse struct {
	ID                   string            `json:"id"`
	WalletAddress        string            `json:"wallet_address"`
	SolReference         *string           `json:"sol_reference,omitempty"`
	AmountCrypto         float64           `json:"amount_crypto"`
	CurrencyCryptoName   string            `json:"currency_crypto_name"`
	CurrencyCryptoSymbol string            `json:"currency_crypto_symbol"`
	AmountFiat           float64           `json:"amount_fiat"`
	CurrencyFiat         string            `json:"currency_fiat"`
	Status               entpayment.Status `json:"status"`
	ExpiresAt            time.Time         `json:"expires_at"`
	LastTransactionHash  string            `json:"last_transaction_hash,omitempty"`
	Confirmations        *int              `json:"confirmations"`
	RequiredConfirmations int              `json:"required_confirmations"`
}

func toResponse(p *ent.Payment, cryptoName, cryptoSymbol, txHash string, confs *int, reqConfs int) *GetPaymentResponse {
	return &GetPaymentResponse{
		ID:                    p.ID.String(),
		WalletAddress:         p.WalletAddress,
		SolReference:          p.SolReference,
		AmountCrypto:          p.AmountCrypto,
		CurrencyCryptoName:    cryptoName,
		CurrencyCryptoSymbol:  cryptoSymbol,
		AmountFiat:            p.AmountFiat,
		CurrencyFiat:          p.CurrencyFiat,
		Status:                p.Status,
		ExpiresAt:             p.ExpiresAt,
		LastTransactionHash:   txHash,
		Confirmations:         confs,
		RequiredConfirmations: reqConfs,
	}
}

func Get(ctx *gin.Context) {
	p, err := db.PaymentById(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "payment not found"})
		return
	}

	// Terminal states — skip on-chain check
	if p.Status == entpayment.StatusPAID || p.Status == entpayment.StatusEXPIRED {
		utils.SendJSON(ctx, http.StatusOK, toResponse(p, p.CurrencyCrypto, p.CurrencyCrypto, "", nil, 0))
		return
	}

	if p.ExpiresAt.Before(time.Now()) {
		if p.Status != entpayment.StatusEXPIRED {
			if _, err := db.UpdatePayment(p.ID, entpayment.StatusEXPIRED); err != nil {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			p.Status = entpayment.StatusEXPIRED
		}
		utils.SendJSON(ctx, http.StatusOK, toResponse(p, p.CurrencyCrypto, p.CurrencyCrypto, "", nil, 0))
		return
	}

	chain, ok := cryptocurrency.GetBySymbol(p.CurrencyCrypto)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	addr := &cryptocurrency.PaymentAddress{
		Address:   p.WalletAddress,
		Index:     uint32(p.AddressIndex),
	}
	if p.SolReference != nil {
		addr.Reference = *p.SolReference
	}

	status, err := chain.CheckPayment(addr, p.AmountCrypto)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var newStatus entpayment.Status
	var confs *int

	switch {
	case status.ReceivedAmount >= p.AmountCrypto && status.Confirmations >= chain.RequiredConfirmations():
		newStatus = entpayment.StatusPAID
	case status.ReceivedAmount > 0 || status.IsPending:
		newStatus = entpayment.StatusCONFIRMING
		confs = utils.Ptr(status.Confirmations)
	default:
		newStatus = entpayment.StatusPENDING
	}

	if newStatus != p.Status {
		if _, err := db.UpdatePaymentWithHash(p.ID, newStatus, status.TxHash, status.ReceivedAmount); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		p.Status = newStatus
	}

	utils.SendJSON(ctx, http.StatusOK, toResponse(p, chain.Name(), chain.Symbol(), status.TxHash, confs, chain.RequiredConfirmations()))
}
