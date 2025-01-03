package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// CompetitorRepository is the interface that provides competitor operations
// This is used to interact with the competitor repository

type CompetitorRepository interface {
	// CreateCompetitors creates competitors in a workspace
	CreateCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, error)

	// GetCompetitor gets a competitor by its ID
	GetCompetitor(ctx context.Context, competitorID uuid.UUID) (models.Competitor, error)

	// ListWorkspaceCompetitors lists all competitors in a workspace
	ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]models.Competitor, error)

	// RemoveWorkspaceCompetitor removes a competitor from a workspace
	// When competitorIDs are nil, all competitors are removed from the workspace
	RemoveWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) []error

	// WorkspaceCompetitorExists checks if a competitor exists in a workspace
	// This is optimized for quick lookups over the competitor table
	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)
}
