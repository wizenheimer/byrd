package ai

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type openAIService struct {
	apiKey     string
	httpClient interfaces.HTTPClient
}

func NewOpenAIService(apiKey string, httpClient interfaces.HTTPClient) (interfaces.AIService, error) {
	return &openAIService{
		apiKey:     apiKey,
		httpClient: httpClient,
	}, nil
}

// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
func (s *openAIService) AnalyzeContentDifferences(ctx context.Context, content1, content2 string) (*models.URLDiffAnalysis, error) {
	// Implementation
	return nil, nil
}

// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
func (s *openAIService) AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 []byte) (*models.URLDiffAnalysis, error) {
	// Implementation
	return nil, nil
}

// EnrichReport enriches a weekly report with AI-generated summaries
func (s *openAIService) EnrichReport(ctx context.Context, report *models.WeeklyReport) error {
	// Implementation
	return nil
}
