// ./src/internal/service/screenshot/interface.go
package screenshot

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type ScreenshotService interface {
	// Refresh refreshes the screenshot and html content for the given URL
	// opts are the options to use for the screenshot request
	// backDate is a flag to indicate if the responses should be backdated
	// Prior to refresh it attempts to retrieve
	Refresh(ctx context.Context, opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImage, *models.ScreenshotContent, error)

	// Retrieve retrieves the screenshot and html content for the given URL
	// opts are the options to use for the screenshot request
	// when backDate is true, it will attempt to retrieve the last available screenshot
	// when backDate is false, it will attempt to retrieve the latest screenshot
	Retrieve(ctx context.Context, opts models.ScreenshotRequestOptions, backDate bool) (*models.ScreenshotImage, *models.ScreenshotContent, error)
}
