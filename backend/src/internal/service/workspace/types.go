package workspace

import (
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceService struct {
	workspaceRepo     repo.WorkspaceRepository
	competitorService svc.CompetitorService
	userService       svc.UserService
	logger            *logger.Logger
}

// compile time check if the interface is implemented
var _ svc.WorkspaceService = (*workspaceService)(nil)
