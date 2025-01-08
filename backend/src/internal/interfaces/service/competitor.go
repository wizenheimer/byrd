package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
)

// CompetitorService embeds competitor repository and page service
// It holds the business logic for competitor management
// PageService is embedded to manage pages within the context of a competitor
type CompetitorService interface {
	// <---------  Competitor Management --------->
	// Create, read, update, delete operations for competitors

	// CreateCompetitor creates a competitor in a workspace.
	// It returns the created competitor and any errors that occurred.
	// This is triggered when a user creates a competitor during workspace creation
	// or otherwise.
	CreateCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorReq api.CreateCompetitorRequest) (api.CreateWorkspaceCompetitorResponse, err.Error)

	// GetCompetitor gets a competitor by ID.
	// It returns the competitor if it exists, otherwise it returns an error.
	// This is triggered when a user wants to get a competitor by ID.
	// Pagination is used to limit the number of pages returned per competitor.
	GetCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pagePaginationParam api.PaginationParams) (api.GetWorkspaceCompetitorResponse, err.Error)

	// CompetitorExists checks if a competitor exists in a workspace.
	// Its optimized for quick lookups over the competitor table.
	CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, err.Error)

	// PageExists checks if a page exists in a competitor.
	// Its optimized for quick lookups over the page table.
	PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, err.Error)

	// ListWorkspaceCompetitors lists the competitors of a workspace.
	// It returns the competitors of the workspace.
	// This is triggered when a user wants to list all competitors in a workspace.
	// Pagination is used to limit the number of competitors returned.
	// The pages are returned in their entirety for every competitor requested.
	ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorListingParam api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, err.Error)

	// RemoveCompetitors removes a list of competitors from a workspace.
	// Removing a competitor also removes all its pages from the workspace.
	// It returns an error if the competitor could not be removed.
	// When competitorIDs are nil, all competitors and their pages
	// are removed from the workspace.
	RemoveCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) err.Error

	// <--------- Page Management --------->
	// Page management is operational within the context of a competitor
	// Creat, Read, Update, Delete operations for pages

	// AddPagesToCompetitor adds a list of pages to an existing competitor.
	// It returns the created page and any errors that occurred.
	// This is triggered when a user adds a page to an existing competitor.
	AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, err.Error)

	// GetCompetitorPage gets a competitor's page by ID.
	// It returns the page if it exists, otherwise it returns an error.
	// This is triggered when a user wants to get a page by ID.
	// Pagination is used to limit the number of page history returned per page.
	GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) (api.GetPageResponse, err.Error)

	// UpdatePage updates a page.
	// It returns the updated page and any errors that occurred.
	// This is triggered when a user updates a page in a competitor.
	UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, err.Error)

	// RemovePagesFromCompetitor removes pages from a competitor.
	// It returns an error if the page could not be removed.
	// This is triggered when a user removes a page from a competitor.
	// When pageIDs are nil, all pages are removed from the competitor.
	RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) err.Error

	// ListCompetitorPages lists the pages of a competitor.
	// It returns the pages of the competitor.
	// This is triggered when a user wants to list all pages in a competitor.
	// This is helpful for reporting workflow.
	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID) ([]models.Page, err.Error)
}

var (
	ErrFailedToListWorkspaceCompetitors   = errors.New("failed to list workspace competitors")
	ErrFailedToCreateWorkspaceCompetitors = errors.New("failed to create workspace competitors")
	ErrFailedToAddPagesToCompetitor       = errors.New("failed to add pages to competitor")
	ErrFailedToGetCompetitorForWorkspace  = errors.New("failed to get competitor for workspace")
	ErrFailedToListPagesForCompetitor     = errors.New("failed to list pages for competitor")

	ErrFailedToRemoveWorkspaceCompetitors = errors.New("failed to remove workspace competitors")
)
