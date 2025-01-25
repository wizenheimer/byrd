// ./src/internal/service/history/service.go
package history

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/history"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// compile time check if the interface is implemented
var _ PageHistoryService = (*pageHistoryService)(nil)

type pageHistoryService struct {
	pageHistoryRepo history.PageHistoryRepository
	logger          *logger.Logger
}

func NewPageHistoryService(pageHistoryRepo history.PageHistoryRepository, logger *logger.Logger) PageHistoryService {
	return &pageHistoryService{
		pageHistoryRepo: pageHistoryRepo,
		logger:          logger.WithFields(map[string]interface{}{"module": "page_history_service"}),
	}
}

// CreatePageHistory creates a page history for a page.
// This is trigger during page creation by the page service and by workflow service.
// It returns true if the new page history was created or it returns false if the page history already exists.
// Error is returned if there was an issue creating the page history.
func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID, diff *models.DynamicChanges, prevURL, currURL string) error {
	ph.logger.Debug("creating page history", zap.Any("pageID", pageID), zap.Any("diff", diff))
	return ph.pageHistoryRepo.CreateHistoryForPage(ctx, pageID, diff, prevURL, currURL)
}

// ListPageHistory lists the history of a page, paginated by pageHistoryPaginationParam
// This is triggered when a user wants to list all page histories of a page
func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	ph.logger.Debug("listing page history", zap.Any("pageID", pageID), zap.Any("limit", limit), zap.Any("offset", offset))
	return ph.pageHistoryRepo.BatchGetPageHistory(ctx, pageID, limit, offset)
}

// ClearPageHistory clears the history of a page.
func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) error {
	ph.logger.Debug("clearing page history", zap.Any("pageIDs", pageIDs))
	return ph.pageHistoryRepo.BatchRemovePageHistory(ctx, pageIDs)
}
