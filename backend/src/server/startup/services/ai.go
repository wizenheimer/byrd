// ./src/server/startup/services/ai.go
package services

import (
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

func SetupAIService(cfg *config.Config, logger *logger.Logger) (ai.AIService, error) {
	logger.Debug("setting up AI service", zap.Any("service_config", cfg.Services))

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, logger)
	if err != nil {
		logger.Fatal("Failed to initialize AI service", zap.Error(err))
		return nil, err
	}

	return aiService, nil
}
