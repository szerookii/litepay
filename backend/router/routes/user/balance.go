package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/router/middleware"
	"github.com/szerookii/litepay/backend/utils"
)

func Balance(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString(middleware.UserIDKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	balance, err := db.UserBalance(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch balance"})
		return
	}

	utils.SendJSON(c, http.StatusOK, balance)
}
