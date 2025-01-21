// ./src/internal/service/screenshot/service.go
package screenshot

import (
	"context"
	"errors"
	"fmt"

	// "errors"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"strconv"

	"github.com/wizenheimer/byrd/src/internal/client"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

type screenshotService struct {
	storage    screenshot.ScreenshotRepository
	httpClient client.HTTPClient
	config     *models.ScreenshotServiceConfig
	logger     *logger.Logger
}

// compile time check if the interface is implemented
// TODO: reduce overhead by passing stuff by reference
var _ ScreenshotService = (*screenshotService)(nil)

// NewScreenshotService creates a new screenshot service with the given options
func NewScreenshotService(logger *logger.Logger, opts ...ScreenshotServiceOption) (ScreenshotService, error) {
	logger.Debug("creating new screenshot service")
	s := &screenshotService{
		config: defaultConfig(),
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
	s.logger.Debug("refreshing screenshot", zap.String("url", url), zap.Any("opts", opts))
	imgResp, err := s.GetCurrentImage(ctx, true, false, opts)
	if err != nil {
		return nil, nil, err
	}

	htmlOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   opts.URL,
		RenderedURL: imgResp.Metadata.RenderedURL,
	}

	htmlContentResp, err := s.GetCurrentHTMLContent(ctx, true, false, htmlOpts)
	if err != nil {
		return nil, nil, err
	}

	return imgResp, htmlContentResp, nil
}

// Initiate initializes the screenshot repository with the latest screenshot and html content for a given URL by back dating the current capture
func (s *screenshotService) Initiate(ctx context.Context, url string, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error) {
	s.logger.Debug("initiating screenshot", zap.String("url", url), zap.Any("opts", opts))
	imgResp, err := s.GetCurrentImage(ctx, true, true, opts)
	if err != nil {
		return nil, nil, err
	}

	htmlOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   opts.URL,
		RenderedURL: imgResp.Metadata.RenderedURL,
	}

	htmlContentResp, err := s.GetCurrentHTMLContent(ctx, true, true, htmlOpts)
	if err != nil {
		return nil, nil, err
	}

	return imgResp, htmlContentResp, nil
}

// Retrieve retrieves the previous screenshot and html content for a given URL
func (s *screenshotService) Retrieve(ctx context.Context, url string) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error) {
	s.logger.Debug("retrieving previous screenshot", zap.String("url", url))
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

func (s *screenshotService) createScreenshot(opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImageResponse, error) {
	// Prepare screenshot
	resp, err := s.prepareScreenshot(opts)
	if err != nil {
		return nil, err
	}

	var currentDayString string
	var currentWeek int
	var currentYear int
	if backDate {
		currentYear, currentWeek, currentDayString = utils.GetPreviousTimeComponents(true)
	} else {
		currentYear, currentWeek, currentDayString = utils.GetCurrentTimeComponents(true)
	}
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, fmt.Errorf("failed to convert current day to int, %s, date: %s", err.Error(), currentDayString)
	}

	// Parse the response
	imgResp, err := s.prepareScreenshotImageResponse(resp, opts.URL, currentYear, currentWeek, currentDay)
	if err != nil {
		return nil, err
	}
	if imgResp == nil {
		return nil, errors.New("failed to prepare screenshot image response")
	}

	return imgResp, nil
}

// GetCurrentImage retrieves the current screenshot from the storage if present
// Or it will take a new screenshot and store it as an image
func (s *screenshotService) GetCurrentImage(ctx context.Context, save bool, backDate bool, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, error) {
	s.logger.Debug("getting current image", zap.Any("opts", opts))

	// Get screenshot if it exists
	screenshotResponse, err := s.getExistingScreenshotImage(ctx, opts.URL, backDate)
	if err == nil && !backDate {
		return screenshotResponse, nil
	}

	// Incase the screenshot does not exist, create a new one
	if screenshotResponse == nil || err != nil {
		// Create a new screenshot
		screenshotResponse, err = s.createScreenshot(opts, backDate)
		if err != nil {
			return nil, err
		}
		// Check if the screenshot was created
		if screenshotResponse == nil {
			return nil, errors.New("failed to create screenshot")
		}
	}

	// Save the screenshot if required
	if save {
		path, err := utils.GetCurrentScreenshotPath(opts.URL, backDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get current content path, %s", err.Error())
		}
		if err := s.storage.StoreScreenshotImage(ctx, *screenshotResponse, path); err != nil {
			return nil, err
		}
	}

	return screenshotResponse, nil
}

