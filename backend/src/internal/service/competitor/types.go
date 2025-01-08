package competitor

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type competitorService struct {
	competitorRepository repo.CompetitorRepository
	pageService          svc.PageService
	nameFinder           *CompanyNameFinder
	logger               *logger.Logger
}

var _ svc.CompetitorService = (*competitorService)(nil)
