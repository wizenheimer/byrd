// ./src/internal/service/history/service.go
package history

import (
	"context"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/history"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// compile time check if the interface is implemented
var _ PageHistoryService = (*pageHistoryService)(nil)

type pageHistoryService struct {
	pageHistoryRepo   history.PageHistoryRepository
	screenshotService screenshot.ScreenshotService
	diffService       diff.DiffService
	logger            *logger.Logger
}

func NewPageHistoryService(pageHistoryRepo history.PageHistoryRepository, screenshotService screenshot.ScreenshotService, diffService diff.DiffService, logger *logger.Logger) PageHistoryService {
	return &pageHistoryService{
		pageHistoryRepo:   pageHistoryRepo,
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger,
	}
}

// CreatePageHistory creates a page history for a page.
// This is trigger during page creation by the page service and by workflow service.
// It returns true if the new page history was created or it returns false if the page history already exists.
// Error is returned if there was an issue creating the page history.
func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, error) {
	return false, nil
}

// ListPageHistory lists the history of a page, paginated by pageHistoryPaginationParam
// This is triggered when a user wants to list all page histories of a page
func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, error) {
	return nil, nil
}

// ClearPageHistory clears the history of a page.
func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) error {
	return nil
}
