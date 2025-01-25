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
	"go.uber.org/zap"
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
		logger:             logger.WithFields(map[string]interface{}{"module": "page_service"}),
	}
}

func (ps *pageService) CreatePage(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	ps.logger.Debug("creating pages", zap.Any("competitorID", competitorID), zap.Any("pages", pages))
	if len(pages) > maxPageBatchSize {
		return nil, errors.New("non-fatal: batch size exceeds the maximum limit")
	}

	if len(pages) == 0 {
		return nil, errors.New("non-fatal: pages unspecified for creating competitors")
	}
	var createdPages []models.Page
	var err error

	if len(pages) == 1 {
		// If there is only one page, perform a single add
		createdPage, err := ps.pageRepo.AddPageToCompetitor(
			ctx,
			competitorID,
			pages[0],
		)
		if err != nil {
			return nil, err
		}
		if createdPage == nil {
			return nil, errors.New("failed to create page")
		}
		createdPages = append(createdPages, *createdPage)
	} else {
		// If there are multiple pages, perform a batch add
		createdPages, err = ps.pageRepo.BatchAddPageToCompetitor(
			ctx,
			competitorID,
			pages,
		)
		if err != nil {
			return nil, err
		}
	}

	if len(createdPages) != len(pages) {
		return createdPages, errors.New("non-fatal: failed to create all pages")
	}

	ps.backdateRefresh(createdPages)

	return createdPages, nil
}

func (ps *pageService) backdateRefresh(pages []models.Page) {
	for _, page := range pages {
		go func(page models.Page) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if ir, hr, err := ps.screenshotService.Refresh(ctx, page.CaptureProfile, true); err != nil {
				ps.logger.Error("failed to refresh page", zap.Any("pageID", page.ID), zap.Error(err))
			} else {
				ps.logger.Debug("refreshed page", zap.Any("pageID", page.ID), zap.Any("imagePath", ir.StoragePath), zap.Any("contentPath", hr.StoragePath))
			}
		}(page)
	}
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error) {
	ps.logger.Debug("getting page", zap.Any("competitorID", competitorID), zap.Any("pageID", pageID))
	return ps.pageRepo.GetCompetitorPageByID(ctx, competitorID, pageID)
}

func (ps *pageService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	ps.logger.Debug("listing page history", zap.Any("pageID", pageID), zap.Any("limit", limit), zap.Any("offset", offset))
	return ps.pageHistoryService.ListPageHistory(ctx, pageID, limit, offset)
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	ps.logger.Debug("updating page", zap.Any("competitorID", competitorID), zap.Any("pageID", pageID), zap.Any("page", page))
	captureProfileRequiresUpdate := page.CaptureProfile != nil || page.URL != ""
	diffProfileRequiresUpdate := len(page.DiffProfile) > 0
	urlRequiresUpdate := page.URL != ""

	var updatedPage *models.Page
	var err error
	if captureProfileRequiresUpdate && diffProfileRequiresUpdate && urlRequiresUpdate {
		updatedPage, err = ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, page)
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

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	ps.logger.Debug("listing competitor pages", zap.Any("competitorID", competitorID), zap.Any("limit", limit), zap.Any("offset", offset))
	return ps.pageRepo.GetCompetitorPages(ctx, competitorID, limit, offset)
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []uuid.UUID, <-chan error) {
	ps.logger.Debug("listing active pages", zap.Any("batchSize", batchSize), zap.Any("lastPageID", lastPageID))
	pagesChan := make(chan []uuid.UUID)
	errorsChan := make(chan error)

	go func() {
		defer close(pagesChan)
		defer close(errorsChan)

		hasMore := true
		for hasMore {
			activePages, err := ps.pageRepo.GetActivePages(ctx, batchSize, lastPageID)
			if err != nil {
				errorsChan <- err
				return
			}

			hasMore = activePages.HasMore
			pagesChan <- activePages.PageIDs
		}
	}()

	return pagesChan, errorsChan
}

func (ps *pageService) RefreshPage(ctx context.Context, pageID uuid.UUID) error {
	ps.logger.Debug("refreshing page", zap.Any("pageID", pageID))
	urlContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	page, err := ps.pageRepo.GetPageByPageID(ctx, pageID)
	if err != nil {
		return err
	}

	currentImgResp, currentHTMLContentResp, err := ps.screenshotService.Refresh(urlContext, page.CaptureProfile, false)
	if err != nil {
		return err
	}

	prevImgResp, previousHtmlContentResp, err := ps.screenshotService.Retrieve(ctx, page.CaptureProfile, false)
	if err != nil {
		return err
	}

	diff, err := ps.diffService.Compare(ctx, previousHtmlContentResp, currentHTMLContentResp, page.DiffProfile)
	if err != nil {
		return err
	}

	return ps.pageHistoryService.CreatePageHistory(ctx, pageID, diff, prevImgResp.StoragePath, currentImgResp.StoragePath)
}

func (ps *pageService) RemovePage(ctx context.Context, competitorIDs []uuid.UUID, pageIDs []uuid.UUID) error {
	ps.logger.Debug("removing page", zap.Any("competitorIDs", competitorIDs), zap.Any("pageIDs", pageIDs))
	if len(pageIDs) > maxPageBatchSize {
		return errors.New("non-fatal: page batch size exceeds the maximum limit")
	}
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
	ps.logger.Debug("checking if page exists", zap.Any("competitorID", competitorID), zap.Any("pageID", pageID))
	page, err := ps.pageRepo.GetCompetitorPageByID(ctx, competitorID, pageID)
	if err != nil {
		return false, err
	}
	return page != nil, nil
}
