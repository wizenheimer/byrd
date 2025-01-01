package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/api"
	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
	"go.uber.org/zap"
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	screenshotResult, err := h.screenshotService.GetCurrentImage(c.Context(), true, sOpts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	hOpts := models.ScreenshotHTMLRequestOptions{
		SourceURL:   sOpts.URL,
		RenderedURL: screenshotResult.Metadata.RenderedURL,
	}

	contentResult, err := h.screenshotService.GetCurrentHTMLContent(c.Context(), true, hOpts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": map[string]interface{}{
			"image":   screenshotResult,
			"content": contentResult,
		},
	})
}

// GetScreenshotImage returns a screenshot image
func (h *ScreenshotHandler) GetScreenshotImage(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	h.logger.Debug("getting screenshot image", zap.Any("url", opts.URL), zap.Any("year", opts.Year), zap.Any("week_number", opts.WeekNumber), zap.Any("week_day", opts.WeekDay))

	result, err := h.screenshotService.GetImage(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return h.sendPNGResponse(c, result)
}

// GetScreenshotContent retrieves the screenshot content from the screenshot service
func (h *ScreenshotHandler) GetScreenshotContent(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// handleDefaults sets the year, week number, and week day to the current values if they are not set
	h.handleTimeDefaults(&opts)

	h.logger.Debug("getting screenshot image", zap.Any("url", opts.URL), zap.Any("year", opts.Year), zap.Any("week_number", opts.WeekNumber), zap.Any("week_day", opts.WeekDay))

	result, err := h.screenshotService.GetHTMLContent(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": result.Status,
		"data":   result,
	})
}

func (h *ScreenshotHandler) ListScreenshots(c *fiber.Ctx) error {
	var opts models.ListScreenshotsOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	h.logger.Debug("listing screenshots", zap.Any("url", opts.URL), zap.Any("content_type", opts.ContentType))

	result, err := h.screenshotService.ListScreenshots(c.Context(), opts.URL, opts.ContentType, opts.MaxItems)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
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

// sendPNGResponse sends a PNG response with the screenshot image
func (h *ScreenshotHandler) sendPNGResponse(c *fiber.Ctx, result *models.ScreenshotImageResponse) error {
	// If the result is nil, return a 404 Not Found error
	if result == nil || result.Image == nil {
		return fiber.NewError(fiber.StatusNotFound, "Screenshot not found")
	}

	// WritePNGResponse writes the image to a PNG byte array
	pngBytes, err := api.WritePNGResponse(
		result.Image,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header to image/png
	c.Set("Content-Type", "image/png")

	// Add the screenshot metadata to the response headers
	h.addScreenshotMetadataToHeaders(c, result)

	// Send the PNG byte array as the response
	return c.Send(pngBytes)
}

// handleDefaults sets the year, week number, and week day to the current values if they are not set
func (h *ScreenshotHandler) handleTimeDefaults(opts *models.GetScreenshotOptions) {
	now := time.Now()

	// If the year, week number, or week day are not set, use the current year, week number, and week day
	currentYear, currentWeek := now.ISOWeek()
	currentWeekDay := int(now.Weekday())

	if opts.Year == nil {
		opts.Year = ptr.To(currentYear)
	}
	if opts.WeekNumber == nil {
		opts.WeekNumber = ptr.To(currentWeek)
	}
	if opts.WeekDay == nil {
		opts.WeekDay = ptr.To(currentWeekDay)
	}
}
