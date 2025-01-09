package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/auth"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type UserHandler struct {
	userService svc.UserService
	logger      *logger.Logger
}

func NewUserHandler(userService svc.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (uh *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	e := uh.userService.DeleteUser(c.Context(), clerkUser)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not delete user", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

func (uh *UserHandler) ValidateToken(c *fiber.Ctx) error {
	clerkClaims, err := auth.GetClerkClaimsFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user claims from context", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Couldn't get user from context", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "User is authenticated", map[string]interface{}{
		"user":   clerkUser,
		"claims": clerkClaims,
	})
}
