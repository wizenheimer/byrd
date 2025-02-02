package highlightzap

import (
	"context"
	"fmt"

	"github.com/highlight/highlight/sdk/highlight-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

var (
	logSeverityKey = attribute.Key(highlight.LogSeverityAttribute)
	logMessageKey  = attribute.Key(highlight.LogMessageAttribute)
)

type HighlightCore struct {
	zapcore.LevelEnabler

	coreFields map[string]any
	// syncOnWrite bool
}

// NewHighlightCore creates a new core to transmit logs to highlight.
// Highlight token and other options should be set before creating a new core
func NewHighlightCore(minLevel zapcore.Level) *HighlightCore {
	return &HighlightCore{
		LevelEnabler: minLevel,
		coreFields:   make(map[string]any),
	}
}

// With provides structure
func (c *HighlightCore) With(fields []zapcore.Field) zapcore.Core {
	fieldMap := fieldsToMap(fields)
	for k, v := range fieldMap {
		c.coreFields[k] = v
	}

	return c
}

// Check determines if this should be sent to roll bar based on LevelEnabler
func (c *HighlightCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}

	return checkedEntry
}

func (c *HighlightCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	ctx := context.TODO()

	span, _ := highlight.StartTrace(ctx, "highlightzap/log")
	defer highlight.EndTrace(span)

	attrs := []attribute.KeyValue{
		logSeverityKey.String(entry.Level.String()),
		logMessageKey.String(entry.Message),
	}

	if entry.Caller.Function != "" {
		attrs = append(attrs, semconv.CodeFunctionKey.String(entry.Caller.Function))
	}

	if entry.Caller.File != "" {
		attrs = append(attrs, semconv.CodeFilepathKey.String(entry.Caller.File))
		attrs = append(attrs, semconv.CodeLineNumberKey.Int(entry.Caller.Line))
	}

	for k, v := range fieldsToMap(fields) {
		if entry.Level == zapcore.ErrorLevel {
			attrs = append(attrs, attribute.String(k, fmt.Sprintf("%+v", v)))
			errMap := extractError(fields)
			if errMap != nil {
				attrs = append(attrs, attribute.String(k, fmt.Sprintf("%s", errMap)))
			}
		} else {
			attrs = append(attrs, attribute.String(k, fmt.Sprintf("%+v", v)))
		}
	}

	span.AddEvent(highlight.LogEvent, trace.WithAttributes(attrs...))

	if entry.Level <= zapcore.ErrorLevel {
		span.SetStatus(codes.Error, entry.Message)
	}

	return nil
}

func (c *HighlightCore) Sync() error {
	return nil
}

func extractError(fields []zapcore.Field) error {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	var foundError error
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			foundError = f.Interface.(error)
		}
	}
	return foundError
}

func fieldsToMap(fields []zapcore.Field) map[string]any {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	m := make(map[string]any)
	for k, v := range enc.Fields {
		m[k] = v
	}
	return m
}
