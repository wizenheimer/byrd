package diff

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type diffService struct {
	diffRepo   interfaces.DiffRepository
	aiService  interfaces.AIService
	// screenshot interfaces.ScreenshotService
	logger     *logger.Logger
}

func NewDiffService(
	diffRepo interfaces.DiffRepository,
	aiService interfaces.AIService,
	// screenshot interfaces.ScreenshotService,
	logger *logger.Logger,
) (interfaces.DiffService, error) {
	logger.Debug("creating new diff service")

	// Initialize diff service
	diffService := &diffService{
		diffRepo:   diffRepo,
		aiService:  aiService,
		// screenshot: screenshot,
		logger:     logger.WithFields(map[string]interface{}{"module": "diff_service"}),
	}
	return diffService, nil
}

// CreateDiff creates a diff for between any 2 versions of a URL
// This is used when both versions are available in the storage via screenshotService
func (d *diffService) GetDiffAnalysis(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error) {
	d.logger.Debug("creating diff", zap.Any("url", req.URL), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_number_1", req.WeekNumber1), zap.Any("week_day_2", req.WeekDay2), zap.Any("week_number_2", req.WeekNumber2))
	// Implementation
	return nil, nil
}

// CreateReport creates a weekly report for the given competitor
// This is used when the competitor has multiple URLs and each URL has multiple available versions
func (d *diffService) CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error) {
	d.logger.Debug("creating report", zap.Any("week_number", req.WeekNumber), zap.Any("week_day_1", req.WeekDay1), zap.Any("week_day_2", req.WeekDay2), zap.Any("urls", req.URLs))
	// Implementation
	return nil, nil
}

// CompareImageResponse compares 2 image responses and returns the ImageResponseDiff
func (d *diffService) CompareImageResponse(ctx context.Context, img1, img2 *models.ScreenshotImageResponse) (*models.ImageResponseDiff, error) {
	d.logger.Debug("comparing image responses", zap.Any("image_1", img1), zap.Any("image_2", img2))
	// Implementation
	return nil, nil
}

// CompareHTMLContentResponse compares 2 HTML content responses and returns the HTMLContentResponseDiff
func (d *diffService) CompareHTMLContentResponse(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse) (*models.HTMLContentResponseDiff, error) {
	d.logger.Debug("comparing HTML content responses", zap.Any("content_1", content1), zap.Any("content_2", content2))
	// Implementation
	return nil, nil
}

// CompareScreenshotContents compares the contents of 2 screenshots and returns the URLDiffAnalysis
func (d *diffService) CompareScreenshotContents(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse) (*models.URLDiffAnalysis, error) {
	d.logger.Debug("comparing screenshot contents", zap.Any("content_1", content1), zap.Any("content_2", content2))
	// Implementation
	return nil, nil
}

// CompareScreenshotImages compares the images of 2 screenshots and returns the URLDiffAnalysis
func (d *diffService) CompareScreenshotImages(ctx context.Context, img1, img2 *models.ScreenshotImageResponse) (*models.URLDiffAnalysis, error) {
	d.logger.Debug("comparing screenshot images", zap.Any("image_1", img1), zap.Any("image_2", img2))
	// Implementation
	return nil, nil
}
