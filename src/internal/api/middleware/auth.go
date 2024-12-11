package middleware


import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	// Validate token
	// Implementation here

	return c.Next()
}