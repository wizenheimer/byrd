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
	var req models.URLDiffRequest
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

func (h *DiffHandler) CreateReport(c *fiber.Ctx) error {
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
