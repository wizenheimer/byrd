// ./src/server/startup/services/screenshot.go
package services

import (
	"fmt"

	"github.com/wizenheimer/byrd/src/internal/client"
	"github.com/wizenheimer/byrd/src/internal/config"
	screenshot_repo "github.com/wizenheimer/byrd/src/internal/repository/screenshot"
	screenshot_svc "github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

func SetupScreenshotService(cfg *config.Config, screenshotHTTPClient *client.HTTPClient, logger *logger.Logger) (screenshot_svc.ScreenshotService, error) {
	if logger == nil {
		return nil, fmt.Errorf("can't initialize screenshot service, logger is required")
	}

	logger.Debug("setting up screenshot service", zap.Any("storage_config", cfg.Storage))
	validateStorageConfig(cfg, logger)

	storageRepo, err := createStorageRepository(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize storage", zap.Error(err))
		return nil, err
	}

	validateServiceConfig(cfg, logger)

	screenshotServiceOptions := []screenshot_svc.ScreenshotServiceOption{
		screenshot_svc.WithStorage(storageRepo),
		screenshot_svc.WithHTTPClient(screenshotHTTPClient),
		screenshot_svc.WithKey(cfg.Services.ScreenshotServiceAPIKey),
		screenshot_svc.WithOrigin(cfg.Services.ScreenshotServiceOrigin),
	}

	return screenshot_svc.NewScreenshotService(logger, screenshotServiceOptions...)
}

func validateStorageConfig(cfg *config.Config, logger *logger.Logger) {
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
}

func validateServiceConfig(cfg *config.Config, logger *logger.Logger) {
	if cfg.Services.ScreenshotServiceAPIKey == "" {
		logger.Warn("API key is empty", zap.String("service", "screenshot"))
	}
	if cfg.Services.ScreenshotServiceOrigin == "" {
		logger.Warn("Origin is empty", zap.String("service", "screenshot"))
	}
}

func createStorageRepository(cfg *config.Config, logger *logger.Logger) (screenshot_repo.ScreenshotRepository, error) {
	switch cfg.Storage.Type {
	case "r2":
		return screenshot_repo.NewR2ScreenshotRepo(
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.AccountId,
			logger,
		)
	case "s3":
		return screenshot_repo.NewS3ScreenshotRepo(
			cfg.Storage.BaseEndpoint,
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.Region,
			logger,
		)
	case "local":
		return screenshot_repo.NewLocalScreenshotRepo(cfg.Storage.Bucket, logger)
	default:
		logger.Warn("Unknown storage type, defaulting to local storage", zap.String("type", cfg.Storage.Type))
		return screenshot_repo.NewLocalScreenshotRepo(cfg.Storage.Bucket, logger)
	}
}
