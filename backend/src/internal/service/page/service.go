// ./src/internal/service/page/service.go
package page

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

func NewPageService(pageRepo repo.PageRepository, pageHistoryService svc.PageHistoryService, logger *logger.Logger) svc.PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		logger:             logger,
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, errs.Error) {
	cErr := errs.New()
	pages, pErr := ps.pageRepo.AddPagesToCompetitor(ctx, competitorID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return nil, cErr
	}

	return pages, nil
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, errs.Error) {
	cErr := errs.New()
	page, pErr := ps.pageRepo.GetCompetitorPage(ctx, competitorID, pageID)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return models.Page{}, cErr
	}

	return page, nil
}

func (ps *pageService) GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, errs.Error) {
	cErr := errs.New()
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

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, errs.Error) {
	cErr := errs.New()
	updatedPage, pErr := ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		return models.Page{}, cErr
	}

	return updatedPage, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, errs.Error) {
	cErr := errs.New()
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

func (ps *pageService) ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, errs.Error) {
	cErr := errs.New()
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

func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) errs.Error {
	cErr := errs.New()
	pErr := ps.pageRepo.RemovePagesFromCompetitor(
		ctx,
		competitorID,
		pageIDs,
	)
	if pErr != nil && pErr.HasErrors() {
		cErr.Merge(pErr)
		ps.logger.Error("Failed to remove pages from competitor", zap.Error(pErr))
		return cErr
	}
	return nil
}
