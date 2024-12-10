package db

import (
	"context"
	"database/sql"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type competitorRepository struct {
	db *sql.DB
}

func NewCompetitorRepository(db *sql.DB) (interfaces.CompetitorRepository, error) {
	return &competitorRepository{db: db}, nil
}

func (r *competitorRepository) Create(ctx context.Context, competitor *models.Competitor) error {
	// Implementation
	return nil
}

func (r *competitorRepository) Update(ctx context.Context, competitor *models.Competitor) error {
	// Implementation
	return nil
}

func (r *competitorRepository) Delete(ctx context.Context, id int) error {
	// Implementation
	return nil
}
func (r *competitorRepository) GetByID(ctx context.Context, id int) (*models.Competitor, error) {
	// Implementation
	return nil, nil
}

func (r *competitorRepository) List(ctx context.Context, limit, offset int) ([]models.Competitor, int, error) {
	// Implementation
	return nil, 0, nil
}

func (r *competitorRepository) FindByURLHash(ctx context.Context, hash string) ([]models.Competitor, error) {
	// Implementation
	return nil, nil
}

func (r *competitorRepository) AddURL(ctx context.Context, competitorID int, url string) error {
	// Implementation
	return nil
}

func (r *competitorRepository) RemoveURL(ctx context.Context, competitorID int, url string) error {
	// Implementation
	return nil
}
