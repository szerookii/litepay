package user

import (
	"context"
	"crypto/sha512"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
	"github.com/szerookii/litepay/backend/router/middleware"
	"github.com/szerookii/litepay/backend/utils"
	"golang.org/x/crypto/pbkdf2"
)

type cashoutRequest struct {
	Symbol      string `json:"symbol"      binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type cashoutTx struct {
	FromAddress string  `json:"from_address"`
	TxHash      string  `json:"tx_hash"`
	Amount      float64 `json:"amount"`
}

func Cashout(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString(middleware.UserIDKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	var req cashoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	driver, ok := cryptocurrency.GetBySymbol(req.Symbol)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "unsupported symbol"})
		return
	}
	spendable, ok := driver.(cryptocurrency.Spendable)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cashout not supported for " + req.Symbol})
		return
	}

	masterSeed, err := DeriveMasterSeed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "seed unavailable"})
		return
	}
	defer utils.Zeroize(masterSeed)

	user, err := db.UserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	// Fetch all PAID payments for this user+coin
	payments, err := db.Client().Payment.Query().
		Where(
			entpayment.UserID(uid),
			entpayment.CurrencyCrypto(req.Symbol),
			entpayment.StatusEQ(entpayment.StatusPAID),
		).
		All(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to query payments"})
		return
	}
	if len(payments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no paid payments to cashout"})
		return
	}

	type skipReason struct {
		Address string `json:"address"`
		Reason  string `json:"reason"`
	}
	var results []cashoutTx
	var skipped []skipReason

	for _, p := range payments {
		privKey, err := spendable.DerivePrivKey(masterSeed, uint32(user.AccountIndex), uint32(p.AddressIndex))
		if err != nil {
			skipped = append(skipped, skipReason{p.WalletAddress, "key derivation: " + err.Error()})
			continue
		}

		balance, err := spendable.GetOnChainBalance(p.WalletAddress)
		if err != nil {
			utils.Zeroize(privKey)
			skipped = append(skipped, skipReason{p.WalletAddress, "balance check: " + err.Error()})
			continue
		}
		if balance < 0.000006 {
			utils.Zeroize(privKey)
			skipped = append(skipped, skipReason{p.WalletAddress, fmt.Sprintf("balance too low: %.9f", balance)})
			continue
		}

		txHash, err := spendable.SendFunds(privKey, p.WalletAddress, req.Destination, 0)
		utils.Zeroize(privKey)
		if err != nil {
			skipped = append(skipped, skipReason{p.WalletAddress, "send failed: " + err.Error()})
			continue
		}

		db.UpdatePaymentWithHash(p.ID, entpayment.StatusCASHED_OUT, txHash, balance)

		results = append(results, cashoutTx{
			FromAddress: p.WalletAddress,
			TxHash:      txHash,
			Amount:      balance,
		})
	}

	if len(results) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no funds available to cashout", "details": skipped})
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{"transactions": results})
}

// DeriveMasterSeed converts MASTER_SEED mnemonic → 64-byte seed via BIP39 PBKDF2-HMAC-SHA512.
// Must match the derivation used in payment/post.go.
func DeriveMasterSeed() ([]byte, error) {
	mnemonic := os.Getenv("MASTER_SEED")
	if mnemonic == "" {
		return nil, nil
	}
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
	return seed, nil
}
