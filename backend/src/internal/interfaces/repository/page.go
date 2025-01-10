package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// PageRepository is the interface that provides page operations
// This is used to interact with the page repository

type PageRepository interface {
	// AddPagesToCompetitor adds pages to a competitor
	AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, errs.Error)

	// RemovePagesFromCompetitor removes pages from a competitor
	// When pageIDs are nil, all pages are removed from the competitor
	RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) errs.Error

	// GetCompetitorPages gets the pages for a competitor
	// This is used to get the pages that belong to a competitor
	// When limit and offset are nil, all pages are returned
	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, errs.Error)

	// ListActivePages lists all active pages in batches
	// This is triggered when a batch of active pages is requested by page service
	// lastPageID is use to seek to the last page in the previous response
	// The table is ordered by created at so the ordering is consistent across runs
	ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, errs.Error)

	// GetCompetitorPage gets a page for a competitor
	GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (models.Page, errs.Error)

	// UpdateCompetitorPage updates a page for a competitor
	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (models.Page, errs.Error)
}
