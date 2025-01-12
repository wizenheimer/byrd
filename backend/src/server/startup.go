// ./src/server/startup.go
package main

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/client"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"

	// ----- clf imports -----

	// ----- repo imports -----
	competitor_repo "github.com/wizenheimer/byrd/src/internal/repository/competitor"
	history_repo "github.com/wizenheimer/byrd/src/internal/repository/history"
	page_repo "github.com/wizenheimer/byrd/src/internal/repository/page"
	screenshot_repo "github.com/wizenheimer/byrd/src/internal/repository/screenshot"
	user_repo "github.com/wizenheimer/byrd/src/internal/repository/user"
	workspace_repo "github.com/wizenheimer/byrd/src/internal/repository/workspace"

	// ----- service imports -----
	ai_svc "github.com/wizenheimer/byrd/src/internal/service/ai"
	competitor_svc "github.com/wizenheimer/byrd/src/internal/service/competitor"
	diff_svc "github.com/wizenheimer/byrd/src/internal/service/diff"
	history_svc "github.com/wizenheimer/byrd/src/internal/service/history"
	page_svc "github.com/wizenheimer/byrd/src/internal/service/page"
	screenshot_svc "github.com/wizenheimer/byrd/src/internal/service/screenshot"
	user_svc "github.com/wizenheimer/byrd/src/internal/service/user"
	workspace_svc "github.com/wizenheimer/byrd/src/internal/service/workspace"
)

func initializer(cfg *config.Config, tm *transaction.TxManager, logger *logger.Logger) (*routes.HandlerContainer, workspace_svc.WorkspaceService, error) {
	// Initialize HTTP client for services
	screenshotClientOpts := []client.ClientOption{
		client.WithLogger(logger),
		client.WithAuth(client.BearerAuth{
			Token: cfg.Services.ScreenshotServiceAPIKey,
		}),
	}
	screenshotHttpClient, err := client.NewClient(screenshotClientOpts...)
	if err != nil {
		return nil, nil, err
	}

	// Create a new rate limited client
	// This client will be used to make requests to the screenshot service
	screenshotRateLimitedClient := client.NewRateLimitedClient(screenshotHttpClient, cfg.Services.ScreenshotServiceQPS)

	// Initialize validator
	utils.InitializeValidator()

	// Initialize services
	screenshotService, err := setupScreenshotService(cfg, screenshotRateLimitedClient, logger)
	if err != nil {
		return nil, nil, err
	}

	aiService, err := setupAIService(cfg, logger)
	if err != nil {
		return nil, nil, err
	}

	diffService, err := diff_svc.NewDiffService(aiService, logger)
	if err != nil {
		return nil, nil, err
	}

	// Intialize repository
	// Repositories are responsible for running transactions
	competitorRepo := competitor_repo.NewCompetitorRepository(tm, logger)
	workspaceRepo := workspace_repo.NewWorkspaceRepository(tm, logger)
	userRepo := user_repo.NewUserRepository(tm, logger)
	pageRepo := page_repo.NewPageRepository(tm, logger)
	historyRepo := history_repo.NewPageHistoryRepository(tm, logger)

	// Initialize services
	// Services are responsible for setting transaction boundaries
	historyService := history_svc.NewPageHistoryService(historyRepo, screenshotService, diffService, logger)
	pageService := page_svc.NewPageService(pageRepo, historyService, logger)
	competitorService := competitor_svc.NewCompetitorService(competitorRepo, pageService, tm, logger)
	userService := user_svc.NewUserService(userRepo, logger)
	workspaceService := workspace_svc.NewWorkspaceService(workspaceRepo, competitorService, userService, tm, logger)

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		aiService,
		userService,
		workspaceService,
		tm,
		logger,
	)

	return handlers, workspaceService, nil
}

func setupScreenshotService(cfg *config.Config, screenshotHTTPClient client.HTTPClient, logger *logger.Logger) (screenshot_svc.ScreenshotService, error) {
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

	var storageRepo screenshot_repo.ScreenshotRepository
	var err error
	switch cfg.Storage.Type {
	case "r2":
		storageRepo, err = screenshot_repo.NewR2ScreenshotRepo(
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.AccountId,
			logger,
		)
	case "local":
		storageRepo, err = screenshot_repo.NewLocalScreenshotRepo(cfg.Storage.Bucket, logger)
	case "s3":
		storageRepo, err = screenshot_repo.NewS3ScreenshotRepo(
			cfg.Storage.BaseEndpoint,
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.Region,
			logger,
		)
	default:
		logger.Warn("Unknown storage type, defaulting to local storage", zap.String("type", cfg.Storage.Type))
		storageRepo, err = screenshot_repo.NewLocalScreenshotRepo(cfg.Storage.Bucket, logger)
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
	screenshotServiceOptions := []screenshot_svc.ScreenshotServiceOption{
		screenshot_svc.WithStorage(storageRepo),
		screenshot_svc.WithHTTPClient(screenshotHTTPClient),
		screenshot_svc.WithKey(cfg.Services.ScreenshotServiceAPIKey),
		screenshot_svc.WithOrigin(cfg.Services.ScreenshotServiceOrigin),
	}

	// Create a new screenshot service
	return screenshot_svc.NewScreenshotService(
		logger,
		screenshotServiceOptions...,
	)
}

func setupAIService(cfg *config.Config, logger *logger.Logger) (ai_svc.AIService, error) {
	logger.Debug("setting up AI service", zap.Any("service_config", cfg.Services))

	aiService, err := ai_svc.NewOpenAIService(cfg.Services.OpenAIKey, logger)
	if err != nil {
		logger.Fatal("Failed to initialize AI service: %v", zap.Any("error", err))
	}

	return aiService, nil
}

func setupDB(cfg *config.Config) (*pgxpool.Pool, error) {
	// Prepare connection string
	connString := prepareConnectionString(cfg)

	// Create connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// You can configure pool settings here if needed
	// poolConfig.MaxConns = 10
	// poolConfig.MinConns = 1
	// poolConfig.MaxConnLifetime = time.Hour
	// poolConfig.MaxConnIdleTime = time.Minute * 30

	// Initialize the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return pool, nil
}

func prepareConnectionString(cfg *config.Config) string {
	if cfg.Database.ConnectionString != "" {
		return cfg.Database.ConnectionString
	}

	// Determine the environment
	var sslMode string
	switch cfg.Environment.EnvProfile {
	case "development":
		sslMode = "disable"
	case "production":
		sslMode = "verify-full"
	default:
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
