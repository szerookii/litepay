package payment

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/cryptocurrency"
	"github.com/szerookii/litepay/db"
	prisma "github.com/szerookii/litepay/prisma/db"
	"github.com/szerookii/litepay/utils"
	"time"
)

type CreatePaymentRequest struct {
	Symbol   string  `json:"symbol" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type CreatePaymentResponse struct {
	Id             string               `json:"id"`
	WalletAddress  string               `json:"wallet_address"`
	AmountCrypto   float64              `json:"amount_crypto"`
	CurrencyCrypto string               `json:"currency_crypto"`
	AmountFiat     float64              `json:"amount_fiat"`
	CurrencyFiat   string               `json:"currency_fiat"`
	Status         prisma.PaymentStatus `json:"status"`
	ExpiresAt      time.Time            `json:"expires_at"`
}

func Post(ctx fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if utils.Validate(req) != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	c := cryptocurrency.GetBySymbol(req.Symbol)
	if c == nil {
		return fmt.Errorf("cryptocurrency with symbol %s not found", req.Symbol)
	}

	if !c.Synced() {
		return fmt.Errorf("cryptocurrency with symbol %s already synced", req.Symbol)
	}

	exchangeRate, err := c.Price(req.Currency)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get price in %s", req.Currency)
		return fmt.Errorf("failed to get price in %s", req.Currency)
	}

	fees, err := c.EstimateFees()
	if err != nil {
		log.Info().Msg("failed to estimate fees, not using fees")
	}

	cryptoAmount := (req.Amount / exchangeRate) + fees

	address, err := c.GetNewAddress("")
	if err != nil {
		return fmt.Errorf("failed to get new address")
	}

	log.Debug().Msgf("address=%s, amount=%f %s, exchange_rate=%f %s, fees=%f %s, crypto_amount=%f %s", address, req.Amount, req.Currency, exchangeRate, req.Currency, fees, req.Currency, cryptoAmount, c.Symbol())
	payment, err := db.CreatePayment(address, cryptoAmount, c.Symbol(), req.Amount, req.Currency, time.Now().Add(time.Hour))
	if err != nil {
		return fmt.Errorf("failed to create payment")
	}

	return utils.SendJSON(ctx, 200, &CreatePaymentResponse{
		Id:             payment.ID,
		WalletAddress:  payment.WalletAddress,
		AmountCrypto:   payment.AmountCrypto,
		CurrencyCrypto: payment.CurrencyCrypto,
		AmountFiat:     payment.AmountFiat,
		CurrencyFiat:   payment.CurrencyFiat,
		Status:         payment.Status,
		ExpiresAt:      payment.ExpiresAt,
	})
}
