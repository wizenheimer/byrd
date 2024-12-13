package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type ScreenshotService interface {
	// TakeScreenshot takes a screenshot of the given URL
	TakeScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, error)

	// GetScreenshot retrieves a screenshot from the storage
	GetScreenshotImage(ctx context.Context, url, weekNumber, weekDay string) (*models.ScreenshotImageResponse, error)

	// GetContent retrieves the content of a screenshot from the storage
	GetScreenshotContent(ctx context.Context, url, weekNumber, weekDay string) (*models.ScreenshotContentResponse, error)
}

type DiffService interface {
	// CreateDiff creates a diff for between versions of a URL
	CreateDiff(ctx context.Context, req models.URLDiffRequest) (*models.URLDiffAnalysis, error)

	// CreateReport creates a weekly report for the given competitor
	CreateReport(ctx context.Context, req models.WeeklyReportRequest) (*models.WeeklyReport, error)
}

type CompetitorService interface {
	Create(ctx context.Context, input models.CompetitorInput) (*models.Competitor, error)
	Update(ctx context.Context, id int, input models.CompetitorInput) (*models.Competitor, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Competitor, error)
	List(ctx context.Context, limit, offset int) ([]models.Competitor, int, error)
	FindByURLHash(ctx context.Context, hash string) ([]models.Competitor, error)
	AddURL(ctx context.Context, id int, url string) error
	RemoveURL(ctx context.Context, id int, url string) error
}

type NotificationService interface {
	SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResults, error)
}

type AIService interface {
	// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
	AnalyzeContentDifferences(ctx context.Context, content1, content2 string) (*models.URLDiffAnalysis, error)
	// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
	AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 []byte) (*models.URLDiffAnalysis, error)
	// EnrichReport enriches a weekly report with AI-generated summaries
	EnrichReport(ctx context.Context, report *models.WeeklyReport) error
}
