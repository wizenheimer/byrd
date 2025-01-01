package handlers

import (
	"github.com/gofiber/fiber/v2"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

// URLHandler is the handler that provides URL management operations
type URLHandler struct {
	// urlService is an interface that provides URL management operations
	urlService interfaces.URLService
	// logger is a structured logger for logging
	logger *logger.Logger
}

// NewURLHandler creates a new URLHandler
// It returns an error if the logger or urlService is nil
func NewURLHandler(urlService interfaces.URLService, logger *logger.Logger) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		logger:     logger.WithFields(map[string]interface{}{"module": "url_handler"}),
	}
}

// AddURL adds a new URL if it does not exist
func (h *URLHandler) AddURL(c *fiber.Ctx) error {
	h.logger.Debug("adding new URL")

	var input models.URLInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.urlService.AddURL(c.Context(), input.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

// ListURLs lists all URLs in batches
func (h *URLHandler) ListURLs(c *fiber.Ctx) error {
	h.logger.Debug("listing URLs")

	var input models.URLListInput
	if err := c.BodyParser(&input); err != nil {
		input.BatchSize = 10
	}

	if input.BatchSize <= 0 {
		input.BatchSize = 10
	}

	urlChan, errChan := h.urlService.ListURLs(c.Context(), input.BatchSize, input.LastSeenID)

	var urls []models.URL
	var errorStrings []string
	for url := range urlChan {
		urls = append(urls, url.URLs...)
	}

	if err := <-errChan; err != nil {
		errorStrings = append(errorStrings, err.Error())
	}

	data := map[string]interface{}{
		"urls":   urls,
		"errors": errorStrings,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	h.logger.Debug("deleting URL")

	var input models.URLInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := h.urlService.DeleteURL(c.Context(), input.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
