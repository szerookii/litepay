package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	setAuthCookie(c, "", -1)
	c.Status(http.StatusNoContent)
}
