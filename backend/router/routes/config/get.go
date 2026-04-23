package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/utils"
)

type ConfigResponse struct {
	AllowRegister bool   `json:"allow_register"`
	Version       string `json:"version"`
}

func Get(c *gin.Context) {
	allowRegister := strings.ToLower(os.Getenv("ALLOW_REGISTER")) != "false"
	utils.SendJSON(c, http.StatusOK, &ConfigResponse{
		AllowRegister: allowRegister,
		Version:       utils.Version,
	})
}
