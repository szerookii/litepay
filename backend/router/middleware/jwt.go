package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtutil "github.com/szerookii/litepay/backend/utils/jwt"
)

const UserIDKey = "user_id"

func JWT(c *gin.Context) {
	var tokenStr string

	if cookie, err := c.Cookie("token"); err == nil {
		tokenStr = cookie
	} else if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		tokenStr = strings.TrimPrefix(h, "Bearer ")
	}

	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	claims, err := jwtutil.ValidateToken(tokenStr)
	if err != nil {
		if errors.Is(err, jwtutil.ErrExpiredToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token expired"})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	c.Set(UserIDKey, claims.UserID)
	c.Next()
}
