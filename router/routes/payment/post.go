package payment

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/crypto"
	"github.com/szerookii/litepay/db"
	prisma "github.com/szerookii/litepay/prisma/db"
	"github.com/szerookii/litepay/utils"
	"os"
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

	c := crypto.GetBySymbol(req.Symbol)
	if c == nil {
		return fmt.Errorf("crypto with symbol %s not found", req.Symbol)
	}

	info, err := c.Info()
	if err != nil {
		log.Error().Err(err).Msg("failed to get blockchain info")
		return fmt.Errorf("failed to get blockchain info")
	}

	if info.VerificationProgress < 0.95 {
		log.Error().Msg("blockchain is not fully synced")
		return fmt.Errorf("blockchain is not fully synced")
	}

	exchangeRate, err := c.Price(req.Currency)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get price in %s", req.Currency)
		return fmt.Errorf("failed to get price in %s", req.Currency)
	}

	wallets, err := c.ListWallets()
	if err != nil {
		log.Error().Err(err).Msg("failed to list wallets")
		return fmt.Errorf("failed to list wallets")
	}

	var walletLoaded bool
	for _, wallet := range wallets {
		if wallet == os.Getenv("WALLET_NAME") {
			walletLoaded = true
			break
		}
	}

	if !walletLoaded {
		log.Info().Msgf("wallet %s not found, loading...", os.Getenv("WALLET_NAME"))
		if err := c.LoadWallet(os.Getenv("WALLET_NAME")); err != nil {
			return fmt.Errorf("failed to load wallet")
		}
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
