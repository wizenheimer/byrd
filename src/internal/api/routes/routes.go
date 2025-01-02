package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	"github.com/wizenheimer/iris/src/internal/api/middleware"
	"github.com/wizenheimer/iris/src/internal/config"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type HandlerContainer struct {
	AIHandler         *handlers.AIHandler
	ScreenshotHandler *handlers.ScreenshotHandler
	UserHandler       *handlers.UserHandler
	WorkspaceHandler  *handlers.WorkspaceHandler
	WorkflowHandler   *handlers.WorkflowHandler
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
		// Handlers for user management
		UserHandler: handlers.NewUserHandler(),
		// Handlers for workspace management
		WorkspaceHandler: handlers.NewWorkspaceHandler(),
		// Handlers for workflow management
		WorkflowHandler: handlers.NewWorkflowHandler(nil, logger),
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
func setupPublicRoutes(app *fiber.App, h *HandlerContainer) {
	// Public routes for production and development
	public := app.Group("/api/v1")

	// Workspace routes
	wh := h.WorkspaceHandler
	workspace := public.Group("/workspace", middleware.ClerkAuthenticationMiddleware)

	// <------- Workspace Management Routes ------->
	// Create a new workspace
	workspace.Post("/", wh.CreateWorkspace)
	// List workspaces for a user
	workspace.Get("/", wh.ListWorkspaces)
	// Get a workspace by ID
	workspace.Get("/:id", wh.GetWorkspace)
	// Update a workspace by ID
	workspace.Put("/:id", wh.UpdateWorkspace)
	// Delete a workspace by ID
	workspace.Delete("/:id", wh.DeleteWorkspace)

	// <------- Workspace User Management Routes ------->
	// Exit a workspace by ID
	workspace.Post("/:id/exit", wh.ExitWorkspace)
	// Join a workspace by ID
	workspace.Post("/:id/join", wh.JoinWorkspace)
	// List users for a workspace
	workspace.Get("/:id/users", wh.ListWorkspaceUsers)
	// Add user to a workspace
	workspace.Post("/:id/users", wh.AddUserToWorkspace)
	// Remove user from a workspace
	workspace.Delete("/:id/users/:userId", wh.RemoveUserFromWorkspace)
	// Update user role in a workspace
	workspace.Put("/:id/users/:userId", wh.UpdateUserRoleInWorkspace)

	// <------- Workspace Competitor Management Routes ------->
	// Create competitor for a workspace
	workspace.Post("/:id/competitors", wh.CreateCompetitorForWorkspace)
	// Add page to a competitor
	workspace.Post("/:id/competitors/:competitorId/pages", wh.AddPageToCompetitor)
	// List workspace competitors
	workspace.Get("/:id/competitors", wh.ListWorkspaceCompetitors)
	// List page history
	workspace.Get("/:id/competitors/:competitorId/pages/:pageId/history", wh.ListPageHistory)
	// Remove page from a competitor
	workspace.Delete("/:id/competitors/:competitorId/pages/:pageId", wh.RemovePageFromCompetitor)
	// Remove competitor from a workspace
	workspace.Delete("/:id/competitors/:competitorId", wh.RemoveCompetitorFromWorkspace)
	// Update page in a competitor
	workspace.Put("/:id/competitors/:competitorId/pages/:pageId", wh.UpdatePageInCompetitor)

	// <------- Workflow Management Routes ------->
	// Workflow routes
	workflow := public.Group("/workflow")
	workflow.Post("/", h.WorkflowHandler.StartWorkflow)
	workflow.Delete("/", h.WorkflowHandler.StopWorkflow)
	workflow.Get("/", h.WorkflowHandler.GetWorkflow)
	workflow.Get("/list", h.WorkflowHandler.ListWorkflows)

	// <------- User Management Routes ------->
	// User routes
	uh := h.UserHandler
	user := public.Group("/user", middleware.ClerkAuthenticationMiddleware)
	// Delete Account
	user.Delete("/", uh.DeleteAccount)
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
