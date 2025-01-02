package interfaces

import (
	"context"
	"image"

	models "github.com/wizenheimer/iris/src/internal/models/core"
)

type AIService interface {
	// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
	AnalyzeContentDifferences(ctx context.Context, version1, version2 string, fields []string) (*models.DynamicChanges, error)

	// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
	AnalyzeVisualDifferences(ctx context.Context, version1, version2 image.Image, fields []string) (*models.DynamicChanges, error)
}
