package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	"github.com/wizenheimer/iris/src/internal/api/middleware"
	"github.com/wizenheimer/iris/src/internal/config"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
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
	screenshotService svc.ScreenshotService,
	aiService svc.AIService,
	userService svc.UserService,
	workspaceService svc.WorkspaceService,
	logger *logger.Logger,
) *HandlerContainer {
	return &HandlerContainer{
		// Handlers for screenshot management
		ScreenshotHandler: handlers.NewScreenshotHandler(screenshotService, logger),
		// Handlers for AI management
		AIHandler: handlers.NewAIHandler(aiService, logger),
		// Handlers for user management
		UserHandler: handlers.NewUserHandler(userService, logger),
		// Handlers for workspace management
		WorkspaceHandler: handlers.NewWorkspaceHandler(
			workspaceService,
			logger,
		),
		// Handlers for workflow management
		WorkflowHandler: handlers.NewWorkflowHandler(nil, logger),
	}
}

// SetupRoutes sets up the routes for the application
// This includes public and private routes
func SetupRoutes(app *fiber.App, handlers *HandlerContainer, ws svc.WorkspaceService, logger *logger.Logger) {
	authMiddleware := middleware.NewAuthenticatedMiddleware(logger)
    authorizationMiddleware := middleware.NewAuthorizationMiddleware(ws, logger)
    pathMiddleware := middleware.NewWorkspacePathValidationMiddleware(ws, logger)

	setupPublicRoutes(app, handlers, authMiddleware, authorizationMiddleware, pathMiddleware)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers, authMiddleware)
	}

}

// setupPublicRoutes sets up the public routes for the application
func setupPublicRoutes(app *fiber.App, h *HandlerContainer, authMiddleware *middleware.AuthenticatedMiddleware, authorization *middleware.AuthorizationMiddleware, pathMiddleware *middleware.WorkspacePathValidationMiddleware) {
	// Public routes for production and development
	public := app.Group("/api/public/v1")

	// Workspace routes
	wh := h.WorkspaceHandler

    // -------------------------------------------------
	// Base workspace routes require
    // valid authentication token
	workspace := public.Group("/workspace", authMiddleware.AuthenticationMiddleware)

	// Create a new workspace
	workspace.Post("/", wh.CreateWorkspace)

	// List workspaces for a user
	workspace.Get("/", wh.ListWorkspaces)

    // -------------------------------------------------
	// workspace admin are routes that require
	// valid workspace path
    // valid workspace admin role
    // inherits authentication middleware
	workspaceAdmin := workspace.Group("", pathMiddleware.ValidateWorkspacePath, authorization.RequireWorkspaceAdmin)

	// Update a workspace by ID
	workspaceAdmin.Put("/:workspaceID/", wh.UpdateWorkspace)

	// Delete a workspace by ID
	workspaceAdmin.Delete("/:workspaceID/", wh.DeleteWorkspace)

	// Remove user from a workspace
	workspaceAdmin.Delete("/:workspaceID/users/:userId", wh.RemoveUserFromWorkspace)

	// Update user role in a workspace
	workspaceAdmin.Put("/:workspaceID/users/:userId", wh.UpdateUserRoleInWorkspace)

    // -------------------------------------------------
	// workspaceMember routes are routes that require
	// valid workspace path
    // valid workspace membership role
    // inherits authentication middleware
	workspaceMember := workspace.Group("", pathMiddleware.ValidateWorkspacePath, authorization.RequireWorkspaceMembership)

	// Get a workspace by ID
	workspaceMember.Get("/:workspaceID/", wh.GetWorkspace)

	// Exit a workspace by ID
	workspaceMember.Post("/:workspaceID/exit", wh.ExitWorkspace)

	// Join a workspace by ID
	workspaceMember.Post("/:workspaceID/join", wh.JoinWorkspace)

	// List users for a workspace
	workspaceMember.Get("/:workspaceID/users", wh.ListWorkspaceUsers)

	// Add user to a workspace
	workspaceMember.Post("/:workspaceID/users", wh.AddUserToWorkspace)

	// Create competitor for a workspace
	workspaceMember.Post("/:workspaceID/competitors", wh.CreateCompetitorForWorkspace)

	// List workspace competitors
	workspaceMember.Get("/:workspaceID/competitors", wh.ListWorkspaceCompetitors)

    // -------------------------------------------------
	// Competitor routes are the routes that require
	// a valid workspace path
    // a valid competitor path
    // inheritance of workspace membership role
	competitorManagement := workspaceMember.Group("", pathMiddleware.ValidateCompetitorPath)

	// Add page to a competitor
	competitorManagement.Post("/:workspaceID/competitors/:competitorID/pages", wh.AddPageToCompetitor)

	// Remove competitor from a workspace
	competitorManagement.Delete("/:workspaceID/competitors/:competitorID", wh.RemoveCompetitorFromWorkspace)

    // -------------------------------------------------
    // Page Management routes are the routes that require
    // a valid page path
    // inherits a valid workspace path
    // inherits a valid competitor path
    pageManagement := competitorManagement.Group("", pathMiddleware.ValidatePagePath)

	// List page history
	pageManagement.Get("/:workspaceID/competitors/:competitorID/pages/:pageID/history", wh.ListPageHistory)

	// Remove page from a competitor
	pageManagement.Delete("/:workspaceID/competitors/:competitorID/pages/:pageID", wh.RemovePageFromCompetitor)

	// Update page in a competitor
	pageManagement.Put("/:workspaceID/competitors/:competitorID/pages/:pageID", wh.UpdatePageInCompetitor)


    // -------------------------------------------------
    // User routes are the routes that require
    // a valid user path
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
