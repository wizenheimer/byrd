// ./src/internal/service/page/service.go
package page

import (
	"context"
	"errors"
	"fmt"
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
	if len(pages) > maxPageBatchSize {
		return nil, errors.New("batch size exceeds the maximum limit")
	}

	if len(pages) == 0 {
		return nil, errors.New("pages unspecified for creating competitors")
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
		return createdPages, errors.New("failed to create all pages")
	}

	ps.backdateRefresh(createdPages)

	return createdPages, nil
}

func (ps *pageService) backdateRefresh(pages []models.Page) {
	for _, page := range pages {
		go func(page models.Page) {
			ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
			defer cancel()

			screenshotRequestOptions := models.GetScreenshotRequestOptions(page.URL, page.CaptureProfile)
			ir, _, err := ps.screenshotService.Refresh(ctx, screenshotRequestOptions, true)

			if err != nil {
				ps.logger.Error("failed to refresh page", zap.Any("pageID", page.ID), zap.Error(err))
			}

			if ir == nil {
				ps.logger.Error("failed to get image response for page", zap.Any("pageID", page.ID))
				ir = &models.ScreenshotImage{
					StoragePath: "",
				}
			}

			diff, err := models.NewEmptyDynamicChanges(page.DiffProfile)
			if err != nil {
				ps.logger.Error("failed to create empty dynamic changes", zap.Any("pageID", page.ID), zap.Error(err))
			}

			if diff == nil {
				ps.logger.Error("failed to get diff for page, defaulting to empty", zap.Any("pageID", page.ID))
				diff = &models.DynamicChanges{}
			}

			if err := ps.pageHistoryService.CreatePageHistory(
				context.Background(),
				page.ID,
				diff,
				ir.StoragePath,
				ir.StoragePath,
			); err != nil {
				ps.logger.Error("failed to create page history", zap.Any("pageID", page.ID), zap.Error(err))
			}
		}(page)
	}
}

func (ps *pageService) GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error) {
	return ps.pageRepo.GetCompetitorPageByID(ctx, competitorID, pageID)
}

func (ps *pageService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	return ps.pageHistoryService.ListPageHistory(ctx, pageID, limit, offset)
}

func (ps *pageService) UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	captureProfileRequiresUpdate := page.CaptureProfile != nil
	diffProfileRequiresUpdate := len(page.DiffProfile) > 0
	urlRequiresUpdate := page.URL != ""

	var updatedPage *models.Page
	var err error

	// If all three fields require an update, update the page
	if captureProfileRequiresUpdate && diffProfileRequiresUpdate && urlRequiresUpdate {
		updatedPage, err = ps.pageRepo.UpdateCompetitorPage(ctx, competitorID, pageID, page)
		if err != nil {
			return nil, err
		}
		return updatedPage, nil
	}

	// If only one field requires an update, update the page with that field one at a time
	// This is done to avoid updating the page with a nil value
	if captureProfileRequiresUpdate {
		updatedPage, err = ps.pageRepo.UpdateCompetitorCaptureProfile(ctx, competitorID, pageID, page.CaptureProfile, page.URL)
		if err != nil {
			return nil, err
		}
	}

	// If only one field requires an update, update the page with that field one at a time
	// This is done to avoid updating the page with an empty value
	if diffProfileRequiresUpdate {
		updatedPage, err = ps.pageRepo.UpdateCompetitorDiffProfile(ctx, competitorID, pageID, page.DiffProfile)
		if err != nil {
			return nil, err
		}
	}

	// If only one field requires an update, update the page with that field one at a time
	// This is done to avoid updating the page with an empty value
	if urlRequiresUpdate {
		updatedPage, err = ps.pageRepo.UpdateCompetitorPageURL(ctx, competitorID, pageID, page.URL)
		if err != nil {
			return nil, err
		}
	}

	// Return the updated page
	return updatedPage, nil
}

func (ps *pageService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	return ps.pageRepo.GetCompetitorPages(ctx, competitorID, limit, offset)
}

