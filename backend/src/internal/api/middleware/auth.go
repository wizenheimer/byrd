package middleware

import (
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/auth"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type AuthorizationMiddleware struct {
	workspaceService svc.WorkspaceService
	logger           *logger.Logger
}

func NewAuthorizationMiddleware(ws svc.WorkspaceService, logger *logger.Logger) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		workspaceService: ws,
		logger:           logger,
	}
}

// RequireWorkspaceAdmin checks if the user is an admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceAdmin(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error()})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceAdmin(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
	}

	return c.Next()
}

// RequireWorkspaceMembership checks if the user is a member or admin in the workspace
func (m *AuthorizationMiddleware) RequireWorkspaceMembership(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if clerkUser == nil || err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error()})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
	}

	if exists, err := m.workspaceService.ClerkUserIsWorkspaceMember(c.Context(), workspaceUUID, clerkUser); !exists || err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", map[string]interface{}{"error": err.Error(), "workspaceID": workspaceID})
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
		return sendErrorResponse(c, fiber.StatusUnauthorized, "No authorization header", map[string]interface{}{"status": fiber.StatusUnauthorized})
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		m.logger.Debug("Invalid authorization header format", zap.Any("status", fiber.StatusUnauthorized), zap.Any("tokenParts", tokenParts), zap.Any("len", len(tokenParts)))
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format", map[string]interface{}{"status": fiber.StatusUnauthorized, "tokenParts": tokenParts, "len": len(tokenParts)})
	}

	// Verify token using jwt.Verify
	claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
		Token: tokenParts[1],
	})
	if err != nil {
		m.logger.Debug("Invalid token", zap.Any("status", fiber.StatusUnauthorized), zap.Any("error", err.Error()))
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", map[string]interface{}{"status": fiber.StatusUnauthorized, "error": err.Error()})
	}

	// Store session info in context
	auth.StoreSessionInfoInContext(c, claims)

	// Continue to next middleware
	return c.Next()
}
