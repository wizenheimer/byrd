package screenshot

import (
	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// ScreenshotServiceOption is a function type that modifies ScreenshotService
type ScreenshotServiceOption func(*screenshotService)

// defaultConfig returns the default configuration
func defaultConfig() *models.ScreenshotServiceConfig {
	return &models.ScreenshotServiceConfig{
		// 40 requests per minute
		QPS: 0.6,
		// Default origin
		Origin: "https://api.screenshotone.com",
	}
}

// WithStorage sets the storage repository
func WithStorage(storage repo.ScreenshotRepository) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.storage = storage
	}
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(client clf.HTTPClient) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.httpClient = client
	}
}

// WithQPS sets the QPS limit
func WithQPS(qps float64) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.config.QPS = qps
	}
}

// WithOrigin sets the origin for making requests
func WithOrigin(origin string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.config.Origin = origin
	}
}

// WithKey sets the API key for making requests
func WithKey(key string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.config.Key = key
	}
}

// WithSignature sets the signature for making requests
func WithSignature(signature string) ScreenshotServiceOption {
	return func(s *screenshotService) {
		s.config.Signature = signature
	}
}
