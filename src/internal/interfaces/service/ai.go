package interfaces

import (
	"context"
	"image"

	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

type AIService interface {
	// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
	AnalyzeContentDifferences(ctx context.Context, version1, version2 string, fields []string) (*core_models.DynamicChanges, error)

	// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
	AnalyzeVisualDifferences(ctx context.Context, version1, version2 image.Image, fields []string) (*core_models.DynamicChanges, error)

	// EnrichReport enriches a weekly report with AI-generated summaries
	EnrichReport(ctx context.Context, report *core_models.WeeklyReport) error
}