func (ps *pageService) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []uuid.UUID, <-chan error) {
	pagesChan := make(chan []uuid.UUID)
	errorsChan := make(chan error)

	go func() {
		defer close(pagesChan)
		defer close(errorsChan)

		currentPageID := lastPageID
		for {
			select {
			case <-ctx.Done():
				return
			default:
				activePages, err := ps.pageRepo.GetActivePages(ctx, batchSize, currentPageID)
				if err != nil {
					errorsChan <- err
					return
				}

				if len(activePages.PageIDs) == 0 {
					// No more pages to process
					return
				}

				// Send the current batch
				select {
				case pagesChan <- activePages.PageIDs:
				case <-ctx.Done():
					return
				}

				if !activePages.HasMore {
					// No more pages after this batch
					return
				}

				// Update cursor for next iteration
				currentPageID = activePages.LastSeen
			}
		}
	}()

	return pagesChan, errorsChan
}

// RefreshPage with the given pageID using best effort strategy
func (ps *pageService) RefreshPage(ctx context.Context, pageID uuid.UUID) error {
	urlContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	page, err := ps.pageRepo.GetPageByPageID(ctx, pageID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	screenshotOptions := models.GetScreenshotRequestOptions(page.URL, page.CaptureProfile)

	// Handle current screenshot refresh with fallback
	var currentPath string
	var currentHTMLContent *models.ScreenshotContent
	currentImgResp, currentHTMLContentResp, err := ps.screenshotService.Refresh(urlContext, screenshotOptions, false)
	if err != nil {
		ps.logger.Error("failed to refresh screenshot", zap.Error(err), zap.Any("pageID", pageID))
		currentPath = ""
		currentHTMLContent = &models.ScreenshotContent{}
	} else {
		currentPath = currentImgResp.StoragePath
		currentHTMLContent = currentHTMLContentResp
	}

	// Handle previous screenshot retrieval with fallback
	var previousPath string
	var previousHTMLContent *models.ScreenshotContent
	prevImgResp, previousHtmlContentResp, err := ps.screenshotService.Retrieve(ctx, screenshotOptions, false)
	if err != nil {
		ps.logger.Error("failed to retrieve previous screenshot", zap.Error(err), zap.Any("pageID", pageID))
		previousPath = ""
		previousHTMLContent = &models.ScreenshotContent{}
	} else {
		previousPath = prevImgResp.StoragePath
		previousHTMLContent = previousHtmlContentResp
	}

	var diff *models.DynamicChanges
	// Only perform diff if both paths are non-empty
	if currentPath != "" && previousPath != "" {
		diffResult, err := ps.diffService.Compare(ctx, previousHTMLContent, currentHTMLContent, page.DiffProfile)
		if err != nil {
			ps.logger.Error("failed to compare contents", zap.Error(err), zap.Any("pageID", pageID))
			diff, err = models.NewEmptyDynamicChanges(page.DiffProfile)
			if err != nil {
				diff = &models.DynamicChanges{}
			}
		} else {
			diff = diffResult
		}
	} else {
		ps.logger.Warn("skipping diff due to missing screenshots",
			zap.Any("pageID", pageID),
			zap.String("currentPath", currentPath),
			zap.String("previousPath", previousPath))
		diff, err = models.NewEmptyDynamicChanges(page.DiffProfile)
		if err != nil {
			diff = &models.DynamicChanges{}
		}
	}

	// Create history with best effort approach
	return ps.pageHistoryService.CreatePageHistory(ctx, pageID, diff, previousPath, currentPath)
}

func (ps *pageService) RemovePage(ctx context.Context, competitorIDs []uuid.UUID, pageIDs []uuid.UUID) error {
	if len(pageIDs) > maxPageBatchSize {
		return errors.New("page batch size exceeds the maximum limit")
	}
	if competitorIDs == nil {
		return errors.New("competitorIDs unspecified for removing pages")
	}

	if len(competitorIDs) > 1 {
		// Perform batch delete if multiple competitorIDs are provided

		if pageIDs == nil {
			// Remove all pages for all competitors
			return ps.pageRepo.BatchDeleteAllCompetitorPages(ctx, competitorIDs)
		}

		return errors.New("pageIDs ambiguous for removing pages")
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

func (ps *pageService) GetLatestPageHistory(ctx context.Context, pageID []uuid.UUID) ([]models.PageHistory, error) {
	return ps.pageHistoryService.GetLatestPageHistory(ctx, pageID)
}

func (ps *pageService) CountActivePagesForCompetitors(ctx context.Context, competitorIDs []uuid.UUID) (int, error) {
	pageCountMap, err := ps.pageRepo.GetActivePageCountsByCompetitors(ctx, competitorIDs)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, c := range pageCountMap {
		count += c
	}

	return count, nil
}
