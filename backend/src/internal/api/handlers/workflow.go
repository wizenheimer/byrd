// ./src/internal/api/handlers/workflow.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/pkg/logger"
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
	return nil
}

func (wh *WorkflowHandler) StopWorkflow(c *fiber.Ctx) error {
	return nil
}

func (wh *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	return nil
}

func (wh *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	return nil
}
