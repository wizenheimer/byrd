// ./src/internal/api/middleware/core.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wizenheimer/byrd/src/internal/config"
)

func SetupMiddleware(cfg *config.Config, app *fiber.App) {
	// Recover from panics
	app.Use(recover.New())
	// Handle shutdown
	app.Use(rejectRequestsDuringShutdown)
	// Log requests
	app.Use(logger.New())
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
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
