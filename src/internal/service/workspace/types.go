package workspace

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

type workspaceService struct {
	workspaceRepo     repo.WorkspaceRepository
	competitorService svc.CompetitorService
	userService       svc.UserService
}

// compile time check if the interface is implemented
var _ svc.WorkspaceService = (*workspaceService)(nil)
