package payment

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/ent"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
	"github.com/szerookii/litepay/backend/router/middleware"
	"github.com/szerookii/litepay/backend/secrets"
	"github.com/szerookii/litepay/backend/utils"
)

type CreatePaymentRequest struct {
	Symbol   string  `json:"symbol"   validate:"required"`
	Amount   float64 `json:"amount"   validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type CreatePaymentResponse struct {
	ID             string            `json:"id"`
	WalletAddress  string            `json:"wallet_address"`
	AmountCrypto   float64           `json:"amount_crypto"`
	CurrencyCrypto string            `json:"currency_crypto"`
	AmountFiat     float64           `json:"amount_fiat"`
	CurrencyFiat   string            `json:"currency_fiat"`
	Status         entpayment.Status `json:"status"`
	ExpiresAt      time.Time         `json:"expires_at"`
}

func Post(ctx *gin.Context) {
	u := ctx.MustGet(middleware.AuthedUserKey).(*ent.User)

	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if utils.Validate(req) != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	chain, ok := cryptocurrency.GetBySymbol(strings.ToUpper(req.Symbol))
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "cryptocurrency not supported"})
		return
	}

	exchangeRate, err := chain.Price(req.Currency)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get %s price in %s", chain.Symbol(), req.Currency)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get price"})
		return
	}

	cryptoAmount := req.Amount / exchangeRate

	masterSeed, err := secrets.DeriveMasterSeed()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "seed unavailable"})
		return
	}
	defer func() {
		for i := range masterSeed {
			masterSeed[i] = 0
		}
	}()

	dbTx, err := db.Client().Tx(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to start transaction"})
		return
	}

	paymentIndex, err := db.NextPaymentIndex(context.Background(), dbTx, u.ID, chain.Symbol())
	if err != nil {
		_ = dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get payment index"})
		return
	}

	addr, err := chain.NewPaymentAddress(masterSeed, uint32(u.AccountIndex), paymentIndex)
	if err != nil {
		_ = dbTx.Rollback()
		log.Error().Err(err).Msgf("failed to derive %s address", chain.Symbol())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to derive address"})
		return
	}

	p, err := db.CreatePayment(dbTx, u.ID, addr.Address, nil, int(paymentIndex), cryptoAmount, chain.Symbol(), req.Amount, req.Currency, time.Now().Add(time.Hour))
	if err != nil {
		_ = dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create payment"})
		return
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to commit transaction"})
		return
	}

	log.Debug().Msgf("payment created: id=%s addr=%s amount=%f %s", p.ID, addr.Address, cryptoAmount, chain.Symbol())

	utils.SendJSON(ctx, http.StatusOK, &CreatePaymentResponse{
		ID:             p.ID.String(),
		WalletAddress:  p.WalletAddress,
		AmountCrypto:   p.AmountCrypto,
		CurrencyCrypto: p.CurrencyCrypto,
		AmountFiat:     p.AmountFiat,
		CurrencyFiat:   p.CurrencyFiat,
		Status:         p.Status,
		ExpiresAt:      p.ExpiresAt,
	})
}
