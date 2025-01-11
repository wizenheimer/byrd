// ./src/internal/api/handlers/screenshot.go
package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type ScreenshotHandler struct {
	screenshotService screenshot.ScreenshotService
	logger            *logger.Logger
}

// NewScreenshotHandler creates a new screenshot handler
func NewScreenshotHandler(screenshotService screenshot.ScreenshotService, logger *logger.Logger) *ScreenshotHandler {
	logger.Debug("creating new screenshot handler")

	return &ScreenshotHandler{
		screenshotService: screenshotService,
		logger:            logger.WithFields(map[string]interface{}{"module": "screenshot_handler"}),
	}
}

// CreateScreenshot creates a new screenshot
func (h *ScreenshotHandler) CreateScreenshot(c *fiber.Ctx) error {
	h.logger.Debug("creating new screenshot")

	var sOpts models.ScreenshotRequestOptions
	if err := c.BodyParser(&sOpts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&sOpts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	screenshotResult, err := h.screenshotService.GetCurrentImage(c.Context(), true, sOpts)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create screenshot", err)
	}

	hOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   sOpts.URL,
		RenderedURL: screenshotResult.Metadata.RenderedURL,
	}

	contentResult, err := h.screenshotService.GetCurrentHTMLContent(c.Context(), true, hOpts)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create screenshot content", err)
	}

	return sendDataResponse(c, fiber.StatusCreated, "Screenshot created successfully", map[string]interface{}{
		"image":   screenshotResult,
		"content": contentResult,
	})
}

// GetScreenshotImage returns a screenshot image
func (h *ScreenshotHandler) GetScreenshotImage(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	result, err := h.screenshotService.GetImage(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get screenshot image", err)
	}

	return h.sendPNGResponse(c, result)
}

// GetScreenshotContent retrieves the screenshot content from the screenshot service
func (h *ScreenshotHandler) GetScreenshotContent(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	result, err := h.screenshotService.GetHTMLContent(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get screenshot content", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "Fetched screenshot content successfully", result)
}

func (h *ScreenshotHandler) ListScreenshots(c *fiber.Ctx) error {
	var opts models.ListScreenshotsOptions
	if err := c.BodyParser(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	result, err := h.screenshotService.ListScreenshots(c.Context(), opts.URL, opts.ContentType, opts.MaxItems)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list screenshots", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed screenshots successfully", result)
}

// addScreenshotMetadataToHeaders adds the screenshot metadata to the response headers
func (h *ScreenshotHandler) addScreenshotMetadataToHeaders(c *fiber.Ctx, screenshotResult *models.ScreenshotImageResponse) {
	if screenshotResult == nil {
		return
	}
	c.Set("x-image-height", fmt.Sprintf("%v", screenshotResult.ImageHeight))
	c.Set("x-image-width", fmt.Sprintf("%v", screenshotResult.ImageWidth))

	metadata := screenshotResult.Metadata
	if metadata == nil {
		return
	}
	c.Set("x-source-url", metadata.SourceURL)
	c.Set("x-rendered-url", metadata.RenderedURL)
	c.Set("x-screenshot-year", fmt.Sprintf("%v", metadata.Year))
	c.Set("x-screenshot-week-number", fmt.Sprintf("%v", metadata.WeekNumber))
	c.Set("x-screenshot-week-day", fmt.Sprintf("%v", metadata.WeekDay))
}

// handleDefaults sets the year, week number, and week day to the current values if they are not set
func (h *ScreenshotHandler) handleTimeDefaults(opts *models.GetScreenshotOptions) {
	now := time.Now()

	// If the year, week number, or week day are not set, use the current year, week number, and week day
	currentYear, currentWeek := now.ISOWeek()
	currentWeekDay := int(now.Weekday())

	if opts.Year == nil {
		opts.Year = utils.ToPtr(currentYear)
	}
	if opts.WeekNumber == nil {
		opts.WeekNumber = utils.ToPtr(currentWeek)
	}
	if opts.WeekDay == nil {
		opts.WeekDay = utils.ToPtr(currentWeekDay)
	}
}
