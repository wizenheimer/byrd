// ./src/internal/api/handlers/user.go
package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type UserHandler struct {
	userService      user.UserService
	workspaceService workspace.WorkspaceService
	logger           *logger.Logger
}

func NewUserHandler(userService user.UserService, workspaceService workspace.WorkspaceService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService:      userService,
		workspaceService: workspaceService,
		logger:           logger.WithFields(map[string]interface{}{"module": "user_handler"}),
	}
}

// GetCurrentUser returns the current user
func (uh *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusBadRequest, "Couldn't get user email for user", err.Error())
	}

	ctx := c.Context()
	user, err := uh.userService.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not get user", err.Error())
	}

	// Synchronize the user with the database if the user is first time user
	firstTimeUser := user.Status != models.AccountStatusActive
	if firstTimeUser {
		user, err = uh.userService.ActivateUser(ctx, userEmail)
		if err != nil {
			return err
		}
	}

	return sendDataResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

// DeleteCurrentUser deletes the current user
func (uh *UserHandler) DeleteCurrentUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusBadRequest, "Couldn't get user email for user", err.Error())
	}

	// Find if the user has any workspaces
	membershipStatus := models.ActiveMember
	workspaces, _, err := uh.workspaceService.ListWorkspacesForUser(c.Context(), userEmail, &membershipStatus, nil, nil)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not list user workspaces", err.Error())
	}

	if len(workspaces) > 0 {
		return sendErrorResponse(c, uh.logger, fiber.StatusForbidden, "Exit the existing workspaces prior to deleting account", "User has active membership to workspaces")
	}

	if err := uh.userService.DeleteUserByEmail(c.Context(), userEmail); err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not delete user", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

// ListWorkspaces lists workspaces for a user
func (uh *UserHandler) ListWorkspacesForUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusUnauthorized, "User is not authorized to list the workspace", err.Error())
	}

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusBadRequest, "Couldn't get user email for user", err.Error())
	}

	user, err := uh.userService.GetUserByEmail(c.Context(), userEmail)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not get user", err.Error())
	}

	ctx := c.Context()
	// Synchronize the user with the database if the user is first time user
	firstTimeUser := user.Status != models.AccountStatusActive
	if firstTimeUser {
		user, err = uh.userService.ActivateUser(ctx, userEmail)
		if err != nil {
			return err
		}
	}

	membershipStatusString := strings.ToLower(c.Query("membership_status", ""))
	var membershipStatus *models.MembershipStatus
	switch membershipStatusString {
	case "active":
		m := models.ActiveMember
		membershipStatus = &m
	case "pending":
		m := models.PendingMember
		membershipStatus = &m
	default:
		membershipStatus = nil
	}

	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	workspaces, hasMore, err := uh.workspaceService.ListWorkspacesForUser(ctx, userEmail, membershipStatus, &limits, &offsets)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Workspace couldn't be listed for the user", err.Error())
	}

	listResponse := map[string]any{
		"workspaces": workspaces,
		"has_more":   hasMore,
	}

	includeUserProfile := strings.ToLower(c.Query("include_profile", "false")) == "true"
	if includeUserProfile {
		listResponse["user"] = user
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspaces successfully", listResponse)
}
