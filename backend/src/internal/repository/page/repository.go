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

func (r *pageRepo) AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) BatchAddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	return nil, nil
}

func (r *pageRepo) GetCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) BatchGetCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return nil, nil
}

func (r *pageRepo) GetCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return nil, nil
}

func (r *pageRepo) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) UpdateCompetitorPageURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) UpdateCompetitorCaptureProfile(ctx context.Context, competitorID, pageID uuid.UUID, captureProfile *models.ScreenshotRequestOptions, url string) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) UpdateCompetitorDiffProfile(ctx context.Context, competitorID, pageID uuid.UUID, diffProfile []string) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) UpdateCompetitorURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error) {
	return nil, nil
}

func (r *pageRepo) DeleteCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) error {
	return nil
}

func (r *pageRepo) BatchDeleteCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) error {
	return nil
}

func (r *pageRepo) DeleteAllCompetitorPages(ctx context.Context, competitorID uuid.UUID) error {
	return nil
}

func (r *pageRepo) BatchDeleteAllCompetitorPages(ctx context.Context, competitorIDs []uuid.UUID) error {
	return nil
}

func (r *pageRepo) GetActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error) {
	return models.ActivePageBatch{}, nil
}

func (r *pageRepo) GetPageByPageID(ctx context.Context, pageID uuid.UUID) (*models.Page, error) {
	return nil, nil
}
