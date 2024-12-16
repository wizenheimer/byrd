package interfaces

import (
	"context"
	"image"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type ScreenshotService interface {
	// CaptureScreenshot takes a screenshot of the given URL
	CaptureScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, image.Image, string, error)

	// GetPreviousScreenshot retrieves the previous screenshot from the storage
	GetPreviousScreenshotImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error)

	// GetPreviousScreenshotContent retrieves the previous content of a screenshot from the storage
	GetPreviousScreenshotContent(ctx context.Context, url string) (*models.ScreenshotContentResponse, error)

	// GetScreenshot retrieves a screenshot from the storage
	GetScreenshotImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error)

	// GetContent retrieves the content of a screenshot from the storage
	GetScreenshotContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotContentResponse, error)

	// ListScreenshots lists the latest content (images or text) for a given URL
	ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error)
}
