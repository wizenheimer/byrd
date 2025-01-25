package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

func sendErrorResponse(c *fiber.Ctx, logger *logger.Logger, status int, message string, details ...any) error {
	logger.Debug(message, zap.Int("status", status), zap.Any("details", details))
	return commons.SendErrorResponse(c, status, message, details)
}
