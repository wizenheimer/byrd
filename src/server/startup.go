package main

import (
	"database/sql"
	"fmt"
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
	"github.com/wizenheimer/iris/src/internal/service/url"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

func initializer(cfg *config.Config, sqlDb *sql.DB, logger *logger.Logger) (*routes.HandlerContainer, error) {
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

	urlService, err := setupURLService(sqlDb, logger)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		urlService,
		diffService,
		competitorService,
		notificationService,
		logger,
	)

	return handlers, nil
}

func setupURLService(sqlDb *sql.DB, logger *logger.Logger) (interfaces.URLService, error) {
	logger.Debug("setting up URL service", zap.Any("database", sqlDb.Stats()))

	urlRepo := db.NewURLRepository(sqlDb, logger)

	return url.NewURLService(urlRepo, logger)
}

func setupScreenshotService(cfg *config.Config, screenshotHTTPClient client.HTTPClient, logger *logger.Logger) (interfaces.ScreenshotService, error) {
	logger.Debug("setting up screenshot service", zap.Any("storage_config", cfg.Storage))

	if logger == nil {
		return nil, fmt.Errorf("can't initialize screenshot service, logger is required")
	}

	if cfg.Storage.AccessKey == "" {
		logger.Warn("Access key is empty", zap.String("type", cfg.Storage.Type))
	}

	if cfg.Storage.SecretKey == "" {
		logger.Warn("Secret key is empty", zap.String("type", cfg.Storage.Type))
	}

	if cfg.Storage.Bucket == "" {
		logger.Warn("Bucket is empty", zap.String("type", cfg.Storage.Type))
	}

	if cfg.Storage.AccountId == "" {
		logger.Warn("Account ID is empty", zap.String("type", cfg.Storage.Type))
	}

	if cfg.Storage.Region == "" {
		logger.Warn("Region is empty", zap.String("type", cfg.Storage.Type))
	}

	var storageRepo interfaces.StorageRepository
	var err error
	switch cfg.Storage.Type {
	case "r2":
		storageRepo, err = storage.NewR2Storage(
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.AccountId,
			logger,
		)
	case "local":
		storageRepo, err = storage.NewLocalStorage(cfg.Storage.Bucket, logger)
	case "s3":
		storageRepo, err = storage.NewS3Storage(
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.AccountId,
			"", // session is empty
			cfg.Storage.Region,
			logger,
		)
	default:
		logger.Warn("Unknown storage type, defaulting to local storage", zap.String("type", cfg.Storage.Type))
		storageRepo, err = storage.NewLocalStorage(cfg.Storage.Bucket, logger)
	}

	if err != nil {
		logger.Fatal("Failed to initialize storage", zap.Error(err))
	}

	if cfg.Services.ScreenshotServiceAPIKey == "" {
		logger.Warn("API key is empty", zap.String("service", "screenshot"))
	}

	if cfg.Services.ScreenshotServiceOrigin == "" {
		logger.Warn("Origin is empty", zap.String("service", "screenshot"))
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
	logger.Debug("setting up notification service", zap.Any("notification_config", cfg.Services))

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
	logger.Debug("setting up diff service", zap.Any("database", sqlDb.Stats()), zap.Any("service_config", cfg.Services), zap.Any("storage_config", cfg.Storage))

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
	logger.Debug("setting up competitor service", zap.Any("database", sqlDb.Stats()))

	competitorRepo, err := db.NewCompetitorRepository(sqlDb, logger)
	if err != nil {
		log.Fatalf("Failed to initialize competitor repository: %v", err)
	}

	return competitor.NewCompetitorService(competitorRepo, logger)
}
