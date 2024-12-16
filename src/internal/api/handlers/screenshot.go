package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
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

	var opts models.ScreenshotRequestOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, _, _, err := h.screenshotService.CaptureScreenshot(c.Context(), opts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
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

	result, err := h.screenshotService.GetScreenshotImage(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
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

	result, err := h.screenshotService.GetScreenshotContent(c.Context(), opts.URL, *opts.Year, *opts.WeekNumber, *opts.WeekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.addScreenshotMetadataToHeaders(c, result.Metadata)

	return c.JSON(fiber.Map{
		"status": result.Status,
		"data":   result.Content,
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
func (h *ScreenshotHandler) addScreenshotMetadataToHeaders(c *fiber.Ctx, metadata *models.ScreenshotMetadata) {
	if metadata == nil {
		return
	}
	c.Set("x-source-url", metadata.SourceURL)
	c.Set("x-fetched-at", metadata.FetchedAt)
	c.Set("x-screenshot-service", metadata.ScreenshotService)
	c.Set("x-image-height", fmt.Sprintf("%v", metadata.ImageHeight))
	c.Set("x-image-width", fmt.Sprintf("%v", metadata.ImageWidth))
	if metadata.PageTitle != nil {
		c.Set("x-page-title", *metadata.PageTitle)
	}
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
	h.addScreenshotMetadataToHeaders(c, result.Metadata)

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
