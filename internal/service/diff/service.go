package diff

import (
	"context"
	"fmt"

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

func (s *diffService) CreateDiff(ctx context.Context, req models.DiffRequest) (*models.DiffAnalysis, error) {
	// Implementation
	return nil, nil
}

func (s *diffService) GenerateReport(ctx context.Context, req models.ReportRequest) (*models.AggregatedReport, error) {
	report := &models.AggregatedReport{}

	// Process each URL and aggregate the results
	// Implementation here

	if req.Enriched {
		if err := s.aiService.EnrichReport(ctx, report); err != nil {
			return nil, fmt.Errorf("failed to enrich report: %w", err)
		}
	}

	return report, nil
}

func (s *diffService) GetDiffHistory(ctx context.Context, params models.DiffHistoryParams) (*models.DiffHistoryResponse, error) {
	// Implementation
	return nil, nil
}
