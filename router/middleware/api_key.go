package middleware

import (
	"github.com/gofiber/fiber/v3"
	"os"
	"strings"
)

func APIKey(c fiber.Ctx) error {
	if strings.ReplaceAll(c.Get("authorization"), "Bearer ", "") != os.Getenv("API_KEY") {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}
