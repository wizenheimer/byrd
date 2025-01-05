package history

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
	ErrFailedToCreatePageHistory = errors.New("failed to create page history")
	ErrFailedToListPageHistory   = errors.New("failed to list page history")
	ErrFailedToClearPageHistory  = errors.New("failed to clear page history")
    
)

func NewPageHistoryService(pageHistoryRepo repo.PageHistoryRepository, screenshotService svc.ScreenshotService, diffService svc.DiffService, logger *logger.Logger) svc.PageHistoryService {
	return &pageHistoryService{
		pageHistoryRepo:   pageHistoryRepo,
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger,
	}
}

func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, error) {
	// TODO: TBD
	return true, nil
}

func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, error) {
	limit, offset := pageHistoryPaginationParam.GetLimit(), pageHistoryPaginationParam.GetOffset()

	history, errs := ph.pageHistoryRepo.ListPageHistory(ctx, pageID, &limit, &offset)
	if errs != nil {
		return nil, ErrFailedToListPageHistory
	}

	return history, nil
}

func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) []error {
	errs := ph.pageHistoryRepo.RemovePageHistory(ctx, pageIDs)
	if errs != nil {
		return []error{ErrFailedToClearPageHistory}
	}

	return nil
}
