package screenshot

import (
	"context"
	"fmt"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type screenshotService struct {
	storage    interfaces.StorageRepository
	httpClient interfaces.HTTPClient
	config     *models.ScreenshotServiceConfig
}

func NewScreenshotService(
	storage interfaces.StorageRepository,
	httpClient interfaces.HTTPClient,
	config *models.ScreenshotServiceConfig,
) (interfaces.ScreenshotService, error) {
	return &screenshotService{
		storage:    storage,
		httpClient: httpClient,
		config:     config,
	}, nil
}

func (s *screenshotService) TakeScreenshot(ctx context.Context, opts models.ScreenshotOptions) (*models.ScreenshotResponse, error) {
	// Implementation
	return nil, nil
}

func (s *screenshotService) GetContent(ctx context.Context, hash, weekNumber, runID string) (*models.ScreenshotResponse, error) {
	screenshotPath := s.getScreenshotPath(hash, weekNumber, runID)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}

func (s *screenshotService) GetScreenshot(ctx context.Context, hash, weekNumber, runID string) (*models.ScreenshotResponse, error) {
	screenshotPath := s.getScreenshotPath(hash, weekNumber, runID)
	_, _, err := s.storage.Get(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotResponse{}, nil
}

func (s *screenshotService) getScreenshotPath(hash, weekNumber, runID string) string {
	return fmt.Sprintf("screenshots/%s/%s/%s", hash, weekNumber, runID)
}
