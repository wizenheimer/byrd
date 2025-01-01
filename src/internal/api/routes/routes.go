package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type HandlerContainer struct {
	ScreenshotHandler   *handlers.ScreenshotHandler
	DiffHandler         *handlers.DiffHandler
	CompetitorHandler   *handlers.CompetitorHandler
	NotificationHandler *handlers.NotificationHandler
	URLHandler          *handlers.URLHandler
	WorkflowHandler     *handlers.WorkflowHandler
	AIHandler           *handlers.AIHandler
}

func NewHandlerContainer(
	screenshotService interfaces.ScreenshotService,
	urlService interfaces.URLService,
	aiService interfaces.AIService,
	diffService interfaces.DiffService,
	competitorService interfaces.CompetitorService,
	notificationService interfaces.NotificationService,
	workflowService interfaces.WorkflowService,
	logger *logger.Logger,
) *HandlerContainer {
	return &HandlerContainer{
		// Handlers for screenshot management
		ScreenshotHandler: handlers.NewScreenshotHandler(screenshotService, logger),
		// Handlers for URL management
		URLHandler: handlers.NewURLHandler(urlService, logger),
		// Handlers for AI management
		AIHandler:           handlers.NewAIHandler(aiService, logger),
		DiffHandler:         handlers.NewDiffHandler(diffService, logger),
		CompetitorHandler:   handlers.NewCompetitorHandler(competitorService, logger),
		NotificationHandler: handlers.NewNotificationHandler(notificationService, logger),
		// Handler for workflow management
		WorkflowHandler: handlers.NewWorkflowHandler(workflowService, logger),
	}
}

func SetupRoutes(app *fiber.App, handlers *HandlerContainer) {
	api := app.Group("/api/v1")

	// Screenshot routes
	screenshot := api.Group("/screenshot")
	screenshot.Post("/", handlers.ScreenshotHandler.CreateScreenshot)
	screenshot.Get("/", handlers.ScreenshotHandler.ListScreenshots)
	screenshot.Get("/image", handlers.ScreenshotHandler.GetScreenshotImage)
	screenshot.Get("/content", handlers.ScreenshotHandler.GetScreenshotContent)

	// AI routes
	ai := api.Group("/ai")
	ai.Post("/content", handlers.AIHandler.AnalyzeContentDifferences)
	ai.Post("/visual", handlers.AIHandler.AnalyzeVisualDifferences)

	// URL routes
	url := api.Group("/url")
	url.Post("/", handlers.URLHandler.AddURL)
	url.Get("/", handlers.URLHandler.ListURLs)
	url.Delete("/", handlers.URLHandler.DeleteURL)

	// Diff routes
	diff := api.Group("/diff")
	diff.Post("/get", handlers.DiffHandler.Get)
	diff.Get("/compare", handlers.DiffHandler.Compare)

	// Workflow routes
	workflow := api.Group("/workflow")
	workflow.Post("/", handlers.WorkflowHandler.StartWorkflow)
	workflow.Delete("/", handlers.WorkflowHandler.StopWorkflow)
	workflow.Get("/", handlers.WorkflowHandler.GetWorkflow)
	workflow.Get("/list", handlers.WorkflowHandler.ListWorkflows)

	// Competitor routes
	competitors := api.Group("/competitors")
	competitors.Post("/", handlers.CompetitorHandler.CreateCompetitor)
	competitors.Get("/", handlers.CompetitorHandler.ListCompetitors)
	competitors.Get("/id/:id", handlers.CompetitorHandler.GetCompetitor)

	// Notification routes
	api.Post("/notify", handlers.NotificationHandler.SendNotification)
}
