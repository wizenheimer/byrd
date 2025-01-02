package history

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

// compile time check if the interface is implemented
var _ svc.PageHistoryService = (*pageHistoryService)(nil)

type pageHistoryService struct {
	pageHistoryRepo repo.PageHistoryRepository
}
