package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

type ScreenshotHandler struct {
	screenshotService interfaces.ScreenshotService
	logger            *logger.Logger
}

// NewScreenshotHandler creates a new screenshot handler
func NewScreenshotHandler(screenshotService interfaces.ScreenshotService, logger *logger.Logger) *ScreenshotHandler {
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

	screenshotResult, e := h.screenshotService.GetCurrentImage(c.Context(), true, sOpts)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create screenshot", e)
	}

	hOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   sOpts.URL,
		RenderedURL: screenshotResult.Metadata.RenderedURL,
	}

	contentResult, e := h.screenshotService.GetCurrentHTMLContent(c.Context(), true, hOpts)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create screenshot content", e)
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

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	result, e := h.screenshotService.GetImage(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get screenshot image", e)
	}

	return h.sendPNGResponse(c, result)
}

// GetScreenshotContent retrieves the screenshot content from the screenshot service
func (h *ScreenshotHandler) GetScreenshotContent(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	result, e := h.screenshotService.GetHTMLContent(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get screenshot content", e)
	}

	return sendDataResponse(c, fiber.StatusOK, "Fetched screenshot content successfully", result)
}

func (h *ScreenshotHandler) ListScreenshots(c *fiber.Ctx) error {
	var opts models.ListScreenshotsOptions
	if err := c.BodyParser(&opts); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	result, e := h.screenshotService.ListScreenshots(c.Context(), opts.URL, opts.ContentType, opts.MaxItems)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list screenshots", e)
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
