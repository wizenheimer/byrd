package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/api/middleware"
	"github.com/wizenheimer/iris/internal/api/routes"
	"github.com/wizenheimer/iris/internal/config"
	"github.com/wizenheimer/iris/pkg/logger"
	"github.com/wizenheimer/iris/pkg/utils/database"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
		return
	}

	// Initialize logger
	loggerConfig := logger.PrepareLoggerConfig(cfg)
	logger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		return
	} else {
		defer logger.Sync()
	}

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

	// Initialize database
	sqlDb, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}

	// Initialize handlers
	handlers, err := initializer(cfg, sqlDb, logger)
	if err != nil {
		logger.Fatal("Failed to initialize handlers", zap.Error(err))
		return
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: middleware.CustomErrorHandler,
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Setup routes
	routes.SetupRoutes(app, handlers)

	// Start server in a goroutine
	serverError := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in server goroutine",
					zap.Any("error", r),
					zap.String("stack", string(debug.Stack())),
				)
				serverError <- fiber.NewError(fiber.StatusInternalServerError, "server panic")
			}
		}()

		if err := app.Listen(cfg.Server.Port); err != nil {
			logger.Error("Server error", zap.Error(err))
			serverError <- err
		}
	}()

	// Setup shutdown handler
	shutdownHandler := NewShutdownHandler(app, sqlDb, logger, ShutdownConfig{
		Timeout:     cfg.Server.ShutdownTimeout,
		MaxAttempts: cfg.Server.ShutdownMaxAttempts,
	})
	shutdownHandler.HandleGracefulShutdown()
}
