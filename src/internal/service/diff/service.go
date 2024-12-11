package diff

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/logger"
	"go.uber.org/zap"
)

type diffService struct {
	diffRepo   interfaces.DiffRepository
	aiService  interfaces.AIService
	screenshot interfaces.ScreenshotService
	logger     *logger.Logger
}

func NewDiffService(
	diffRepo interfaces.DiffRepository,
	aiService interfaces.AIService,
	screenshot interfaces.ScreenshotService,
	logger *logger.Logger,
) (interfaces.DiffService, error) {
	logger.Debug("creating new diff service")

	// Initialize diff service
	diffService := &diffService{
		diffRepo:   diffRepo,
		aiService:  aiService,
		screenshot: screenshot,
		logger:     logger.WithFields(map[string]interface{}{"module": "diff_service"}),
	}
	return diffService, nil
}

func (s *diffService) CreateDiff(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error) {
	s.logger.Debug("creating diff", zap.Any("url", req.URL), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_number_1", req.WeekNumber1), zap.Any("week_day_2", req.WeekDay2), zap.Any("week_number_2", req.WeekNumber2))
	// Implementation
	return nil, nil
}

func (s *diffService) CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error) {
	s.logger.Debug("creating report", zap.Any("week_number", req.WeekNumber), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_day_2", req.WeekDay2), zap.Any("urls", req.URLs))
	// Implementation
	return nil, nil
}
