// ./src/internal/api/middleware/core.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	highlightFiber "github.com/highlight/highlight/sdk/highlight-go/middleware/fiber"
	"github.com/wizenheimer/byrd/src/internal/config"
)

func SetupMiddleware(cfg *config.Config, app *fiber.App, rc *RateLimiters, lm *LoggingMiddleware) {
	// Rwecover from panics
	app.Use(recover.New())
	// Highlight middleware for Fiber, only in non-development environments
	if cfg.Environment.EnvProfile != "development" {
		app.Use(highlightFiber.Middleware())
	}
	// Handle shutdown
	app.Use(rejectRequestsDuringShutdown)
	// Log requests
	if cfg.Environment.EnvProfile == "development" {
		app.Use(lm.RequestResponseLogger())
	}
	// Rate limiter - Global rate limiting
	app.Use(rc.GlobalLimiter)
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
