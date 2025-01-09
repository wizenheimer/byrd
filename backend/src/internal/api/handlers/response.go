package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

func sendDataResponse(c *fiber.Ctx, status int, message string, data ...any) error {
	return commons.SendDataResponse(c, status, message, data)
}

func sendErrorResponse(c *fiber.Ctx, status int, message string, details ...any) error {
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
