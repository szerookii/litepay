package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/utils"
)

type updateWebhookRequest struct {
	WebhookURL *string `json:"webhook_url"`
}

func UpdateWallets(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	var req updateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	updated, err := db.UpdateUserWebhook(u.ID, req.WebhookURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update webhook"})
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{
		"webhook_url": updated.WebhookURL,
	})
}
