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

// CreatePageHistory creates a new page history.
// The page history is created with the provided page ID and page history.
func (r *historyRepo) CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, error) {
	return models.PageHistory{}, nil
}

// BatchGetPageHistory lists page history for a page ordered by created at.
// When limit and offset are nil, all page history is returned.
func (r *historyRepo) BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
	return []models.PageHistory{}, nil
}

// RemovePageHistory removes page history for a list of pages.
// Returns an error if pageIDs are nil.
func (r *historyRepo) BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error {
	return nil
}
