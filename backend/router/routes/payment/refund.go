package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	entpayment "github.com/szerookii/litepay/backend/ent/payment"
	"github.com/szerookii/litepay/backend/router/middleware"
	"github.com/szerookii/litepay/backend/secrets"
	"github.com/szerookii/litepay/backend/utils"
)

func Refund(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString(middleware.UserIDKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	paymentID := c.Param("id")
	p, err := db.PaymentById(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "payment not found"})
		return
	}
	if p.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	if p.Status != entpayment.StatusPAID {
		if p.Status == entpayment.StatusREFUNDED {
			c.JSON(http.StatusBadRequest, gin.H{"message": "payment already refunded"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "only PAID payments can be refunded"})
		}
		return
	}
	if p.TransactionHash == nil || *p.TransactionHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no transaction hash on this payment"})
		return
	}

	driver, ok := cryptocurrency.GetBySymbol(p.CurrencyCrypto)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "unsupported symbol"})
		return
	}
	spendable, ok := driver.(cryptocurrency.Spendable)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "refund not supported for " + p.CurrencyCrypto})
		return
	}

	senderAddr, err := spendable.GetSender(*p.TransactionHash, p.WalletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not determine sender: " + err.Error()})
		return
	}

	u, err := db.UserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	masterSeed, err := secrets.DeriveMasterSeed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "seed unavailable"})
		return
	}
	defer utils.Zeroize(masterSeed)

	privKey, err := spendable.DerivePrivKey(masterSeed, uint32(u.AccountIndex), uint32(p.AddressIndex))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "key derivation failed"})
		return
	}
	defer utils.Zeroize(privKey)

	balance, err := spendable.GetOnChainBalance(p.WalletAddress)
	if err != nil || balance < 0.000006 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no funds available to refund"})
		return
	}

	txHash, err := spendable.SendFunds(privKey, p.WalletAddress, senderAddr, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "broadcast failed: " + err.Error()})
		return
	}

	if _, err := db.UpdatePaymentWithHash(p.ID, entpayment.StatusREFUNDED, txHash, balance); err != nil {
		// Log but don't fail — tx already broadcast
		c.JSON(http.StatusOK, gin.H{
			"tx_hash":    txHash,
			"to":         senderAddr,
			"amount":     balance,
			"payment_id": p.ID,
			"warning":    "tx broadcast but status update failed",
		})
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{
		"tx_hash":    txHash,
		"to":         senderAddr,
		"amount":     balance,
		"payment_id": p.ID,
	})
}
