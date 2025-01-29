// ./src/internal/api/handlers/user.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
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

	ctx := c.Context()
	user, err := uh.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not get user", err.Error())
	}

	// Synchronize the user with the database if the user is first time user
	firstTimeUser := user.ClerkID == nil || user.Status != models.AccountStatusActive
	if firstTimeUser {
		if err := uh.userService.ActivateUser(ctx, user.ID, clerkUser); err != nil {
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

	// Find if the user has any workspaces
	workspaces, err := uh.workspaceService.ListUserWorkspaces(c.Context(), clerkUser, models.ActiveMember)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not list user workspaces", err.Error())
	}

	if len(workspaces) > 0 {
		return sendErrorResponse(c, uh.logger, fiber.StatusForbidden, "Exit the existing workspaces prior to deleting account", "User has active membership to workspaces")
	}

	if err := uh.userService.DeleteUser(c.Context(), clerkUser); err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not delete user", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

// CreateOrUpdateUser creates or updates a user
func (uh *UserHandler) CreateOrUpdateUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	ctx := c.Context()
	user, err := uh.userService.GetOrCreateUser(ctx, clerkUser)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusInternalServerError, "Could not create or update user", err.Error())
	}

	// Synchronize the user with the database if the user is first time user
	firstTimeUser := user.ClerkID == nil || user.Status != models.AccountStatusActive
	if firstTimeUser {
		if err := uh.userService.ActivateUser(ctx, user.ID, clerkUser); err != nil {
			return err
		}
	}

	return sendDataResponse(c, fiber.StatusOK, "User created or updated successfully", user)
}
