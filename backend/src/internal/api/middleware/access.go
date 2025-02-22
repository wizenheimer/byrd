// ./src/internal/api/middleware/access.go
package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	user_svc "github.com/wizenheimer/byrd/src/internal/service/user"
	workspace_svc "github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// Determines if user has rights to access a resource
type AccessMiddleware struct {
	workspaceService workspace_svc.WorkspaceService
	userService      user_svc.UserService
	tokenManager     *utils.TokenManager
	logger           *logger.Logger
}

// NewAccessMiddleware creates a new AccessMiddleware
func NewAccessMiddleware(workspaceService workspace_svc.WorkspaceService, userService user_svc.UserService, tokenManager *utils.TokenManager, logger *logger.Logger) *AccessMiddleware {
	return &AccessMiddleware{
		workspaceService: workspaceService,
		userService:      userService,
		tokenManager:     tokenManager,
		logger: logger.WithFields(
			map[string]interface{}{
				"module": "access_middleware",
			},
		),
	}
}

// Checks if the user is authenticated
// This prevents unauthenticated users from accessing the resource
func (m *AccessMiddleware) RequiresClerkToken(c *fiber.Ctx) error {
	// Validate user
	if err := createSession(c); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "User Authentication Failed", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

// Checks if the user is authenticated with a private token
func (m *AccessMiddleware) RequiresPrivateToken(c *fiber.Ctx) error {
	// Validate token
	if err := m.validatePrivateToken(c); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Token Authentication Failed", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

// Checks if the user has an active membership in the workspace with role admin
func (m *AccessMiddleware) RequiresWorkspaceAdmin(c *fiber.Ctx) error {
	// Validate workspace membership
	allowedRoles := []models.WorkspaceRole{models.RoleAdmin}
	allowedStatus := []models.MembershipStatus{models.ActiveMember}

	if err := m.validateWorkspaceMembership(c, allowedRoles, allowedStatus); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusForbidden, "User Access Denied", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

// Checks if the user has an active membership in the workspace with role admin or user
func (m *AccessMiddleware) RequiresWorkspaceMember(c *fiber.Ctx) error {
	// Validate workspace membership
	allowedRoles := []models.WorkspaceRole{models.RoleAdmin, models.RoleUser}
	allowedStatus := []models.MembershipStatus{models.ActiveMember}

	if err := m.validateWorkspaceMembership(c, allowedRoles, allowedStatus); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusForbidden, "User Access Denied", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

func (m *AccessMiddleware) RequiresActiveOrPendingWorkspaceMembership(c *fiber.Ctx) error {
	// Validate workspace membership
	allowedRoles := []models.WorkspaceRole{models.RoleAdmin, models.RoleUser}
	allowedStatus := []models.MembershipStatus{models.ActiveMember, models.PendingMember}

	if err := m.validateWorkspaceMembership(c, allowedRoles, allowedStatus); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusForbidden, "User Access Denied", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

// Checks if the user has an pending membership in the workspace with role admin or user
func (m *AccessMiddleware) RequiresPendingWorkspaceMember(c *fiber.Ctx) error {
	// Validate workspace membership
	allowedRoles := []models.WorkspaceRole{models.RoleAdmin, models.RoleUser}
	allowedStatus := []models.MembershipStatus{models.PendingMember}

	if err := m.validateWorkspaceMembership(c, allowedRoles, allowedStatus); err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusForbidden, "User Access Denied", err.Error())
	}

	// Continue to next middleware
	return c.Next()
}

// validatePrivateToken validates the private token
// This function returns an error if the token is invalid
func (m *AccessMiddleware) validatePrivateToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("authorization header not found")
	}

	// Extract token from "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return errors.New("invalid authorization header")
	}
	token := authHeader[7:]

	if !m.tokenManager.ValidateToken(token) {
		return errors.New("invalid token")
	}

	return nil
}

// validateWorkspaceMembership validates the user's membership in the workspace
func (m *AccessMiddleware) validateWorkspaceMembership(c *fiber.Ctx, allowedRoles []models.WorkspaceRole, allowedStatus []models.MembershipStatus) error {
	// Get workspace ID from context
	workspaceID, err := getWorkspaceIDFromContext(c)
	if err != nil {
		return err
	}

	// Get clerk user from context
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return err
	}

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return err
	}

	// Check if user is an admin
	workspaceUser, err := m.workspaceService.GetWorkspaceUser(c.Context(), workspaceID, userEmail)
	if err != nil {
		return err
	}

	// Check if user has the required role
	if !utils.Contains(allowedRoles, workspaceUser.Role) {
		return errors.New("user does not have the required role")
	}

	// Check if user has the required status
	if !utils.Contains(allowedStatus, workspaceUser.MembershipStatus) {
		return errors.New("user does not have the required status")
	}

	return nil
}

func (m *AccessMiddleware) RequiresSlackSignature(c *fiber.Ctx) error {
	// Verify the request signature
	verifier, err := slack.NewSecretsVerifier(http.Header(c.GetReqHeaders()), os.Getenv("SLACK_SIGNATURE_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("error initializing signature verifier: %s", err.Error()))
	}

	// Get the raw body from the request
	bodyBytes := c.Body()

	// Create a copy of the body bytes
	bodyBytesCopy := make([]byte, len(bodyBytes))
	copy(bodyBytesCopy, bodyBytes)

	// Write the request body to the verifier
	if _, err = verifier.Write(bodyBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error writing request body bytes for verification: %s", err.Error()))
	}

	// Ensure the request signature is valid
	if err = verifier.Ensure(); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("error verifying slack signature: %s", err.Error()))
	}

	// Store the body copy in locals if needed for subsequent middleware/handlers
	c.Locals("rawBody", bodyBytesCopy)

	// Continue to next middleware/handler
	return c.Next()
}
