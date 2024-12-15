package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type WorkflowHandler struct {
	workflowService interfaces.WorkflowService
	logger          *logger.Logger
}

func NewWorkflowHandler(workflowService interfaces.WorkflowService, logger *logger.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
		logger:          logger.WithFields(map[string]interface{}{"module": "workflow_handler"}),
	}
}

// StartWorkflow starts a new workflow
func (h *WorkflowHandler) StartWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	switch workflowRequest.Type {
	case "screenshot", "report":
		// Validate the request
		if err := workflowRequest.Validate(true); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unknown workflow type. Allowed types are 'screenshot' and 'report'",
		})
	}

	// Start the workflow
	workflowResponse, err := h.workflowService.StartWorkflow(c.Context(), workflowRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workflowResponse)
}

func (h *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	switch workflowRequest.Type {
	case "screenshot", "report":
		// Validate the request
		if err := workflowRequest.Validate(true); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unknown workflow type. Allowed types are 'screenshot' and 'report'",
		})
	}

	// Start the workflow
	workflowResponse, err := h.workflowService.GetWorkflow(c.Context(), workflowRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workflowResponse)
}

func (h *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	limit := c.Query("limit", "100")
	status := c.Query("status", "running")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 100
	}
	// List all workflows
	workflows, count, err := h.workflowService.ListWorkflows(c.Context(), status, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := models.WorkflowListResponse{
		Workflows: workflows,
		Total:     count,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
