package competitor

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

type competitorService struct {
	competitorRepository repo.CompetitorRepository
	pageService          svc.PageService
}

var _ svc.CompetitorService = (*competitorService)(nil)
