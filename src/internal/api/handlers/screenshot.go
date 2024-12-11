package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type ScreenshotHandler struct {
	screenshotService interfaces.ScreenshotService
	logger            *logger.Logger
}

func NewScreenshotHandler(screenshotService interfaces.ScreenshotService, logger *logger.Logger) *ScreenshotHandler {
	logger.Debug("creating new screenshot handler")

	return &ScreenshotHandler{
		screenshotService: screenshotService,
		logger:            logger.WithFields(map[string]interface{}{"module": "notification_handler"}),
	}
}

func (h *ScreenshotHandler) CreateScreenshot(c *fiber.Ctx) error {
	h.logger.Debug("creating new screenshot")

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
	h.logger.Debug("getting screenshot", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

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
	h.logger.Debug("getting screenshot content", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

	result, err := h.screenshotService.GetContent(c.Context(), url, weekNumber, weekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
