package history

import (
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// compile time check if the interface is implemented
var _ svc.PageHistoryService = (*pageHistoryService)(nil)

type pageHistoryService struct {
	pageHistoryRepo   repo.PageHistoryRepository
	screenshotService svc.ScreenshotService
	diffService       svc.DiffService
	logger            *logger.Logger
}
