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

func (s *openAIService) AnalyzeDifferences(ctx context.Context, content1, content2 string) (*models.DiffAnalysis, error) {
	// Implementation
	return nil, nil
}

func (s *openAIService) EnrichReport(ctx context.Context, report *models.AggregatedReport) error {
	// Implementation
	return nil
}
