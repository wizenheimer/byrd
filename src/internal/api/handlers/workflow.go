package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type WorkflowHandler struct {
	workflowService interfaces.WorkflowService
	logger          *logger.Logger
}

func NewWorkflowHandler(workflowService interfaces.WorkflowService, logger *logger.Logger) *WorkflowHandler {
	wh := WorkflowHandler{
		workflowService: workflowService,
		logger:          logger.WithFields(map[string]interface{}{"module": "workflow_handler"}),
	}
	return &wh
}

func (wh *WorkflowHandler) StartWorkflow(c *fiber.Ctx) error {
	var workflowRequest api.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		wh.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request body"})
	}

	workflow, err := wh.workflowService.StartWorkflow(context.Background(), workflowRequest)
	if err != nil {
		wh.logger.Error("failed to start workflow", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to start workflow"})
	}
	return c.JSON(workflow)
}

func (wh *WorkflowHandler) StopWorkflow(c *fiber.Ctx) error {
	var workflowRequest api.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		wh.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request body"})
	}

	if err := wh.workflowService.StopWorkflow(context.Background(), workflowRequest); err != nil {
		wh.logger.Error("failed to stop workflow", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to stop workflow"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (wh *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	var workflowRequest api.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		wh.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request body"})
	}

	workflow, err := wh.workflowService.GetWorkflow(context.Background(), workflowRequest)
	if err != nil {
		wh.logger.Error("failed to get workflow", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get workflow"})
	}
	return c.JSON(workflow)
}

func (wh *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	workflowStatusString := c.Query("workflow_status")
	workflowStatus, err := models.ParseWorkflowStatus(workflowStatusString)
	if err != nil {
		wh.logger.Error("failed to parse workflow status", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse workflow status"})
	}

	workflowTypeString := c.Query("workflow_type")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		wh.logger.Error("failed to parse workflow type", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse workflow type"})
	}

	wh.logger.Debug("list workflows", zap.String("workflow_status", workflowStatusString), zap.String("workflow_type", workflowTypeString))

	workflows, err := wh.workflowService.ListWorkflows(context.Background(), workflowStatus, workflowType)
	if err != nil {
		wh.logger.Error("failed to list workflows", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list workflows"})
	}
	return c.JSON(workflows)
}
