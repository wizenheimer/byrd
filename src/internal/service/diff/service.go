package diff

import (
	"context"
	"fmt"

	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type diffService struct {
	diffRepo  repo.DiffRepository
	aiService svc.AIService
	processor *MarkdownProcessor
	logger    *logger.Logger
}

func NewDiffService(
	diffRepo repo.DiffRepository,
	aiService svc.AIService,
	// screenshot interfaces.ScreenshotService,
	logger *logger.Logger,
) (svc.DiffService, error) {
	logger.Debug("creating new diff service")

	// Initialize markdown processor
	processor, err := NewMarkdownProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create new markdown processor: %w", err)
	}

	// Initialize diff service
	diffService := &diffService{
		diffRepo:  diffRepo,
		aiService: aiService,
		processor: processor,
		logger:    logger.WithFields(map[string]interface{}{"module": "diff_service"}),
	}
	return diffService, nil
}

// Get creates a diff for between any 2 versions of a URL
// This is used when both versions are available in the storage via screenshotService
func (d *diffService) Get(ctx context.Context, req api_models.URLDiffRequest) (*core_models.DynamicChanges, error) {
	d.logger.Debug("creating diff", zap.Any("url", req.URL), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_number_1", req.WeekNumber1), zap.Any("week_day_2", req.WeekDay2), zap.Any("week_number_2", req.WeekNumber2))
	return d.diffRepo.Get(ctx, req)
}

// Compare compares the contents of 2 screenshots
// This is used when the screenshots are available in memory
// If the AI analysis is not available, it falls back to the default diff analysis
func (d *diffService) Compare(ctx context.Context, content1, content2 *core_models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*core_models.DynamicChanges, error) {
	d.logger.Debug("comparing screenshot contents", zap.Any("content_1", len(content1.HTMLContent)), zap.Any("content_2", len(content2.HTMLContent)))
	// Attempt to get the AI analysis for the screenshot contents if available
	req := api_models.URLDiffRequest{
		URL:         content1.Metadata.SourceURL,
		WeekDay1:    content1.Metadata.WeekDay,
		WeekNumber1: content1.Metadata.WeekNumber,
		WeekDay2:    content2.Metadata.WeekDay,
		WeekNumber2: content2.Metadata.WeekNumber,
	}
	aiAnalysis, err := d.Get(ctx, req)
	if err == nil {
		return aiAnalysis, nil
	}
	// If not available, fallback to the default diff analysis
	markdownContent1, err := d.processor.Process(content1.HTMLContent)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown content 1: %w", err)
	}

	markdownContent2, err := d.processor.Process(content2.HTMLContent)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown content 2: %w", err)
	}

	profileFields := []string{
		"customers",
		"messaging",
		"product",
		"pricing",
		"partnerships",
		"roadmap",
	}

	aiAnalysis, err = d.aiService.AnalyzeContentDifferences(ctx, markdownContent1, markdownContent2, profileFields)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze content differences: %w", err)
	}

	if persist {
		// Save the AI analysis for future use
		if err = d.diffRepo.Set(ctx, req, aiAnalysis); err != nil {
			return nil, fmt.Errorf("failed to save diff analysis: %w", err)
		}
	}

	return aiAnalysis, nil
}
