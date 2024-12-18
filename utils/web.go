package utils

import (
	"github.com/gofiber/fiber/v3"
)

func SendJSON(c fiber.Ctx, status int, data interface{}) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	if err := c.SendStatus(status); err != nil {
		return err
	}

	return c.JSON(data)
}
