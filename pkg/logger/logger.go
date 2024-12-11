package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wizenheimer/iris/internal/config"
	"github.com/wizenheimer/iris/pkg/utils/parser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultLoggerConfig returns a default logger configuration
func PrepareLoggerConfig(cfg *config.Config) LoggerConfig {
	return LoggerConfig{
		Level:         LogLevel(cfg.Environment.EnvLogLevel),
		Development:   cfg.Environment.EnvProfile == "development",
		Environment:   cfg.Environment.EnvProfile,
		OutputPaths:   []string{"stdout"},
		ErrorPaths:    []string{"stderr"},
		ServiceName:   "service",
		InitialFields: make(map[string]interface{}),
		LogDir:        "logs", // Default log directory
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(cfg LoggerConfig) (*Logger, error) {
	// Handle development mode specific configuration
	if cfg.Development {
		// Create a new log file
		logFile, err := createLogFile(cfg.LogDir)
		if err != nil {
			return nil, err
		}

		// Always include stdout and the log file in development mode
		cfg.OutputPaths = append([]string{"stdout", logFile}, cfg.OutputPaths...)
		cfg.ErrorPaths = append([]string{"stderr", logFile}, cfg.ErrorPaths...)

		// Remove duplicates from OutputPaths and ErrorPaths
		cfg.OutputPaths = parser.DeduplicateElements(cfg.OutputPaths)
		cfg.ErrorPaths = parser.DeduplicateElements(cfg.ErrorPaths)
	}

	// Convert LogLevel to zapcore.Level
	var zapLevel zapcore.Level
	switch cfg.Level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Customize for development
	if cfg.Development {
		// Use colorized output in development mode
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
	}

	// Create basic configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      cfg.Development,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorPaths,
		InitialFields:    cfg.InitialFields,
	}

	if cfg.Development {
		// Use console encoding for development
		zapConfig.Encoding = "console"
	}

	// Create logger
	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Add basic fields
	logger = logger.With(
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment),
	)

	return &Logger{
		log:   logger,
		level: cfg.Level,
	}, nil
}

func createLogFile(logDir string) (string, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(logDir, fmt.Sprintf("%s.log", timestamp))

	// Create the log file
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	return filename, nil
}

// Methods for logging at different levels
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &Logger{
		log:   l.log.With(zapFields...),
		level: l.level,
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.log.Sync()
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}
