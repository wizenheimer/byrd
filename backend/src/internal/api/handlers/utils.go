// ./src/internal/api/handlers/response.go
package handlers

import (
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

const (
	ClaimsContextKey    = "claims"
	UserIDContextKey    = "userId"
	WorkspaceIDParamKey = "workspaceId"
)

func sendDataResponse(c *fiber.Ctx, status int, message string, data any) error {
	return commons.SendDataResponse(c, status, message, data)
}

func sendErrorResponse(c *fiber.Ctx, status int, message string, details any) error {
	return commons.SendErrorResponse(c, status, message, details)
}

// sendPNGResponse sends a PNG response with the screenshot image
func (h *ScreenshotHandler) sendPNGResponse(c *fiber.Ctx, result *models.ScreenshotImageResponse) error {
	// If the result is nil, return a 404 Not Found error
	if result == nil || result.Image == nil {
		return sendErrorResponse(c, fiber.StatusNotFound, "Screenshot not found", nil)
	}

	// WritePNGResponse writes the image to a PNG byte array
	pngBytes, err := utils.WritePNGResponse(
		result.Image,
	)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not write PNG response", err)
	}

	// Set the Content-Type header to image/png
	c.Set("Content-Type", "image/png")

	// Add the screenshot metadata to the response headers
	h.addScreenshotMetadataToHeaders(c, result)

	// Send the PNG byte array as the response
	return c.Send(pngBytes)
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
func getClerkClaimsFromContext(c *fiber.Ctx) (*clerk.SessionClaims, error) {
	clerkSessionClaims, ok := c.Locals(ClaimsContextKey).(*clerk.SessionClaims)
	if !ok || clerkSessionClaims == nil {
		return nil, fmt.Errorf("clerk session claims not found in context")
	}

	return clerkSessionClaims, nil
}
