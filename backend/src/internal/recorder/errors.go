// pkg/recorder/recorder.go
package recorder

import (
	"context"
	"runtime"
	"time"

	"github.com/highlight/highlight/sdk/highlight-go"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrorRecorder handles error recording based on environment
type ErrorRecorder struct {
	logger      *logger.Logger
	isDev       bool
	serviceName string
}

// NewErrorRecorder creates a new ErrorRecorder
func NewErrorRecorder(logger *logger.Logger, isDev bool, serviceName string) *ErrorRecorder {
	return &ErrorRecorder{
		logger:      logger,
		isDev:       isDev,
		serviceName: serviceName,
	}
}

// RecordError records an error with zap fields
func (er *ErrorRecorder) RecordError(ctx context.Context, err error, fields ...zap.Field) {
	if err == nil {
		return
	}

	// Get the caller information
	_, file, line, _ := runtime.Caller(1)

	// Add default fields
	fields = append(fields,
		zap.Error(err),
		zap.String("file", file),
		zap.Int("line", line),
		zap.String("service", er.serviceName),
	)

	if er.isDev {
		// In development, just use zap logging directly
		er.logger.Error("Error occurred", fields...)
	} else {
		// In production, convert zap fields to Highlight attributes
		attrs := zapFieldsToAttributes(fields)
		highlight.RecordError(ctx, err, attrs...)
	}
}

// zapFieldsToAttributes converts zap fields to OpenTelemetry attributes
func zapFieldsToAttributes(fields []zap.Field) []attribute.KeyValue {
	encoder := zapcore.NewMapObjectEncoder()
	attrs := make([]attribute.KeyValue, 0, len(fields))

	for _, field := range fields {
		// Special handling for Error type
		if field.Type == zapcore.ErrorType {
			if field.Interface != nil {
				attrs = append(attrs, attribute.String(field.Key, field.Interface.(error).Error()))
			}
			continue
		}

		// Add the field to our encoder
		field.AddTo(encoder)

		// Convert based on the field type
		switch field.Type {
		case zapcore.StringType:
			attrs = append(attrs, attribute.String(field.Key, field.String))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			attrs = append(attrs, attribute.Int64(field.Key, field.Integer))
		case zapcore.Float64Type, zapcore.Float32Type:
			attrs = append(attrs, attribute.Float64(field.Key, field.Interface.(float64)))
		case zapcore.BoolType:
			attrs = append(attrs, attribute.Bool(field.Key, field.Integer == 1))
		case zapcore.DurationType:
			attrs = append(attrs, attribute.Int64(field.Key, int64(field.Interface.(time.Duration))))
		default:
			// For complex types, convert to string representation
			attrs = append(attrs, attribute.String(field.Key, encoder.Fields[field.Key].(string)))
		}
	}

	return attrs
}
