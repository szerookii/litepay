package db

import (
	"context"
	"github.com/szerookii/litepay/prisma/db"
	"strings"
	"time"
)

func CreatePayment(walletAddr string, amountCrypto float64, currencyCrypto string, amountFiat float64, currencyFiat string, expiresAt time.Time) (*db.PaymentModel, error) {
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Hour)
	}

	return Client().Payment.CreateOne(db.Payment.WalletAddress.Set(walletAddr), db.Payment.AmountCrypto.Set(amountCrypto), db.Payment.CurrencyCrypto.Set(strings.ToUpper(currencyCrypto)), db.Payment.AmountFiat.Set(amountFiat), db.Payment.CurrencyFiat.Set(strings.ToUpper(currencyFiat)), db.Payment.ExpiresAt.Set(expiresAt)).Exec(context.Background())
}

func PaymentById(id string) (*db.PaymentModel, error) {
	return Client().Payment.FindUnique(db.Payment.ID.Equals(id)).Exec(context.Background())
}

func UpdatePayment(id string, status db.PaymentStatus) (*db.PaymentModel, error) {
	return Client().Payment.FindUnique(db.Payment.ID.Equals(id)).Update(db.Payment.Status.Set(status)).Exec(context.Background())
}

func DeletePayment(id string) error {
	_, err := Client().Payment.FindUnique(db.Payment.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}
