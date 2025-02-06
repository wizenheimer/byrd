// ./src/server/main.go
package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/highlight/highlight/sdk/highlight-go"
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/server/shutdown"
	"github.com/wizenheimer/byrd/src/server/startup"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration", zap.Error(err))
		return
	}

	// Highlight is only started in non-development environments
	if cfg.Environment.EnvProfile != "development" {
		highlight.SetProjectID(cfg.Services.HightlightProjectID)
		highlight.Start(
			highlight.WithServiceName(cfg.ServiceName),
			highlight.WithServiceVersion("git-sha"),
		)
		defer highlight.Stop()
	}

	// Initialize Clerk with the secret key
	clerk.SetKey(cfg.Services.ClerkAPIKey)

	// Initialize logger
	loggerConfig := logger.PrepareLoggerConfig(cfg)
	logger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
		return
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}()

	// Recover from panics in main
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic in main",
				zap.Any("error", r),
				zap.String("stack", string(debug.Stack())),
			)
			os.Exit(1)
		}
	}()

	// ErrorRecorder is initialized with the logger and service name
	// The ErrorRecorder is used to record errors in the application
	errorRecorder := recorder.NewErrorRecorder(
		logger,
		cfg.Environment.EnvProfile == "development",
		cfg.ServiceName,
	)

	// Initialize handlers using the new modular initializer
	handlers, rm, am, err := startup.Initialize(context.Background(), cfg, logger, errorRecorder)
	if err != nil {
		logger.Fatal("Failed to initialize handlers", zap.Error(err))
		return
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: middleware.CustomErrorHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// Initialize rate limiters
	ratelimiters := middleware.NewRateLimiters(cfg)

	// Initialize logging middleware
	logSkipPaths := []string{"/health", "/metrics"}
	lm, err := middleware.NewLoggingMiddleware(logger, logSkipPaths)
	if err != nil {
		logger.Fatal("Failed to initialize logging middleware", zap.Error(err))
		return
	}

	// Setup middleware with rate limiters
	middleware.SetupMiddleware(cfg, app, ratelimiters, lm)

	// Setup routes with handlers and rate limiters
	routes.SetupRoutes(app, handlers, ratelimiters, am, rm)

	logger.Info("server started", zap.Any("port", cfg.Server.Port))

	// Start server in a goroutine
	serverError := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errorRecorder.RecordError(context.TODO(), r.(error), zap.String("stack", string(debug.Stack())))
				serverError <- fiber.NewError(fiber.StatusInternalServerError, "server panic")
			}
		}()

		if err := app.Listen(cfg.Server.Port); err != nil {
			logger.Error("Server error", zap.Error(err))
			serverError <- err
		}
	}()

	// Setup shutdown handler
	shutdownHandler := shutdown.NewShutdownHandler(app, cfg.Server.ShutdownTimeout, cfg.Server.ShutdownMaxAttempts, logger)
	shutdownHandler.HandleGracefulShutdown()
}
