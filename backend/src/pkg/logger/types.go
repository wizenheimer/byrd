package logger

import "go.uber.org/zap"

// Logger is a wrapper around zap.Logger
type Logger struct {
	log   *zap.Logger
	level LogLevel
}

// Config holds the logger configuration
type LoggerConfig struct {
	Level         LogLevel
	Development   bool
	Environment   string
	OutputPaths   []string // e.g., "stdout", "/var/log/app.log"
	ErrorPaths    []string // paths for error output
	ServiceName   string
	InitialFields map[string]interface{}
	LogDir        string // Directory to store log files
}

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)
