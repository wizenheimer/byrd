package handlers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/api"
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
		logger:            logger.WithFields(map[string]interface{}{"module": "screenshot_handler"}),
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

// GetScreenshotImage returns a screenshot image
func (h *ScreenshotHandler) GetScreenshotImage(c *fiber.Ctx) error {
	var opts models.GetScreenshotOptions
	if err := c.BodyParser(&opts); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	url, err := url.Parse(opts.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid URL")
	}

	if opts.WeekNumber < 1 || opts.WeekNumber > 52 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid week number, must be between 1 and 52")
	}
	weekNumber := strconv.Itoa(opts.WeekNumber)

	if opts.WeekDay < 1 || opts.WeekDay > 7 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid week day, must be between 1 and 7")
	}
	weekDay := strconv.Itoa(opts.WeekDay)

	h.logger.Debug("getting screenshot", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

	result, err := h.screenshotService.GetScreenshotImage(c.Context(), url.String(), weekNumber, weekDay)
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

	url, err := url.Parse(opts.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid URL")
	}

	if opts.WeekNumber < 1 || opts.WeekNumber > 52 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid week number, must be between 1 and 52")
	}
	weekNumber := strconv.Itoa(opts.WeekNumber)

	if opts.WeekDay < 1 || opts.WeekDay > 7 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid week day, must be between 1 and 7")
	}
	weekDay := strconv.Itoa(opts.WeekDay)

	h.logger.Debug("getting screenshot content", zap.Any("url", url), zap.Any("week_number", weekNumber), zap.Any("week_day", weekDay))

	result, err := h.screenshotService.GetScreenshotContent(c.Context(), url.String(), weekNumber, weekDay)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.addScreenshotMetadataToHeaders(c, result.Metadata)

	return c.JSON(fiber.Map{
		"status": result.Status,
		"data":   result.Content,
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
