package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// PageHistoryService is a service for page history operations
// It holds the business logic for page history management and operations
// It embeds PageHistoryRepository to interact with the database
// PageHistoryRepository inturn interacts with the page_histories table
// It embeds DiffService to perform diff operations
// And ScreenshotsService to manage screenshot operations

type PageHistoryService interface {
	// CreatePageHistory creates a page history for a page.
	// This is trigger during page creation by the page service and by workflow service.
	// It returns true if the new page history was created or it returns false if the page history already exists.
	// Error is returned if there was an issue creating the page history.
	CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, errs.Error)

	// ListPageHistory lists the history of a page, paginated by pageHistoryPaginationParam
	// This is triggered when a user wants to list all page histories of a page
	ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, errs.Error)

	// ClearPageHistory clears the history of a page.
	ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) errs.Error
}

var (
	ErrFailedToCreatePageHistory = errors.New("failed to create page history")
	ErrFailedToListPageHistory   = errors.New("failed to list page history")
	ErrFailedToClearPageHistory  = errors.New("failed to clear page history")
)
