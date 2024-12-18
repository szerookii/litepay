package payment

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/szerookii/litepay/crypto"
	"github.com/szerookii/litepay/db"
	prisma "github.com/szerookii/litepay/prisma/db"
	"github.com/szerookii/litepay/utils"
	"time"
)

type GetPaymentResponse struct {
	Id                   string               `json:"id"`
	WalletAddress        string               `json:"wallet_address"`
	AmountCrypto         float64              `json:"amount_crypto"`
	CurrencyCryptoName   string               `json:"currency_crypto_name"`
	CurrencyCryptoSymbol string               `json:"currency_crypto_symbol"`
	AmountFiat           float64              `json:"amount_fiat"`
	CurrencyFiat         string               `json:"currency_fiat"`
	Status               prisma.PaymentStatus `json:"status"`
	ExpiresAt            time.Time            `json:"expires_at"`

	LastTransactionHash   string `json:"last_transaction_hash,omitempty"`
	Confirmations         *int   `json:"confirmations"`
	RequiredConfirmations int    `json:"required_confirmations"`
}

func Get(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	payment, err := db.PaymentById(id)
	if err != nil {
		return fmt.Errorf("payment not found")
	}

	// check if payment is expired
	// TODO: maybe move this to a cron job to check for expired payments
	if payment.ExpiresAt.Before(time.Now()) && payment.Status != prisma.PaymentStatusPaid {
		if payment.Status != prisma.PaymentStatusExpired {
			_, err := db.UpdatePayment(payment.ID, prisma.PaymentStatusExpired)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}

			payment.Status = prisma.PaymentStatusExpired
		}

		return utils.SendJSON(ctx, fiber.StatusOK, &GetPaymentResponse{
			Id:                   payment.ID,
			WalletAddress:        payment.WalletAddress,
			AmountCrypto:         payment.AmountCrypto,
			CurrencyCryptoName:   payment.CurrencyCrypto,
			CurrencyCryptoSymbol: payment.CurrencyCrypto,
			AmountFiat:           payment.AmountFiat,
			CurrencyFiat:         payment.CurrencyFiat,
			Status:               prisma.PaymentStatusExpired,
			ExpiresAt:            payment.ExpiresAt,
		})
	}

	c := crypto.GetBySymbol(payment.CurrencyCrypto)
	if c == nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	transactions, err := c.ListUnspent(payment.WalletAddress)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	response := &GetPaymentResponse{
		Id:                    payment.ID,
		WalletAddress:         payment.WalletAddress,
		AmountCrypto:          payment.AmountCrypto,
		CurrencyCryptoName:    c.Name(),
		CurrencyCryptoSymbol:  c.Symbol(),
		AmountFiat:            payment.AmountFiat,
		CurrencyFiat:          payment.CurrencyFiat,
		Status:                payment.Status,
		ExpiresAt:             payment.ExpiresAt,
		RequiredConfirmations: c.RequiredConfirmations(),
	}

	if len(transactions) >= 1 {
		var amountWaiting float64
		var totalConfirmedAmount float64
		var totalConfirmations int

		for _, transaction := range transactions {
			if transaction.Confirmations >= c.RequiredConfirmations() {
				totalConfirmedAmount += transaction.Amount
			}

			amountWaiting += transaction.Amount
			totalConfirmations += transaction.Confirmations
		}

		// check if amount is paid
		if totalConfirmedAmount >= payment.AmountCrypto {
			_, err := db.UpdatePayment(payment.ID, prisma.PaymentStatusPaid)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}

			response.Status = prisma.PaymentStatusPaid
		} else {
			averageConfirmations := totalConfirmations / len(transactions)
			response.Confirmations = utils.Ptr(averageConfirmations)
		}

		response.LastTransactionHash = transactions[0].Txid
	}

	return utils.SendJSON(ctx, fiber.StatusOK, response)
}
