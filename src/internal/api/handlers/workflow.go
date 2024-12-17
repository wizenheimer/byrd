package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
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
	var workflowRequest models.WorkflowRequest
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
	var workflowRequest models.WorkflowRequest
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
	var workflowRequest models.WorkflowRequest
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
	var workflowStatus models.WorkflowStatus
	var workflowType models.WorkflowType
	if err := c.QueryParser(&workflowStatus); err != nil {
		wh.logger.Error("failed to parse request query", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request query"})
	}
	if err := c.QueryParser(&workflowType); err != nil {
		wh.logger.Error("failed to parse request query", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request query"})
	}

	workflows, err := wh.workflowService.ListWorkflows(context.Background(), workflowStatus, workflowType)
	if err != nil {
		wh.logger.Error("failed to list workflows", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list workflows"})
	}
	return c.JSON(workflows)
}
