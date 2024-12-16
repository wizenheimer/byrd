package diff

import (
	"context"
	"image"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type diffService struct {
	diffRepo   interfaces.DiffRepository
	aiService  interfaces.AIService
	screenshot interfaces.ScreenshotService
	logger     *logger.Logger
}

func NewDiffService(
	diffRepo interfaces.DiffRepository,
	aiService interfaces.AIService,
	screenshot interfaces.ScreenshotService,
	logger *logger.Logger,
) (interfaces.DiffService, error) {
	logger.Debug("creating new diff service")

	// Initialize diff service
	diffService := &diffService{
		diffRepo:   diffRepo,
		aiService:  aiService,
		screenshot: screenshot,
		logger:     logger.WithFields(map[string]interface{}{"module": "diff_service"}),
	}
	return diffService, nil
}

// CreateDiff creates a diff for between any 2 versions of a URL
// This is used when both versions are available in the storage via screenshotService
func (d *diffService) CreateDiff(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error) {
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

// CreateCurrentDiffFromScreenshotImages creates a diff between two screenshots
// This is used when the screenshots are available in memory
// Once the diff is created, it is saved to the storage
func (d *diffService) CreateCurrentDiffFromScreenshotImages(ctx context.Context, url string, screenshot1, screenshot2 image.Image) (*models.URLDiffAnalysis, error) {
	d.logger.Debug("creating diff from screenshot images", zap.Any("url", url))
	d.logger.Debug("saving diff to storage", zap.Any("url", url))
	// Implementation
	return &models.URLDiffAnalysis{
		Branding:    make([]string, 0),
		Integration: make([]string, 0),
		Pricing:     make([]string, 0),
		Product:     make([]string, 0),
		Positioning: make([]string, 0),
		Partnership: make([]string, 0),
	}, nil
}

// CreateCurrentDiffFromScreenshotContents creates a diff between two screenshots
// This is used when the screenshots are available in memory
// Once the diff is created, it is saved to the storage
func (d *diffService) CreateCurrentDiffFromScreenshotContents(ctx context.Context, url, content1, content2 string) (*models.URLDiffAnalysis, error) {
	d.logger.Debug("creating diff from screenshot contents", zap.Any("url", url))
	d.logger.Debug("saving diff to storage", zap.Any("url", url))
	// Implementation
	return &models.URLDiffAnalysis{
		Branding:    make([]string, 0),
		Integration: make([]string, 0),
		Pricing:     make([]string, 0),
		Product:     make([]string, 0),
		Positioning: make([]string, 0),
		Partnership: make([]string, 0),
	}, nil
}
