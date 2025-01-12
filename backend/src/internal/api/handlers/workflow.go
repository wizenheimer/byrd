// ./src/internal/api/handlers/workflow.go
package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/api"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type WorkflowHandler struct {
	workflowService workflow.WorkflowService
	logger          *logger.Logger
}

func NewWorkflowHandler(workflowService workflow.WorkflowService, logger *logger.Logger) *WorkflowHandler {
	wh := WorkflowHandler{
		workflowService: workflowService,
		logger:          logger.WithFields(map[string]interface{}{"module": "workflow_handler"}),
	}
	return &wh
}

func (wh *WorkflowHandler) StartWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse request body", err)
	}
	wh.logger.Info("Starting workflow", zap.Any("workflow_request", workflowRequest))

	workflow, err := wh.workflowService.StartWorkflow(context.Background(), workflowRequest)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to start workflow", err)
	}
	return sendDataResponse(c, fiber.StatusCreated, "workflow started successfully", workflow)
}

func (wh *WorkflowHandler) StopWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	if err := wh.workflowService.StopWorkflow(context.Background(), workflowRequest); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to stop workflow", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "workflow stopped successfully", workflowRequest)
}

func (wh *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	var workflowRequest models.WorkflowRequest
	if err := c.BodyParser(&workflowRequest); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	workflow, err := wh.workflowService.GetWorkflow(context.Background(), workflowRequest)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to get workflow", err)
	}

	return sendDataResponse(c, fiber.StatusOK, "workflow retrieved successfully", workflow)
}

func (wh *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	workflowStatusString := c.Query("workflow_status")
	workflowStatus, err := models.ParseWorkflowStatus(workflowStatusString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse workflow status", err)
	}

	workflowTypeString := c.Query("workflow_type")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse workflow type", err)
	}

	workflows, err := wh.workflowService.ListWorkflows(context.Background(), workflowStatus, workflowType)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to list workflows", err)
	}
	return c.JSON(workflows)
}
