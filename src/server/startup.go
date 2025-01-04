package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/wizenheimer/iris/src/internal/api/routes"
	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/config"
	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/internal/repository/db"
	"github.com/wizenheimer/iris/src/internal/repository/storage"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/internal/service/ai"
	"github.com/wizenheimer/iris/src/internal/service/competitor"
	"github.com/wizenheimer/iris/src/internal/service/diff"
	"github.com/wizenheimer/iris/src/internal/service/history"
	"github.com/wizenheimer/iris/src/internal/service/page"
	"github.com/wizenheimer/iris/src/internal/service/screenshot"
	"github.com/wizenheimer/iris/src/internal/service/user"
	"github.com/wizenheimer/iris/src/internal/service/workspace"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
	"go.uber.org/zap"
)

func initializer(cfg *config.Config, sqldb *sql.DB, logger *logger.Logger) (*routes.HandlerContainer, error) {
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

	// Initialize validator
	utils.InitializeValidator()

	// Initialize services
	screenshotService, err := setupScreenshotService(cfg, screenshotRateLimitedClient, logger)
	if err != nil {
		return nil, err
	}

	aiService, err := setupAIService(cfg, logger)
	if err != nil {
		return nil, err
	}

	diffService, err := diff.NewDiffService(aiService, logger)
	if err != nil {
		return nil, err
	}

	// Intialize transaction manager
	tm := transaction.NewTxManager(sqldb)

	// Intialize repository
	competitorRepo := db.NewCompetitorRepository(tm, logger)
	workspaceRepo := db.NewWorkspaceRepository(tm, logger)
	userRepo := db.NewUserRepository(tm, logger)
	pageRepo := db.NewPageRepository(tm, logger)
	historyRepo := db.NewPageHistoryRepository(tm, logger)

	// Initialize services
	historyService := history.NewPageHistoryService(historyRepo, screenshotService, diffService, logger)
	pageService := page.NewPageService(pageRepo, historyService, logger)
	competitorService := competitor.NewCompetitorService(competitorRepo, pageService, logger)
	userService := user.NewUserService(userRepo, logger)
	workspaceService := workspace.NewWorkspaceService(workspaceRepo, competitorService, userService, logger)

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		aiService,
		userService,
		workspaceService,
		logger,
	)

	return handlers, nil
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

func setupAIService(cfg *config.Config, logger *logger.Logger) (svc.AIService, error) {
	logger.Debug("setting up AI service", zap.Any("service_config", cfg.Services))

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, logger)
	if err != nil {
		logger.Fatal("Failed to initialize AI service: %v", zap.Any("error", err))
	}

	return aiService, nil
}

func setupDB(cfg *config.Config) (*sql.DB, error) {
	// Prepare connection string
	connString := prepareConnectionString(cfg)

	// initialize db connection
	return sql.Open(cfg.Database.Driver, connString)
}

func prepareConnectionString(cfg *config.Config) string {
	if cfg.Database.ConnectionString != "" {
		return cfg.Database.ConnectionString
	}

	// Determine the environment
	var sslMode string
	switch cfg.Environment.EnvProfile {
	case "development":
		// Set up development environment
		sslMode = "disable"
	case "production":
		// Set up production environment
		sslMode = "require"
	default:
		// Set up default environment
		sslMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		sslMode,
	)
}
