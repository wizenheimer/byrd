package interfaces

import (
	"context"

	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

type DiffService interface {
	// CompareScreenshotContents compares the contents of two screenshots using AI services
	Compare(ctx context.Context, content1, content2 *core_models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*core_models.DynamicChanges, error)

	// Get returns the diff analysis for the given URL
	Get(ctx context.Context, req api_models.URLDiffRequest) (*core_models.DynamicChanges, error)
}
