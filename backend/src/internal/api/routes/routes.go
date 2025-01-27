// ./src/internal/api/routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/handlers"
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/internal/service/notification"
	"github.com/wizenheimer/byrd/src/internal/service/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type HandlerContainer struct {
	AIHandler           *handlers.AIHandler
	ScreenshotHandler   *handlers.ScreenshotHandler
	UserHandler         *handlers.UserHandler
	WorkspaceHandler    *handlers.WorkspaceHandler
	WorkflowHandler     *handlers.WorkflowHandler
	ScheduleHandler     *handlers.ScheduleHandler
	NotificationHandler *handlers.NotificationHandler
}

func NewHandlerContainer(
	screenshotService screenshot.ScreenshotService,
	aiService ai.AIService,
	userService user.UserService,
	workspaceService workspace.WorkspaceService,
	workflowService workflow.WorkflowService,
	schedulerService scheduler.SchedulerService,
	notificationService notification.NotificationService,
	library template.TemplateLibrary,
	tx *transaction.TxManager,
	logger *logger.Logger,
) (*HandlerContainer, error) {
	notificationHandler, err := handlers.NewNotificationHandler(logger, notificationService, library)
	if err != nil {
		logger.Error("Couldn't create notification handler", zap.Error(err))
		return nil, err
	}

	hc := HandlerContainer{
		// Handlers for screenshot management
		ScreenshotHandler: handlers.NewScreenshotHandler(screenshotService, logger),
		// Handlers for AI management
		AIHandler: handlers.NewAIHandler(aiService, logger),
		// Handlers for user management
		UserHandler: handlers.NewUserHandler(userService, workspaceService, logger),
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
		// Handlers for notification management
		NotificationHandler: notificationHandler,
	}
	return &hc, nil
}

// SetupRoutes sets up the routes for the application
// This includes public and private routes
func SetupRoutes(
	app *fiber.App,
	handlers *HandlerContainer,
	m *middleware.AccessMiddleware,
	r *middleware.ResourceMiddleware,
) {

	setupPublicRoutes(app, handlers, m, r)

	if config.IsDevelopment() {
		setupPrivateRoutes(app, handlers, m)
	}
}

// setupPublicRoutes configures all public API endpoints
func setupPublicRoutes(
	app *fiber.App,
	h *HandlerContainer,
	m *middleware.AccessMiddleware,
	r *middleware.ResourceMiddleware,
) {
	// Base public API group
	public := app.Group("/api/public/v1", m.RequiresClerkToken)

	// Token validation
	public.Get("/token", h.UserHandler.ValidateClerkToken)

	// User management routes
	setupUserRoutes(public, h.UserHandler)

	// Workspace and related routes
	setupWorkspaceRoutes(public, h.WorkspaceHandler, m)

	// Member management routes
	setupMemberRoutes(public, h.WorkspaceHandler, m)

	// Competitor management routes
	setupCompetitorRoutes(public, h.WorkspaceHandler, m, r)

	// Page management routes
	setupPageRoutes(public, h.WorkspaceHandler, m, r)
}

// setupUserRoutes configures user management routes
func setupUserRoutes(
	router fiber.Router,
	handler *handlers.UserHandler,
) {
	// Delete the current user account
	router.Delete("/users", handler.DeleteCurrentUser)
	// Get the current user account
	router.Get("/users", handler.GetCurrentUser)
	// Create or update a user account
	router.Post("/users", handler.CreateOrUpdateUser)
}

// setupWorkspaceRoutes configures workspace and related resource management routes
func setupWorkspaceRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	m *middleware.AccessMiddleware,
) {
	// Create a new workspace for a user
	router.Post("/workspace",
		workspaceHandler.CreateWorkspaceForUser)

	// List all workspaces for a user
	router.Get("/workspace",
		workspaceHandler.ListWorkspacesForUser)

	// Get a workspace by ID
	router.Get("/workspace/:workspaceID",
		m.RequiresWorkspaceMember,
		workspaceHandler.GetWorkspaceByID)

	// Update a workspace by ID
	router.Put("/workspace/:workspaceID",
		m.RequiresWorkspaceMember,
		workspaceHandler.UpdateWorkspaceByID)

	// Delete a workspace by ID
	router.Delete("/workspace/:workspaceID",
		m.RequiresWorkspaceAdmin,
		workspaceHandler.DeleteWorkspaceByID)

	// Join a workspace by ID
	router.Post("/workspace/:workspaceID/join",
		m.RequiresPendingWorkspaceMember,
		workspaceHandler.JoinWorkspaceByID)

	// Exit a workspace by ID
	router.Post("/workspace/:workspaceID/exit",
		m.RequiresActiveOrPendingWorkspaceMembership,
		workspaceHandler.ExitWorkspaceByID)
}

func setupMemberRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	m *middleware.AccessMiddleware,
) {
	// List all users in a workspace
	router.Get("/workspace/:workspaceID/users",
		m.RequiresWorkspaceMember,
		workspaceHandler.ListUsersForWorkspace)

	// Invite a user to a workspace
	router.Post("/workspace/:workspaceID/users",
		m.RequiresWorkspaceMember,
		workspaceHandler.InviteUsersToWorkspace)

	// Update a user's role in a workspace
	router.Put("/workspace/:workspaceID/users/:userID",
		m.RequiresWorkspaceAdmin,
		workspaceHandler.UpdateUserRoleInWorkspace)

	// Remove a user from a workspace
	router.Delete("/workspace/:workspaceID/users/:userID",
		m.RequiresWorkspaceAdmin,
		workspaceHandler.RemoveUserFromWorkspace)
}

func setupCompetitorRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	m *middleware.AccessMiddleware,
	r *middleware.ResourceMiddleware,
) {
	// List all competitors in a workspace
	router.Get("/workspace/:workspaceID/competitors",
		m.RequiresWorkspaceMember,
		workspaceHandler.ListCompetitorsForWorkspace)

	// Add a competitor to a workspace
	router.Post("/workspace/:workspaceID/competitors",
		m.RequiresWorkspaceMember,
		workspaceHandler.CreateCompetitorForWorkspace)

	// Get a competitor in a workspace
	router.Get("/workspace/:workspaceID/competitors/:competitorID",
		m.RequiresWorkspaceMember,
		r.ValidateCompetitorResource,
		workspaceHandler.GetCompetitorForWorkspace)

	// Update a competitor in a workspace
	router.Put("/workspace/:workspaceID/competitors/:competitorID",
		m.RequiresWorkspaceMember,
		r.ValidateCompetitorResource,
		workspaceHandler.UpdateCompetitorForWorkspace)

	// Delete a competitor from a workspace
	router.Delete("/workspace/:workspaceID/competitors/:competitorID",
		m.RequiresWorkspaceMember,
		r.ValidateCompetitorResource,
		workspaceHandler.RemoveCompetitorFromWorkspace)
}

func setupPageRoutes(
	router fiber.Router,
	workspaceHandler *handlers.WorkspaceHandler,
	m *middleware.AccessMiddleware,
	r *middleware.ResourceMiddleware,
) {
	// List all pages in a workspace
	router.Get("/workspace/:workspaceID/competitors/:competitorID/pages",
		m.RequiresWorkspaceMember,
		r.ValidateCompetitorResource,
		workspaceHandler.ListPagesForCompetitor)

	// Add a page to a competitor
	router.Post("/workspace/:workspaceID/competitors/:competitorID/pages",
		m.RequiresWorkspaceMember,
		r.ValidateCompetitorResource,
		workspaceHandler.AddPagesToCompetitor)

	// Get a page in a competitor
	router.Get("/workspace/:workspaceID/competitors/:competitorID/pages/:pageID",
		m.RequiresWorkspaceMember,
		r.ValidatePageResource,
		workspaceHandler.GetPageForCompetitor)

	// Update a page in a competitor
	router.Put("/workspace/:workspaceID/competitors/:competitorID/pages/:pageID",
		m.RequiresWorkspaceMember,
		r.ValidatePageResource,
		workspaceHandler.UpdatePageForCompetitor)

	// Delete a page from a competitor
	router.Delete("/workspace/:workspaceID/competitors/:competitorID/pages/:pageID",
		m.RequiresWorkspaceMember,
		r.ValidatePageResource,
		workspaceHandler.RemovePageFromCompetitor)

	// List page history for a page
	router.Get("/workspace/:workspaceID/competitors/:competitorID/pages/:pageID/history",
		m.RequiresWorkspaceMember,
		r.ValidatePageResource,
		workspaceHandler.ListPageHistory)
}

// setupPrivateRoutes configures all private API endpoints
func setupPrivateRoutes(app *fiber.App, h *HandlerContainer, m *middleware.AccessMiddleware) {
	private := app.Group("/api/private/v1", m.RequiresPrivateToken)

	// Token validation
	private.Get("/token", h.UserHandler.ValidateManagementToken)

	// Workflow management routes
	setupWorkflowRoutes(private, h.WorkflowHandler)

	// Screenshot management routes
	setupScreenshotRoutes(private, h.ScreenshotHandler)

	// AI analysis routes
	setupAIRoutes(private, h.AIHandler)

	// Schedule management routes
	setupScheduleRoutes(private, h.ScheduleHandler)

	// Notification management routes
	setupNotificationRoutes(private, h.NotificationHandler)
}

// setupWorkflowRoutes configures workflow management endpoints
func setupWorkflowRoutes(router fiber.Router, handler *handlers.WorkflowHandler) {
	// Job management
	router.Post("/workflow/:workflowType/job", handler.StartWorkflow)
	router.Delete("/workflow/:workflowType/job/:jobID", handler.StopWorkflow)
	router.Get("/workflow/:workflowType/job/:jobID", handler.GetWorkflow)

	// Workflow monitoring
	router.Get("/workflow/checkpoint", handler.ListCheckpoint)
	router.Get("/workflow/history", handler.ListHistory)
}

// setupScreenshotRoutes configures screenshot management endpoints
func setupScreenshotRoutes(router fiber.Router, handler *handlers.ScreenshotHandler) {
	// Screenshot operations endpoints
	router.Post("/screenshot/refresh", handler.Refresh)
	router.Post("/screenshot/retrieve", handler.Retrieve)
}

// setupAIRoutes configures AI analysis endpoints
func setupAIRoutes(router fiber.Router, handler *handlers.AIHandler) {
	// Analysis endpoints
	router.Post("/ai/content", handler.AnalyzeContentDifferences)
	router.Post("/ai/visual", handler.AnalyzeVisualDifferences)
}

// setupScheduleRoutes configures schedule management endpoints
func setupScheduleRoutes(router fiber.Router, handler *handlers.ScheduleHandler) {
	// CRUD operations for schedules
	router.Post("/schedule", handler.CreateSchedule)
	router.Get("/schedule", handler.ListSchedules)
	router.Get("/schedule/:scheduleID", handler.GetSchedule)
	router.Delete("/schedule/:scheduleID", handler.DeleteSchedule)
	router.Put("/schedule/:scheduleID", handler.UpdateSchedule)
}

func setupNotificationRoutes(router fiber.Router, handler *handlers.NotificationHandler) {
	router.Post("/notification", handler.SendNotification)
}
