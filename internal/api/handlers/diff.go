package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type DiffHandler struct {
	diffService interfaces.DiffService
}

func NewDiffHandler(diffService interfaces.DiffService) *DiffHandler {
	return &DiffHandler{
		diffService: diffService,
	}
}

func (h *DiffHandler) CreateDiff(c *fiber.Ctx) error {
	var req models.DiffRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.diffService.CreateDiff(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *DiffHandler) GetDiffHistory(c *fiber.Ctx) error {
	params := models.DiffHistoryParams{
		URL:        c.Query("url"),
		FromRunID:  c.Query("fromRunId"),
		ToRunID:    c.Query("toRunId"),
		WeekNumber: c.Query("weekNumber"),
		Limit:      c.QueryInt("limit", 10),
	}

	result, err := h.diffService.GetDiffHistory(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *DiffHandler) CreateReport(c *fiber.Ctx) error {
	var req models.ReportRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.diffService.GenerateReport(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
