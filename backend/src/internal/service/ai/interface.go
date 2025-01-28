// ./src/internal/service/ai/interface.go
package ai

import (
	"context"
	"errors"
	"image"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type AIService interface {
	// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
	AnalyzeContentDifferences(ctx context.Context, version1, version2 string, fields []string) (*models.DynamicChanges, error)

	// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
	AnalyzeVisualDifferences(ctx context.Context, version1, version2 image.Image, fields []string) (*models.DynamicChanges, error)

	// SummarizeChanges summarizes the changes in a report
	SummarizeChanges(ctx context.Context, changes []*models.DynamicChanges) ([]models.CategoryChange, error)
}

var (
	ErrBuildingProfile         = errors.New("failed to build profile")
	ErrPreparingChatCompletion = errors.New("failed to prepare chat completion")
	ErrConvertingImageToBase64 = errors.New("failed to convert image to base64")
	ErrEncounteredRefusal      = errors.New("refusal encountered in chat completion")
	ErrParsingChanges          = errors.New("failed to parse changes")
	ErrEncodingImage           = errors.New("failed to encode image")
)

var (
	ErrProfileParsing      = errors.New("couldn't parse profile")
	ErrProfileNameMissing  = errors.New("profile name is required")
	ErrProfileFieldParsing = errors.New("couldn't parse profile field")
)

var (
	ErrProfileFieldNotFound = errors.New("profile field not found")
)
