package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/ent"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
)

const (
	maxRetries    = 3
	retryDelay    = 2 * time.Second
	clientTimeout = 10 * time.Second
)

var httpClient = &http.Client{Timeout: clientTimeout}

type Event string

const (
	EventPaymentPaid       Event = "payment.paid"
	EventPaymentConfirming Event = "payment.confirming"
	EventPaymentExpired    Event = "payment.expired"
)

type PaymentPayload struct {
	ID              string   `json:"id"`
	Status          string   `json:"status"`
	AmountCrypto    float64  `json:"amount_crypto"`
	CurrencyCrypto  string   `json:"currency_crypto"`
	AmountFiat      float64  `json:"amount_fiat"`
	CurrencyFiat    string   `json:"currency_fiat"`
	ReceivedAmount  *float64 `json:"received_amount"`
	TransactionHash *string  `json:"transaction_hash"`
	WalletAddress   string   `json:"wallet_address"`
	ExpiresAt       string   `json:"expires_at"`
	CreatedAt       string   `json:"created_at"`
}

type Payload struct {
	Event     Event          `json:"event"`
	Payment   PaymentPayload `json:"payment"`
	Timestamp int64          `json:"timestamp"`
}

func sign(secret, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func paymentToPayload(p *ent.Payment) PaymentPayload {
	pp := PaymentPayload{
		ID:             p.ID.String(),
		Status:         string(p.Status),
		AmountCrypto:   p.AmountCrypto,
		CurrencyCrypto: p.CurrencyCrypto,
		AmountFiat:     p.AmountFiat,
		CurrencyFiat:   p.CurrencyFiat,
		WalletAddress:  p.WalletAddress,
		ExpiresAt:      p.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:      p.CreateTime.UTC().Format(time.RFC3339),
	}
	if p.ReceivedAmount != nil {
		pp.ReceivedAmount = p.ReceivedAmount
	}
	if p.TransactionHash != nil {
		pp.TransactionHash = p.TransactionHash
	}
	return pp
}

func eventForStatus(s entpayment.Status) (Event, bool) {
	switch s {
	case entpayment.StatusPAID:
		return EventPaymentPaid, true
	case entpayment.StatusCONFIRMING:
		return EventPaymentConfirming, true
	case entpayment.StatusEXPIRED:
		return EventPaymentExpired, true
	default:
		return "", false
	}
}

// Dispatch sends a webhook for the given payment status.
// webhookURL and webhookSecret come from the user record.
// Fires-and-forgets with up to maxRetries attempts.
func Dispatch(webhookURL, webhookSecret string, p *ent.Payment) {
	event, ok := eventForStatus(p.Status)
	if !ok {
		return
	}

	payload := Payload{
		Event:     event,
		Payment:   paymentToPayload(p),
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Str("payment_id", p.ID.String()).Msg("Webhook: failed to marshal payload")
		return
	}

	sig := sign([]byte(webhookSecret), body)

	go func() {
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err := send(webhookURL, sig, body)
			if err == nil {
				log.Info().
					Str("payment_id", p.ID.String()).
					Str("event", string(event)).
					Str("url", webhookURL).
					Msg("Webhook: delivered")
				return
			}

			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Str("payment_id", p.ID.String()).
				Msg("Webhook: delivery failed")

			if attempt < maxRetries {
				time.Sleep(retryDelay * time.Duration(attempt))
			}
		}
		log.Error().
			Str("payment_id", p.ID.String()).
			Str("url", webhookURL).
			Msg("Webhook: all retries exhausted")
	}()
}

func send(url, signature string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-LitePay-Signature", signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
