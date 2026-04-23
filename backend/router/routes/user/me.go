package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/ent"
	"github.com/szerookii/litepay/backend/router/middleware"
	"github.com/szerookii/litepay/backend/utils"
)

type meResponse struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	APIKey       string   `json:"api_key"`
	WebhookURL   *string  `json:"webhook_url"`
	AccountIndex int      `json:"account_index"`
	SupportedCoins []string `json:"supported_coins"`
}

func getUser(c *gin.Context) (*ent.User, error) {
	uid, err := uuid.Parse(c.GetString(middleware.UserIDKey))
	if err != nil {
		return nil, err
	}
	return db.UserByID(uid)
}

func Me(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	coins := make([]string, 0)
	for _, b := range cryptocurrency.All() {
		coins = append(coins, b.Symbol())
	}

	utils.SendJSON(c, http.StatusOK, &meResponse{
		ID:             u.ID.String(),
		Email:          u.Email,
		APIKey:         u.APIKey,
		WebhookURL:     u.WebhookURL,
		AccountIndex:   u.AccountIndex,
		SupportedCoins: coins,
	})
}
