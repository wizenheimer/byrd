package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/api/auth"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
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
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidRequest", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidRequest", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	workspace, e := wh.workspaceService.CreateWorkspace(c.Context(), clerkUser, req)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create workspace", e)
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created workspace successfully", workspace)

}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	workspaces, e := wh.workspaceService.ListUserWorkspaces(c.Context(), clerkUser)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspaces", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspaces successfully", workspaces)
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	workspace, e := wh.workspaceService.GetWorkspace(c.Context(), workspaceID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Fetched workspace successfully", workspace)
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req api.WorkspaceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if e := wh.workspaceService.UpdateWorkspace(c.Context(), workspaceID, req); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated workspace successfully",
		map[string]any{
			"workspaceId":  workspaceID,
			"billingEmail": req.BillingEmail,
			"name":         req.Name,
		})
}

// DeleteWorkspace deletes a workspace by ID
func (wh *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	status, e := wh.workspaceService.DeleteWorkspace(c.Context(), workspaceID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not delete workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Deleted workspace successfully", map[string]any{
		"workspaceId": workspaceID,
		"status":      status,
	})
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	if e := wh.workspaceService.LeaveWorkspace(c.Context(), clerkUser, workspaceID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not exit workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Exited workspace successfully", map[string]any{
		"workspaceId": workspaceID,
	})
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	if e := wh.workspaceService.JoinWorkspace(c.Context(), clerkUser, workspaceID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not join workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Joined workspace successfully", map[string]any{
		"workspaceId": workspaceID,
	})
}

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	params := api.WorkspaceMembersListingParams{
		IncludeMembers: c.QueryBool("includeMembers", true),
		IncludeAdmins:  c.QueryBool("includeAdmins", true),
	}

	users, e := wh.workspaceService.ListWorkspaceMembers(c.Context(), workspaceID, params)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace users", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace users successfully", users)
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var reqs []api.InviteUserToWorkspaceRequest
	if err := c.BodyParser(&reqs); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	for index := range reqs {
		if err := utils.SetDefaultsAndValidate(&reqs[index]); err != nil {
			return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
		}
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	responses, e := wh.workspaceService.InviteUsersToWorkspace(c.Context(), clerkUser, workspaceID, reqs)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not invite users to workspace", e)
	}

	return sendDataResponse(c, fiber.StatusCreated, "Invited users to workspace successfully", responses)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	e := wh.workspaceService.RemoveUserFromWorkspace(c.Context(), workspaceID, userID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove user from workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed user from workspace successfully", map[string]any{
		"userId":      userID,
		"workspaceId": workspaceID,
	})
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	var req api.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	e := wh.workspaceService.UpdateWorkspaceMemberRole(c.Context(), workspaceID, userID, req.Role)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update user role in workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated user role in workspace successfully", map[string]any{
		"role": req.Role,
	})
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	e := wh.workspaceService.CreateWorkspaceCompetitor(c.Context(), clerkUser, workspaceID, req)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create competitor", e)
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created competitor successfully", req)
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var pages []api.CreatePageRequest
	if err := c.BodyParser(&pages); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	for index := range pages {
		if err := utils.SetDefaultsAndValidate(&pages[index]); err != nil {
			return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
		}
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	createdPages, e := wh.workspaceService.AddPageToCompetitor(c.Context(), clerkUser, competitorID.String(), pages)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not add page to competitor", e)
	}

	return sendDataResponse(c, fiber.StatusCreated, "Added page to competitor successfully", createdPages)
}

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
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
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace competitors", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace competitors successfully", response)
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid page ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	pageNumber := c.QueryInt("pageNumber", 10)
	pageSize := c.QueryInt("pageSize", 0)

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	history, e := wh.workspaceService.ListWorkspacePageHistory(c.Context(), clerkUser, workspaceID, competitorID, pageID, params)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list page history", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed page history successfully", history)
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	if e := wh.workspaceService.RemovePageFromWorkspace(c.Context(), clerkUser, competitorID, pageID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove page from competitor", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed page from competitor successfully", nil)
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidWorkspaceID", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	if e := wh.workspaceService.RemoveCompetitorFromWorkspace(c.Context(), clerkUser, workspaceID, competitorID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove competitor from workspace", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed competitor from workspace successfully", nil)
}

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	var req api.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if e := wh.workspaceService.UpdateCompetitorPage(c.Context(), competitorID, pageID, req); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update page in competitor", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated page in competitor successfully", nil)
}
