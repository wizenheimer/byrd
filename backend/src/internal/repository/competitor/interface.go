package competitor

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// CompetitorRepository is the interface that provides competitor operations
// This is used to interact with the competitor repository

type CompetitorRepository interface {

	// ------ CRUD operations for competitor ------

	// ---- Create operations for competitor ----

	// CreateCompetitor creates a competitor
	CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorName string) (models.Competitor, error)

	// BatchCreateCompetitorsForWorkspace creates multiple competitors
	BatchCreateCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, error)

	// ---- Read operations for competitor ----

	// GetCompetitorForWorkspace gets a competitor by its ID
	GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (models.Competitor, error)

	// BatchGetCompetitors gets competitors by their IDs
	BatchGetCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error)

	// ListCompetitorsForWorkspace lists competitors for a workspace
	ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, error)

	// ---- Update operations for competitor ----

	// UpdateCompetitorForWorkspace updates a competitor
	UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (models.Competitor, error)

	// ---- Delete operations for competitor ----

	// RemoveCompetitorForWorkspace removes a competitor
	RemoveCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error

	// ---- Optimized Lookup for Competitors ----

	// WorkspaceCompetitorExists checks if a competitor exists in a workspace
	// This is optimized for quick lookups over the competitor table
	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)
}
