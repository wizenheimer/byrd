package history

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// PageHistoryRepository is the interface that provides page history operations
// This is used to interact with the page history repository

type PageHistoryRepository interface {
	CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (*models.PageHistory, error)

	BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error)

	BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error
}
