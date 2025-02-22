// ./src/internal/api/handlers/workspace_user.go
package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListUsersForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	roleFilterString := strings.ToLower(c.Query("role", ""))

	var roleFilter *models.WorkspaceRole
	switch roleFilterString {
	case "admin":
		adminRole := models.RoleAdmin
		roleFilter = &adminRole
	case "user":
		userRole := models.RoleUser
		roleFilter = &userRole
	default:
		// Note: This would result in all roles being returned
		roleFilter = nil
	}

	ctx := c.Context()
	users, hasMore, err := wh.workspaceService.ListWorkspaceMembers(ctx, workspaceID, &limits, &offsets, roleFilter)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list workspace users", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace users successfully", map[string]any{
		"users":    users,
		"has_more": hasMore,
	})
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) InviteUsersToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var users api.AddUsersToWorkspaceRequest
	if err := c.BodyParser(&users); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Unprocessable request body", err.Error())
	}

	err = utils.SetDefaultsAndValidate(&users)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}
	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Couldn't locate user email", err.Error())
	}

	ctx := c.Context()

	workspaceUsers, err := wh.workspaceService.AddUsersToWorkspace(ctx, userEmail, workspaceID, users.Emails)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not invite users to workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Invited users to workspace successfully", workspaceUsers)
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	var req api.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()

	if err := wh.workspaceService.UpdateWorkspaceMemberRole(ctx, workspaceID, userID, req.Role); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update user role in workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated user role in workspace successfully", map[string]any{
		"role": req.Role,
	})
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.RemoveUserFromWorkspace(ctx, workspaceID, userID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not remove user from workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed user from workspace successfully", map[string]any{
		"user_id":      userID,
		"workspace_id": workspaceID,
	})
}
