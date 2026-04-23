package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setAuthCookie(c *gin.Context, token string, maxAge int) {
	secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
