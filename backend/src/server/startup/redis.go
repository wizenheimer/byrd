// ./src/server/startup/redis.go
package startup

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

func SetupRedis(cfg *config.Config, logger *logger.Logger) (*redis.Client, error) {
	if cfg.Workflow.RedisAddr == "" {
		logger.Warn("Redis URL is empty")
	}

	logger.Debug("Setting up workflow service", zap.Any("redis_config", cfg.Workflow))

	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Workflow.RedisAddr,
		Password: cfg.Workflow.RedisPassword,
		DB:       cfg.Workflow.RedisDB,
	})

	_, err := c.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return c, nil
}
