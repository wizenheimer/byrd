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
	"time"

	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/parser"
	"github.com/wizenheimer/iris/src/pkg/utils/path"
	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
	"go.uber.org/zap"
)

type screenshotService struct {
	storage    interfaces.ScreenshotRepository
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

func (s *screenshotService) GetExistingScreenshot(ctx context.Context, url string) (*models.ScreenshotResponse, image.Image, string, error) {
	screenshotPath, err := path.GetCurrentScreenshotPath(url)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get current screenshot path: %v", err)
	}

	contentPath, err := path.GetCurrentContentPath(url)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get current content path: %v", err)
	}

	var screenshotImage image.Image
	var screenshotContent string
	var screenshotMetadata models.ScreenshotMetadata
	if screenshotImage, screenshotMetadata, err = s.storage.GetScreenshotImage(ctx, screenshotPath); err != nil {
		s.logger.Warn("failed to get image", zap.Error(err))
		return nil, nil, "", fmt.Errorf("failed to get existing screenshot: %v", err)
	}
	// Get content if it exists
	if screenshotContent, screenshotMetadata, err = s.storage.GetScreenshotContent(ctx, contentPath); err != nil {
		s.logger.Warn("failed to get content", zap.Error(err))
		return nil, nil, "", fmt.Errorf("failed to get existing screenshot: %v", err)
	}

	screenshotResponse := models.ScreenshotResponse{
		Status: "success",
		Paths: &models.ScreenshotPaths{
			Screenshot: screenshotPath,
			Content:    contentPath,
		},
		Metadata: &models.ScreenshotMeta{
			ImageWidth:  screenshotMetadata.ImageWidth,
			ImageHeight: screenshotMetadata.ImageHeight,
			PageTitle:   screenshotMetadata.PageTitle,
			ContentSize: ptr.To(len(screenshotContent)),
		},
		URL: ptr.To(url),
	}

	return &screenshotResponse, screenshotImage, screenshotContent, nil
}

// CaptureScreenshot takes a screenshot of a given URL
func (s *screenshotService) CaptureScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, image.Image, string, error) {
	// Get screenshot if it exists
	if screenshotResponse, screenshotImage, screenshotContent, err := s.GetExistingScreenshot(ctx, opts.URL); err == nil {
		return screenshotResponse, screenshotImage, screenshotContent, nil
	}

	resp, err := s.prepareScreenshot(opts)
	if err != nil {
		return nil, nil, "", err
	}

	img, width, height, cleanText, title, err := s.parseScreenshotRespose(resp)
	if err != nil {
		return nil, nil, "", err
	}

	metadata := models.ScreenshotMetadata{
		SourceURL:         opts.URL,
		FetchedAt:         time.Now().String(),
		ScreenshotService: "screenshotone",
		ImageWidth:        width,
		ImageHeight:       height,
		PageTitle:         ptr.To(title),
	}

	screenshotPath, err := path.GetCurrentScreenshotPath(opts.URL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get current screenshot path: %v", err)
	}
	if err := s.storage.StoreScreenshotImage(ctx, img, screenshotPath, metadata); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store image: %v", err)
	}

	contentPath, err := path.GetCurrentContentPath(opts.URL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get current content path: %v", err)
	}
	if err := s.storage.StoreScreenshotContent(ctx, cleanText, contentPath, metadata); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store content: %v", err)
	}

	screenshotResponse := models.ScreenshotResponse{
		Status: "success",
		Paths: &models.ScreenshotPaths{
			Screenshot: screenshotPath,
			Content:    contentPath,
		},
		Metadata: &models.ScreenshotMeta{
			ImageWidth:  width,
			ImageHeight: height,
			PageTitle:   ptr.To(title),
			ContentSize: ptr.To(len(cleanText)),
		},
		URL: ptr.To(opts.URL),
	}
	return &screenshotResponse, img, cleanText, nil
}

// GetPreviousScreenshotImage retrieves previous screenshot image from the storage
func (s *screenshotService) GetPreviousScreenshotImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error) {
	s.logger.Debug("getting previous screenshot", zap.Any("url", url))

	screenshotPath, err := path.GetPreviousScreenshotPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous screenshot path: %v", err)
	}

	// s.logger.Debug("getting previous screenshot", zap.Any("path", screenshotPath))

	img, metadata, err := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous screenshot: %v", err)
	}

	screenshotResponse := models.ScreenshotImageResponse{
		Status:   "success",
		Image:    img,
		Metadata: ptr.To(metadata),
		URL:      ptr.To(url),
		Path:     ptr.To(screenshotPath),
	}

	return &screenshotResponse, nil
}

