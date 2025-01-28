package logger

import (
	"fmt"
)

// PosthogLogger is an adapter that implements the posthog.Logger interface
// using our zap-based Logger
type PosthogLogger struct {
	logger *Logger
}

// NewPosthogLogger creates a new PostHog logger adapter
func NewPosthogLogger(logger *Logger) *PosthogLogger {
	return &PosthogLogger{
		logger: logger,
	}
}

// Logf implements the posthog.Logger interface for info-level logging
func (l *PosthogLogger) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Info(msg)
}

// Errorf implements the posthog.Logger interface for error-level logging
func (l *PosthogLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Error(msg)
}

// You can also add a constructor function in your main logger package:
func (l *Logger) AsPosthogLogger() *PosthogLogger {
	return NewPosthogLogger(l)
}
