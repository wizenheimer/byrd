// ./src/internal/service/competitor/interface.go
package competitor

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

var (
	maxCompetitorBatchSize int = 25
)

// CompetitorService embeds competitor repository and page service
// It holds the business logic for competitor management
// PageService is embedded to manage pages within the context of a competitor
type CompetitorService interface {
	// BatchCreateCompetitorsForWorkspace creates multiple competitors for a workspace
	// Each page in the pages slice is used to create a competitor
	BatchCreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error)

	// CreateCompetitorForWorkspace creates a competitor for a workspace
	// Here all the pages are used to create a single competitor
	CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (models.Competitor, error)

	GetCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error)

	ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error)

	UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error)

	RemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error

	CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)

	PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error)

	AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error)

	UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)

	RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) error

	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error)

	ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error)

	// ListReports lists the reports for a competitor.
	ListReports(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error)

	// CreateReport creates a report for a competitor.
	CreateReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) (*models.Report, error)

	// DispatchReport dispatches a report for a competitor.
	DispatchReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID, subscriberEmails []string) error

	CountPagesForCompetitors(ctx context.Context, competitorIDs []uuid.UUID) (int, error)

	CountCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID) (int, error)
}
