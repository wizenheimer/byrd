// ./src/internal/service/screenshot/service.go
package screenshot

import (
	"context"
	"errors"

	// "errors"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	"github.com/wizenheimer/byrd/src/internal/client"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type screenshotService struct {
	storage    screenshot.ScreenshotRepository
	httpClient client.HTTPClient
	qps        float64
	origin     string
	key        string
	signature  string
	logger     *logger.Logger
}

// NewScreenshotService creates a new screenshot service with the given options
func NewScreenshotService(logger *logger.Logger, opts ...ScreenshotServiceOption) (ScreenshotService, error) {
	logger.Debug("creating new screenshot service")
	s := &screenshotService{
		logger: logger.WithFields(
			map[string]interface{}{
				"module": "screenshot_service",
			},
		),
	}

	// Apply all options
	for _, opt := range opts {
		opt(s)
	}

	// Validate required dependencies
	if s.storage == nil {
		return nil, errors.New("storage repository is required")
	}
	if s.httpClient == nil {
		return nil, errors.New("HTTP client is required")
	}
	if s.key == "" {
		return nil, errors.New("screenshot key is required")
	}
	if s.origin == "" {
		return nil, errors.New("screenshot origin is required")
	}

	return s, nil
}

func (s *screenshotService) Refresh(ctx context.Context, opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImage, *models.ScreenshotContent, error) {
	return nil, nil, nil
}

func (s *screenshotService) Retrieve(ctx context.Context, opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImage, *models.ScreenshotContent, error) {
	return nil, nil, nil
}
