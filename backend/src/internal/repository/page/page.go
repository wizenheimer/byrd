package page

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type pageRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageRepository(tm *transaction.TxManager, logger *logger.Logger) PageRepository {
	return &pageRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "page_repository"}),
	}
}

// ---- Create operations for page ----

// AddPageToCompetitor adds a page to a competitor.
// This is used to add a page to a competitor.
func (pr *pageRepo) AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, page models.PageProps) (models.Page, error) {
	return models.Page{}, nil
}

// AddPagesToCompetitor adds pages to a competitor.
// This is used to add multiple pages to a competitor.
func (pr *pageRepo) BatchAddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	return []models.Page{}, nil
}

// ---- Read operations for page ----

// GetPageByID gets a page by its ID.
// This is used to get the page details.
func (pr *pageRepo) GetCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) (models.Page, error) {
	return models.Page{}, nil
}

// BatchGetPagesByIDs gets pages by their IDs.
// This is used to get the page details.
func (pr *pageRepo) BatchGetCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return []models.Page{}, nil
}

// GetCompetitorPages gets active pages for a competitor.
// This is used to get the active pages that belong to a competitor.
func (pr *pageRepo) GetCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return []models.Page{}, nil
}

// ---- Update operations for page ----

// UpdateCompetitorPage updates a page for a competitor
func (pr *pageRepo) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (models.Page, error) {
	return models.Page{}, nil
}

// UpdateCompetitorPageURL updates a page URL for a competitor
func (pr *pageRepo) UpdateCompetitorPageURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (models.Page, error) {
	return models.Page{}, nil
}

// UpdateCompetitorCaptureProfile updates a page capture profile for a competitor
func (pr *pageRepo) UpdateCompetitorCaptureProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (models.Page, error) {
	return models.Page{}, nil
}

// UpdateCompetitorDiffProfile updates a page diff profile for a competitor
func (pr *pageRepo) UpdateCompetitorDiffProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (models.Page, error) {
	return models.Page{}, nil
}

// UpdateCompetitorURL updates a page URL for a competitor
func (pr *pageRepo) UpdateCompetitorURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (models.Page, error) {
	return models.Page{}, nil
}

// ---- Delete operations for page ----

// DeleteCompetitorPageByID deletes a page by its ID.
// This is used to delete the page.
func (pr *pageRepo) DeleteCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) error {
	return nil
}

// ---- Workflow operations for page ----

// GetActivePages lists all active pages in batches
// This is triggered when a batch of active pages is requested by page service
// lastPageID is use to seek to the last page in the previous response
// The table is ordered by created at so the ordering is consistent across runs
func (pr *pageRepo) GetActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error) {
	return models.ActivePageBatch{}, nil
}
