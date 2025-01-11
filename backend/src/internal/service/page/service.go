// ./src/internal/service/page/service.go
package page

import (
	"context"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/page"
	"github.com/wizenheimer/byrd/src/internal/service/history"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// compile time check if the interface is implemented
var _ PageService = (*pageService)(nil)

type pageService struct {
	pageRepo           page.PageRepository
	pageHistoryService history.PageHistoryService
	logger             *logger.Logger
}

func NewPageService(pageRepo page.PageRepository, pageHistoryService history.PageHistoryService, logger *logger.Logger) PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		logger:             logger,
	}
}

// CreatePage creates a page for a competitor
// It returns the created page and any errors that occurred
// This is triggered when a user creates a page for a competitor
func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, error) {
	return nil, nil
}

// GetPage gets a page by ID.
// It returns the page if it exists, otherwise it returns an error.
func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, error) {
	return models.Page{}, nil
}

// GetPageWithHistory gets a page along with its history.
// It returns the page if it exists, otherwise it returns an error.
func (ps *pageService) GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, error) {
	return api.GetPageResponse{}, nil
}

// UpdatePage updates a page.
// It returns the updated page and any errors that occurred.
// This is triggered when a user updates a page.
func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	return models.Page{}, nil
}

// ListCompetitorPages lists the pages of a competitor.
// It returns the pages of the competitor
// This is triggered when a user wants to list all pages in a competitor.
// Pagination is used to limit the number of pages returned.
// When pagination param is nil all the pages are returned for the competitor.
func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, error) {
	return nil, nil
}

// ListActivePages lists all active pages.
// It returns the active pages.
// This is triggered by workflow to list all active pages.
// lastPageID serves as a checkpoint and is used to seek to the last page in the previous batch.
// batchSize is used to limit the number of pages returned.
func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	return nil, nil
}

// ListPageHistory lists the history of a page
// It returns the history of the page.
// This is triggered when a user wants to list the history of a page.
// Pagination is used to limit the number of history returned.
func (ps *pageService) ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error) {
	return nil, nil
}

// RemovePage removes pages from a competitor
// It returns any errors that occurred
// This is triggered when a user wants to remove pages from a competitor
// When pageIDs are nil all pages are removed from the competitor.
func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) error {
	return nil
}
