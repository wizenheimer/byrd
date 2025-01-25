package diff

import (
	"context"
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

type diffService struct {
	aiService ai.AIService
	processor *utils.MarkdownProcessor
	logger    *logger.Logger
}

// NewDiffService creates a new diff service
func NewDiffService(aiService ai.AIService, logger *logger.Logger) (DiffService, error) {
	processor, err := utils.NewMarkdownProcessor()
	if err != nil {
		return nil, err
	}
	return &diffService{
		aiService: aiService,
		processor: processor,
		logger:    logger.WithFields(map[string]interface{}{"module": "diff_service"}),
	}, nil
}

func (d *diffService) Compare(ctx context.Context, content1, content2 *models.ScreenshotContent, profileFields []string) (*models.DynamicChanges, error) {
	d.logger.Debug("comparing content", zap.Any("profile_fields", profileFields))
	markdownContent1, err := d.processor.Process(content1.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown content 1: %w", err)
	}

	markdownContent2, err := d.processor.Process(content2.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown content 2: %w", err)
	}

	aiAnalysis, err := d.aiService.AnalyzeContentDifferences(ctx, markdownContent1, markdownContent2, profileFields)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze content differences: %w", err)
	}

	return aiAnalysis, nil
}
