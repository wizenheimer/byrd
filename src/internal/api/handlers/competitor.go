package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/logger"
)

type CompetitorHandler struct {
	competitorService interfaces.CompetitorService
	logger            *logger.Logger
}

func NewCompetitorHandler(competitorService interfaces.CompetitorService, logger *logger.Logger) *CompetitorHandler {
	logger.Debug("creating new competitor handler")

	return &CompetitorHandler{
		competitorService: competitorService,
		logger:            logger.WithFields(map[string]interface{}{"module": "competitor_handler"}),
	}
}

func (h *CompetitorHandler) CreateCompetitor(c *fiber.Ctx) error {
	h.logger.Debug("creating new competitor")

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
	h.logger.Debug("listing competitors")

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
	h.logger.Debug("getting competitor by ID")

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
