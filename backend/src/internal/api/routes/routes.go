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

	setupPublicRoutes(app, handlers, authMiddleware, authorizationMiddleware, pathMiddleware)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers, authMiddleware)
	}
}

// setupPublicRoutes configures all public API endpoints
func setupPublicRoutes(
	app *fiber.App,
	h *HandlerContainer,
	authMiddleware *middleware.AuthenticatedMiddleware,
	authorization *middleware.AuthorizationMiddleware,
	pathMiddleware *middleware.WorkspacePathValidationMiddleware,
) {
	// Base public API group
	public := app.Group("/api/public/v1")

	// Authentication routes
	setupAuthRoutes(public, h.UserHandler, authMiddleware)

	// User management routes
	setupUserRoutes(public, h.UserHandler, authMiddleware)

	// Workspace and related routes
	setupWorkspaceRoutes(public, h.WorkspaceHandler, h.UserHandler,
		authMiddleware, authorization, pathMiddleware)
}

// setupAuthRoutes configures authentication-related routes
func setupAuthRoutes(
	router fiber.Router,
	handler *handlers.UserHandler,
	authMiddleware *middleware.AuthenticatedMiddleware,
) {
	router.Get("/auth",
		authMiddleware.AuthenticationMiddleware,
		handler.ValidateClerkToken,
	)
}

// setupUserRoutes configures user management routes
func setupUserRoutes(
	router fiber.Router,
	handler *handlers.UserHandler,
	authMiddleware *middleware.AuthenticatedMiddleware,
) {
	user := router.Group("/user",
		authMiddleware.AuthenticationMiddleware,
	)
	user.Delete("/", handler.DeleteAccount)
}

// setupWorkspaceRoutes configures workspace and related resource management routes
func setupWorkspaceRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthenticatedMiddleware,
	authorization *middleware.AuthorizationMiddleware,
	pathMiddleware *middleware.WorkspacePathValidationMiddleware,
) {
	// Base workspace group with authentication
	workspace := router.Group("/workspace",
		authMiddleware.AuthenticationMiddleware,
	)

	// Basic workspace operations
	workspace.Post("/", workspaceHandler.CreateWorkspace)
	workspace.Get("/", workspaceHandler.ListWorkspaces)

	// Workspace-specific routes
	workspaceBase := workspace.Group("/:workspaceID",
		pathMiddleware.ValidateWorkspacePath,
	)

	setupWorkspaceAdminRoutes(workspaceBase, workspaceHandler, authorization)
	setupWorkspaceMemberRoutes(workspaceBase, workspaceHandler, userHandler, authorization)
	setupCompetitorRoutes(workspaceBase, workspaceHandler, authorization, pathMiddleware)
}

// setupWorkspaceAdminRoutes configures admin-only workspace management routes
func setupWorkspaceAdminRoutes(
	router fiber.Router,
	handler *handlers.WorkspaceHandler,
	authorization *middleware.AuthorizationMiddleware,
) {
	admin := router.Group("",
		authorization.RequireWorkspaceAdmin,
	)

	admin.Put("/", handler.UpdateWorkspace)
	admin.Delete("/", handler.DeleteWorkspace)
	admin.Delete("/users/:userId", handler.RemoveUserFromWorkspace)
	admin.Put("/users/:userId", handler.UpdateUserRoleInWorkspace)
}

// setupWorkspaceMemberRoutes configures member-specific workspace routes
func setupWorkspaceMemberRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	userHandler *handlers.UserHandler,
	authorization *middleware.AuthorizationMiddleware,
) {
	// Pending member routes
	router.Post("/join", authorization.RequirePendingWorkspaceMembership, userHandler.Sync, workspaceHandler.JoinWorkspace)

	// Active member routes
	member := router.Group("",
		authorization.RequireWorkspaceMembership,
	)
	member.Post("/exit", userHandler.Sync, workspaceHandler.ExitWorkspace)
	member.Get("/", workspaceHandler.GetWorkspace)

	// Extended member privileges
	activeMember := router.Group("",
		authorization.RequireActiveWorkspaceMembership,
	)
	activeMember.Get("/users", workspaceHandler.ListWorkspaceUsers)
	activeMember.Post("/users", workspaceHandler.AddUserToWorkspace)
	activeMember.Post("/competitors", workspaceHandler.CreateCompetitorForWorkspace)
	activeMember.Get("/competitors", workspaceHandler.ListWorkspaceCompetitors)
}

