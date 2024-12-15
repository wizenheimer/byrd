package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type HandlerContainer struct {
	ScreenshotHandler   *handlers.ScreenshotHandler
	WorkflowHandler     *handlers.WorkflowHandler
	DiffHandler         *handlers.DiffHandler
	CompetitorHandler   *handlers.CompetitorHandler
	NotificationHandler *handlers.NotificationHandler
	URLHandler          *handlers.URLHandler
}

func NewHandlerContainer(
	screenshotService interfaces.ScreenshotService,
	urlService interfaces.URLService,
	diffService interfaces.DiffService,
	competitorService interfaces.CompetitorService,
	notificationService interfaces.NotificationService,
	workflowService interfaces.WorkflowService,
	logger *logger.Logger,
) *HandlerContainer {
	return &HandlerContainer{
		// Handlers for screenshot management
		ScreenshotHandler: handlers.NewScreenshotHandler(screenshotService, logger),
		// Handlers for workflow management
		WorkflowHandler: handlers.NewWorkflowHandler(workflowService, logger),
		// Handlers for URL management
		URLHandler:          handlers.NewURLHandler(urlService, logger),
		DiffHandler:         handlers.NewDiffHandler(diffService, logger),
		CompetitorHandler:   handlers.NewCompetitorHandler(competitorService, logger),
		NotificationHandler: handlers.NewNotificationHandler(notificationService, logger),
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

	// URL routes
	url := api.Group("/url")
	url.Post("/", handlers.URLHandler.AddURL)
	url.Get("/", handlers.URLHandler.ListURLs)
	url.Delete("/", handlers.URLHandler.DeleteURL)

	// Workflow routes
	workflow := api.Group("/workflow")
	workflow.Post("/", handlers.WorkflowHandler.StartWorkflow)
	workflow.Get("/", handlers.WorkflowHandler.GetWorkflow)
	workflow.Delete("/", handlers.WorkflowHandler.StopWorkflow)
	workflow.Get("/list", handlers.WorkflowHandler.ListWorkflows)
	workflow.Post("/recover", handlers.WorkflowHandler.RecoverWorkflow)

	// Diff routes
	diff := api.Group("/diff")
	diff.Post("/create", handlers.DiffHandler.CreateDiff)
	diff.Get("/report", handlers.DiffHandler.CreateReport)

	// Competitor routes
	competitors := api.Group("/competitors")
	competitors.Post("/", handlers.CompetitorHandler.CreateCompetitor)
	competitors.Get("/", handlers.CompetitorHandler.ListCompetitors)
	competitors.Get("/id/:id", handlers.CompetitorHandler.GetCompetitor)

	// Notification routes
	api.Post("/notify", handlers.NotificationHandler.SendNotification)
}
