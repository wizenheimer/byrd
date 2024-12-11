package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type ScreenshotHandler struct {
	screenshotService interfaces.ScreenshotService
}

func NewScreenshotHandler(screenshotService interfaces.ScreenshotService) *ScreenshotHandler {
	return &ScreenshotHandler{
		screenshotService: screenshotService,
	}
}

func (h *ScreenshotHandler) CreateScreenshot(c *fiber.Ctx) error {
	var opts models.ScreenshotRequestOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.screenshotService.TakeScreenshot(c.Context(), opts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *ScreenshotHandler) GetScreenshot(c *fiber.Ctx) error {
	url := c.Params("hash")
	weekNumber := c.Params("weekNumber")
	weekDay := c.Params("weekDay")

	result, err := h.screenshotService.GetScreenshot(c.Context(), url, weekNumber, weekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *ScreenshotHandler) GetScreenshotContent(c *fiber.Ctx) error {
	url := c.Params("hash")
	weekNumber := c.Params("weekNumber")
	weekDay := c.Params("weekDay")

	result, err := h.screenshotService.GetContent(c.Context(), url, weekNumber, weekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
