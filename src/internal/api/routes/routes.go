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
func SetupRoutes(app *fiber.App, handlers *HandlerContainer, logger *logger.Logger) {
	authMiddleware := middleware.NewAuthenticatedMiddleware(logger)

	setupPublicRoutes(app, handlers, authMiddleware)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers, authMiddleware)
	}

}

// setupPublicRoutes sets up the public routes for the application
func setupPublicRoutes(app *fiber.App, h *HandlerContainer, authMiddleware *middleware.AuthenticatedMiddleware) {
	// Public routes for production and development
	public := app.Group("/api/public/v1")

	// Workspace routes
	wh := h.WorkspaceHandler
	workspace := public.Group("/workspace", authMiddleware.AuthenticationMiddleware)

	// <------- Workspace Management Routes ------->
	// Create a new workspace
	workspace.Post("/", wh.CreateWorkspace)
	// List workspaces for a user
	workspace.Get("/", wh.ListWorkspaces)
	// Get a workspace by ID
	workspace.Get("/:workspaceID", wh.GetWorkspace)
	// Update a workspace by ID
	workspace.Put("/:workspaceID", wh.UpdateWorkspace)
	// Delete a workspace by ID
	workspace.Delete("/:workspaceID", wh.DeleteWorkspace)

	// <------- Workspace User Management Routes ------->
	// Exit a workspace by ID
	workspace.Post("/:workspaceID/exit", wh.ExitWorkspace)
	// Join a workspace by ID
	workspace.Post("/:workspaceID/join", wh.JoinWorkspace)
	// List users for a workspace
	workspace.Get("/:workspaceID/users", wh.ListWorkspaceUsers)
	// Add user to a workspace
	workspace.Post("/:workspaceID/users", wh.AddUserToWorkspace)
	// Remove user from a workspace
	workspace.Delete("/:workspaceID/users/:userId", wh.RemoveUserFromWorkspace)
	// Update user role in a workspace
	workspace.Put("/:workspaceID/users/:userId", wh.UpdateUserRoleInWorkspace)

	// <------- Workspace Competitor Management Routes ------->
	// Create competitor for a workspace
	workspace.Post("/:workspaceID/competitors", wh.CreateCompetitorForWorkspace)
	// Add page to a competitor
	workspace.Post("/:workspaceID/competitors/:competitorID/pages", wh.AddPageToCompetitor)
	// List workspace competitors
	workspace.Get("/:workspaceID/competitors", wh.ListWorkspaceCompetitors)
	// List page history
	workspace.Get("/:workspaceID/competitors/:competitorID/pages/:pageID/history", wh.ListPageHistory)
	// Remove page from a competitor
	workspace.Delete("/:workspaceID/competitors/:competitorID/pages/:pageID", wh.RemovePageFromCompetitor)
	// Remove competitor from a workspace
	workspace.Delete("/:workspaceID/competitors/:competitorID", wh.RemoveCompetitorFromWorkspace)
	// Update page in a competitor
	workspace.Put("/:workspaceID/competitors/:competitorID/pages/:pageID", wh.UpdatePageInCompetitor)

	// <------- User Management Routes ------->
	// User routes
	uh := h.UserHandler
	user := public.Group("/user", authMiddleware.AuthenticationMiddleware)
	// Delete Account
	user.Delete("/", uh.DeleteAccount)
}

// setupPrivateRoutes sets up the private routes for the application
func setupPrivateRoutes(app *fiber.App, h *HandlerContainer, authMiddleware *middleware.AuthenticatedMiddleware) {
	// Private routes for development
	private := app.Group("/api/private/v1")

	// <------- Auth validation routes ------->
	// User routes
	uh := h.UserHandler
	user := private.Group("/auth", authMiddleware.AuthenticationMiddleware)
	// Validate token
	user.Get("/validate", uh.ValidateToken)

	// <------- Workflow Management Routes ------->
	// Workflow routes
	workflow := private.Group("/workflow")
	// Start a new workflow
	workflow.Post("/", h.WorkflowHandler.StartWorkflow)
	// Stop a workflow
	workflow.Delete("/", h.WorkflowHandler.StopWorkflow)
	// List workflows
	workflow.Get("/", h.WorkflowHandler.ListWorkflows)
	// Get a workflow
	workflow.Get("/:id", h.WorkflowHandler.GetWorkflow)

	// <------- Screenshot Management Routes ------->
	// Screenshot routes
	screenshot := private.Group("/screenshot")
	// Create a new screenshot
	screenshot.Post("/", h.ScreenshotHandler.CreateScreenshot)
	// List screenshots
	screenshot.Get("/", h.ScreenshotHandler.ListScreenshots)
	// Get existing screenshot image
	screenshot.Get("/image", h.ScreenshotHandler.GetScreenshotImage)
	// Get existing screenshot content
	screenshot.Get("/content", h.ScreenshotHandler.GetScreenshotContent)

	// AI routes
	ai := private.Group("/ai")
	// Analyze content differences
	ai.Post("/content", h.AIHandler.AnalyzeContentDifferences)
	// Analyze visual differences
	ai.Post("/visual", h.AIHandler.AnalyzeVisualDifferences)

}
