package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type RateLimiters struct {
	// Global base limiter
	GlobalLimiter fiber.Handler
	// Specific operation limiters
	WorkspaceCDLimiter  fiber.Handler
	CompetitorCDLimiter fiber.Handler
	PageCDLimiter       fiber.Handler
	UserCDLimiter       fiber.Handler
}

func NewRateLimiters(cfg *config.Config, logger *logger.Logger) *RateLimiters {
	rLogger := logger.WithFields(
		map[string]interface{}{
			"module": "rate_limiters",
		},
	)

	gl := limiter.New(limiter.Config{
		Max:        cfg.Server.GlobalRequestsPerMinute,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Global rate limit exceeded",
				"details": "Exceeds the authorized number of requests per minute",
			})
		},
	})

	rLogger.Debug("created global rate limiter", zap.Any("max_requests", cfg.Server.GlobalRequestsPerMinute))

	wl := limiter.New(limiter.Config{
		Max:        cfg.Server.WorkspaceCDRequestsPerMinute,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":workspace"
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Workspace rate limit exceeded",
				"details": "Exceeds the authorized number of requests per second",
			})
		},
	})

	rLogger.Debug("created workspace rate limiter", zap.Any("max_requests", cfg.Server.WorkspaceCDRequestsPerMinute))

	cl := limiter.New(limiter.Config{
		Max:        cfg.Server.CompetitorCDRequestsPerSecond,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":competitor"
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Competitor rate limit exceeded",
				"details": "Exceeds the authorized number of requests per second",
			})
		},
	})

	rLogger.Debug("created competitor rate limiter", zap.Any("max_requests", cfg.Server.CompetitorCDRequestsPerSecond))

	pl := limiter.New(limiter.Config{
		Max:        cfg.Server.PageCDRequestsPerSecond,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":page"
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Page rate limit exceeded",
				"details": "Exceeds the authorized number of requests per second",
			})
		},
	})

	rLogger.Debug("created page rate limiter", zap.Any("max_requests", cfg.Server.PageCDRequestsPerSecond))

	ul := limiter.New(limiter.Config{
		Max:        cfg.Server.UserCDRequestsPerSecond,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":user"
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "User rate limit exceeded",
				"details": "Exceeds the authorized number of requests per second",
			})
		},
	})

	rLogger.Debug("created user rate limiter", zap.Any("max_requests", cfg.Server.UserCDRequestsPerSecond))

	return &RateLimiters{
		// Global base limiter
		GlobalLimiter: gl,

		WorkspaceCDLimiter:  wl,
		CompetitorCDLimiter: cl,
		PageCDLimiter:       pl,
		UserCDLimiter:       ul,
	}
}
