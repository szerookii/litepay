package utils

import "github.com/gin-gonic/gin"

func SendJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
