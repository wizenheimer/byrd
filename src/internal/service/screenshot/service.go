package screenshot

import (
	"context"
	"errors"
	"fmt"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"strconv"

	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/path"
	"go.uber.org/zap"
)

// NewScreenshotService creates a new screenshot service with the given options
func NewScreenshotService(logger *logger.Logger, opts ...ScreenshotServiceOption) (svc.ScreenshotService, error) {
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
	if s.config.Origin == "" {
		return nil, errors.New("screenshot origin is required")
	}

	s.logger.Debug("created new screenshot service", zap.Any("config", s.config))
	return s, nil
}

// Refresh retrieves the current screenshot and html content for a given URL
func (s *screenshotService) Refresh(ctx context.Context, url string, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error) {
	imgResp, err := s.GetCurrentImage(ctx, true, opts)
	if err != nil {
		return nil, nil, err
	}

	htmlOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   opts.URL,
		RenderedURL: imgResp.Metadata.RenderedURL,
	}

	htmlContentResp, err := s.GetCurrentHTMLContent(ctx, true, htmlOpts)
	if err != nil {
		return nil, nil, err
	}

	return imgResp, htmlContentResp, nil
}

// Retrieve retrieves the previous screenshot and html content for a given URL
func (s *screenshotService) Retrieve(ctx context.Context, url string) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error) {
	imgResp, err := s.GetPreviousImage(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	htmlContentResp, err := s.GetPreviousHTMLContent(ctx, imgResp.Metadata.RenderedURL)
	if err != nil {
		return nil, nil, err
	}

	return imgResp, htmlContentResp, nil
}

// GetCurrentImage retrieves the current screenshot from the storage if present
// Or it will take a new screenshot and store it as an image
func (s *screenshotService) GetCurrentImage(ctx context.Context, save bool, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, error) {
	// Get screenshot if it exists
	if screenshotResponse, err := s.getExistingScreenshotImage(ctx, opts.URL); err == nil {
		return screenshotResponse, nil
	}

	// Prepare screenshot
	resp, err := s.prepareScreenshot(opts)
	if err != nil {
		return nil, err
	}

	currentYear, currentWeek, currentDayString := path.GetCurrentTimeComponents(false)
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, fmt.Errorf("failed to convert current day to int: %v", err)
	}

	// Parse the response
	imgResp, err := s.prepareScreenshotImageResponse(resp, opts.URL, currentYear, currentWeek, currentDay)
	if err != nil {
		return nil, err
	}

	// Save the screenshot if required
	if save {
		currentPath, err := path.GetCurrentScreenshotPath(opts.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to get current screenshot path: %v", err)
		}
		if err := s.storage.StoreScreenshotImage(ctx, *imgResp, currentPath); err != nil {
			return nil, fmt.Errorf("failed to store screenshot: %v", err)
		}
	}

	return imgResp, nil
}

// GetCurrentHTMLContent retrieves the current html content from the storage if present
// Or it will take a new screenshot and store it as html
func (s *screenshotService) GetCurrentHTMLContent(ctx context.Context, save bool, opts models.ScreenshotHTMLRequestOptions) (*models.ScreenshotHTMLContentResponse, error) {
	// Get screenshot if it exists
	if htmlContentResp, err := s.getExistingHTMLContent(ctx, opts.RenderedURL); err == nil {
		return htmlContentResp, nil
	}

	// Prepare screenshot
	resp, err := s.prepareScreenshotHTML(opts)
	if err != nil {
		return nil, err
	}

	// Get current time components
	currentYear, currentWeek, currentDayString := path.GetCurrentTimeComponents(false)
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, fmt.Errorf("failed to convert current day to int: %v", err)
	}

	// Parse the response
	htmlContentResp, err := s.prepareScreenshotHTMLContentResponse(resp, opts.SourceURL, opts.RenderedURL, currentYear, currentWeek, currentDay)
	if err != nil {
		return nil, err
	}

	// Save the screenshot if required
	if save {
		currentPath, err := path.GetCurrentContentPath(opts.SourceURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get current content path: %v", err)
		}
		if err := s.storage.StoreScreenshotHTMLContent(ctx, *htmlContentResp, currentPath); err != nil {
			return nil, fmt.Errorf("failed to store screenshot: %v", err)
		}
	}

	return nil, nil
}

// GetPreviousImage retrieves previous screenshot image from the storage
func (s *screenshotService) GetPreviousImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error) {
	screenshotPath, err := path.GetPreviousScreenshotPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous screenshot path: %v", err)
	}

	imgResp, errs := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get previous screenshot: %v", errs)
	}

	return &imgResp, nil
}

// GetPreviousScreenshotContent retrieves previous screenshot content from the storage
func (s *screenshotService) GetPreviousHTMLContent(ctx context.Context, url string) (*models.ScreenshotHTMLContentResponse, error) {
	contentPath, err := path.GetPreviousContentPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous content path: %v", err)
	}

	s.logger.Debug("getting previous content", zap.Any("path", contentPath))

	contentResp, errs := s.storage.GetScreenshotHTMLContent(ctx, contentPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get previous content: %v", errs)
	}

	return &contentResp, nil
}

// GetImage retrieves a screenshot image from the storage
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error) {
	screenshotPath, err := path.GetScreenshotPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot path: %v", err)
	}

	imgResp, errs := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get screenshot: %v", errs)
	}

	return &imgResp, nil
}

// GetHTMLContent retrieves the screenshot content from the screenshot service
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetHTMLContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotHTMLContentResponse, error) {
	contentPath, err := path.GetContentPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get content path: %v", err)
	}

	contentResp, errs := s.storage.GetScreenshotHTMLContent(ctx, contentPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get content: %v", errs)
	}

	return &contentResp, nil
}

// ListScreenshots lists the screenshots for a given URL
func (s *screenshotService) ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error) {
	prefix, err := path.GetListingPrefixFromContentType(url, contentType)
	if err != nil || prefix == nil {
		return nil, fmt.Errorf("failed to get listing prefix: %v", err)
	}

	// Get the list of screenshots
	return s.storage.List(ctx, *prefix, maxItems)
}
