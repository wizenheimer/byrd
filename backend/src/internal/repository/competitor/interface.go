// ./src/internal/repository/competitor/interface.go
package competitor

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// CompetitorRepository is the interface that provides competitor operations
// This is used to interact with the competitor repository

type CompetitorRepository interface {
	CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorName string) (*models.Competitor, error)

	BatchCreateCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, error)

	GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error)

	BatchGetCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error)

	ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error)

	UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error)

	RemoveCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error

	BatchRemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorID []uuid.UUID) error

	RemoveAllCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)
}
