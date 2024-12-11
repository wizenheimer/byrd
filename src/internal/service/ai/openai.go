package ai

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/logger"
	"go.uber.org/zap"
)

type openAIService struct {
	apiKey     string
	httpClient interfaces.HTTPClient
	logger     *logger.Logger
}

func NewOpenAIService(apiKey string, httpClient interfaces.HTTPClient, logger *logger.Logger) (interfaces.AIService, error) {
	logger.Debug("creating new openAI service")
	return &openAIService{
		apiKey:     apiKey,
		httpClient: httpClient,
		logger:     logger.WithFields(map[string]interface{}{"module": "ai_service"}),
	}, nil
}

// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
func (s *openAIService) AnalyzeContentDifferences(ctx context.Context, content1, content2 string) (*models.URLDiffAnalysis, error) {
	s.logger.Debug("analyzing content differences", zap.Any("content1_len", len(content1)), zap.Any("content2_len", len(content2)))
	// Implementation
	return nil, nil
}

// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
func (s *openAIService) AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 []byte) (*models.URLDiffAnalysis, error) {
	s.logger.Debug("analyzing visual differences")
	// Implementation
	return nil, nil
}

// EnrichReport enriches a weekly report with AI-generated summaries
func (s *openAIService) EnrichReport(ctx context.Context, report *models.WeeklyReport) error {
	s.logger.Debug("enriching report", zap.Any("report", report))
	// Implementation
	return nil
}
