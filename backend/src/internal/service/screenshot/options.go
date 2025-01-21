// ./src/internal/service/screenshot/options.go
package screenshot

import (
	"github.com/wizenheimer/byrd/src/internal/client"
	"github.com/wizenheimer/byrd/src/internal/repository/screenshot"
)

// ScreenshotServiceOption is a function type that modifies ScreenshotService
type ScreenshotServiceOption func(*screenshotService)

// WithStorage sets the storage repository
func WithStorage(storage screenshot.ScreenshotRepository) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.storage = storage
	}
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(client client.HTTPClient) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.httpClient = client
	}
}

// WithQPS sets the QPS limit
func WithQPS(qps float64) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.qps = qps
	}
}

// WithOrigin sets the origin for making requests
func WithOrigin(origin string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.origin = origin
	}
}

// WithKey sets the API key for making requests
func WithKey(key string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.key = key
	}
}

// WithSignature sets the signature for making requests
func WithSignature(signature string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.signature = signature
	}
}
