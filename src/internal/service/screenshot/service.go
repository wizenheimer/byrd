package screenshot

import (
	"context"
	"errors"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/logger"
	"github.com/wizenheimer/iris/pkg/utils/competitor"
	"go.uber.org/zap"
)

type screenshotService struct {
	storage    interfaces.StorageRepository
	httpClient interfaces.HTTPClient
	config     *models.ScreenshotServiceConfig
	logger     *logger.Logger
}

// NewScreenshotService creates a new screenshot service with the given options
func NewScreenshotService(logger *logger.Logger, opts ...ScreenshotServiceOption) (interfaces.ScreenshotService, error) {
	logger.Debug("creating new screenshot service")
	s := &screenshotService{
		config: defaultConfig(),
		logger: logger.WithFields(map[string]interface{}{"module": "screenshot_service"}),
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
		return nil, errors.New("screenshot key is required")
	}

	return s, nil
}

func (s *screenshotService) TakeScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, error) {
	s.logger.Debug("taking screenshot", zap.Any("url", opts.URL))

	// Implementation
	return nil, nil
}

func (s *screenshotService) GetContent(ctx context.Context, hash, weekNumber, weekDay string) (*models.ScreenshotResponse, error) {
	s.logger.Debug("getting content", zap.Any("hash", hash), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

	screenshotPath := competitor.GetScreenshotPath(hash, weekNumber, weekDay)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}

func (s *screenshotService) GetScreenshot(ctx context.Context, hash, weekNumber, weekDay string) (*models.ScreenshotResponse, error) {
	s.logger.Debug("getting screenshot", zap.Any("hash", hash), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

	screenshotPath := competitor.GetScreenshotPath(hash, weekNumber, weekDay)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}
