package competitor

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

func NewCompetitorService(competitorRepository repo.CompetitorRepository, pageService svc.PageService, logger *logger.Logger) svc.CompetitorService {
	return &competitorService{
		competitorRepository: competitorRepository,
		pageService:          pageService,
		logger:               logger,
		nameFinder:           NewCompanyNameFinder(logger),
	}
}

func (cs *competitorService) CreateCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorReq api.CreateCompetitorRequest) (api.CreateWorkspaceCompetitorResponse, err.Error) {
	cErr := err.New()
	competitors := make([]models.Competitor, 0)

	// Iterate over the list of pages in the request
	for _, page := range competitorReq.Pages {
		// Get the page URL and find the competitor name
		urls := []string{page.URL}
		competitorName := cs.nameFinder.ProcessURLs(urls)

		// Create competitor using the competitor name
		createdCompetitors, err := cs.competitorRepository.CreateCompetitors(ctx, workspaceID, []string{competitorName})
		if err != nil && err.HasErrors() {
			cErr.Merge(err)
			continue
		}

		// Add the page to the competitor
		if _, err := cs.pageService.CreatePage(ctx, createdCompetitors[0].ID, []api.CreatePageRequest{page}); err != nil && err.HasErrors() {
			cErr.Merge(err)
			continue
		}

		// Append the created competitor to the list of competitors
		competitors = append(competitors, createdCompetitors...)
	}

	// Return the list of competitors and errors that occurred during the creation
	return api.CreateWorkspaceCompetitorResponse{
		Competitors: competitors,
	}, cErr

}

func (cs *competitorService) GetCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pagePaginationParam api.PaginationParams) (api.GetWorkspaceCompetitorResponse, err.Error) {
	cErr := err.New()
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil && err.HasErrors() {
		cErr.Merge(err)
		return api.GetWorkspaceCompetitorResponse{}, cErr
	}

	pages, err := cs.pageService.ListCompetitorPages(ctx, competitorID, &pagePaginationParam)
	if err != nil {
		cErr.Merge(err)
		return api.GetWorkspaceCompetitorResponse{}, cErr
	}

	return api.GetWorkspaceCompetitorResponse{
		Competitor: competitor,
		Pages:      pages,
	}, cErr
}

func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, err.Error) {
	cErr := err.New()
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil && err.HasErrors() {
		cErr.Merge(err)
		return false, cErr
	}

	return competitor.WorkspaceID == workspaceID, nil
}

func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, err.Error) {
	cErr := err.New()
	if _, err := cs.pageService.GetPage(ctx, competitorID, pageID); err != nil {
		cErr.Merge(err)
		return false, cErr
	}

	return true, nil
}

func (cs *competitorService) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, param api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, err.Error) {
	cErr := err.New()
	competitors, err := cs.competitorRepository.ListWorkspaceCompetitors(ctx, workspaceID, param.GetLimit(), param.GetOffset())
	if err != nil && err.HasErrors() {
		cErr.Merge(err)
		return nil, cErr
	}

	var competitorResponses []api.GetWorkspaceCompetitorResponse
	for _, competitor := range competitors {
		// Setting pagination as nil would list all pages for the competitor
		pages, err := cs.pageService.ListCompetitorPages(ctx, competitor.ID, nil)
		if err != nil && err.HasErrors() {
			cErr.Merge(err)
			continue
		}

		competitorResponses = append(competitorResponses, api.GetWorkspaceCompetitorResponse{
			Competitor: competitor,
			Pages:      pages,
		})
	}

	return competitorResponses, cErr
}

func (cs *competitorService) RemoveCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) err.Error {
	cErr := err.New()
	// Remove workspace competitors

	if err := cs.competitorRepository.RemoveWorkspaceCompetitors(ctx, workspaceID, competitorIDs); err != nil && err.HasErrors() {
		cErr.Merge(err)
		return cErr
	}

	// Remove pages for the competitors
	for _, competitorID := range competitorIDs {
		// Remove pages for the competitor
		err := cs.pageService.RemovePage(ctx, competitorID, nil)
		if err != nil && err.HasErrors() {
			cErr.Merge(err)
		}
	}

	// TODO: make clean ups atomic and transactional
	return cErr
}

func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, err.Error) {
	cErr := err.New()
	page, err := cs.pageService.CreatePage(ctx, competitorID, pageReq)
	if err != nil && err.HasErrors() {
		cErr.Merge(err)
		return nil, cErr
	}
	return page, nil
}

func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, err.Error) {
	cErr := err.New()
	page, err := cs.pageService.GetPageWithHistory(ctx, competitorID, pageID, historyPaginationParams)
	if err != nil && err.HasErrors() {
		return api.GetPageResponse{}, cErr
	}
	return page, nil
}

func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, err.Error) {
	cErr := err.New()
	page, err := cs.pageService.UpdatePage(ctx, competitorID, pageID, pageReq)
	if err != nil {
		cErr.Merge(err)
		return models.Page{}, cErr
	}

	return page, nil
}

func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) err.Error {
	cErr := err.New()
	// Remove pages for the competitor
	errs := cs.pageService.RemovePage(ctx, competitorID, pageID)
	if len(errs) > 0 {
		cErr.Merge(errs)
		return cErr
	}
	return nil
}

func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID) ([]models.Page, err.Error) {
	cErr := err.New()
	pages, err := cs.pageService.ListCompetitorPages(ctx, competitorID, nil)
	if err != nil {
		cErr.Merge(err)
		return nil, cErr
	}
	return pages, nil
}