func (s *screenshotService) createHTMLContent(opts models.ScreenshotHTMLRequestOptions, backDate bool) (*models.ScreenshotHTMLContentResponse, error) {
	// Prepare screenshot
	resp, err := s.prepareScreenshotHTML(opts)
	if err != nil {
		return nil, err
	}

	// Get current time components
	var currentDayString string
	var currentWeek int
	var currentYear int
	if backDate {
		currentYear, currentWeek, currentDayString = utils.GetPreviousTimeComponents(true)
	} else {
		currentYear, currentWeek, currentDayString = utils.GetCurrentTimeComponents(true)
	}
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, fmt.Errorf("failed to convert current day to int, %s, date: %s", err.Error(), currentDayString)
	}

	// Parse the response
	htmlContentResp, err := s.prepareScreenshotHTMLContentResponse(resp, opts.SourceURL, opts.RenderedURL, currentYear, currentWeek, currentDay)
	if err != nil {
		return nil, err
	}
	if htmlContentResp == nil {
		return nil, errors.New("failed to prepare screenshot html content response")
	}

	return htmlContentResp, nil
}

// GetCurrentHTMLContent retrieves the current html content from the storage if present
// Or it will take a new screenshot and store it as html
func (s *screenshotService) GetCurrentHTMLContent(ctx context.Context, save bool, backDate bool, opts models.ScreenshotHTMLRequestOptions) (*models.ScreenshotHTMLContentResponse, error) {
	s.logger.Debug("getting current html content", zap.Any("opts", opts))
	// Get screenshot if it exists
	htmlContentResp, err := s.getExistingHTMLContent(ctx, opts.RenderedURL, backDate)
	if err == nil && !backDate {
		return htmlContentResp, nil
	}
	if htmlContentResp == nil || err != nil {
		// Create a new screenshot
		htmlContentResp, err = s.createHTMLContent(opts, backDate)
		if err != nil {
			return nil, err
		}
		// Check if the screenshot was created
		if htmlContentResp == nil {
			return nil, errors.New("failed to create screenshot")
		}
	}

	// Save the screenshot if required
	if save {
		var path string
		path, err := utils.GetCurrentContentPath(opts.SourceURL, backDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get current content path, %s", err.Error())
		}
		if err := s.storage.StoreScreenshotHTMLContent(ctx, *htmlContentResp, path); err != nil {
			return nil, err
		}
	}

	return htmlContentResp, nil
}

// GetPreviousImage retrieves previous screenshot image from the storage
func (s *screenshotService) GetPreviousImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error) {
	s.logger.Debug("getting previous image", zap.String("url", url))
	screenshotPath, err := utils.GetPreviousScreenshotPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous screenshot path, %s", err.Error())
	}

	imgResp, errs := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if errs != nil {
		return nil, errors.New("failed to get previous screenshot image")
	}

	return &imgResp, nil
}

// GetPreviousScreenshotContent retrieves previous screenshot content from the storage
func (s *screenshotService) GetPreviousHTMLContent(ctx context.Context, url string) (*models.ScreenshotHTMLContentResponse, error) {
	s.logger.Debug("getting previous html content", zap.String("url", url))
	contentPath, err := utils.GetPreviousContentPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous content path, %s", err.Error())
	}

	s.logger.Debug("getting previous content", zap.Any("path", contentPath))

	contentResp, errs := s.storage.GetScreenshotHTMLContent(ctx, contentPath)
	if errs != nil {
		return nil, errors.New("failed to get previous content")
	}

	return &contentResp, nil
}

// GetImage retrieves a screenshot image from the storage
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error) {
	s.logger.Debug("getting image", zap.String("url", url), zap.Int("year", year), zap.Int("week_number", weekNumber), zap.Int("week_day", weekDay))

	screenshotPath, err := utils.GetScreenshotPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot path, %s", err.Error())
	}

	imgResp, errs := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if errs != nil {
		return nil, errors.New("failed to get screenshot image")
	}

	return &imgResp, nil
}

// GetHTMLContent retrieves the screenshot content from the screenshot service
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetHTMLContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotHTMLContentResponse, error) {
	s.logger.Debug("getting html content", zap.String("url", url), zap.Int("year", year), zap.Int("week_number", weekNumber), zap.Int("week_day", weekDay))

	contentPath, err := utils.GetContentPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get content path, %s", err.Error())
	}

	contentResp, errs := s.storage.GetScreenshotHTMLContent(ctx, contentPath)
	if errs != nil {
		return nil, errors.New("failed to get previous content")
	}

	return &contentResp, nil
}

// ListScreenshots lists the screenshots for a given URL
func (s *screenshotService) ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error) {
	s.logger.Debug("listing screenshots", zap.String("url", url), zap.String("content_type", contentType), zap.Int("max_items", maxItems))

	prefix, err := utils.GetListingPrefixFromContentType(url, contentType)
	if err != nil || prefix == nil {
		return nil, errors.New("failed to get listing prefix from content type")
	}

	// Get the list of screenshots

	resp, err := s.storage.List(ctx, *prefix, maxItems)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
