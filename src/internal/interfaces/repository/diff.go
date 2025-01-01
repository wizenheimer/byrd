package interfaces

import (
	"context"

	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

type DiffRepository interface {
	// Set saves the diff analysis for the given URL
	Set(ctx context.Context, req api_models.URLDiffRequest, diff *core_models.DynamicChanges) error

	// Get retrieves the diff analysis for the given URL, week day, and week number
	Get(ctx context.Context, req api_models.URLDiffRequest) (*core_models.DynamicChanges, error)
}
