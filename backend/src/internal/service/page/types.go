// ./src/internal/service/page/types.go
package page

import (
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// compile time check if the interface is implemented
var _ svc.PageService = (*pageService)(nil)

type pageService struct {
	pageRepo           repo.PageRepository
	pageHistoryService svc.PageHistoryService
	logger             *logger.Logger
}
