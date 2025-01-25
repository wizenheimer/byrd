// ./src/internal/api/handlers/user.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type UserHandler struct {
	userService user.UserService
	logger      *logger.Logger
}

func NewUserHandler(userService user.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger.WithFields(map[string]interface{}{"module": "user_handler"}),
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

	return sendDataResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

// DeleteCurrentUser deletes the current user
func (uh *UserHandler) DeleteCurrentUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, uh.logger, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
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

	return sendDataResponse(c, fiber.StatusOK, "User created or updated successfully", user)
}
