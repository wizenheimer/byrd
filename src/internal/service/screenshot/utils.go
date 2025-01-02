package screenshot

import (
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"reflect"
	"strconv"

	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/utils"
	"go.uber.org/zap"
)

// prepareImageResponse prepares the image response from the HTTP response
func (s *screenshotService) prepareScreenshotImageResponse(resp *http.Response, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	img, width, height, err := parseImageFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image from response: %v", err)
	}

	renderedURL := resp.Header.Get("X-ScreenshotOne-Content-URL")
	if renderedURL == "" {
		return nil, fmt.Errorf("no content URL found in headers")
	}

	screenshotResponse := models.ScreenshotImageResponse{
		Status: "success",
		Image:  img,
		Metadata: &models.ScreenshotMetadata{
			SourceURL:   url,
			RenderedURL: renderedURL,
			Year:        year,
			WeekNumber:  weekNumber,
			WeekDay:     weekDay,
		},
		ImageHeight: utils.ToPtr(height),
		ImageWidth:  utils.ToPtr(width),
	}

	return &screenshotResponse, nil
}

// prepareHTMLContentResponse prepares the HTML content response from the HTTP response
func (s *screenshotService) prepareScreenshotHTMLContentResponse(resp *http.Response, sourceURL, renderedURL string, year int, weekNumber int, weekDay int) (*models.ScreenshotHTMLContentResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the HTML content
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTML body: %v", err)
	}

	htmlContent := string(htmlBytes)

	htmlResponse := models.ScreenshotHTMLContentResponse{
		Status:      "success",
		HTMLContent: htmlContent,
		Metadata: &models.ScreenshotMetadata{
			SourceURL:   sourceURL,
			RenderedURL: renderedURL,
			Year:        year,
			WeekNumber:  weekNumber,
			WeekDay:     weekDay,
		},
	}

	return &htmlResponse, nil
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

	if !utils.Contains(imageContentTypes, resp.Header.Get("Content-Type")) {
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

// getExistingScreenshotImage retrieves the existing screenshot image from the storage
// and returns the image and metadata
func (s *screenshotService) getExistingScreenshotImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error) {
	screenshotPath, err := utils.GetCurrentScreenshotPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get current screenshot path: %v", err)
	}

	var screenshotImageResp models.ScreenshotImageResponse

	screenshotImageResp, errs := s.storage.GetScreenshotImage(ctx, screenshotPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get existing screenshot: %v", errs)
	}
	return &screenshotImageResp, nil
}

// getExistingHTMLContent retrieves the existing screenshot content from the storage
// and returns the content and metadata
func (s *screenshotService) getExistingHTMLContent(ctx context.Context, url string) (*models.ScreenshotHTMLContentResponse, error) {
	contentPath, err := utils.GetCurrentContentPath(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get current content path: %v", err)
	}

	var screenshotContentResponse models.ScreenshotHTMLContentResponse

	screenshotContentResponse, errs := s.storage.GetScreenshotHTMLContent(ctx, contentPath)
	if errs != nil {
		return nil, fmt.Errorf("failed to get existing screenshot: %v", errs)
	}

	return &screenshotContentResponse, nil
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

// prepareScreenshotHTML creates a request for the screenshot HTML
func (s *screenshotService) prepareScreenshotHTML(opts models.ScreenshotHTMLRequestOptions) (*http.Response, error) {
	// Get HTML content
	htmlResp, err := http.Get(opts.RenderedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML content: %v", err)
	}

	if htmlResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", htmlResp.StatusCode)
	}

	// Prepare the request
	return htmlResp, nil
}

// getDefaultScreenshotRequestOptions returns the default options for the screenshot request
func getDefaultScreenshotRequestOptions() models.ScreenshotRequestOptions {
	// Get default options
	defaultOpt := models.ScreenshotRequestOptions{
		// Capture options
		Format:                utils.ToPtr("png"),
		ImageQuality:          utils.ToPtr(80),
		CaptureBeyondViewport: utils.ToPtr(true),
		FullPage:              utils.ToPtr(true),
		FullPageAlgorithm:     utils.ToPtr(models.FullPageAlgorithmDefault),

		// Resource blocking options
		BlockAds:                 utils.ToPtr(true),
		BlockCookieBanners:       utils.ToPtr(true),
		BlockBannersByHeuristics: utils.ToPtr(true),
		BlockTrackers:            utils.ToPtr(true),
		BlockChats:               utils.ToPtr(true),

		// Wait and delay options
		Delay:             utils.ToPtr(0),
		Timeout:           utils.ToPtr(60),
		NavigationTimeout: utils.ToPtr(30),
		WaitUntil: []models.WaitUntilOption{
			models.WaitUntilNetworkIdle2,
			models.WaitUntilNetworkIdle0,
		},

		// Styling options
		DarkMode:      utils.ToPtr(false),
		ReducedMotion: utils.ToPtr(true),

		// Response options
		MetadataImageSize:      utils.ToPtr(true),
		MetadataPageTitle:      utils.ToPtr(true),
		MetadataContent:        utils.ToPtr(true),
		MetadataHttpStatusCode: utils.ToPtr(true),
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
