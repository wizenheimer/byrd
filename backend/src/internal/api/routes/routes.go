// ./src/internal/api/routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/handlers"
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/internal/service/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type HandlerContainer struct {
	AIHandler         *handlers.AIHandler
	ScreenshotHandler *handlers.ScreenshotHandler
	UserHandler       *handlers.UserHandler
	WorkspaceHandler  *handlers.WorkspaceHandler
	WorkflowHandler   *handlers.WorkflowHandler
	ScheduleHandler   *handlers.ScheduleHandler
}

func NewHandlerContainer(
	screenshotService screenshot.ScreenshotService,
	aiService ai.AIService,
	userService user.UserService,
	workspaceService workspace.WorkspaceService,
	workflowService workflow.WorkflowService,
	schedulerService scheduler.SchedulerService,
	tx *transaction.TxManager,
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
			tx,
			logger,
		),
		// Handlers for workflow management
		WorkflowHandler: handlers.NewWorkflowHandler(workflowService, logger),
		// Handlers for schedule management
		ScheduleHandler: handlers.NewScheduleHandler(schedulerService, logger),
	}
}

// SetupRoutes sets up the routes for the application
// This includes public and private routes
func SetupRoutes(app *fiber.App, handlers *HandlerContainer,
	pathMiddleware *middleware.WorkspacePathValidationMiddleware,
	authorizationMiddleware *middleware.AuthorizationMiddleware,
	authMiddleware *middleware.AuthenticatedMiddleware) {
	// authMiddleware := middleware.NewAuthenticatedMiddleware(logger)
	// authorizationMiddleware := middleware.NewAuthorizationMiddleware(ws, logger)
	// pathMiddleware := middleware.NewWorkspacePathValidationMiddleware(ws, logger)

	setupPublicRoutes(app, handlers, authMiddleware, authorizationMiddleware, pathMiddleware)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers, authMiddleware)
	}
}

// setupPublicRoutes sets up public routes for the application
func setupPublicRoutes(app *fiber.App, h *HandlerContainer, authMiddleware *middleware.AuthenticatedMiddleware, authorization *middleware.AuthorizationMiddleware, pathMiddleware *middleware.WorkspacePathValidationMiddleware) {
	// Public routes for production and development
	public := app.Group("/api/public/v1")

	// Workspace routes
	wh := h.WorkspaceHandler

	// -------------------------------------------------
	// Base workspace routes require valid authentication token
	workspace := public.Group("/workspace", authMiddleware.AuthenticationMiddleware)

	// Base workspace endpoints (no ID needed)
	workspace.Post("/", wh.CreateWorkspace)
	workspace.Get("/", wh.ListWorkspaces)

	// -------------------------------------------------
	// All routes that need workspace validation
	workspaceBase := workspace.Group("/:workspaceID",
		pathMiddleware.ValidateWorkspacePath)

	// Workspace admin routes
	workspaceAdmin := workspaceBase.Group("",
		authorization.RequireWorkspaceAdmin)

	workspaceAdmin.Put("/", wh.UpdateWorkspace)
	workspaceAdmin.Delete("/", wh.DeleteWorkspace)
	workspaceAdmin.Delete("/users/:userId", wh.RemoveUserFromWorkspace)
	workspaceAdmin.Put("/users/:userId", wh.UpdateUserRoleInWorkspace)

	// -------------------------------------------------
	// Workspace member routes
	workspaceMember := workspaceBase.Group("",
		authorization.RequireWorkspaceMembership)

	workspaceMember.Get("/", wh.GetWorkspace)
	workspaceMember.Post("/exit", wh.ExitWorkspace)
	workspaceMember.Post("/join", wh.JoinWorkspace)
	workspaceMember.Get("/users", wh.ListWorkspaceUsers)
	workspaceMember.Post("/users", wh.AddUserToWorkspace)
	workspaceMember.Post("/competitors", wh.CreateCompetitorForWorkspace)
	workspaceMember.Get("/competitors", wh.ListWorkspaceCompetitors)

	// -------------------------------------------------
	// Competitor management routes
	competitorManagement := workspaceMember.Group("/competitors/:competitorID",
		pathMiddleware.ValidateCompetitorPath)

	competitorManagement.Post("/pages", wh.AddPageToCompetitor)
  competitorManagement.Get("/pages", wh.ListPagesForCompetitor)
	competitorManagement.Delete("/", wh.RemoveCompetitorFromWorkspace)

	// -------------------------------------------------
	// Page management routes
	pageManagement := competitorManagement.Group("/pages/:pageID",
		pathMiddleware.ValidatePagePath)

	pageManagement.Get("/history", wh.ListPageHistory)
	pageManagement.Delete("/", wh.RemovePageFromCompetitor)
	pageManagement.Put("/", wh.UpdatePageInCompetitor)

	// -------------------------------------------------
	// User routes
	uh := h.UserHandler
	user := public.Group("/user", authMiddleware.AuthenticationMiddleware)
	user.Delete("/", uh.DeleteAccount)
	user.Get("/", uh.Sync)
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
	workflow.Post("/:workflowType/job", h.WorkflowHandler.StartWorkflow)
	// Stop a workflow
	workflow.Delete("/:workflowType/job/:jobID", h.WorkflowHandler.StopWorkflow)
	// List workflows
	workflow.Get("/checkpoint", h.WorkflowHandler.ListCheckpoint)
	workflow.Get("/history", h.WorkflowHandler.ListHistory)
	// Get a workflow
	workflow.Get("/:workflowType/job/:jobID", h.WorkflowHandler.GetWorkflow)

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

	// -------- Schedule Management Routes --------
	// schedule routes
	schedule := private.Group("/schedule")

	// Schedule a new workflow
	schedule.Post("/", h.ScheduleHandler.CreateSchedule)
	// List all scheduled workflows
	schedule.Get("/", h.ScheduleHandler.ListSchedules)
	// Get a scheduled workflow
	schedule.Get("/:scheduleID", h.ScheduleHandler.GetSchedule)
	// Delete a scheduled workflow
	schedule.Delete("/:scheduleID", h.ScheduleHandler.DeleteSchedule)
	// Update a scheduled workflow
	schedule.Put("/:scheduleID", h.ScheduleHandler.UpdateSchedule)
}
