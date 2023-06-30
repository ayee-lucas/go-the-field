package middleware

import (
	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func EnsureAuth(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	_, err := auth.GetSession(sessionId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get session", "message": err.Error()})
	}

	return c.Next()

}
