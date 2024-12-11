package db

import (
	"context"
	"database/sql"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type diffRepository struct {
	db *sql.DB
}

func NewDiffRepository(db *sql.DB) (interfaces.DiffRepository, error) {
	return &diffRepository{
		db: db,
	}, nil
}

// SaveDiff saves the diff analysis for a URL
func (r *diffRepository) SaveDiff(ctx context.Context, url string, diff *models.URLDiffAnalysis) error {
	// Implementation
	return nil
}

// GetDiff retrieves the diff analysis for a URL, week day, and week number
func (r *diffRepository) GetDiff(ctx context.Context, url, weekDay, weekNumber string) (*models.URLDiffAnalysis, error) {
	// Implementation
	return nil, nil
}
