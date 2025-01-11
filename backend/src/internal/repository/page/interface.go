package page

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// PageRepository is the interface that provides page operations
// This is used to interact with the page repository

type PageRepository interface {
	// ------ CRUD operations for page ------

	// ---- Create operations for page ----

	// AddPageToCompetitor adds a page to a competitor.
	// This is used to add a page to a competitor.
	AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, page models.PageProps) (models.Page, error)

	// AddPagesToCompetitor adds pages to a competitor.
	// This is used to add multiple pages to a competitor.
	BatchAddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	// ---- Read operations for page ----

	// GetPageByID gets a page by its ID.
	// This is used to get the page details.
	GetCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) (models.Page, error)

	// BatchGetPagesByIDs gets pages by their IDs.
	// This is used to get the page details.
	BatchGetCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID, limit, offset *int) ([]models.Page, error)

	// GetCompetitorPages gets active pages for a competitor.
	// This is used to get the active pages that belong to a competitor.
	GetCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error)

	// ---- Update operations for page ----

	// UpdateCompetitorPage updates a page for a competitor
	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (models.Page, error)

	// UpdateCompetitorPageURL updates a page URL for a competitor
	UpdateCompetitorPageURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (models.Page, error)

	// UpdateCompetitorCaptureProfile updates a page capture profile for a competitor
	UpdateCompetitorCaptureProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (models.Page, error)

	// UpdateCompetitorDiffProfile updates a page diff profile for a competitor
	UpdateCompetitorDiffProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (models.Page, error)

	// UpdateCompetitorURL updates a page URL for a competitor
	UpdateCompetitorURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (models.Page, error)

	// ---- Delete operations for page ----

	// DeleteCompetitorPageByID deletes a page by its ID.
	// This is used to delete the page.
	DeleteCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) error

	// ---- Workflow operations for page ----

	// GetActivePages lists all active pages in batches
	// This is triggered when a batch of active pages is requested by page service
	// lastPageID is use to seek to the last page in the previous response
	// The table is ordered by created at so the ordering is consistent across runs
	GetActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error)
}
