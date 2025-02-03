package middleware

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	// Logger instance
	Logger *logger.Logger

	// Skip logging for specific paths
	SkipPaths []string
}

func NewLoggingMiddleware(logger *logger.Logger, skipPaths []string) (*LoggingMiddleware, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	// Deduplicate the skip paths
	skipPaths = deduplicate(skipPaths)

	lm := LoggingMiddleware{
		Logger:    logger,
		SkipPaths: skipPaths,
	}

	return &lm, nil
}

// deduplicate removes duplicate elements from a slice
func deduplicate(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for v := range slice {
		slice[v] = strings.ToLower(slice[v])
		if !encountered[slice[v]] {
			// Add element to result
			result = append(result, slice[v])
			// Record this element as an encountered element
			encountered[slice[v]] = true
		}
	}
	return result
}

// RequestResponseLogger returns a middleware that logs HTTP requests/responses
func (lm LoggingMiddleware) RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if path should be skipped
		for _, path := range lm.SkipPaths {
			if path == c.Path() {
				return c.Next()
			}
		}

		start := time.Now()

		// Log request before processing
		reqFields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", string(c.Request().Header.UserAgent())),
			zap.Any("headers", getHeaders(c)),
			zap.Any("query_params", c.Queries()),
		}

		// Add parsed request body fields
		reqFields = append(reqFields, getBodyAsFields("req_body", c.Body())...)

		// Log request
		lm.Logger.Debug("Incoming request", reqFields...)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response body - need to copy as it will be sent to client
		responseBody := string(c.Response().Body())

		// Log response after processing
		fields := []zap.Field{
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", duration),
			zap.String("response_body", responseBody),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		lm.Logger.Debug("Outgoing response", fields...)

		return err
	}
}

// getHeaders returns a map of request headers
func getHeaders(c *fiber.Ctx) map[string]string {
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

// parseJSON attempts to parse a string into a map[string]interface{}
func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if len(data) == 0 {
		return nil, nil
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getBodyAsFields converts body to zap fields, parsing JSON if possible
func getBodyAsFields(prefix string, body []byte) []zap.Field {
	parsed, err := parseJSON(body)
	if err != nil {
		// If parsing fails, log as raw string
		return []zap.Field{
			zap.String(prefix+"_raw", string(body)),
		}
	}

	if parsed == nil {
		return []zap.Field{
			zap.String(prefix+"_raw", ""),
		}
	}

	// Convert each key-value pair to a zap field
	fields := make([]zap.Field, 0, len(parsed))
	for key, value := range parsed {
		fields = append(fields, zap.Any(prefix+"_"+key, value))
	}
	return fields
}
