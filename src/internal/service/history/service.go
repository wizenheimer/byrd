package history

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
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

func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, error) {
	// TODO: TBD
	return true, nil
}

func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, error) {
	limit, offset := pageHistoryPaginationParam.GetLimit(), pageHistoryPaginationParam.GetOffset()

	return ph.pageHistoryRepo.ListPageHistory(ctx, pageID, &limit, &offset)
}

func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) []error {
	return ph.pageHistoryRepo.RemovePageHistory(ctx, pageIDs)
}
