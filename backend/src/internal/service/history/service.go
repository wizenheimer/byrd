package history

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

func NewPageHistoryService(pageHistoryRepo repo.PageHistoryRepository, screenshotService svc.ScreenshotService, diffService svc.DiffService, logger *logger.Logger) svc.PageHistoryService {
	return &pageHistoryService{
		pageHistoryRepo:   pageHistoryRepo,
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger,
	}
}

func (ph *pageHistoryService) CreatePageHistory(ctx context.Context, pageID uuid.UUID) (bool, errs.Error) {
	// TODO: TBD
	return true, nil
}

func (ph *pageHistoryService) ListPageHistory(ctx context.Context, pageID uuid.UUID, pageHistoryPaginationParam api.PaginationParams) ([]models.PageHistory, errs.Error) {
	cErr := errs.New()
	limit, offset := pageHistoryPaginationParam.GetLimit(), pageHistoryPaginationParam.GetOffset()

	history, hErr := ph.pageHistoryRepo.ListPageHistory(ctx, pageID, &limit, &offset)
	if hErr != nil && hErr.HasErrors() {
		cErr.Merge(hErr)
		ph.logger.Info(hErr.Error())
		// Initialize an empty slice instead of returning nil
		return make([]models.PageHistory, 0), cErr
	}

	// If history is nil, return an empty slice
	if history == nil {
		return make([]models.PageHistory, 0), nil
	}

	return history, nil
}

func (ph *pageHistoryService) ClearPageHistory(ctx context.Context, pageIDs []uuid.UUID) errs.Error {
	cErr := errs.New()
	hErr := ph.pageHistoryRepo.RemovePageHistory(ctx, pageIDs)
	if hErr != nil && hErr.HasErrors() {
		cErr.Merge(hErr)
		return cErr
	}

	return nil
}
