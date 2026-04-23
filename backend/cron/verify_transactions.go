package cron

import (
	"context"
	"time"

	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
	"github.com/szerookii/litepay/backend/webhook"
)

func StartTransactionVerifier(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Info().Msg("Starting transaction verifier cron...")
	
	// Run once immediately on start
	verifyPendingPayments()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping transaction verifier cron...")
			return
		case <-ticker.C:
			verifyPendingPayments()
		}
	}
}

func verifyPendingPayments() {
	// 1. Get all pending or confirming payments that haven't expired
	payments, err := db.Client().Payment.Query().
		Where(
			entpayment.StatusIn(entpayment.StatusPENDING, entpayment.StatusCONFIRMING),
			entpayment.ExpiresAtGT(time.Now()),
		).
		All(context.Background())

	if err != nil {
		log.Error().Err(err).Msg("Cron: Failed to query pending payments")
		return
	}

	if len(payments) == 0 {
		log.Debug().Msg("Cron: No pending payments to check")
		return
	}

	log.Debug().Int("count", len(payments)).Msg("Cron: Checking pending payments")

	for _, p := range payments {
		log.Debug().
			Str("id", p.ID.String()).
			Str("symbol", p.CurrencyCrypto).
			Str("addr", p.WalletAddress).
			Msg("Cron: Verifying payment")

		// 2. Find the corresponding crypto driver
		driver, ok := cryptocurrency.GetBySymbol(p.CurrencyCrypto)
		if !ok {
			log.Error().Str("symbol", p.CurrencyCrypto).Msg("Cron: Driver not found for payment")
			continue
		}

		// 3. Check the blockchain status
		solRef := ""
		if p.SolReference != nil {
			solRef = *p.SolReference
		}

		addr := &cryptocurrency.PaymentAddress{
			Address:   p.WalletAddress,
			Index:     uint32(p.AddressIndex),
			Reference: solRef,
		}

		status, err := driver.CheckPayment(addr, p.AmountCrypto)
		if err != nil {
			log.Error().Err(err).Str("id", p.ID.String()).Msg("Cron: Failed to check payment status")
			continue
		}

		// 4. Update status based on confirmations
		if status.ReceivedAmount >= p.AmountCrypto {
			newStatus := entpayment.StatusCONFIRMING
			if status.Confirmations >= driver.RequiredConfirmations() {
				newStatus = entpayment.StatusPAID
				log.Info().
					Str("id", p.ID.String()).
					Str("hash", status.TxHash).
					Float64("amount", status.ReceivedAmount).
					Msg("Cron: Payment fully confirmed!")
			} else {
				log.Info().
					Str("id", p.ID.String()).
					Int("confs", status.Confirmations).
					Int("required", driver.RequiredConfirmations()).
					Msg("Cron: Payment seen, waiting for confirmations...")
			}

			updated, err := db.UpdatePaymentWithHash(p.ID, newStatus, status.TxHash, status.ReceivedAmount)
			if err != nil {
				log.Error().Err(err).Str("id", p.ID.String()).Msg("Cron: Failed to update payment status")
				continue
			}

			user, err := db.UserByID(p.UserID)
			if err != nil {
				log.Error().Err(err).Str("id", p.ID.String()).Msg("Cron: Failed to load user for webhook")
				continue
			}
			if user.WebhookURL != nil && *user.WebhookURL != "" && user.WebhookSecret != nil {
				webhook.Dispatch(*user.WebhookURL, *user.WebhookSecret, updated)
			}
		} else {
			log.Debug().
				Str("id", p.ID.String()).
				Float64("received", status.ReceivedAmount).
				Float64("required", p.AmountCrypto).
				Msg("Cron: No funds received yet")
		}
	}

	// 5. Mark old payments as expired
	n, err := db.Client().Payment.Update().
		Where(
			entpayment.StatusEQ(entpayment.StatusPENDING),
			entpayment.ExpiresAtLT(time.Now()),
		).
		SetStatus(entpayment.StatusEXPIRED).
		Save(context.Background())

	if err != nil {
		log.Error().Err(err).Msg("Cron: Failed to expire old payments")
	} else if n > 0 {
		log.Info().Int("count", n).Msg("Cron: Expired old payments")
	}
}
