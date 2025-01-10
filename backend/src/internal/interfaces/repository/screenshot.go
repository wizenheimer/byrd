// ./src/internal/interfaces/repository/screenshot.go
package interfaces

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// ScreenshotRepository is the interface that provides screenshot operations
// This is used to interact with the screenshot repository

type ScreenshotRepository interface {
	// StoreScreenshotImage stores screenshot image in the storage
	StoreScreenshotImage(ctx context.Context, data models.ScreenshotImageResponse, path string) error

	// StoreScreenshotContent stores screenshot content in the storage
	StoreScreenshotHTMLContent(ctx context.Context, data models.ScreenshotHTMLContentResponse, path string) error

	// GetContent retrieves a text content from the storage
	// Serialize the content to a string and return it
	GetScreenshotHTMLContent(ctx context.Context, path string) (models.ScreenshotHTMLContentResponse, []error)

	// GetScreenshot retrieves a screenshot from the storage
	// Deserialize the content to an image and return it
	GetScreenshotImage(ctx context.Context, path string) (models.ScreenshotImageResponse, []error)

	// Get retrieves a binary from the storage
	// Return the binary content and the metadata
	Get(ctx context.Context, path string) ([]byte, map[string]string, error)

	// Delete deletes a file from the storage
	// Return an error if the file does not exist or cannot be deleted
	Delete(ctx context.Context, path string) error

	// List lists the latest content matching the prefix
	// Return a list of ScreenshotListResponse objects or an error
	List(ctx context.Context, prefix string, maxItems int) ([]models.ScreenshotListResponse, error)
}