// setupCompetitorRoutes configures competitor and page management routes
func setupCompetitorRoutes(
	router fiber.Router,
	handler *handlers.WorkspaceHandler,
	authorization *middleware.AuthorizationMiddleware,
	pathMiddleware *middleware.WorkspacePathValidationMiddleware,
) {
	competitor := router.Group("/competitors/:competitorID",
		authorization.RequireActiveWorkspaceMembership,
		pathMiddleware.ValidateCompetitorPath,
	)

	// Competitor management
	competitor.Post("/pages", handler.AddPageToCompetitor)
	competitor.Get("/pages", handler.ListPagesForCompetitor)
	competitor.Delete("/", handler.RemoveCompetitorFromWorkspace)

	// Page management
	page := competitor.Group("/pages/:pageID",
		pathMiddleware.ValidatePagePath,
	)
	page.Get("/history", handler.ListPageHistory)
	page.Delete("/", handler.RemovePageFromCompetitor)
	page.Put("/", handler.UpdatePageInCompetitor)
}

// setupPrivateRoutes configures all private API endpoints
func setupPrivateRoutes(app *fiber.App, h *HandlerContainer, authMiddleware *middleware.AuthenticatedMiddleware) {
	private := app.Group("/api/private/v1",
		authMiddleware.PrivateRouteAuthenticationMiddleware)

	// Management token validation
	private.Get("/validate", h.UserHandler.ValidateManagementToken)

	setupWorkflowRoutes(private, h.WorkflowHandler)
	setupScreenshotRoutes(private, h.ScreenshotHandler)
	setupAIRoutes(private, h.AIHandler)
	setupScheduleRoutes(private, h.ScheduleHandler)
}

// setupWorkflowRoutes configures workflow management endpoints
func setupWorkflowRoutes(router fiber.Router, handler *handlers.WorkflowHandler) {
	workflow := router.Group("/workflow")

	// Job management
	workflow.Post("/:workflowType/job", handler.StartWorkflow)
	workflow.Delete("/:workflowType/job/:jobID", handler.StopWorkflow)
	workflow.Get("/:workflowType/job/:jobID", handler.GetWorkflow)

	// Workflow monitoring
	workflow.Get("/checkpoint", handler.ListCheckpoint)
	workflow.Get("/history", handler.ListHistory)
}

// setupScreenshotRoutes configures screenshot management endpoints
func setupScreenshotRoutes(router fiber.Router, handler *handlers.ScreenshotHandler) {
	screenshot := router.Group("/screenshot")

	// Screenshot operations
	screenshot.Post("/", handler.CreateScreenshot)
	screenshot.Get("/", handler.ListScreenshots)

	// Screenshot content retrieval
	screenshot.Get("/image", handler.GetScreenshotImage)
	screenshot.Get("/content", handler.GetScreenshotContent)
}

// setupAIRoutes configures AI analysis endpoints
func setupAIRoutes(router fiber.Router, handler *handlers.AIHandler) {
	ai := router.Group("/ai")

	// Analysis endpoints
	ai.Post("/content", handler.AnalyzeContentDifferences)
	ai.Post("/visual", handler.AnalyzeVisualDifferences)
}

// setupScheduleRoutes configures schedule management endpoints
func setupScheduleRoutes(router fiber.Router, handler *handlers.ScheduleHandler) {
	schedule := router.Group("/schedule")

	// CRUD operations for schedules
	schedule.Post("/", handler.CreateSchedule)
	schedule.Get("/", handler.ListSchedules)
	schedule.Get("/:scheduleID", handler.GetSchedule)
	schedule.Delete("/:scheduleID", handler.DeleteSchedule)
	schedule.Put("/:scheduleID", handler.UpdateSchedule)
}
