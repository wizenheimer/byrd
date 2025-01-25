package screenshot

import (
	"context"
	"errors"

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
	if opts.URL == "" {
		return nil, nil, errors.New("URL is required for generating screenshot")
	}

	// Step 1: Check if the screenshot exists in storage
	existingScreenshot, existingContent, err := s.Retrieve(ctx, opts, backDate)
	if err == nil {
		return existingScreenshot, existingContent, nil
	}

	// Step 2: If the screenshot does not exist in storage, create a new screenshot
	// Step 2.1: Determine the screenshot path
	screenshotPath, err := DeterminePath(opts.URL, ContentTypeImage, backDate)
	if err != nil {
		return nil, nil, err
	}
	contentPath, err := DeterminePath(opts.URL, ContentTypeContent, backDate)
	if err != nil {
		return nil, nil, err
	}

	// Step 2.2: Refresh the screenshot image and content
	img, content, err := s.refreshScreenshot(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// Step 2.3: Get the screenshot metadata
	metadata, err := getScreenshotMetadata(backDate)
	if err != nil {
		return nil, nil, err
	}

	// Step 2.4: Create the screenshot image and content
	screenshotImage := &models.ScreenshotImage{
		StoragePath: screenshotPath,
		Image:       *img,
		Metadata:    metadata,
	}
	screenshotContent := &models.ScreenshotContent{
		StoragePath: contentPath,
		Content:     *content,
		Metadata:    metadata,
	}

	// Step 3: Store the screenshot in storage
	if err := s.storage.StoreScreenshotImage(ctx, screenshotImage); err != nil {
		return nil, nil, err
	}

	// Step 4: Store the screenshot content in storage
	if err := s.storage.StoreScreenshotContent(ctx, screenshotContent); err != nil {
		return nil, nil, err
	}

	return screenshotImage, screenshotContent, nil
}

func (s *screenshotService) Retrieve(ctx context.Context, opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImage, *models.ScreenshotContent, error) {
	// Step 1: Determine the screenshot path
	screenshotPath, err := DeterminePath(opts.URL, ContentTypeImage, backDate)
	if err != nil {
		return nil, nil, err
	}

	// Step 2: Retrieve the screenshot from storage
	screenshotImage, err := s.storage.RetrieveScreenshotImage(ctx, screenshotPath)
	if err != nil {
		return nil, nil, err
	}

	// Step 3: Determine the screenshot content path
	contentPath, err := DeterminePath(opts.URL, ContentTypeContent, backDate)
	if err != nil {
		return nil, nil, err
	}

	// Step 4: Retrieve the screenshot content from storage
	screenshotContent, err := s.storage.RetrieveScreenshotContent(ctx, contentPath)
	if err != nil {
		return nil, nil, err
	}

	return screenshotImage, screenshotContent, nil
}
