package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type DiffService interface {
	// CompareImageResponse compares the contents of two image responses
	CompareImageResponse(ctx context.Context, img1, img2 *models.ScreenshotImageResponse) (*models.ImageResponseDiff, error)

	// CompareHTMLContentResponse compares the contents of two HTML responses
	CompareHTMLContentResponse(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse) (*models.HTMLContentResponseDiff, error)

	// CompareScreenshotContents compares the contents of two screenshots using AI services
	CompareScreenshotContents(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse) (*models.URLDiffAnalysis, error)

	// CompareScreenshotImages compares the images of two screenshots using AI services
	CompareScreenshotImages(ctx context.Context, img1, img2 *models.ScreenshotImageResponse) (*models.URLDiffAnalysis, error)

	// GetDiffAnalysis returns the diff analysis for the given URL
	GetDiffAnalysis(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error)

	// CreateReport creates a weekly report for the given competitor
	CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error)
}
