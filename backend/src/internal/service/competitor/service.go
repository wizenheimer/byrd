// ./src/internal/service/competitor/service.go
package competitor

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type competitorService struct {
	competitorRepository competitor.CompetitorRepository
	pageService          page.PageService
	nameFinder           *CompanyNameFinder
	logger               *logger.Logger
}

var _ CompetitorService = (*competitorService)(nil)

func NewCompetitorService(competitorRepository competitor.CompetitorRepository, pageService page.PageService, logger *logger.Logger) CompetitorService {
	return &competitorService{
		competitorRepository: competitorRepository,
		pageService:          pageService,
		logger:               logger,
		nameFinder:           NewCompanyNameFinder(logger),
	}
}

func (cs *competitorService) AddCompetitorsToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error) {
	if len(pages) == 0 {
		return nil, errors.New("pages unspecified for creating competitors")
	}

	var competitors []models.Competitor
	for _, page := range pages {
		// Figure out competitor's name using the url
		competitorName := cs.nameFinder.FindCompanyName([]string{
			page.URL,
		})

		// Create a new competitor using the competitor's name
		competitor, err := cs.competitorRepository.CreateCompetitorForWorkspace(
			ctx,
			workspaceID,
			competitorName,
		)
		if err != nil {
			cs.logger.Debug("failed to create competitor")
			continue
		}

		// Create a page, and associate it with the created competitor
		_, err = cs.pageService.CreatePage(
			ctx,
			competitor.ID,
			[]models.PageProps{
				page,
			},
		)
		if err != nil {
			cs.logger.Debug("failed to create page for competitor")
		}

		competitors = append(competitors, *competitor)
	}

	totalPages := len(pages)
	totalCompetitors := len(competitors)

	if totalCompetitors == 0 {
		return competitors, errors.New("failed to create competitors")
	}

	if totalPages > totalCompetitors {
		return competitors, errors.New("failed to create some competitors")
	}

	return competitors, nil
}

func (cs *competitorService) GetCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error) {
	return cs.competitorRepository.BatchGetCompetitorsForWorkspace(
		ctx,
		workspaceID,
		competitorIDs,
	)
}

func (cs *competitorService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, error) {
	return cs.competitorRepository.ListCompetitorsForWorkspace(
		ctx,
		workspaceID,
		limit,
		offset,
	)
}

func (cs *competitorService) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return cs.competitorRepository.UpdateCompetitorForWorkspace(
		ctx,
		workspaceID,
		competitorID,
		competitorName,
	)
}


func (cs *competitorService) RemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error {
	if competitorIDs == nil {
		return cs.competitorRepository.RemoveAllCompetitorsForWorkspace(
			ctx,
			workspaceID,
		)
	}

	if len(competitorIDs) == 1 {
		return cs.competitorRepository.RemoveCompetitorForWorkspace(
			ctx,
			workspaceID,
			competitorIDs[0],
		)
	}

	return cs.competitorRepository.BatchRemoveCompetitorForWorkspace(
		ctx,
		workspaceID,
		competitorIDs,
	)
}


func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	return cs.competitorRepository.WorkspaceCompetitorExists(
		ctx,
		workspaceID,
		competitorID,
	)
}

func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	return cs.pageService.PageExists(
		ctx,
		competitorID,
		pageID,
	)
}

func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	return cs.pageService.CreatePage(
		ctx,
		competitorID,
		pages,
	)
}

func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	return cs.pageService.GetPage(
		ctx,
		competitorID,
		pageID,
	)
}

func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return cs.pageService.UpdatePage(
		ctx,
		competitorID,
		pageID,
		page,
	)
}

func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) error {
	return cs.pageService.RemovePage(
		ctx,
		competitorID,
		pageID,
	)
}

func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return cs.pageService.ListCompetitorPages(
		ctx,
		competitorID,
		limit,
        offset,
	)
}

func (cs *competitorService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
    return cs.pageService.ListPageHistory(
        ctx,
        pageID,
        limit,
        offset,
    )
}