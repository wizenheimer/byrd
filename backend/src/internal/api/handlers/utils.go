// ./src/internal/api/handlers/response.go
package handlers

import (
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	ClaimsContextKey    = "claims"
	UserIDContextKey    = "userId"
	WorkspaceIDParamKey = "workspaceId"
)

func sendDataResponse(c *fiber.Ctx, status int, message string, data any) error {
	return commons.SendDataResponse(c, status, message, data)
}

func sendErrorResponse(c *fiber.Ctx, logger *logger.Logger, status int, message string, details any) error {
	logger.Debug(message, zap.Int("status", status), zap.Any("details", details))
	return commons.SendErrorResponse(c, status, message, details)
}

// getClerkUserFromContext gets the Clerk user from the context
// This function returns an error if the Clerk user is not found in the context
func getClerkUserFromContext(c *fiber.Ctx) (*clerk.User, error) {
	userID, err := getClerkUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	clerkUser, err := user.Get(c.Context(), userID)
	if err != nil {
		return nil, err
	}

	if clerkUser == nil {
		return nil, fmt.Errorf("clerk user not found")
	}

	return clerkUser, nil
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

// getClerkClaimsFromContext gets the Clerk session claims from the context
// This function returns an error if the Clerk session claims are not found in the context
// func getClerkClaimsFromContext(c *fiber.Ctx) (*clerk.SessionClaims, error) {
// 	clerkSessionClaims, ok := c.Locals(ClaimsContextKey).(*clerk.SessionClaims)
// 	if !ok || clerkSessionClaims == nil {
// 		return nil, fmt.Errorf("clerk session claims not found in context")
// 	}

// 	return clerkSessionClaims, nil
// }
