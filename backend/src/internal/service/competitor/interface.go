// ./src/internal/interfaces/service/competitor.go
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
	AddCompetitorsToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error)

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
}
