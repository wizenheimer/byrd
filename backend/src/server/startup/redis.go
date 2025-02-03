// ./src/server/startup/redis.go
package startup

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

func SetupRedis(cfg *config.Config, logger *logger.Logger) (*redis.Client, error) {
	var opts *redis.Options
	var err error
	if cfg.Workflow.RedisURL != "" {
		opts, err = redis.ParseURL(cfg.Workflow.RedisURL)
		if err != nil {
			return nil, err
		}
	} else {
		opts = &redis.Options{
			Addr:     cfg.Workflow.RedisAddr,
			Password: cfg.Workflow.RedisPassword,
			DB:       cfg.Workflow.RedisDB,
		}
	}
	c := redis.NewClient(opts)
	_, err = c.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return c, nil
}
