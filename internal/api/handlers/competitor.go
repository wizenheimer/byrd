package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type CompetitorHandler struct {
	competitorService interfaces.CompetitorService
}

func NewCompetitorHandler(competitorService interfaces.CompetitorService) *CompetitorHandler {
	return &CompetitorHandler{
		competitorService: competitorService,
	}
}

func (h *CompetitorHandler) CreateCompetitor(c *fiber.Ctx) error {
	var input models.CompetitorInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.competitorService.Create(c.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *CompetitorHandler) ListCompetitors(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	competitors, total, err := h.competitorService.List(c.Context(), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   competitors,
		"metadata": fiber.Map{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *CompetitorHandler) GetCompetitor(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid competitor ID")
	}

	competitor, err := h.competitorService.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   competitor,
	})
}
