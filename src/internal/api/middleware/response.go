package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/commons"
)

func sendErrorResponse(c *fiber.Ctx, status int, message string, details ...any) error {
	return commons.SendErrorResponse(c, status, message, details)
}
