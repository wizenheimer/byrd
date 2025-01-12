// ./src/internal/api/middleware/go
package middleware

import (
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	ClaimsContextKey = "claims"
	UserIDContextKey = "userId"
)

type AuthorizationMiddleware struct {
	workspaceService workspace.WorkspaceService
	logger           *logger.Logger
}

func NewAuthorizationMiddleware(ws workspace.WorkspaceService, logger *logger.Logger) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		workspaceService: ws,
		logger:           logger.WithFields(map[string]interface{}{"module": "authorization_middleware"}),
	}
}

// RequireWorkspaceAdmin checks if the user is an admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceAdmin(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	clerkUser, err := getClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error()})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceAdmin(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Unauthorized")
	}

	return c.Next()
}

// RequireWorkspaceMembership checks if the user is a member or admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceMembership(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	clerkUser, err := getClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error()})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceMember(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Unauthorized")
	}

	return c.Next()
}

type AuthenticatedMiddleware struct {
	logger *logger.Logger
}

func NewAuthenticatedMiddleware(logger *logger.Logger) *AuthenticatedMiddleware {
	return &AuthenticatedMiddleware{
		logger: logger,
	}
}

// ClerkAuthenticationMiddleware to verify Clerk JWT Token
func (m *AuthenticatedMiddleware) AuthenticationMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		m.logger.Debug("No authorization header", zap.Any("status", fiber.StatusUnauthorized))
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "No authorization header", map[string]interface{}{"status": fiber.StatusUnauthorized})
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		m.logger.Debug("Invalid authorization header format", zap.Any("status", fiber.StatusUnauthorized), zap.Any("tokenParts", tokenParts), zap.Any("len", len(tokenParts)))
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Invalid authorization header format", map[string]interface{}{"status": fiber.StatusUnauthorized, "tokenParts": tokenParts, "len": len(tokenParts)})
	}

	// Verify token using jwt.Verify
	claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
		Token: tokenParts[1],
	})
	if err != nil {
		m.logger.Debug("Invalid token", zap.Any("status", fiber.StatusUnauthorized), zap.Any("error", err.Error()))
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Invalid token", map[string]interface{}{"status": fiber.StatusUnauthorized, "error": err.Error()})
	}

	// Store session info in context
	storeSessionInfoInContext(c, claims)

	// Continue to next middleware
	return c.Next()
}

// getClerkUserIDFromContext gets the Clerk user ID from the context
// This function returns an error if the Clerk user ID is not found in the context
func getClerkUserIDFromContext(c *fiber.Ctx) (string, error) {
	clerkUserID, ok := c.Locals(UserIDContextKey).(string)
	if !ok || clerkUserID == "" {
		return "", fmt.Errorf("clerk user ID not found in context")
	}

	return clerkUserID, nil
}

// storeSessionInfoInContext stores the session info in the context
func storeSessionInfoInContext(c *fiber.Ctx, claims *clerk.SessionClaims) {
	c.Locals(UserIDContextKey, claims.Subject)
	c.Locals(ClaimsContextKey, claims)
}

// getClerkUserFromContext gets the Clerk user from the context
// This function returns an error if the Clerk user is not found in the context
func getClerkUserFromContext(c *fiber.Ctx) (*clerk.User, error) {
	userID, err := getClerkUserIDFromContext(c)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	clerkUser, err := user.Get(c.Context(), userID)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return clerkUser, nil
}
