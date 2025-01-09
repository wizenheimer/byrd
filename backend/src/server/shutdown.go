package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// ShutdownConfig holds configuration for graceful shutdown
type ShutdownConfig struct {
	Timeout     time.Duration
	MaxAttempts int
}

// ShutdownHandler encapsulates shutdown logic
type ShutdownHandler struct {
	app      *fiber.App
	db       *sql.DB
	logger   *logger.Logger
	config   ShutdownConfig
	shutdown chan os.Signal
}

// NewShutdownHandler creates a new shutdown handler
func NewShutdownHandler(app *fiber.App, db *sql.DB, logger *logger.Logger, config ShutdownConfig) *ShutdownHandler {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	return &ShutdownHandler{
		app:      app,
		db:       db,
		logger:   logger,
		config:   config,
		shutdown: shutdown,
	}
}

// HandleGracefulShutdown coordinates the shutdown process
func (h *ShutdownHandler) HandleGracefulShutdown() {
	// Wait for shutdown signal
	<-h.shutdown
	h.logger.Info("Shutdown signal received")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), h.config.Timeout)
	defer cancel()

	// Set shutdown flag to reject new requests
	middleware.IsShuttingDown.Store(true)
	h.logger.Info("Server is no longer accepting new requests")

	// Perform cleanup with retry mechanism
	for attempt := 1; attempt <= h.config.MaxAttempts; attempt++ {
		if err := h.performCleanup(ctx); err != nil {
			if attempt == h.config.MaxAttempts {
				h.logger.Error("Final cleanup attempt failed",
					zap.Error(err),
					zap.Int("attempt", attempt),
				)
				break
			}
			h.logger.Warn("Cleanup attempt failed, retrying",
				zap.Error(err),
				zap.Int("attempt", attempt),
			)
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}
		h.logger.Info("Cleanup successful",
			zap.Int("attempt", attempt),
		)
		break
	}

	// Ensure logs are flushed before exiting
	if err := h.logger.Sync(); err != nil {
		log.Fatalf("Failed to sync logger: %v", err)
	}
}

// performCleanup handles the actual cleanup tasks
func (h *ShutdownHandler) performCleanup(ctx context.Context) error {
	cleanup := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("Panic during cleanup",
					zap.Any("error", r),
					zap.String("stack", string(debug.Stack())),
				)
				cleanup <- fiber.NewError(fiber.StatusInternalServerError, "cleanup panic")
			}
		}()

		// Shutdown the server
		if err := h.app.Shutdown(); err != nil {
			cleanup <- fmt.Errorf("server shutdown error: %w", err)
			return
		}

		// Close database connections
		if err := h.db.Close(); err != nil {
			cleanup <- fmt.Errorf("database closure error: %w", err)
			return
		}

		cleanup <- nil
	}()

	// Wait for either cleanup completion or timeout
	select {
	case err := <-cleanup:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Helper function to check if we're shutting down
func IsShuttingDown() bool {
	return middleware.IsShuttingDown.Load()
}
