// pkg/recorder/recorder.go
package recorder

import (
	"context"
	"fmt"
	"runtime"

	"github.com/highlight/highlight/sdk/highlight-go"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
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

// RecordError records an error either to logs in development or Highlight in production
// attributes are optional and can be provided as key-value pairs
func (er *ErrorRecorder) RecordError(ctx context.Context, err error, attributes ...any) {
	if err == nil {
		return
	}

	// Get the caller information
	_, file, line, _ := runtime.Caller(1)

	if er.isDev {
		// In development, log the error with details
		fields := []zap.Field{
			zap.Error(err),
			zap.String("file", file),
			zap.Int("line", line),
			zap.String("service", er.serviceName),
		}

		// Convert attributes to zap fields
		for i := 0; i < len(attributes); i += 2 {
			if i+1 < len(attributes) {
				key, ok := attributes[i].(string)
				if ok {
					fields = append(fields, zap.Any(key, attributes[i+1]))
				}
			}
		}

		er.logger.Error("Error occurred", fields...)
	} else {
		// In production, send to Highlight
		attrs := []attribute.KeyValue{
			attribute.String("file", file),
			attribute.Int("line", line),
			attribute.String("service", er.serviceName),
		}

		// Convert provided attributes to Highlight attributes
		for i := 0; i < len(attributes); i += 2 {
			if i+1 < len(attributes) {
				key, ok := attributes[i].(string)
				if ok {
					attrs = append(attrs, attribute.String(key, fmt.Sprint(attributes[i+1])))
				}
			}
		}

		highlight.RecordError(ctx, err, attrs...)
	}
}
