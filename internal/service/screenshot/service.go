package screenshot

import (
	"context"
	"errors"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/utils/competitor"
)

type screenshotService struct {
	storage    interfaces.StorageRepository
	httpClient interfaces.HTTPClient
	config     *models.ScreenshotServiceConfig
}

// NewScreenshotService creates a new screenshot service with the given options
func NewScreenshotService(opts ...ScreenshotServiceOption) (interfaces.ScreenshotService, error) {
	s := &screenshotService{
		config: defaultConfig(),
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
	if s.config.Key == "" {
		return nil, errors.New("API key is required")
	}

	return s, nil
}

func (s *screenshotService) TakeScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, error) {
	// Implementation
	return nil, nil
}

func (s *screenshotService) GetContent(ctx context.Context, hash, weekNumber, weekDay string) (*models.ScreenshotResponse, error) {
	screenshotPath := competitor.GetScreenshotPath(hash, weekNumber, weekDay)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}

func (s *screenshotService) GetScreenshot(ctx context.Context, hash, weekNumber, weekDay string) (*models.ScreenshotResponse, error) {
	screenshotPath := competitor.GetScreenshotPath(hash, weekNumber, weekDay)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}
