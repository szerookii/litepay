package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/utils"
)

func RotateWebhookSecret(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	secret, err := db.RotateWebhookSecret(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to rotate webhook secret"})
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{"webhook_secret": secret})
}
