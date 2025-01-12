// ./src/internal/service/page/service.go
package page

import (
	"context"

	"github.com/google/uuid"
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

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	return nil, nil
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error) {
	return nil, nil
}

func (ps *pageService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
	return nil, nil
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return nil, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return nil, nil
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	return nil, nil
}

func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) error {
	return nil
}

func (ps *pageService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	return false, nil
}
