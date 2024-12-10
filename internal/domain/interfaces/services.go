package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/models"
)

type ScreenshotService interface {
	TakeScreenshot(ctx context.Context, opts models.ScreenshotOptions) (*models.ScreenshotResponse, error)
	GetScreenshot(ctx context.Context, hash, weekNumber, runID string) (*models.ScreenshotResponse, error)
	GetContent(ctx context.Context, hash, weekNumber, runID string) (*models.ScreenshotResponse, error)
}

type DiffService interface {
	CreateDiff(ctx context.Context, req models.DiffRequest) (*models.DiffAnalysis, error)
	GetDiffHistory(ctx context.Context, params models.DiffHistoryParams) (*models.DiffHistoryResponse, error)
	GenerateReport(ctx context.Context, req models.ReportRequest) (*models.AggregatedReport, error)
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
	AnalyzeDifferences(ctx context.Context, content1, content2 string) (*models.DiffAnalysis, error)
	EnrichReport(ctx context.Context, report *models.AggregatedReport) error
}
