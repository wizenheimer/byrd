package screenshot

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/competitor"
	"github.com/wizenheimer/iris/src/pkg/utils/parser"
	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
	"go.uber.org/zap"
)

type screenshotService struct {
	storage    interfaces.StorageRepository
	httpClient client.HTTPClient
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
	if s.config.Origin == "" {
		return nil, errors.New("screenshot origin is required")
	}

	s.logger.Debug("created new screenshot service", zap.Any("config", s.config))
	return s, nil
}

func (s *screenshotService) TakeScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, error) {
	s.logger.Debug("taking screenshot", zap.Any("url", opts.URL))

	resp, err := s.createScreenshotRequest(opts)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to take screenshot - %v status code", resp.StatusCode)
	}

	_, width, height, err := parseImageFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image from response: %v", err)
	}

	cleanText, title, err := parseContentFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content from response: %v", err)
	}

	// Implementation
	return &models.ScreenshotResponse{
		Status: "success",
		Paths: &models.ScreenshotPaths{
			Screenshot: "/temp",
			Content:    "/temp",
		},
		Metadata: &models.ScreenshotMeta{
			ImageWidth:  width,
			ImageHeight: height,
			PageTitle:   ptr.To(title),
			ContentSize: ptr.To(len(cleanText)),
		},
		URL: ptr.To(opts.URL),
	}, nil
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

// createScreenshotRequest creates a request for the screenshot API
func (s *screenshotService) createScreenshotRequest(opts models.ScreenshotRequestOptions) (*http.Response, error) {
	// Get default options
	defaultOpt := getDefaultScreenshotRequestOptions()

	// Merge the provided options with default options
	mergedOpt := mergeScreenshotRequestOptions(defaultOpt, opts)

	s.logger.Debug("creating screenshot request", zap.Any("opts", opts), zap.Any("merged_opts", mergedOpt), zap.Any("default_opts", defaultOpt))

	// Prepare the request
	return s.httpClient.NewRequest().
		BaseURL(s.config.Origin).
		Method(http.MethodGet).
		Path("/take").
		QueryParam("access_key", s.config.Key).
		AddQueryParamsFromStruct(mergedOpt).
		Execute(s.httpClient)
}

func getDefaultScreenshotRequestOptions() models.ScreenshotRequestOptions {
	// Get default options
	defaultOpt := models.ScreenshotRequestOptions{
		// Capture options
		Format:                ptr.To("png"),
		ImageQuality:          ptr.To(80),
		CaptureBeyondViewport: ptr.To(true),
		FullPage:              ptr.To(true),
		FullPageAlgorithm:     ptr.To(models.FullPageAlgorithmDefault),

		// Resource blocking options
		BlockAds:                 ptr.To(true),
		BlockCookieBanners:       ptr.To(true),
		BlockBannersByHeuristics: ptr.To(true),
		BlockTrackers:            ptr.To(true),
		BlockChats:               ptr.To(true),

		// Wait and delay options
		Delay:             ptr.To(0),
		Timeout:           ptr.To(60),
		NavigationTimeout: ptr.To(30),
		WaitUntil: []models.WaitUntilOption{
			models.WaitUntilNetworkIdle2,
			models.WaitUntilNetworkIdle0,
		},

		// Styling options
		DarkMode:      ptr.To(false),
		ReducedMotion: ptr.To(true),

		// Response options
		MetadataImageSize:      ptr.To(true),
		MetadataPageTitle:      ptr.To(true),
		MetadataContent:        ptr.To(true),
		MetadataHttpStatusCode: ptr.To(true),
	}

	return defaultOpt
}

// MergeOptions merges the provided options with default options
func mergeScreenshotRequestOptions(defaults, override models.ScreenshotRequestOptions) models.ScreenshotRequestOptions {
	result := defaults

	// Use reflection to handle all fields
	rOverride := reflect.ValueOf(override)
	rResult := reflect.ValueOf(&result).Elem()

	for i := 0; i < rOverride.NumField(); i++ {
		field := rOverride.Field(i)
		resultField := rResult.Field(i)

		// Skip if the override field is nil or zero
		if field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				resultField.Set(field)
			}
		case reflect.String:
			if field.String() != "" {
				resultField.Set(field)
			}
		case reflect.Slice:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		case reflect.Map:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		}
	}

	return result
}

// parseImageFromResponse parses an image from an HTTP response
func parseImageFromResponse(resp *http.Response) (image.Image, int, int, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, -1, -1, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	imageContentTypes := []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
	}

	if !parser.Contains(imageContentTypes, resp.Header.Get("Content-Type")) {
		return nil, -1, -1, fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, -1, -1, fmt.Errorf("failed to decode image: %w", err)
	}

	var width, height int

	imageHeightString := resp.Header.Get("X-ScreenshotOne-Image-Height")
	imageWidthString := resp.Header.Get("X-ScreenshotOne-Image-Width")
	bounds := img.Bounds()

	if imageHeightString != "" || imageWidthString != "" {
		width = bounds.Max.X - bounds.Min.X
		height = bounds.Max.Y - bounds.Min.Y
	} else {
		width, err = strconv.Atoi(imageWidthString)
		if err != nil {
			width = bounds.Max.X - bounds.Min.X
		}

		height, err = strconv.Atoi(imageHeightString)
		if err != nil {
			height = bounds.Max.Y - bounds.Min.Y
		}
	}
	return img, width, height, nil
}

// parseContentFromResponse parses the content from an HTTP response
func parseContentFromResponse(resp *http.Response) (string, string, error) {
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Get the content URL from header
	contentURL := resp.Header.Get("X-ScreenshotOne-Content-URL")
	if contentURL == "" {
		return "", "", fmt.Errorf("no content URL found in headers")
	}

	// Make a new request to the content URL
	htmlResp, err := http.Get(contentURL)
	if err != nil {
		return "", "", fmt.Errorf("error fetching HTML content: %v", err)
	}
	defer htmlResp.Body.Close()

	// Read the HTML content
	htmlBytes, err := io.ReadAll(htmlResp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading HTML body: %v", err)
	}

	htmlContent := string(htmlBytes)

	// Extract clean text
	cleanText, err := parser.ParseTextFromHTML(htmlContent)
	if err != nil {
		return "", "", fmt.Errorf("error extracting text: %v", err)
	}

	// Parse the title
	title := resp.Header.Get("X-ScreenshotOne-Page-Title")

	return cleanText, title, nil
}
