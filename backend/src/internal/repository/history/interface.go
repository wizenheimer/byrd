package history

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// PageHistoryRepository is the interface that provides page history operations
// This is used to interact with the page history repository

type PageHistoryRepository interface {
	// ---- CRUD operations for page history ----

	// CreatePageHistory creates a new page history.
	// The page history is created with the provided page ID and page history.
	CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, error)

	// BatchGetPageHistory lists page history for a page ordered by created at.
	// When limit and offset are nil, all page history is returned.
	BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error)

	// RemovePageHistory removes page history for a list of pages.
	// Returns an error if pageIDs are nil.
	BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error
}
