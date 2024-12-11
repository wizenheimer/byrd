package db

import (
	"context"
	"database/sql"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type diffRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDiffRepository(db *sql.DB, logger *logger.Logger) (interfaces.DiffRepository, error) {
	logger.Debug("creating new diff repository")
	return &diffRepository{
		db:     db,
		logger: logger.WithFields(map[string]interface{}{"module": "diff_repository"}),
	}, nil
}

// SaveDiff saves the diff analysis for a URL
func (r *diffRepository) SaveDiff(ctx context.Context, url string, diff *models.URLDiffAnalysis) error {
	r.logger.Debug("saving diff analysis", zap.String("url", url))
	// Implementation
	return nil
}

// GetDiff retrieves the diff analysis for a URL, week day, and week number
func (r *diffRepository) GetDiff(ctx context.Context, url, weekDay, weekNumber string) (*models.URLDiffAnalysis, error) {
	r.logger.Debug("getting diff analysis", zap.String("url", url), zap.String("week_day", weekDay), zap.String("week_number", weekNumber))
	// Implementation
	return nil, nil
}
