// ./src/internal/service/screenshot/utils.go
package screenshot

import (
	"reflect"

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
