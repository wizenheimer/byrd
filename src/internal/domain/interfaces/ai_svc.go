package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type AIService interface {
	// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
	AnalyzeContentDifferences(ctx context.Context, content1, content2 string) (*models.URLDiffAnalysis, error)
	// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
	AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 []byte) (*models.URLDiffAnalysis, error)
	// EnrichReport enriches a weekly report with AI-generated summaries
	EnrichReport(ctx context.Context, report *models.WeeklyReport) error
}
