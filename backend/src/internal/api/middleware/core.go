// ./src/internal/api/middleware/core.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupMiddleware(app *fiber.App) {
	// Recover from panics
	app.Use(recover.New())
	// Handle shutdown
	app.Use(rejectRequestsDuringShutdown)
	// Log requests
	app.Use(logger.New())
	// Handle Liveness
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			// TODO: Implement liveness probe
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			// TODO: Implement readiness probe
			return true
		},
		ReadinessEndpoint: "/ready",
	},
	))
}
