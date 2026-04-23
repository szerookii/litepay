package db

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/ent"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
)

// NextPaymentIndex returns the next derivation index for a user+coin pair.
// Must be called inside a transaction so the count is stable.
func NextPaymentIndex(ctx context.Context, tx *ent.Tx, userID uuid.UUID, coin string) (uint32, error) {
	count, err := tx.Payment.Query().
		Where(
			entpayment.UserID(userID),
			entpayment.CurrencyCrypto(strings.ToUpper(coin)),
		).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func CreatePayment(
	tx *ent.Tx,
	userID uuid.UUID,
	walletAddr string,
	solReference *string,
	addressIndex int,
	amountCrypto float64,
	currencyCrypto string,
	amountFiat float64,
	currencyFiat string,
	expiresAt time.Time,
) (*ent.Payment, error) {
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Hour)
	}
	q := tx.Payment.Create().
		SetUserID(userID).
		SetWalletAddress(walletAddr).
		SetAddressIndex(addressIndex).
		SetAmountCrypto(amountCrypto).
		SetCurrencyCrypto(strings.ToUpper(currencyCrypto)).
		SetAmountFiat(amountFiat).
		SetCurrencyFiat(strings.ToUpper(currencyFiat)).
		SetExpiresAt(expiresAt)
	if solReference != nil && *solReference != "" {
		q.SetSolReference(*solReference)
	}
	return q.Save(context.Background())
}

func PaymentById(id string) (*ent.Payment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return Client().Payment.Get(context.Background(), uid)
}

func UpdatePayment(id uuid.UUID, status entpayment.Status) (*ent.Payment, error) {
	return Client().Payment.UpdateOneID(id).
		SetStatus(status).
		Save(context.Background())
}

func UpdatePaymentWithHash(id uuid.UUID, status entpayment.Status, txHash string, receivedAmount float64) (*ent.Payment, error) {
	q := Client().Payment.UpdateOneID(id).SetStatus(status).SetReceivedAmount(receivedAmount)
	if txHash != "" {
		q.SetTransactionHash(txHash)
	}
	return q.Save(context.Background())
}

func DeletePayment(id uuid.UUID) error {
	return Client().Payment.DeleteOneID(id).Exec(context.Background())
}

func UserPayments(userID uuid.UUID) ([]*ent.Payment, error) {
	return Client().Payment.Query().
		Where(entpayment.UserID(userID)).
		Order(ent.Desc(entpayment.FieldCreateTime)).
		All(context.Background())
}

// UserBalance returns the total received amount per coin (PAID payments only).
func UserBalance(userID uuid.UUID) (map[string]float64, error) {
	payments, err := Client().Payment.Query().
		Where(
			entpayment.UserID(userID),
			entpayment.StatusEQ(entpayment.StatusPAID),
		).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	balance := map[string]float64{}
	for _, p := range payments {
		if p.ReceivedAmount != nil {
			balance[p.CurrencyCrypto] += *p.ReceivedAmount
		} else {
			balance[p.CurrencyCrypto] += p.AmountCrypto
		}
	}
	return balance, nil
}
