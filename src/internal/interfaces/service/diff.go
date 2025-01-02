package interfaces

import (
	"context"

	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// DiffService is the interface that provides diff operations
type DiffService interface {
	// Compare: compares two HTML contents and returns the differences using the given profile
	Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, error)
}
