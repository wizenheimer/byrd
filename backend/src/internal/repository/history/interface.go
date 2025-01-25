// ./src/internal/repository/history/interface.go
package history

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// PageHistoryRepository is the interface that provides page history operations
// This is used to interact with the page history repository

type PageHistoryRepository interface {
	CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, diffContent any, prev, curr string) error

	BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error)

	BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error
}
