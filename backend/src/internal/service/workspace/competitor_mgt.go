package workspace

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"go.uber.org/zap"
)

func (ws *workspaceService) GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error) {
	competitors, err := ws.competitorService.GetCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil {
		return nil, err
	}
	if len(competitors) == 0 {
		return nil, errors.New("competitor not found")
	}

	return &competitors[0], nil
}

func (ws *workspaceService) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return ws.competitorService.UpdateCompetitorForWorkspace(ctx, workspaceID, competitorID, competitorName)
}

func (ws *workspaceService) AddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (*models.Competitor, error) {
	canCreateCompetitor, _, err := ws.CanCreateCompetitor(ctx, workspaceID, 1, len(pages))
	if err != nil {
		return nil, err
	}
	if !canCreateCompetitor {
		return nil, errors.New("user cannot create competitor")
	}

	competitor, err := ws.competitorService.CreateCompetitorForWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	return &competitor, nil
}

func (ws *workspaceService) BatchAddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error) {
	canCreateCompetitor, _, err := ws.CanCreateCompetitor(ctx, workspaceID, len(pages), len(pages))
	if err != nil {
		return nil, err
	}
	if !canCreateCompetitor {
		return nil, errors.New("user cannot create competitor")
	}

	competitors, err := ws.competitorService.BatchCreateCompetitorForWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	if len(competitors) == 0 {
		ws.logger.Error("failed to create competitors", zap.Any("workspaceID", workspaceID), zap.Any("pageProps", pages), zap.Any("numCompetitors", len(competitors)))
		return nil, errors.New("failed to create competitors")
	}

	return competitors, nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pageProps []models.PageProps) ([]models.Page, error) {
	canCreatePage, _, err := ws.CanCreatePage(ctx, workspaceID, len(pageProps))
	if err != nil {
		return nil, err
	}
	if !canCreatePage {
		return nil, errors.New("user cannot create page")
	}

	createdPages, err := ws.competitorService.AddPagesToCompetitor(ctx, competitorID, pageProps)
	if err != nil {
		return nil, err
	}

	if len(pageProps) != len(createdPages) {
		ws.logger.Error("failed to create some pages", zap.Int("numPages", len(pageProps)), zap.Int("numCreatedPages", len(createdPages)), zap.Any("pages", pageProps), zap.Any("createdPages", createdPages))
	}

	return createdPages, nil
}

func (ws *workspaceService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error) {
	competitors, hasMore, err := ws.competitorService.ListCompetitorsForWorkspace(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, false, err
	}

	return competitors, hasMore, nil
}

// ListPagesForCompetitor lists the pages for a competitor
func (ws *workspaceService) ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	pages, hasMore, err := ws.competitorService.ListCompetitorPages(ctx, competitorID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pages, hasMore, nil
}

// ListHistoryForPage lists the history of a page
func (ws *workspaceService) ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	pageHistory, hasMore, err := ws.competitorService.ListPageHistory(ctx, pageID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pageHistory, hasMore, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error {
	return ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, []uuid.UUID{pageID})
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	return ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, pageProps models.PageProps) (*models.Page, error) {
	return ws.competitorService.UpdatePage(ctx, competitorID, pageID, pageProps)
}

func (ws *workspaceService) GetPageForCompetitor(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	return ws.competitorService.GetCompetitorPage(ctx, competitorID, pageID)
}
