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

// GetWorkflow retrieves a workflow
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

// ListWorkflows lists all workflows
func (h *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	limit := c.Query("limit", "100")
	status := c.Query("status", "all")
	workflowTypeString := c.Query("type", "screenshot")

	// Parse the workflow type
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 100
	}

	workflowStatusList := make([]models.WorkflowStatus, 0)

	// Parse the workflow status
	if status == "all" {
		workflowStatusList = append(workflowStatusList, models.WorkflowStatusPending, models.WorkflowStatusRunning, models.WorkflowStatusCompleted, models.WorkflowStatusFailed, models.WorkflowStatusAborted)
	} else {
		workflowStatus, err := models.ParseWorkflowStatus(status)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		workflowStatusList = append(workflowStatusList, workflowStatus)
	}

	// List all workflows
	responses := make([]models.WorkflowListResponse, 0)
	for _, workflowStatus := range workflowStatusList {

		workflows, count, err := h.workflowService.ListWorkflows(c.Context(), workflowStatus, workflowType, limitInt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if len(workflows) == 0 {
			workflows = make([]models.WorkflowResponse, 0)
		}

		response := models.WorkflowListResponse{
			Workflows:      workflows,
			Total:          count,
			WorkflowStatus: workflowStatus,
		}

		responses = append(responses, response)
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// RecoverWorkflow recovers workflows from available checkpoints
func (h *WorkflowHandler) RecoverWorkflow(c *fiber.Ctx) error {
	if err := h.workflowService.RecoverWorkflow(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Workflows recovered successfully",
	})
}

// Shutdown shuts down the workflow service
func (h *WorkflowHandler) Shutdown(c *fiber.Ctx) error {
	if err := h.workflowService.Shutdown(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Workflow service shutdown successfully",
	})
}

// StopWorkflow stops a workflow
func (h *WorkflowHandler) StopWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := workflowRequest.Validate(true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	workflowType := models.WorkflowType(workflowRequest.Type)
	workflowID := models.WorkflowIdentifier{
		Type:         workflowType,
		Year:         *workflowRequest.Year,
		WeekNumber:   *workflowRequest.WeekNumber,
		BucketNumber: *workflowRequest.BucketNumber,
	}

	if err := h.workflowService.StopWorkflow(c.Context(), workflowID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Workflow stopped successfully",
	})
}
