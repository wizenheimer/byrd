package db

import (
	"context"
	"database/sql"

	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
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

// Set saves the diff analysis for a URL
func (r *diffRepository) Set(ctx context.Context, req api_models.URLDiffRequest, diff *core_models.DynamicChanges) error {
	// Implementation
	return nil
}

// Get retrieves the diff analysis for a URL, week day, and week number
func (r *diffRepository) Get(ctx context.Context, req api_models.URLDiffRequest) (*core_models.DynamicChanges, error) {
	r.logger.Debug("compairing diff", zap.Any("url", req.URL), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_number_1", req.WeekNumber1), zap.Any("week_day_2", req.WeekDay2), zap.Any("week_number_2", req.WeekNumber2))
	// Implementation
	return nil, nil
}
