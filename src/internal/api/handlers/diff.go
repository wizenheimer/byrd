package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type DiffHandler struct {
	diffService interfaces.DiffService
	logger      *logger.Logger
}

func NewDiffHandler(diffService interfaces.DiffService, logger *logger.Logger) *DiffHandler {
	logger.Debug("creating new diff handler")

	return &DiffHandler{
		diffService: diffService,
		logger:      logger.WithFields(map[string]interface{}{"module": "diff_handler"}),
	}
}

func (h *DiffHandler) CreateDiff(c *fiber.Ctx) error {
	h.logger.Debug("creating new diff")

	var req models.URLDiffRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.diffService.GetDiffAnalysis(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *DiffHandler) CreateReport(c *fiber.Ctx) error {
	h.logger.Debug("creating new report")

	var req models.WeeklyReportRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.diffService.CreateReport(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
