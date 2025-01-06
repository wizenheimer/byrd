package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/api/auth"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

type WorkspaceHandler struct {
	workspaceService svc.WorkspaceService
	logger           *logger.Logger
}

func NewWorkspaceHandler(
	workspaceService svc.WorkspaceService,
	logger *logger.Logger,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: workspaceService,
		logger:           logger,
	}
}

// CreateWorkspace creates a new workspace
func (wh *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	var req api.WorkspaceCreationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	workspace, e := wh.workspaceService.CreateWorkspace(c.Context(), clerkUser, req)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create workspace",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workspace)

}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	workspaces, e := wh.workspaceService.ListUserWorkspaces(c.Context(), clerkUser)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list workspaces",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list workspaces",
		})
	}

	return c.JSON(workspaces)
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	workspace, e := wh.workspaceService.GetWorkspace(c.Context(), workspaceID)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get workspace",
		})
	}

	return c.JSON(workspace)
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var req api.WorkspaceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if e := wh.workspaceService.UpdateWorkspace(c.Context(), workspaceID, req); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update workspace",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteWorkspace deletes a workspace by ID
func (wh *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	status, e := wh.workspaceService.DeleteWorkspace(c.Context(), workspaceID)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete workspace",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete workspace",
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if e := wh.workspaceService.LeaveWorkspace(c.Context(), clerkUser, workspaceID); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to leave workspace",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if e := wh.workspaceService.JoinWorkspace(c.Context(), clerkUser, workspaceID); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to join workspace",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	params := api.WorkspaceMembersListingParams{
		IncludeMembers: c.QueryBool("includeMembers", true),
		IncludeAdmins:  c.QueryBool("includeAdmins", true),
	}

	users, e := wh.workspaceService.ListWorkspaceMembers(c.Context(), workspaceID, params)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list workspace users",
		})
	}

	return c.JSON(users)
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var reqs []api.InviteUserToWorkspaceRequest
	if err := c.BodyParser(&reqs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	for index := range reqs {
		if err := utils.SetDefaultsAndValidate(&reqs[index]); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "InvalidRequest",
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	responses, e := wh.workspaceService.InviteUsersToWorkspace(c.Context(), clerkUser, workspaceID, reqs)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to invite users to workspace",
		})
	}

	return c.JSON(responses)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidUserID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	e := wh.workspaceService.RemoveUserFromWorkspace(c.Context(), workspaceID, userID)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove user from workspace",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidUserID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	var req struct {
		Role models.UserWorkspaceRole `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate role
	if !isValidRole(req.Role) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRole",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid role specified",
		})
	}

	e := wh.workspaceService.UpdateWorkspaceMemberRole(c.Context(), workspaceID, userID, req.Role)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update user role",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// Helper function to validate workspace roles
func isValidRole(role models.UserWorkspaceRole) bool {
	switch role {
	case models.UserRoleAdmin, models.UserRoleUser, models.UserRoleViewer:
		return true
	default:
		return false
	}
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var req api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	e := wh.workspaceService.CreateWorkspaceCompetitor(c.Context(), clerkUser, workspaceID, req)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "CompetitorCreationFailed",
			Code:    fiber.StatusBadRequest,
			Message: "Failed to create competitor",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	var pages []api.CreatePageRequest
	if err := c.BodyParser(&pages); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	for index := range pages {
		if err := utils.SetDefaultsAndValidate(&pages[index]); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "InvalidRequest",
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	createdPages, e := wh.workspaceService.AddPageToCompetitor(c.Context(), clerkUser, competitorID.String(), pages)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "PageCreationFailed",
			Code:    fiber.StatusBadRequest,
			Message: "Failed to add pages",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPages)
}

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	pageNumber := offset/limit + 1
	pageSize := limit

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	response, e := wh.workspaceService.ListWorkspaceCompetitors(c.Context(), clerkUser, workspaceID, params)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list competitors",
		})
	}

	return c.JSON(response)
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	pageNumber := c.QueryInt("pageNumber", 10)
	pageSize := c.QueryInt("pageSize", 0)

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	history, e := wh.workspaceService.ListWorkspacePageHistory(c.Context(), clerkUser, workspaceID, competitorID, pageID, params)
	if e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list page history",
		})
	}

	return c.JSON(history)
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if e := wh.workspaceService.RemovePageFromWorkspace(c.Context(), clerkUser, competitorID, pageID); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove page",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if e := wh.workspaceService.RemoveCompetitorFromWorkspace(c.Context(), clerkUser, workspaceID, competitorID); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove competitor",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	var req api.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if e := wh.workspaceService.UpdateCompetitorPage(c.Context(), competitorID, pageID, req); e != nil && e.HasErrors() {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update page",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
