package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/utils"
	jwtutil "github.com/szerookii/litepay/backend/utils/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func Register(c *gin.Context) {
	if strings.ToLower(os.Getenv("ALLOW_REGISTER")) == "false" {
		c.JSON(http.StatusForbidden, gin.H{"message": "registration is disabled"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}
	if utils.Validate(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid fields"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	apiKey, err := db.GenerateAPIKey()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := db.CreateUser(req.Email, string(hash), apiKey)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"message": "email already in use"})
			return
		}
		log.Error().Err(err).Msg("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create account"})
		return
	}

	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.SendJSON(c, http.StatusCreated, gin.H{"token": token})
}
