package interfaces

import (
	"context"
	"image"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type ScreenshotService interface {
	// CaptureScreenshot takes a screenshot of the given URL
	CaptureScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*models.ScreenshotResponse, image.Image, string, error)

	// GetPreviousScreenshot retrieves the previous screenshot from the storage
	GetPreviousScreenshotImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, error)

	// GetPreviousScreenshotContent retrieves the previous content of a screenshot from the storage
	GetPreviousScreenshotContent(ctx context.Context, url string) (*models.ScreenshotContentResponse, error)

	// GetScreenshot retrieves a screenshot from the storage
	GetScreenshotImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, error)

	// GetContent retrieves the content of a screenshot from the storage
	GetScreenshotContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotContentResponse, error)

	// ListScreenshots lists the latest content (images or text) for a given URL
	ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, error)
}

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

// URLService is the interface that provides URL management operations
type URLService interface {
	// AddURL: adds a new URL if it does not exist
	AddURL(ctx context.Context, url string) (*models.URL, error)

	// ListURLs: lists all URLs
	ListURLs(ctx context.Context, batchSize int, lastSeenID *uuid.UUID) (<-chan models.URLBatch, <-chan error)

	// DeleteURL: deletes a URL
	DeleteURL(ctx context.Context, url string) error

	// URLExists: checks if a URL exists
	URLExists(ctx context.Context, url string) (bool, error)
}

type WorkflowService interface {
	// StartWorkflow starts a new workflow
	StartWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error)

	// GetWorkflow retrieves a workflow
	GetWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error)

	// ListWorkflows lists of all workflows
	ListWorkflows(context.Context, models.WorkflowStatus, models.WorkflowType, int) ([]models.WorkflowResponse, int, error)

	// RecoverWorkflow recovers a workflow from a checkpoint
	RecoverWorkflow(ctx context.Context) error

	// Shutdown shuts down the workflow service
	Shutdown(ctx context.Context) error

	// StopWorkflow stops a workflow
	StopWorkflow(ctx context.Context, workflowID models.WorkflowIdentifier) error
}
