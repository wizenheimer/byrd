package page

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

func NewPageService(pageRepo repo.PageRepository, pageHistoryService svc.PageHistoryService, logger *logger.Logger) svc.PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		logger:             logger,
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, err.Error) {
	cErr := err.New()
	pages, pErr := ps.pageRepo.AddPagesToCompetitor(ctx, competitorID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return nil, cErr
	}

	return pages, nil
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, err.Error) {
	cErr := err.New()
	page, pErr := ps.pageRepo.GetCompetitorPage(ctx, competitorID, pageID)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return models.Page{}, cErr
	}

	return page, nil
}

func (ps *pageService) GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, err.Error) {
	cErr := err.New()
	page, err := ps.GetPage(ctx, competitorID, pageID)
	if err != nil {
		cErr.Merge(err)
		return api.GetPageResponse{}, cErr
	}

	pageHistory, err := ps.pageHistoryService.ListPageHistory(ctx, pageID, historyPaginationParams)
	if err != nil {
		cErr.Merge(err)
		return api.GetPageResponse{}, cErr
	}

	return api.GetPageResponse{
		Page:    page,
		History: pageHistory,
	}, nil
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, err.Error) {
	cErr := err.New()
	updatedPage, pErr := ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return models.Page{}, cErr
	}

	return updatedPage, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, err.Error) {
	cErr := err.New()
	var limit, offset *int
	if param != nil {
		limit, offset = utils.ToPtr(param.GetLimit()), utils.ToPtr(param.GetOffset())
	}

	page, pErr := ps.pageRepo.ListCompetitorPages(ctx, competitorID, limit, offset)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return nil, cErr
	}

	return page, nil
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	// TODO: TBD
	return nil, nil
}

func (ps *pageService) ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, err.Error) {
	cErr := err.New()
	pageHistory, err := ps.pageHistoryService.ListPageHistory(
		ctx,
		pageID,
		param,
	)
	if err != nil {
		cErr.Merge(err)
		return nil, cErr
	}

	return pageHistory, nil
}

func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) err.Error {
	cErr := err.New()
	pErr := ps.pageRepo.RemovePagesFromCompetitor(
		ctx,
		competitorID,
		pageIDs,
	)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return cErr
	}
	return nil
}
