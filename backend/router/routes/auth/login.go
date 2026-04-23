package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/utils"
	jwtutil "github.com/szerookii/litepay/backend/utils/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}
	if utils.Validate(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid fields"})
		return
	}

	user, err := db.UserByEmail(req.Email)
	if err != nil {
		// Constant-time response — don't leak whether email exists
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$12$placeholder.hash.to.prevent.timing.attack"), []byte(req.Password))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	setAuthCookie(c, token, 86400)
	utils.SendJSON(c, http.StatusOK, gin.H{"token": token})
}