// GetPreviousScreenshotContent retrieves previous screenshot content from the storage
func (s *screenshotService) GetPreviousScreenshotContent(ctx context.Context, url string) (*models.ScreenshotContentResponse, error) {
	s.logger.Debug("getting previous content", zap.Any("url", url))

	contentPath, err := path.GetPreviousContentPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous content path: %v", err)
	}

	s.logger.Debug("getting previous content", zap.Any("path", contentPath))

	content, metadata, err := s.storage.GetScreenshotContent(ctx, contentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous content: %v", err)
	}

	screenshotResponse := models.ScreenshotContentResponse{
		Status:   "success",
		Content:  content,
		Metadata: ptr.To(metadata),
		URL:      ptr.To(url),
		Path:     ptr.To(contentPath),
	}

	return &screenshotResponse, nil
}

// GetScreenshotContent retrieves the screenshot content from the screenshot service
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetScreenshotContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotContentResponse, error) {
	contentPath, err := path.GetContentPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get content path: %v", err)
	}

	s.logger.Debug("getting content", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay), zap.Any("path", contentPath))

	content, metadata, err := s.storage.GetScreenshotContent(ctx, contentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %v", err)
	}

	return &models.ScreenshotContentResponse{
		Status:   "success",
		Content:  content,
		Metadata: ptr.To(metadata),
		URL:      ptr.To(url),
		Path:     ptr.To(contentPath),
	}, nil
}

// GetScreenshotImage retrieves a screenshot image from the storage
// The URL is used to generate the path to the image
// The week number and week day are used to generate the path to the image
func (s *screenshotService) GetScreenshotImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error) {
	screenshotPath, err := path.GetScreenshotPath(url, year, weekNumber, weekDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot path: %v", err)
	}

	s.logger.Debug("getting screenshot", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay), zap.Any("path", screenshotPath))

	img, metadata, err := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot: %v", err)
	}

	return &models.ScreenshotImageResponse{
		Status:   "success",
		Image:    img,
		Metadata: ptr.To(metadata),
		URL:      ptr.To(url),
		Path:     ptr.To(screenshotPath),
	}, nil
}

// ListScreenshots lists the screenshots for a given URL
func (s *screenshotService) ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error) {
	s.logger.Debug("listing screenshots", zap.Any("url", url), zap.Any("content_type", contentType))

	hash, err := path.GenerateURLHash(url)
	if err != nil {
		return nil, fmt.Errorf("failed to generate URL hash: %w", err)
	}

	var prefix string
	switch contentType {
	case "image", "screenshot":
		contentType = "images"
	case "content", "text":
		contentType = "text"
	default:
		return nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	prefix = fmt.Sprintf("%s/%s", contentType, hash)

	s.logger.Debug("listing screenshots", zap.Any("prefix", prefix), zap.Any("max_items", maxItems))

	// Get the list of screenshots
	return s.storage.List(ctx, prefix, maxItems)
}

// prepareScreenshot creates a request for the screenshot API
// and returns the response
func (s *screenshotService) prepareScreenshot(opts models.ScreenshotRequestOptions) (*http.Response, error) {
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

// parseScreenshotRespose parses the screenshot response
// and returns the image, width, height, clean text, and title
func (s *screenshotService) parseScreenshotRespose(resp *http.Response) (image.Image, int, int, string, string, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, -1, -1, "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	img, width, height, err := parseImageFromResponse(resp)
	if err != nil {
		return nil, -1, -1, "", "", fmt.Errorf("failed to parse image from response: %v", err)
	}

	cleanText, title, err := parseContentFromResponse(resp)
	if err != nil {
		return nil, -1, -1, "", "", fmt.Errorf("failed to parse content from response: %v", err)
	}

	return img, width, height, cleanText, title, nil
}

// getDefaultScreenshotRequestOptions returns the default options for the screenshot request
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
		// "image/jpeg", // Not supported
		// "image/jpg", // Not supported
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
