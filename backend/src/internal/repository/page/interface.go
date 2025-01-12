package page

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// PageRepository is the interface that provides page operations
// This is used to interact with the page repository

type PageRepository interface {
	AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, page models.PageProps) (*models.Page, error)

	BatchAddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	GetCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error)

	BatchGetCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID, limit, offset *int) ([]models.Page, error)

	GetCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error)

	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)

	UpdateCompetitorPageURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error)

	UpdateCompetitorCaptureProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (*models.Page, error)

	UpdateCompetitorDiffProfile(ctx context.Context, competitorID, pageID uuid.UUID, name string) (*models.Page, error)

	UpdateCompetitorURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error)

	DeleteCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) error

	GetActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error)
}
