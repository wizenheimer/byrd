package competitor

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type competitorRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewCompetitorRepository(tm *transaction.TxManager, logger *logger.Logger) CompetitorRepository {
	return &competitorRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "competitor_repository"}),
	}
}

func (r *competitorRepo) CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) BatchCreateCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) BatchGetCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return nil, nil
}

func (r *competitorRepo) RemoveCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	return nil
}

func (r *competitorRepo) BatchRemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorID []uuid.UUID) error {
	return nil
}

func (r *competitorRepo) RemoveAllCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	return nil
}

func (r *competitorRepo) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	return false, nil
}
