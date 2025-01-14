package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// cronLogger is a logger for the cron scheduler
type cronLogger struct {
	logger *logger.Logger
}

func (c *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Info(msg, zap.Any("keysAndValues", keysAndValues))
}

func (c *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.Error(msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

func NewCronLogger(logger *logger.Logger) cron.Logger {
	return &cronLogger{
		logger: logger,
	}
}
