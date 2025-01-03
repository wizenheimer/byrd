package page

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

func NewPageService(pageRepo repo.PageRepository, pageHistoryService svc.PageHistoryService, logger *logger.Logger) svc.PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		logger:             logger,
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, []error) {
	return ps.pageRepo.AddPagesToCompetitor(ctx, competitorID, pageReq)
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, error) {
	return ps.pageRepo.GetCompetitorPage(ctx, competitorID, pageID)
}

func (ps *pageService) GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, error) {
	page, err := ps.GetPage(ctx, competitorID, pageID)
	if err != nil {
		return api.GetPageResponse{}, err
	}

	pageHistory, err := ps.pageHistoryService.ListPageHistory(ctx, pageID, historyPaginationParams)
	if err != nil {
		return api.GetPageResponse{}, err
	}

	return api.GetPageResponse{
		Page:    page,
		History: pageHistory,
	}, nil
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	return ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, pageReq)
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, error) {
	limit := param.GetLimit()
	offet := param.GetOffset()
	return ps.pageRepo.ListCompetitorPages(ctx, competitorID, &limit, &offet)
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	// TODO: TBD
	return nil, nil
}

func (ps *pageService) ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error) {
	return ps.pageHistoryService.ListPageHistory(
		ctx,
		pageID,
		param,
	)
}

func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) []error {
	return ps.pageRepo.RemovePagesFromCompetitor(
		ctx,
		competitorID,
		pageIDs,
	)
}
