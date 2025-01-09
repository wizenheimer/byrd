package competitor

import (
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type competitorService struct {
	competitorRepository repo.CompetitorRepository
	pageService          svc.PageService
	nameFinder           *CompanyNameFinder
	logger               *logger.Logger
}

var _ svc.CompetitorService = (*competitorService)(nil)
