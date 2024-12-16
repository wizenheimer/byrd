package interfaces

import (
	"context"
	"image"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type DiffService interface {
	// CreateDiff creates a diff for between versions of a URL
	CreateDiff(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error)

	// CreateDiffFromScreenshotImages creates a diff between two screenshots
	CreateCurrentDiffFromScreenshotImages(ctx context.Context, url string, screenshot1, screenshot2 image.Image) (*models.URLDiffAnalysis, error)

	// CreateDiffFromScreenshotContents creates a diff between two screenshots
	CreateCurrentDiffFromScreenshotContents(ctx context.Context, url, content1, content2 string) (*models.URLDiffAnalysis, error)

	// CreateReport creates a weekly report for the given competitor
	CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error)
}
