package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type DiffService interface {
	// CompareScreenshotContents compares the contents of two screenshots using AI services
	Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, error)

	// Get returns the diff analysis for the given URL
	Get(ctx context.Context, req models.URLDiffRequest) (*models.DynamicChanges, error)
}
