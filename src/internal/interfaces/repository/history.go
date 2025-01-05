package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// PageHistoryRepository is the interface that provides page history operations
// This is used to interact with the page history repository

type PageHistoryRepository interface {
	// CreatePageHistory creates a new page history
	// The page history is created with the provided page ID and page history
	CreatePageHistory(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, error)

	// PageHistoryExists checks if a page history exists
	// This is optimized for quick lookups over the page history table
	// PageHistoryExists(ctx context.Context, pageID string, weekNumber1, weekNumber2, yearNumber1, yearNumber2 int, bucketID1, bucketID2 string) (bool, error)

	// ListPageHistory lists page history for a page ordered by created at
	// When limit and offset are nil, all page history is returned
	ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, []error)

	// RemovePageHistory removes page history for a list of pages
	// Returns an error if pageIDs are nil
	RemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) []error
}
