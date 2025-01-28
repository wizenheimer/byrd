package report

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// ReportRepository is the interface that provides report operations
type ReportRepository interface {
	// Set creates a new report
	Set(ctx context.Context, report *models.Report) error

	// Get returns the report with the given ID.
	Get(ctx context.Context, reportID uuid.UUID) (*models.Report, error)

	// List returns a list of reports for the given workspace and competitor
	List(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, error)

	// GetLatest returns the latest report for the given workspace and competitor
	GetLatest(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Report, error)
}
