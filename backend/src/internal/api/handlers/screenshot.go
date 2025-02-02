// ./src/internal/api/handlers/screenshot.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type ScreenshotHandler struct {
	screenshotService screenshot.ScreenshotService
	logger            *logger.Logger
}

// NewScreenshotHandler creates a new screenshot handler
func NewScreenshotHandler(screenshotService screenshot.ScreenshotService, logger *logger.Logger) *ScreenshotHandler {
	return &ScreenshotHandler{
		screenshotService: screenshotService,
		logger: logger.WithFields(map[string]any{
			"module": "screenshot_handler",
		}),
	}
}

func (h *ScreenshotHandler) Retrieve(c *fiber.Ctx) error {
	type req struct {
		Options  models.ScreenshotRequestOptions `json:"options"`
		BackDate bool                            `json:"backDate" default:"false"`
	}

	var r req
	if err := c.BodyParser(&r); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	// Create a new screenshot
	screenshotImg, screenshotContent, err := h.screenshotService.Retrieve(c.Context(), r.Options, r.BackDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return sendDataResponse(c, fiber.StatusOK, "screenshot retrieved", map[string]any{
		"screenshot": map[string]any{
			"metadata": screenshotImg.Metadata,
			"path":     screenshotImg.StoragePath,
		},
		"content": map[string]any{
			"metadata": screenshotContent.Metadata,
			"path":     screenshotContent.StoragePath,
		},
	})
}

func (h *ScreenshotHandler) Refresh(c *fiber.Ctx) error {
	type req struct {
		Options  models.ScreenshotRequestOptions `json:"options"`
		BackDate bool                            `json:"backDate" default:"false"`
	}

	var r req
	if err := c.BodyParser(&r); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	// Create a new screenshot
	screenshotImg, screenshotContent, err := h.screenshotService.Refresh(c.Context(), r.Options, r.BackDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return sendDataResponse(c, fiber.StatusOK, "screenshot refreshed", map[string]any{
		"screenshot": map[string]any{
			"metadata": screenshotImg.Metadata,
			"path":     screenshotImg.StoragePath,
		},
		"content": map[string]any{
			"metadata": screenshotContent.Metadata,
			"path":     screenshotContent.StoragePath,
		},
	})
}
