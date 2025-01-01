package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/iris/src/internal/api/routes"
	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/config"
	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/db"
	"github.com/wizenheimer/iris/src/internal/repository/storage"
	"github.com/wizenheimer/iris/src/internal/service/ai"
	"github.com/wizenheimer/iris/src/internal/service/alert"
	"github.com/wizenheimer/iris/src/internal/service/competitor"
	"github.com/wizenheimer/iris/src/internal/service/diff"
	"github.com/wizenheimer/iris/src/internal/service/executor"
	"github.com/wizenheimer/iris/src/internal/service/notification"
	"github.com/wizenheimer/iris/src/internal/service/screenshot"
	"github.com/wizenheimer/iris/src/internal/service/url"
	"github.com/wizenheimer/iris/src/internal/service/workflow"
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

	aiService, err := setupAIService(cfg, logger)
	if err != nil {
		return nil, err
	}

	diffService, err := setupDiffService(cfg, sqlDb, aiService, logger)
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

	workflowService, err := setupWorkflowService(cfg, screenshotService, diffService, urlService, logger)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		urlService,
		aiService,
		diffService,
		competitorService,
		notificationService,
		workflowService,
		logger,
	)

	return handlers, nil
}

func setupWorkflowService(cfg *config.Config, screenshotService svc.ScreenshotService, diffService svc.DiffService, urlService svc.URLService, logger *logger.Logger) (svc.WorkflowService, error) {
	logger.Debug("setting up workflow service", zap.Any("workflow_config", cfg.Workflow))

	if cfg.Workflow.RedisAddr == "" {
		logger.Warn("Redis URL is empty")
	}

	// Create a new redis client
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Workflow.RedisAddr,
			Password: cfg.Workflow.RedisPassword,
			DB:       cfg.Workflow.RedisDB,
		},
	)

	// Create a new workflow repository
	workflowRepo, err := db.NewWorkflowRepository(redisClient, logger)
	if err != nil {
		return nil, err
	}

	// Create a new workflow alert client
	clientConfig := core_models.DefaultSlackConfig()
	clientConfig.Token = cfg.Workflow.SlackAlertToken
	clientConfig.ChannelID = cfg.Workflow.SlackWorkflowChannelId

	var alertClient clf.AlertClient
	if cfg.Environment.EnvProfile == "development" {
		logger.Debug("Using local workflow alert client")
		alertClient = alert.NewLocalWorkflowClient(clientConfig, logger)
	} else {
		slackWorkflowClient, err := alert.NewSlackAlertClient(clientConfig, logger)
		if err != nil {
			return nil, err
		}
		alertClient = slackWorkflowClient
	}

	workflowAlertClient, err := alert.NewWorkflowAlertClient(alertClient, logger)
	if err != nil {
		return nil, err
	}

	screenshotTaskExecutor, err := executor.NewScreenshotTaskExecutor(urlService, screenshotService, diffService, logger)
	if err != nil {
		return nil, err
	}

	screenshotWorkflowExecutor, err := executor.NewWorkflowExecutor(core_models.ScreenshotWorkflowType, workflowRepo, workflowAlertClient, screenshotTaskExecutor, logger)
	if err != nil {
		return nil, err
	}

	workflowService, err := workflow.NewWorkflowService(logger, workflowRepo, screenshotWorkflowExecutor, screenshotWorkflowExecutor)
	if err != nil {
		return nil, err
	}

	return workflowService, nil
}

func setupURLService(sqlDb *sql.DB, logger *logger.Logger) (svc.URLService, error) {
	logger.Debug("setting up URL service", zap.Any("database", sqlDb.Stats()))

	urlRepo := db.NewURLRepository(sqlDb, logger)

	return url.NewURLService(urlRepo, logger)
}

func setupScreenshotService(cfg *config.Config, screenshotHTTPClient clf.HTTPClient, logger *logger.Logger) (svc.ScreenshotService, error) {
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

	var storageRepo repo.ScreenshotRepository
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
			cfg.Storage.BaseEndpoint,
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
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

func setupNotificationService(cfg *config.Config, httpClient clf.HTTPClient, logger *logger.Logger) (svc.NotificationService, error) {
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

func setupAIService(cfg *config.Config, logger *logger.Logger) (svc.AIService, error) {
	logger.Debug("setting up AI service", zap.Any("service_config", cfg.Services))

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, logger)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	return aiService, nil
}

func setupDiffService(cfg *config.Config, sqlDb *sql.DB, aiService svc.AIService, logger *logger.Logger) (svc.DiffService, error) {
	logger.Debug("setting up diff service", zap.Any("database", sqlDb.Stats()), zap.Any("service_config", cfg.Services), zap.Any("storage_config", cfg.Storage))

	diffRepo, err := db.NewDiffRepository(sqlDb, logger)
	if err != nil {
		log.Fatalf("Failed to initialize diff repository: %v", err)
	}

	return diff.NewDiffService(diffRepo, aiService, logger)
}

func setupCompetitorService(sqlDb *sql.DB, logger *logger.Logger) (svc.CompetitorService, error) {
	logger.Debug("setting up competitor service", zap.Any("database", sqlDb.Stats()))

	competitorRepo, err := db.NewCompetitorRepository(sqlDb, logger)
	if err != nil {
		log.Fatalf("Failed to initialize competitor repository: %v", err)
	}

	return competitor.NewCompetitorService(competitorRepo, logger)
}
