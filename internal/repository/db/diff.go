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

func (r *diffRepository) SaveDiff(ctx context.Context, diff *models.DiffAnalysis, metadata models.DiffMetadata) error {
	// Implementation
	return nil
}

func (r *diffRepository) GetDiffHistory(ctx context.Context, params models.DiffHistoryParams) ([]models.DiffReport, error) {
	// Implementation
	return nil, nil
}

func (r *diffRepository) GetLatestDiff(ctx context.Context, url string) (*models.DiffReport, error) {
	// Implementation
	return nil, nil
}
