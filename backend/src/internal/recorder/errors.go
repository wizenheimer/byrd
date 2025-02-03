// pkg/recorder/recorder.go
package recorder

import (
	"context"
	"fmt"
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
// zapFieldsToAttributes converts zap fields to OpenTelemetry attributes
func zapFieldsToAttributes(fields []zap.Field) []attribute.KeyValue {
	encoder := zapcore.NewMapObjectEncoder()
	attrs := make([]attribute.KeyValue, 0, len(fields))

	for _, field := range fields {
		// First add to encoder to handle complex types
		field.AddTo(encoder)

		// Try to get the value from encoder first
		if encodedVal, ok := encoder.Fields[field.Key]; ok {
			// Handle the encoded value based on its type
			switch v := encodedVal.(type) {
			case string:
				attrs = append(attrs, attribute.String(field.Key, v))
			case int:
				attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", v)))
			case int64:
				attrs = append(attrs, attribute.Int64(field.Key, v))
			case float64:
				attrs = append(attrs, attribute.Float64(field.Key, v))
			case bool:
				attrs = append(attrs, attribute.Bool(field.Key, v))
			case []interface{}:
				// Convert slice to string representation
				attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", v)))
			case map[string]interface{}:
				// Convert map to string representation
				attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", v)))
			case time.Time:
				attrs = append(attrs, attribute.String(field.Key, v.Format(time.RFC3339)))
			case time.Duration:
				attrs = append(attrs, attribute.Int64(field.Key, int64(v)))
			default:
				// Fallback to string representation for unknown types
				attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", v)))
			}
			continue
		}

		// Handle special cases and native zap types
		switch field.Type {
		case zapcore.ErrorType:
			if field.Interface != nil {
				if err, ok := field.Interface.(error); ok {
					attrs = append(attrs, attribute.String(field.Key, err.Error()))
				} else {
					attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", field.Interface)))
				}
			}
		case zapcore.StringType:
			attrs = append(attrs, attribute.String(field.Key, field.String))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			attrs = append(attrs, attribute.Int64(field.Key, field.Integer))
		case zapcore.BoolType:
			attrs = append(attrs, attribute.Bool(field.Key, field.Integer == 1))
		case zapcore.DurationType:
			if dur, ok := field.Interface.(time.Duration); ok {
				attrs = append(attrs, attribute.Int64(field.Key, int64(dur)))
			}
		case zapcore.TimeType:
			if t, ok := field.Interface.(time.Time); ok {
				attrs = append(attrs, attribute.String(field.Key, t.Format(time.RFC3339)))
			}
		case zapcore.NamespaceType:
			// Skip namespace fields as they're just structural
			continue
		default:
			// For any other types, use string representation
			attrs = append(attrs, attribute.String(field.Key, fmt.Sprintf("%v", field.Interface)))
		}
	}

	return attrs
}
