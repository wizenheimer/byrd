package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type DiffRepository interface {
	// SaveDiff saves the diff analysis for the given URL
	SaveDiff(ctx context.Context, url string, diff *models.URLDiffAnalysis) error

	// GetDiff retrieves the diff analysis for the given URL, week day, and week number
	GetDiff(ctx context.Context, url, weekDay, weekNumber string) (*models.URLDiffAnalysis, error)
}
