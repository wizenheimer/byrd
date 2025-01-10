// ./src/internal/api/middleware/shutdown.go
package middleware

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

// Global flag for server status
var IsShuttingDown atomic.Bool

// Middleware to reject new requests during shutdown
func rejectRequestsDuringShutdown(c *fiber.Ctx) error {
	if IsShuttingDown.Load() {
		return fiber.NewError(
			fiber.StatusServiceUnavailable,
			"server is shutting down",
		)
	}
	return c.Next()
}
