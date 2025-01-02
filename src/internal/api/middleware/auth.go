package middleware

import (
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/api/auth"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	// Validate token
	// Implementation here

	return c.Next()
}

type AuthorizationMiddleware struct {
	workspaceService interfaces.WorkspaceService
}

func NewWorkspaceRoleMiddleware(ws interfaces.WorkspaceService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		workspaceService: ws,
	}
}

// RequireWorkspaceAdmin checks if the user is an admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceAdmin(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceAdmin(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Next()
}

// RequireWorkspaceMembership checks if the user is a member or admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceMembership(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceMember(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Next()
}

// ClerkAuthenticationMiddleware to verify Clerk JWT Token
func ClerkAuthenticationMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "No authorization header"})
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid authorization header format"})
	}

	// Verify token using jwt.Verify
	claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
		Token: tokenParts[1],
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	// Store session info in context
	auth.StoreSessionInfoInContext(c, claims)

	// Continue to next middleware
	return c.Next()
}
