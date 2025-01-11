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
		logger:      logger,
	}
}

func (uh *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	if err := uh.userService.DeleteUser(c.Context(), clerkUser); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not delete user", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

func (uh *UserHandler) ValidateToken(c *fiber.Ctx) error {
	clerkClaims, err := getClerkClaimsFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user claims from context", err.Error())
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "User is authenticated", map[string]interface{}{
		"user":   clerkUser,
		"claims": clerkClaims,
	})
}

func (uh *UserHandler) Sync(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	if err := uh.userService.SyncUser(c.Context(), clerkUser); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not sync user", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "User is synchronized", nil)
}
