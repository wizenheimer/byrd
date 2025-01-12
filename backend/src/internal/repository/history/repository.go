package history

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type historyRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageHistoryRepository(tm *transaction.TxManager, logger *logger.Logger) PageHistoryRepository {
	return &historyRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "history_repository"}),
	}
}

func (r *historyRepo) CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, diffContent any) error {
	return nil
}

func (r *historyRepo) BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	return nil, false, nil
}

func (r *historyRepo) BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error {
	return nil
}
