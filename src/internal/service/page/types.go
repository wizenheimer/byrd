package page

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

// compile time check if the interface is implemented
var _ svc.PageService = (*pageService)(nil)

type pageService struct {
	pageRepo           repo.PageRepository
	pageHistoryService svc.PageHistoryService
}
