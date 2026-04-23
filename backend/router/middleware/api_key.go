package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/db"
)

const AuthedUserKey = "authed_user"

func APIKey(c *gin.Context) {
	key := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if key == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	user, err := db.UserByAPIKey(key)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Set(AuthedUserKey, user)
	c.Next()
}
