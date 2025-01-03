package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/auth"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	return nil
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
