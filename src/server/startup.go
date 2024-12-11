package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/wizenheimer/iris/src/internal/api/routes"
	"github.com/wizenheimer/iris/src/internal/config"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/repository/db"
	"github.com/wizenheimer/iris/src/internal/repository/storage"
	"github.com/wizenheimer/iris/src/internal/service/ai"
	"github.com/wizenheimer/iris/src/internal/service/competitor"
	"github.com/wizenheimer/iris/src/internal/service/diff"
	"github.com/wizenheimer/iris/src/internal/service/notification"
	"github.com/wizenheimer/iris/src/internal/service/screenshot"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

func initializer(cfg *config.Config, sqlDb *sql.DB, logger *logger.Logger) (*routes.HandlerContainer, error) {
	logger.Debug("initializing services", zap.Any("service_config", zap.Any("config", cfg.Services)), zap.Any("storage_config", cfg.Storage), zap.Any("database_config", cfg.Database), zap.Any("server_config", cfg.Server), zap.Any("environment_config", cfg.Environment))

	// Initialize HTTP client
	httpClient := &http.Client{}

	// Initialize services
	screenshotService, err := setupScreenshotService(cfg, httpClient, logger)
	if err != nil {
		return nil, err
	}

	notificationService, err := setupNotificationService(cfg, httpClient, logger)
	if err != nil {
		return nil, err
	}

	diffService, err := setupDiffService(cfg, sqlDb, httpClient, logger)
	if err != nil {
		return nil, err
	}

	competitorService, err := setupCompetitorService(sqlDb, logger)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		diffService,
		competitorService,
		notificationService,
		logger,
	)

	return handlers, nil
}

func setupScreenshotService(cfg *config.Config, httpClient interfaces.HTTPClient, logger *logger.Logger) (interfaces.ScreenshotService, error) {
	storageRepo, err := storage.NewR2Storage(
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.Bucket,
		cfg.Storage.Region,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to initialize storage", zap.Error(err))
	}

	screenshotServiceOptions := []screenshot.ScreenshotServiceOption{
		screenshot.WithStorage(storageRepo),
		screenshot.WithHTTPClient(httpClient),
		screenshot.WithKey(cfg.Services.ScreenshotServiceAPIKey),
	}

	return screenshot.NewScreenshotService(
		logger,
		screenshotServiceOptions...,
	)
}

func setupNotificationService(cfg *config.Config, httpClient *http.Client, logger *logger.Logger) (interfaces.NotificationService, error) {
	emailClient, err := notification.NewResendEmailClient(cfg, httpClient, logger)
	if err != nil {
		log.Fatalf("Failed to initialize email client: %v", err)
	}

	templateManager, err := notification.NewTemplateManager(logger)
	if err != nil {
		log.Fatalf("Failed to initialize template manager: %v", err)
	}

	return notification.NewNotificationService(emailClient, templateManager, logger)
}

func setupDiffService(cfg *config.Config, sqlDb *sql.DB, httpClient *http.Client, logger *logger.Logger) (interfaces.DiffService, error) {
	screenshotService, err := setupScreenshotService(cfg, httpClient, logger)
	if err != nil {
		log.Fatalf("Failed to initialize screenshot service: %v", err)
	}

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, httpClient, logger)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	diffRepo, err := db.NewDiffRepository(sqlDb, logger)
	if err != nil {
		log.Fatalf("Failed to initialize diff repository: %v", err)
	}

	return diff.NewDiffService(diffRepo, aiService, screenshotService, logger)
}

func setupCompetitorService(sqlDb *sql.DB, logger *logger.Logger) (interfaces.CompetitorService, error) {
	competitorRepo, err := db.NewCompetitorRepository(sqlDb, logger)
	if err != nil {
		log.Fatalf("Failed to initialize competitor repository: %v", err)
	}

	return competitor.NewCompetitorService(competitorRepo, logger)
}
