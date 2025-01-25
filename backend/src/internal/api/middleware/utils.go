// ./src/internal/api/middleware/utils.go
package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	CompetitorIDParamKey = "competitorID"
	PageIDParamKey       = "pageID"
	WorkspaceIDParamKey  = "workspaceID"
	ClaimsContextKey     = "claims"
	UserIDContextKey     = "userId"
)

// Validates the user session and stores the session info in the context
func createSession(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("authorization header not found")
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return errors.New("invalid authorization header")
	}

	// Verify token using jwt.Verify
	claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
		Token: tokenParts[1],
	})
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			err = errors.New("authorization token expired, request a new one")
		}
		return err
	}

	// Store session info in context
	storeSessionInfoInContext(c, claims)

	return nil
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

// getWorkspaceIDFromContext gets the workspace ID from the context
func getWorkspaceIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return uuid.Nil, errors.New("workspace ID not found in context")
	}

	// Verify workspace exists
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return uuid.Nil, errors.New("invalid workspace ID")
	}

	return workspaceUUID, nil
}

// getCompetitorIDFromContext gets the competitor ID from the context
func getCompetitorIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	competitorID := c.Params(CompetitorIDParamKey)
	if competitorID == "" {
		return uuid.Nil, errors.New("competitor ID not found in context")
	}

	// Verify competitor exists
	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return uuid.Nil, errors.New("invalid competitor ID")
	}

	return competitorUUID, nil
}

// getPageIDFromContext gets the page ID from the context
func getPageIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	pageID := c.Params(PageIDParamKey)
	if pageID == "" {
		return uuid.Nil, errors.New("page ID not found in context")
	}

	// Verify page exists
	pageUUID, err := uuid.Parse(pageID)
	if err != nil {
		return uuid.Nil, errors.New("invalid page ID")
	}

	return pageUUID, nil
}

// getClerkUserFromContext gets the Clerk user from the context
// This function returns an error if the Clerk user is not found in the context
func getClerkUserFromContext(c *fiber.Ctx) (*clerk.User, error) {
	userID, err := getClerkUserIDFromContext(c)
	if err != nil {
		return nil, errors.New("couldn't parse user credentials")
	}

	clerkUser, err := user.Get(c.Context(), userID)
	if err != nil {
		return nil, errors.New("couldn't process user credentials")
	}

	if clerkUser == nil {
		return nil, errors.New("clerk user is nil")
	}

	return clerkUser, nil
}
