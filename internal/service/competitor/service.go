package competitor

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type competitorService struct {
	repo interfaces.CompetitorRepository
}

func NewCompetitorService(repo interfaces.CompetitorRepository) (interfaces.CompetitorService, error) {
	return &competitorService{
		repo: repo,
	}, nil
}

func (s *competitorService) Create(ctx context.Context, input models.CompetitorInput) (*models.Competitor, error) {
	// TODO: Validate input
	// TODO: Create competitor entity
	// TODO: Save to repository
	return nil, nil
}

func (s *competitorService) Update(ctx context.Context, id int, input models.CompetitorInput) (*models.Competitor, error) {
	// TODO: Validate input
	// TODO: Check if competitor exists
	// TODO: Update competitor
	// TODO: Save changes
	return nil, nil
}

func (s *competitorService) Delete(ctx context.Context, id int) error {
	// TODO: Check if competitor exists
	// TODO: Delete competitor
	return nil
}

func (s *competitorService) Get(ctx context.Context, id int) (*models.Competitor, error) {
	// TODO: Get competitor from repository
	return nil, nil
}

func (s *competitorService) List(ctx context.Context, limit, offset int) ([]models.Competitor, int, error) {
	// TODO: Validate pagination parameters
	// TODO: Get competitors from repository with pagination
	// TODO: Get total count
	return nil, 0, nil
}

func (s *competitorService) FindByURLHash(ctx context.Context, hash string) ([]models.Competitor, error) {
	// TODO: Validate hash
	// TODO: Find competitors by URL hash
	return nil, nil
}

func (s *competitorService) AddURL(ctx context.Context, id int, url string) error {
	// TODO: Validate URL
	// TODO: Check if competitor exists
	// TODO: Add URL to competitor
	return nil
}

func (s *competitorService) RemoveURL(ctx context.Context, id int, url string) error {
	// TODO: Validate URL
	// TODO: Check if competitor exists
	// TODO: Remove URL from competitor
	return nil
}
