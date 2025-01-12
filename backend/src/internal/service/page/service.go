// ./src/internal/service/page/service.go
package page

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/page"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/service/history"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// compile time check if the interface is implemented
var _ PageService = (*pageService)(nil)

type pageService struct {
	pageRepo           page.PageRepository
	pageHistoryService history.PageHistoryService
	diffService        diff.DiffService
	screenshotService  screenshot.ScreenshotService
	logger             *logger.Logger
}

func NewPageService(pageRepo page.PageRepository, pageHistoryService history.PageHistoryService, diffService diff.DiffService, screenshotService screenshot.ScreenshotService, logger *logger.Logger) PageService {
	return &pageService{
		pageRepo:           pageRepo,
		pageHistoryService: pageHistoryService,
		diffService:        diffService,
		screenshotService:  screenshotService,
		logger:             logger,
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	if len(pages) == 0 {
		return nil, errors.New("non-fatal: pages unspecified for creating competitors")
	}
	var createdPages []models.Page
	var err error
	if len(pages) == 1 {
		createdPages = make([]models.Page, 1)
		createdPage, err := ps.pageRepo.AddPageToCompetitor(
			ctx,
			competitorID,
			pages[0],
		)
		if err != nil {
			return nil, err
		}
		createdPages[0] = *createdPage

	}

	createdPages, err = ps.pageRepo.BatchAddPageToCompetitor(
		ctx,
		competitorID,
		pages,
	)
	if err != nil {
		return nil, err
	}

	if len(createdPages) != len(pages) {
		return createdPages, errors.New("non-fatal: failed to create all pages")
	}

	return createdPages, nil
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error) {
	return ps.pageRepo.GetCompetitorPageByID(ctx, competitorID, pageID)
}

func (ps *pageService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
	return ps.pageHistoryService.ListPageHistory(ctx, pageID, limit, offset)
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	captureProfileRequiresUpdate := page.CaptureProfile != nil || page.URL != ""
	diffProfileRequiresUpdate := len(page.DiffProfile) > 0
	urlRequiresUpdate := page.URL != ""

	var updatedPage *models.Page
	var err error
	if captureProfileRequiresUpdate && diffProfileRequiresUpdate && urlRequiresUpdate {
		ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, page)
	} else {
		if captureProfileRequiresUpdate {
			updatedPage, err = ps.pageRepo.UpdateCompetitorCaptureProfile(ctx, competitorID, pageID, page.CaptureProfile, page.URL)
		}
		if diffProfileRequiresUpdate {
			updatedPage, err = ps.pageRepo.UpdateCompetitorDiffProfile(ctx, competitorID, pageID, page.DiffProfile)
		}
		if urlRequiresUpdate {
			updatedPage, err = ps.pageRepo.UpdateCompetitorPageURL(ctx, competitorID, pageID, page.URL)
		}
	}
	if err != nil {
		return nil, err
	}

	return updatedPage, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	return ps.pageRepo.GetCompetitorPages(ctx, competitorID, limit, offset)
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error) {
	return nil, nil
}

func (ps *pageService) RefreshPage(ctx context.Context, pageID uuid.UUID) error {
	urlContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	page, err := ps.pageRepo.GetPageByPageID(ctx, pageID)
	if err != nil {
		return err
	}

	_, currentHTMLContentResp, err := ps.screenshotService.Refresh(urlContext, page.URL, page.CaptureProfile)
	if err != nil {
		return err
	}

	_, previousHtmlContentResp, err := ps.screenshotService.Retrieve(ctx, page.URL)
	if err != nil {
		return err
	}

	diff, err := ps.diffService.Compare(ctx, previousHtmlContentResp, currentHTMLContentResp, page.DiffProfile)
	if err != nil {
		return err
	}

	return ps.pageHistoryService.CreatePageHistory(ctx, pageID, diff)
}

func (ps *pageService) RemovePage(ctx context.Context, competitorIDs []uuid.UUID, pageIDs []uuid.UUID) error {
	if competitorIDs == nil {
		return errors.New("non-fatal: competitorIDs unspecified for removing pages")
	}

	if len(competitorIDs) > 1 {
		// Perform batch delete if multiple competitorIDs are provided

		if pageIDs == nil {
			// Remove all pages for all competitors
			return ps.pageRepo.BatchDeleteAllCompetitorPages(ctx, competitorIDs)
		}

		return errors.New("non-fatal: pageIDs ambiguous for removing pages")
	}

	// Perform single delete if only one competitorID is provided
	if pageIDs == nil {
		// Remove all pages for a competitor
		return ps.pageRepo.DeleteAllCompetitorPages(ctx, competitorIDs[0])
	} else if len(pageIDs) == 1 {
		// Remove a single page for a competitor
		return ps.pageRepo.DeleteCompetitorPageByID(ctx, competitorIDs[0], pageIDs[0])
	} else {
		// Remove multiple pages for a competitor
		return ps.pageRepo.BatchDeleteCompetitorPagesByIDs(ctx, competitorIDs[0], pageIDs)
	}
}

func (ps *pageService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	page, err := ps.pageRepo.GetCompetitorPageByID(ctx, competitorID, pageID)
	if err != nil {
		return false, err
	}
	return page != nil, nil
}
