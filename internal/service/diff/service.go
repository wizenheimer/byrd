package diff

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type diffService struct {
	diffRepo   interfaces.DiffRepository
	aiService  interfaces.AIService
	screenshot interfaces.ScreenshotService
}

func NewDiffService(
	diffRepo interfaces.DiffRepository,
	aiService interfaces.AIService,
	screenshot interfaces.ScreenshotService,
) (interfaces.DiffService, error) {
	diffService := &diffService{
		diffRepo:   diffRepo,
		aiService:  aiService,
		screenshot: screenshot,
	}
	return diffService, nil
}

func (s *diffService) CreateDiff(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error) {
	// Implementation
	return nil, nil
}

func (s *diffService) CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error) {
	// Implementation
	return nil, nil
}
