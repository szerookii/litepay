package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/utils"
)

func RegenerateKey(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}
	newKey, err := db.RegenerateAPIKey(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to regenerate key"})
		return
	}
	utils.SendJSON(c, http.StatusOK, gin.H{"api_key": newKey})
}
