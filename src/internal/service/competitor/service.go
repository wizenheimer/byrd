package competitor

import (
	"context"

	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type competitorService struct {
	repo   repo.CompetitorRepository
	logger *logger.Logger
}

func NewCompetitorService(repo repo.CompetitorRepository, logger *logger.Logger) (svc.CompetitorService, error) {
	logger.Debug("creating new competitor service")
	return &competitorService{
		repo:   repo,
		logger: logger.WithFields(map[string]interface{}{"module": "competitor_service"}),
	}, nil
}

func (s *competitorService) Create(ctx context.Context, input core_models.CompetitorInput) (*core_models.Competitor, error) {
	s.logger.Debug("creating new competitor", zap.Any("input", input))
	// TODO: Validate input
	// TODO: Create competitor entity
	// TODO: Save to repository
	return nil, nil
}

func (s *competitorService) Update(ctx context.Context, id int, input core_models.CompetitorInput) (*core_models.Competitor, error) {
	s.logger.Debug("updating competitor", zap.Int("id", id), zap.Any("input", input))
	// TODO: Validate input
	// TODO: Check if competitor exists
	// TODO: Update competitor
	// TODO: Save changes
	return nil, nil
}

func (s *competitorService) Delete(ctx context.Context, id int) error {
	s.logger.Debug("deleting competitor", zap.Int("id", id))
	// TODO: Check if competitor exists
	// TODO: Delete competitor
	return nil
}

func (s *competitorService) Get(ctx context.Context, id int) (*core_models.Competitor, error) {
	s.logger.Debug("getting competitor by ID", zap.Int("id", id))
	// TODO: Get competitor from repository
	return nil, nil
}

func (s *competitorService) List(ctx context.Context, limit, offset int) ([]core_models.Competitor, int, error) {
	s.logger.Debug("listing competitors", zap.Int("limit", limit), zap.Int("offset", offset))
	// TODO: Validate pagination parameters
	// TODO: Get competitors from repository with pagination
	// TODO: Get total count
	return nil, 0, nil
}

func (s *competitorService) FindByURLHash(ctx context.Context, hash string) ([]core_models.Competitor, error) {
	s.logger.Debug("finding competitor by URL hash", zap.String("hash", hash))
	// TODO: Validate hash
	// TODO: Find competitors by URL hash
	return nil, nil
}

func (s *competitorService) AddURL(ctx context.Context, id int, url string) error {
	s.logger.Debug("adding URL to competitor", zap.Int("id", id), zap.String("url", url))
	// TODO: Validate URL
	// TODO: Check if competitor exists
	// TODO: Add URL to competitor
	return nil
}

func (s *competitorService) RemoveURL(ctx context.Context, id int, url string) error {
	s.logger.Debug("removing URL from competitor", zap.Int("id", id), zap.String("url", url))
	// TODO: Validate URL
	// TODO: Check if competitor exists
	// TODO: Remove URL from competitor
	return nil
}
