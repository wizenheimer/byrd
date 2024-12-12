package main

import (
	"database/sql"
	"log"

	"github.com/wizenheimer/iris/src/internal/api/routes"
	"github.com/wizenheimer/iris/src/internal/client"
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

	// Initialize HTTP client for services
	screenshotClientOpts := []client.ClientOption{
		client.WithLogger(logger),
		client.WithAuth(client.BearerAuth{
			Token: cfg.Services.ScreenshotServiceAPIKey,
		}),
	}
	screenshotHttpClient, err := client.NewClient(screenshotClientOpts...)
	if err != nil {
		return nil, err
	}

	// Create a new rate limited client
	// This client will be used to make requests to the screenshot service
	screenshotRateLimitedClient := client.NewRateLimitedClient(screenshotHttpClient, cfg.Services.ScreenshotServiceQPS)

	aiClientOpts := []client.ClientOption{
		client.WithLogger(logger),
		client.WithAuth(client.BearerAuth{
			Token: cfg.Services.OpenAIKey,
		}),
	}

	aiHttpClient, err := client.NewClient(aiClientOpts...)
	if err != nil {
		return nil, err
	}

	commonClientOpts := []client.ClientOption{
		client.WithLogger(logger),
	}

	commonHttpClient, err := client.NewClient(commonClientOpts...)
	if err != nil {
		return nil, err
	}

	// Initialize services
	screenshotService, err := setupScreenshotService(cfg, screenshotRateLimitedClient, logger)
	if err != nil {
		return nil, err
	}

	notificationService, err := setupNotificationService(cfg, commonHttpClient, logger)
	if err != nil {
		return nil, err
	}

	diffService, err := setupDiffService(cfg, sqlDb, screenshotRateLimitedClient, aiHttpClient, logger)
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

func setupScreenshotService(cfg *config.Config, screenshotHTTPClient client.HTTPClient, logger *logger.Logger) (interfaces.ScreenshotService, error) {
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

	// Create screenshot service option
	screenshotServiceOptions := []screenshot.ScreenshotServiceOption{
		screenshot.WithStorage(storageRepo),
		screenshot.WithHTTPClient(screenshotHTTPClient),
		screenshot.WithKey(cfg.Services.ScreenshotServiceAPIKey),
		screenshot.WithOrigin(cfg.Services.ScreenshotServiceOrigin),
	}

	// Create a new screenshot service
	return screenshot.NewScreenshotService(
		logger,
		screenshotServiceOptions...,
	)
}

func setupNotificationService(cfg *config.Config, httpClient client.HTTPClient, logger *logger.Logger) (interfaces.NotificationService, error) {
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

func setupDiffService(cfg *config.Config, sqlDb *sql.DB, screenshotHTTPClient, aiHTTPClient client.HTTPClient, logger *logger.Logger) (interfaces.DiffService, error) {
	screenshotService, err := setupScreenshotService(cfg, screenshotHTTPClient, logger)
	if err != nil {
		log.Fatalf("Failed to initialize screenshot service: %v", err)
	}

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, aiHTTPClient, logger)
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
