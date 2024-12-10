package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	log = logger
	return nil
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}
