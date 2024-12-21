package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type DiffRepository interface {
	// Set saves the diff analysis for the given URL
	Set(ctx context.Context, req models.URLDiffRequest, diff *models.DynamicChanges) error

	// Get retrieves the diff analysis for the given URL, week day, and week number
	Get(ctx context.Context, req models.URLDiffRequest) (*models.DynamicChanges, error)
}
