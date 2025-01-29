package report

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// ReportService is the interface that provides report generation methods.
type ReportService interface {
	// Get returns the report with the given ID.
	Get(ctx context.Context, reportID uuid.UUID) (*models.Report, error)

	// GetLatest returns the latest report for the given workspace and competitor
	GetLatest(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Report, error)

	// List returns a list of reports for the given workspace and competitor
	List(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error)

	// Create creates a new report for the given workspace and competitor
	Create(ctx context.Context, workspaceID, competitorID uuid.UUID, history []models.PageHistory) (*models.Report, error)

	// Dispatch send the report to it's subscribers.
	Dispatch(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string, subscriberEmails []string) error
}
