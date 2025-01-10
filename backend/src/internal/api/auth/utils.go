// ./src/internal/api/auth/utils.go
package auth

import (
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
)

// GetClerkUserFromContext gets the Clerk user from the context
// This function returns an error if the Clerk user is not found in the context
func GetClerkUserFromContext(c *fiber.Ctx) (*clerk.User, error) {
	userID, err := GetClerkUserIDFromContext(c)
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

// GetClerkUserIDFromContext gets the Clerk user ID from the context
// This function returns an error if the Clerk user ID is not found in the context
func GetClerkUserIDFromContext(c *fiber.Ctx) (string, error) {
	clerkUserID, ok := c.Locals(UserIDContextKey).(string)
	if !ok || clerkUserID == "" {
		return "", fmt.Errorf("clerk user ID not found in context")
	}

	return clerkUserID, nil
}

// GetClerkClaimsFromContext gets the Clerk session claims from the context
// This function returns an error if the Clerk session claims are not found in the context
func GetClerkClaimsFromContext(c *fiber.Ctx) (*clerk.SessionClaims, error) {
	clerkSessionClaims, ok := c.Locals(ClaimsContextKey).(*clerk.SessionClaims)
	if !ok || clerkSessionClaims == nil {
		return nil, fmt.Errorf("clerk session claims not found in context")
	}

	return clerkSessionClaims, nil
}

func StoreSessionInfoInContext(c *fiber.Ctx, claims *clerk.SessionClaims) {
	c.Locals(UserIDContextKey, claims.Subject)
	c.Locals(ClaimsContextKey, claims)
}
