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
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	claims, err := jwtutil.ValidateToken(strings.TrimPrefix(header, "Bearer "))
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
