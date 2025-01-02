package interfaces

import (
	"context"

	models "github.com/wizenheimer/iris/src/internal/models/core"
)

type ScreenshotService interface {
	// Refresh refreshes the screenshot repository with the latest screenshot and html content
	// url is the URL to take a screenshot of
	// opts are the options for the screenshot
	Refresh(ctx context.Context, url string, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error)

	// Retrieve retrieves previous screenshot and html content from the storage
	// url is the URL to retrieve the screenshot from
	Retrieve(ctx context.Context, url string) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, error)

	// GetCurrentImage retrieves the current screenshot from the storage if present
	// Or it will take a new screenshot and store it as an image
	GetCurrentImage(ctx context.Context, save bool, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, error)

	// GetCurrentHTML retrieves the current html content from the storage if present
	// Or it will take a new screenshot and store it as html
	GetCurrentHTMLContent(ctx context.Context, save bool, opts models.ScreenshotHTMLRequestOptions) (*models.ScreenshotHTMLContentResponse, error)

	// GetPreviousImage retrieves the previous screenshot from the storage
	GetPreviousImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error)

	// GetPreviousHTML retrieves the previous content of a screenshot from the storage
	GetPreviousHTMLContent(ctx context.Context, url string) (*models.ScreenshotHTMLContentResponse, error)

	// GetImage retrieves a screenshot from the storage
	GetImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error)

	// GetHTMLContent retrieves the content of a screenshot from the storage
	GetHTMLContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotHTMLContentResponse, error)

	// ListScreenshots lists the latest content (images or text) for a given URL
	ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error)
}
