package competitor

import (
	"context"

	"github.com/google/uuid"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

func (cs *competitorService) CreateCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorReq api.CreateCompetitorRequest) (api.CreateWorkspaceCompetitorResponse, []error) {
	var competitorNames []string

	for _, page := range competitorReq.Pages {
		// TODO: add competitor name
		competitorNames = append(competitorNames, page.URL)
	}

	createdCompetitors, err := cs.competitorRepository.CreateCompetitors(ctx, workspaceID, competitorNames)
	if err != nil {
		return api.CreateWorkspaceCompetitorResponse{}, []error{err}
	}

	return api.CreateWorkspaceCompetitorResponse{
		Competitors: createdCompetitors,
	}, nil

}

func (cs *competitorService) GetCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pagePaginationParam api.PaginationParams) (api.GetWorkspaceCompetitorResponse, error) {
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil {
		return api.GetWorkspaceCompetitorResponse{}, err
	}

	pages, err := cs.pageService.ListCompetitorPages(ctx, competitorID, &pagePaginationParam)
	if err != nil {
		return api.GetWorkspaceCompetitorResponse{}, err
	}

	return api.GetWorkspaceCompetitorResponse{
		Competitor: competitor,
		Pages:      pages,
	}, nil
}

func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil {
		return false, err
	}

	return competitor.WorkspaceID == workspaceID, nil
}

func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	if _, err := cs.pageService.GetPage(ctx, competitorID, pageID); err != nil {
		return false, err
	}

	return true, nil
}

func (cs *competitorService) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, param api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, error) {
	competitors, err := cs.competitorRepository.ListWorkspaceCompetitors(ctx, workspaceID.String(), param.GetLimit(), param.GetOffset())
	if err != nil {
		return nil, err
	}

	var competitorResponses []api.GetWorkspaceCompetitorResponse
	for _, competitor := range competitors {
		// Setting pagination as nil would list all pages for the competitor
		pages, err := cs.pageService.ListCompetitorPages(ctx, competitor.ID, nil)
		if err != nil {
			return nil, err
		}

		competitorResponses = append(competitorResponses, api.GetWorkspaceCompetitorResponse{
			Competitor: competitor,
			Pages:      pages,
		})
	}

	return competitorResponses, nil
}

func (cs *competitorService) RemoveCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) []error {
	// Remove workspace competitors
	errs := cs.competitorRepository.RemoveWorkspaceCompetitors(ctx, workspaceID, competitorIDs)
	if errs != nil {
		return errs
	}

	// Remove pages for the competitors
	var pageErrs []error
	for _, competitorID := range competitorIDs {
		// Remove pages for the competitor
		errs := cs.pageService.RemovePage(ctx, competitorID, nil)
		if len(errs) > 0 {
			pageErrs = append(pageErrs, errs...)
		}
	}

	// TODO: make clean ups atomic and transactional
	return pageErrs
}

func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, []error) {
	page, err := cs.pageService.CreatePage(ctx, competitorID, pageReq)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, error) {
	return cs.pageService.GetPageWithHistory(ctx, competitorID, pageID, historyPaginationParams)
}

func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	return cs.pageService.UpdatePage(ctx, competitorID, pageID, pageReq)
}

func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) []error {
	// Remove pages for the competitor
	errs := cs.pageService.RemovePage(ctx, competitorID, pageID)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID) ([]models.Page, error) {
	return cs.pageService.ListCompetitorPages(ctx, competitorID, nil)
}
