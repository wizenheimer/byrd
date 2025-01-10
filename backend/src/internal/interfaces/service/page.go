// ./src/internal/interfaces/service/page.go
package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// PageService is the service that manages pages
// It holds the business logic for page management
// It embeds PageRepository to interact with the database
// PageRepository inturn interacts primarily with the pages table
// It embeds PageHistoryService to manage page history for a page

type PageService interface {
	// <--------- Page Management --------->
	// CreatePage creates a page for a competitor
	// It returns the created page and any errors that occurred
	// This is triggered when a user creates a page for a competitor
	CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, errs.Error)

	// GetPage gets a page by ID.
	// It returns the page if it exists, otherwise it returns an error.
	GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, errs.Error)

	// GetPageWithHistory gets a page along with its history.
	// It returns the page if it exists, otherwise it returns an error.
	GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, errs.Error)

	// UpdatePage updates a page.
	// It returns the updated page and any errors that occurred.
	// This is triggered when a user updates a page.
	UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, errs.Error)

	// ListCompetitorPages lists the pages of a competitor.
	// It returns the pages of the competitor
	// This is triggered when a user wants to list all pages in a competitor.
	// Pagination is used to limit the number of pages returned.
	// When pagination param is nil all the pages are returned for the competitor.
	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, errs.Error)

	// ListActivePages lists all active pages.
	// It returns the active pages.
	// This is triggered by workflow to list all active pages.
	// lastPageID serves as a checkpoint and is used to seek to the last page in the previous batch.
	// batchSize is used to limit the number of pages returned.
	ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error)

	// ListPageHistory lists the history of a page
	// It returns the history of the page.
	// This is triggered when a user wants to list the history of a page.
	// Pagination is used to limit the number of history returned.
	ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, errs.Error)

	// RemovePage removes pages from a competitor
	// It returns any errors that occurred
	// This is triggered when a user wants to remove pages from a competitor
	// When pageIDs are nil all pages are removed from the competitor.
	RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) errs.Error
}

var (
	ErrFailedToCreatePageForCompetitor  = errors.New("failed to create page")
	ErrFailedToGetPageForCompetitor     = errors.New("failed to get page")
	ErrFailedToUpdatePageForCompetitor  = errors.New("failed to update page")
	ErrFaileToListPagesForCompetitor    = errors.New("failed to list pages")
	ErrFailedToRemovePageFromCompetitor = errors.New("failed to remove page")
)
