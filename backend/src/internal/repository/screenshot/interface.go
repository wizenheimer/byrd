// ./src/internal/repository/screenshot/interface.go
package screenshot

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// ScreenshotRepository is the interface that provides screenshot operations
// This is used to interact with the screenshot repository

type ScreenshotRepository interface {
	// StoreScreenshotImage stores screenshot image in the storage
	StoreScreenshotImage(ctx context.Context, data *models.ScreenshotImage) error

	// StoreScreenshotContent stores screenshot content in the storage
	StoreScreenshotContent(ctx context.Context, data *models.ScreenshotContent) error

	// RetrieveScreenshotImage retrieves screenshot image from the storage
	RetrieveScreenshotImage(ctx context.Context, path string) (*models.ScreenshotImage, error)

	// RetrieveScreenshotContent retrieves screenshot content from the storage
	RetrieveScreenshotContent(ctx context.Context, path string) (*models.ScreenshotContent, error)
}
