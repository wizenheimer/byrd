package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/auth"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't get user from context",
			"error":   err.Error(),
		})
	}

	err = uh.userService.DeleteUser(c.Context(), clerkUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

func (uh *UserHandler) ValidateToken(c *fiber.Ctx) error {
	clerkClaims, err := auth.GetClerkClaimsFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't get claims from context",
			"error":   err.Error(),
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't get user from context",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User is authenticated",
		"data": map[string]interface{}{
			"user":   clerkUser,
			"claims": clerkClaims,
		},
	})
}
