package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	"github.com/wizenheimer/iris/src/internal/config"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type HandlerContainer struct {
	AIHandler         *handlers.AIHandler
	ScreenshotHandler *handlers.ScreenshotHandler
}

func NewHandlerContainer(
	screenshotService interfaces.ScreenshotService,
	aiService interfaces.AIService,
	logger *logger.Logger,
) *HandlerContainer {
	return &HandlerContainer{
		// Handlers for screenshot management
		ScreenshotHandler: handlers.NewScreenshotHandler(screenshotService, logger),
		// Handlers for AI management
		AIHandler: handlers.NewAIHandler(aiService, logger),
	}
}

// SetupRoutes sets up the routes for the application
// This includes public and private routes
func SetupRoutes(app *fiber.App, handlers *HandlerContainer) {
	setupPublicRoutes(app, handlers)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers)
	}

}

// setupPublicRoutes sets up the public routes for the application
func setupPublicRoutes(app *fiber.App, handlers *HandlerContainer) {
	// Public routes for production and development
	// public := app.Group("/api/v1")

	// Workflow routes
	// workflow := public.Group("/workflow")
	// workflow.Post("/", handlers.WorkflowHandler.StartWorkflow)
	// workflow.Delete("/", handlers.WorkflowHandler.StopWorkflow)
	// workflow.Get("/", handlers.WorkflowHandler.GetWorkflow)
	// workflow.Get("/list", handlers.WorkflowHandler.ListWorkflows)
}

// setupPrivateRoutes sets up the private routes for the application
func setupPrivateRoutes(app *fiber.App, handlers *HandlerContainer) {
	// Private routes for development
	private := app.Group("/dev/v1")

	// Screenshot routes
	screenshot := private.Group("/screenshot")
	screenshot.Post("/", handlers.ScreenshotHandler.CreateScreenshot)
	screenshot.Get("/", handlers.ScreenshotHandler.ListScreenshots)
	screenshot.Get("/image", handlers.ScreenshotHandler.GetScreenshotImage)
	screenshot.Get("/content", handlers.ScreenshotHandler.GetScreenshotContent)

	// AI routes
	ai := private.Group("/ai")
	ai.Post("/content", handlers.AIHandler.AnalyzeContentDifferences)
	ai.Post("/visual", handlers.AIHandler.AnalyzeVisualDifferences)

}
