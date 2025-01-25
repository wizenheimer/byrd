// ./src/internal/api/handlers/token.go
package handlers

import "github.com/gofiber/fiber/v2"

func (uh *UserHandler) ValidateClerkToken(c *fiber.Ctx) error {
	// Validation of the clerk token is done by the middleware
	return sendDataResponse(c, fiber.StatusOK, "Clerk token is valid", nil)
}

func (uh *UserHandler) ValidateManagementToken(c *fiber.Ctx) error {
	// Validation of the management token is done by the middleware
	return sendDataResponse(c, fiber.StatusOK, "Token is valid", nil)
}
