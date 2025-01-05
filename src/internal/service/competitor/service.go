package competitor

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	ErrFailedToListWorkspaceCompetitors   = errors.New("failed to list workspace competitors")
	ErrFailedToCreateWorkspaceCompetitors = errors.New("failed to create workspace competitors")
	ErrFailedToAddPagesToCompetitor       = errors.New("failed to add pages to competitor")
	ErrFailedToGetCompetitorForWorkspace  = errors.New("failed to get competitor for workspace")
	ErrFailedToListPagesForCompetitor     = errors.New("failed to list pages for competitor")
	ErrFailedToGetPageForCompetitor       = errors.New("failed to get page for competitor")
	ErrFailedToRemoveWorkspaceCompetitors = errors.New("failed to remove workspace competitors")
	ErrFailedToRemovePageFromCompetitor   = errors.New("failed to remove page from competitor")
	ErrFailedToCreatePageForCompetitor    = errors.New("failed to create page for competitor")
	ErrFailedToUpdatePageForCompetitor    = errors.New("failed to update page for competitor")
)

func NewCompetitorService(competitorRepository repo.CompetitorRepository, pageService svc.PageService, logger *logger.Logger) svc.CompetitorService {
	return &competitorService{
		competitorRepository: competitorRepository,
		pageService:          pageService,
		logger:               logger,
		nameFinder:           NewCompanyNameFinder(logger),
	}
}

func (cs *competitorService) CreateCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorReq api.CreateCompetitorRequest) (api.CreateWorkspaceCompetitorResponse, []error) {
	cerrs := make([]error, 0)
	competitors := make([]models.Competitor, 0)

	// Iterate over the list of pages in the request
	for _, page := range competitorReq.Pages {
		// Get the page URL and find the competitor name
		urls := []string{page.URL}
		competitorName := cs.nameFinder.ProcessURLs(urls)

		// Create competitor using the competitor name
		createdCompetitors, errs := cs.competitorRepository.CreateCompetitors(ctx, workspaceID, []string{competitorName})
		if errs != nil || len(createdCompetitors) == 0 {
			cerrs = append(errs, ErrFailedToCreateWorkspaceCompetitors)
			continue
		}

		// Add the page to the competitor
		if _, errs := cs.pageService.CreatePage(ctx, createdCompetitors[0].ID, []api.CreatePageRequest{page}); len(errs) > 0 {
			cerrs = append(errs, ErrFailedToAddPagesToCompetitor)
			continue
		}

		// Append the created competitor to the list of competitors
		competitors = append(competitors, createdCompetitors...)
	}

	// Return the list of competitors and errors that occurred during the creation
	return api.CreateWorkspaceCompetitorResponse{
		Competitors: competitors,
	}, cerrs

}

func (cs *competitorService) GetCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pagePaginationParam api.PaginationParams) (api.GetWorkspaceCompetitorResponse, error) {
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil {
		return api.GetWorkspaceCompetitorResponse{}, ErrFailedToGetCompetitorForWorkspace
	}

	pages, err := cs.pageService.ListCompetitorPages(ctx, competitorID, &pagePaginationParam)
	if err != nil {
		return api.GetWorkspaceCompetitorResponse{}, ErrFailedToListPagesForCompetitor
	}

	return api.GetWorkspaceCompetitorResponse{
		Competitor: competitor,
		Pages:      pages,
	}, nil
}

func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	competitor, err := cs.competitorRepository.GetCompetitor(ctx, competitorID)
	if err != nil {
		return false, ErrFailedToGetCompetitorForWorkspace
	}

	return competitor.WorkspaceID == workspaceID, nil
}

func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	if _, err := cs.pageService.GetPage(ctx, competitorID, pageID); err != nil {
		return false, ErrFailedToGetPageForCompetitor
	}

	return true, nil
}

func (cs *competitorService) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, param api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, []error) {
	competitors, errs := cs.competitorRepository.ListWorkspaceCompetitors(ctx, workspaceID, param.GetLimit(), param.GetOffset())
	if errs != nil {
		return nil, []error{ErrFailedToListWorkspaceCompetitors}
	}

	var competitorResponses []api.GetWorkspaceCompetitorResponse
	for _, competitor := range competitors {
		// Setting pagination as nil would list all pages for the competitor
		pages, err := cs.pageService.ListCompetitorPages(ctx, competitor.ID, nil)
		if err != nil {
			errs = append(errs, ErrFailedToListPagesForCompetitor)
			continue
		}

		competitorResponses = append(competitorResponses, api.GetWorkspaceCompetitorResponse{
			Competitor: competitor,
			Pages:      pages,
		})
	}

	return competitorResponses, errs
}

func (cs *competitorService) RemoveCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) []error {
	// Remove workspace competitors
	errs := cs.competitorRepository.RemoveWorkspaceCompetitors(ctx, workspaceID, competitorIDs)
	if errs != nil {
		return []error{ErrFailedToRemoveWorkspaceCompetitors}
	}

	// Remove pages for the competitors
	var pageErrs []error
	for _, competitorID := range competitorIDs {
		// Remove pages for the competitor
		errs := cs.pageService.RemovePage(ctx, competitorID, nil)
		if len(errs) > 0 {
			pageErrs = append(pageErrs, ErrFailedToRemovePageFromCompetitor)
		}
	}

	// TODO: make clean ups atomic and transactional
	return pageErrs
}

func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pageReq []api.CreatePageRequest) ([]models.Page, []error) {
	page, errs := cs.pageService.CreatePage(ctx, competitorID, pageReq)
	if len(errs) > 0 {
		return nil, []error{ErrFailedToCreatePageForCompetitor}
	}
	return page, nil
}

func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, historyPaginationParams api.PaginationParams) (api.GetPageResponse, error) {
	page, err := cs.pageService.GetPageWithHistory(ctx, competitorID, pageID, historyPaginationParams)
	if err != nil {
		return api.GetPageResponse{}, ErrFailedToGetPageForCompetitor
	}
	return page, nil
}

func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, pageReq api.UpdatePageRequest) (models.Page, error) {
	page, err := cs.pageService.UpdatePage(ctx, competitorID, pageID, pageReq)
	if err != nil {
		return models.Page{}, ErrFailedToUpdatePageForCompetitor
	}

	return page, nil
}

func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) []error {
	// Remove pages for the competitor
	errs := cs.pageService.RemovePage(ctx, competitorID, pageID)
	if len(errs) > 0 {
		return []error{ErrFailedToRemovePageFromCompetitor}
	}
	return nil
}

func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID) ([]models.Page, error) {
	pages, err := cs.pageService.ListCompetitorPages(ctx, competitorID, nil)
	if err != nil {
		return nil, ErrFailedToListPagesForCompetitor
	}
	return pages, nil
}
