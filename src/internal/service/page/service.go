package page

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

var (
	ErrFailedToCreatePageForCompetitor  = errors.New("failed to create page")
	ErrFailedToGetPageForCompetitor     = errors.New("failed to get page")
	ErrFailedToListPageHistory          = errors.New("failed to list page history")
	ErrFailedToUpdatePageForCompetitor  = errors.New("failed to update page")
	ErrFaileToListPagesForCompetitor    = errors.New("failed to list pages")
	ErrFailedToRemovePageFromCompetitor = errors.New("failed to remove page")
)

func NewPageService(pageRepo repo.PageRepository, pageHistoryService svc.PageHistoryService, logger *logger.Logger) svc.PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		logger:             logger,
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, []error) {
	pages, pErr := ps.pageRepo.AddPagesToCompetitor(ctx, competitorID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		return nil, []error{ErrFailedToCreatePageForCompetitor}
	}

	return pages, nil
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (models.Page, error) {
	page, pErr := ps.pageRepo.GetCompetitorPage(ctx, competitorID, pageID)
	if pErr != nil && pErr.HasErrors() {
		return models.Page{}, ErrFailedToGetPageForCompetitor
	}

	return page, nil
}

func (ps *pageService) GetPageWithHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, error) {
	page, err := ps.GetPage(ctx, competitorID, pageID)
	if err != nil {
		return api.GetPageResponse{}, err
	}

	pageHistory, err := ps.pageHistoryService.ListPageHistory(ctx, pageID, historyPaginationParams)
	if err != nil {
		return api.GetPageResponse{}, ErrFailedToListPageHistory
	}

	return api.GetPageResponse{
		Page:    page,
		History: pageHistory,
	}, nil
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	updatedPage, pErr := ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, pageReq)
	if pErr != nil && pErr.HasErrors() {
		return models.Page{}, ErrFailedToUpdatePageForCompetitor
	}

	return updatedPage, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, param *api.PaginationParams) ([]models.Page, error) {
	var limit, offset *int
	if param != nil {
		limit, offset = utils.ToPtr(param.GetLimit()), utils.ToPtr(param.GetOffset())
	}

	page, pErr := ps.pageRepo.ListCompetitorPages(ctx, competitorID, limit, offset)
	if pErr != nil && pErr.HasErrors() {
		return nil, ErrFaileToListPagesForCompetitor
	}

	return page, nil
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	// TODO: TBD
	return nil, nil
}

func (ps *pageService) ListPageHistory(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error) {
	pageHistory, err := ps.pageHistoryService.ListPageHistory(
		ctx,
		pageID,
		param,
	)
	if err != nil {
		return nil, ErrFailedToListPageHistory
	}

	return pageHistory, nil
}

func (ps *pageService) RemovePage(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) []error {
	pErr := ps.pageRepo.RemovePagesFromCompetitor(
		ctx,
		competitorID,
		pageIDs,
	)
	if pErr != nil && pErr.HasErrors() {
		return []error{ErrFailedToRemovePageFromCompetitor}
	}
	return nil
}
