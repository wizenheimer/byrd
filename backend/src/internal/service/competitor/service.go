// ./src/internal/service/competitor/service.go
package competitor

import (
	"context"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type competitorService struct {
	competitorRepository competitor.CompetitorRepository
	pageService          page.PageService
	nameFinder           *CompanyNameFinder
	logger               *logger.Logger
}

var _ CompetitorService = (*competitorService)(nil)

func NewCompetitorService(competitorRepository competitor.CompetitorRepository, pageService page.PageService, logger *logger.Logger) CompetitorService {
	return &competitorService{
		competitorRepository: competitorRepository,
		pageService:          pageService,
		logger:               logger,
		nameFinder:           NewCompanyNameFinder(logger),
	}
}

// CreateCompetitor creates a competitor in a workspace.
// It returns the created competitor and any errors that occurred.
// This is triggered when a user creates a competitor during workspace creation
// or otherwise.
func (cs *competitorService) CreateCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorReq api.CreateCompetitorRequest) (api.CreateWorkspaceCompetitorResponse, error) {
	return api.CreateWorkspaceCompetitorResponse{}, nil
}

// GetCompetitor gets a competitor by ID.
// It returns the competitor if it exists, otherwise it returns an error.
// This is triggered when a user wants to get a competitor by ID.
// Pagination is used to limit the number of pages returned per competitor.
func (cs *competitorService) GetCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pagePaginationParam api.PaginationParams) (api.GetWorkspaceCompetitorResponse, error) {
	return api.GetWorkspaceCompetitorResponse{}, nil
}

// CompetitorExists checks if a competitor exists in a workspace.
// Its optimized for quick lookups over the competitor table.
func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	return false, nil
}

// PageExists checks if a page exists in a competitor.
// Its optimized for quick lookups over the page table.
func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	return false, nil
}

// ListWorkspaceCompetitors lists the competitors of a workspace.
// It returns the competitors of the workspace.
// This is triggered when a user wants to list all competitors in a workspace.
// Pagination is used to limit the number of competitors returned.
// The pages are returned in their entirety for every competitor requested.
func (cs *competitorService) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorListingParam api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, error) {
	return []api.GetWorkspaceCompetitorResponse{}, nil
}

// RemoveCompetitors removes a list of competitors from a workspace.
// Removing a competitor also removes all its pages from the workspace.
// It returns an error if the competitor could not be removed.
// When competitorIDs are nil, all competitors and their pages
// are removed from the workspace.
func (cs *competitorService) RemoveCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error {
	return nil
}

// <--------- Page Management --------->
// Page management is operational within the context of a competitor
// Creat, Read, Update, Delete operations for pages

// AddPagesToCompetitor adds a list of pages to an existing competitor.
// It returns the created page and any errors that occurred.
// This is triggered when a user adds a page to an existing competitor.
func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, error) {
	return []models.Page{}, nil
}

// GetCompetitorPage gets a competitor's page by ID.
// It returns the page if it exists, otherwise it returns an error.
// This is triggered when a user wants to get a page by ID.
// Pagination is used to limit the number of page history returned per page.
func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) (api.GetPageResponse, error) {
	return api.GetPageResponse{}, nil
}

// UpdatePage updates a page.
// It returns the updated page and any errors that occurred.
// This is triggered when a user updates a page in a competitor.
func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	return models.Page{}, nil
}

// RemovePagesFromCompetitor removes pages from a competitor.
// It returns an error if the page could not be removed.
// This is triggered when a user removes a page from a competitor.
// When pageIDs are nil, all pages are removed from the competitor.
func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) error {
	return nil
}

// ListCompetitorPages lists the pages of a competitor.
// It returns the pages of the competitor.
// This is triggered when a user wants to list all pages in a competitor.
// This is helpful for reporting workflow.
func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID) ([]models.Page, error) {
	return []models.Page{}, nil
}
