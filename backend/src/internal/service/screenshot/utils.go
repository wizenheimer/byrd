// ./src/internal/service/screenshot/utils.go
package screenshot

import (
	"context"
	"errors"
	"image"
	"io"
	"net/http"
	"reflect"
	"strconv"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// getDefaultScreenshotRequestOptions returns the default options for the screenshot request
func GetDefaultScreenshotRequestOptions(url string) models.ScreenshotRequestOptions {
	// Get default options
	defaultOpt := models.ScreenshotRequestOptions{
		URL: url,
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
func MergeScreenshotRequestOptions(defaults, override models.ScreenshotRequestOptions) models.ScreenshotRequestOptions {
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

// refreshScreenshot refreshes the screenshot and html content for the given URL
// it ensures that the screenshot and content are fetched and aren't null before returning
func (s *screenshotService) refreshScreenshot(_ context.Context, opts models.ScreenshotRequestOptions) (*image.Image, *string, error) {
	defaultOpt := GetDefaultScreenshotRequestOptions(opts.URL)
	mergedOpt := MergeScreenshotRequestOptions(defaultOpt, opts)

	resp, err := s.httpClient.NewRequest().
		BaseURL(s.origin).
		Method(http.MethodGet).
		Path("/take").
		QueryParam("access_key", s.key).
		AddQueryParamsFromStruct(mergedOpt).
		Execute(s.httpClient)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("failed to retrieve screenshot")
	}

	img, err := s.getScreenshot(resp)
	if err != nil {
		return nil, nil, err
	} else if img == nil {
		return nil, nil, errors.New("failed to retrieve screenshot")
	}

	content, err := s.getContent(resp)
	if err != nil {
		return nil, nil, err
	} else if content == nil {
		return nil, nil, errors.New("failed to retrieve content")
	}

	return img, content, nil
}

func getScreenshotMetadata(backDate bool) (*models.ScreenshotMetadata, error) {
	var currentDayString string
	var currentWeek int
	var currentYear int
	if backDate {
		currentYear, currentWeek, currentDayString = getPreviousTimeComponents(true)
	} else {
		currentYear, currentWeek, currentDayString = getCurrentTimeComponents(true)
	}
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotMetadata{
		Year:       currentYear,
		WeekNumber: currentWeek,
		WeekDay:    currentDay,
	}, nil
}

func (s *screenshotService) getScreenshot(resp *http.Response) (*image.Image, error) {
	imageContentTypes := []string{
		"image/png",
	}

	if !utils.Contains(imageContentTypes, resp.Header.Get("Content-Type")) {
		return nil, errors.New("received unexpected content type")
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (s *screenshotService) getContent(resp *http.Response) (*string, error) {
	renderedURL := resp.Header.Get("X-ScreenshotOne-Content-URL")
	if renderedURL == "" {
		return nil, errors.New("no content URL found in headers, cannot proceed with rendering")
	}

	htmlResp, err := http.Get(renderedURL)
	if err != nil {
		return nil, err
	}
	defer htmlResp.Body.Close()
	if htmlResp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve content")
	}

	htmlBytes, err := io.ReadAll(htmlResp.Body)
	if err != nil {
		return nil, err
	}

	content := string(htmlBytes)
	return &content, nil
}
