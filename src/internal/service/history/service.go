package history

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

func NewPageHistoryService(pageHistoryRepo repo.PageHistoryRepository, screenshotService svc.ScreenshotService, diffService svc.DiffService, logger *logger.Logger) svc.PageHistoryService {
	return &pageHistoryService{
		pageHistoryRepo:   pageHistoryRepo,
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger,
	}
}

func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, err.Error) {
	// TODO: TBD
	return true, nil
}

func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, err.Error) {
    cErr := err.New()
	limit, offset := pageHistoryPaginationParam.GetLimit(), pageHistoryPaginationParam.GetOffset()

	history, hErr := ph.pageHistoryRepo.ListPageHistory(ctx, pageID, &limit, &offset)
	if hErr != nil && hErr.HasErrors() {
        cErr.Merge(hErr)
		return nil, cErr
	}

	return history, nil
}

func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) err.Error {
    cErr := err.New()
	hErr := ph.pageHistoryRepo.RemovePageHistory(ctx, pageIDs)
	if hErr != nil && hErr.HasErrors() {
        cErr.Merge(hErr)
        return cErr
	}

	return nil
}
