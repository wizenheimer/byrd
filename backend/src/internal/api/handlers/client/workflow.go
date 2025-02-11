// ./src/internal/api/handlers/workflow.go
package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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
	workflowTypeString := c.Params("workflowType")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse workflow type", err.Error())
	}

	jobID, err := wh.workflowService.Submit(c.Context(), workflowType)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to start workflow", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "workflow started successfully", map[string]any{
		"workflowType": workflowType,
		"jobID":        jobID,
	})
}

func (wh *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	workflowTypeString := c.Params("workflowType")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse workflow type", err.Error())
	}

	jobIDString := c.Params("jobID")
	jobID, err := uuid.Parse(jobIDString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse job ID", err.Error())
	}

	job, err := wh.workflowService.State(c.Context(), workflowType, jobID)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to get workflow", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "workflow retrieved successfully", map[string]any{
		"workflowType": workflowType,
		"jobID":        jobID,
		"job":          job,
	})
}

func (wh *WorkflowHandler) StopWorkflow(c *fiber.Ctx) error {
	workflowTypeString := c.Params("workflowType")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse workflow type", err.Error())
	}

	jobIDString := c.Params("jobID")
	jobID, err := uuid.Parse(jobIDString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "failed to parse job ID", err.Error())
	}

	if err := wh.workflowService.Stop(context.Background(), workflowType, jobID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to stop workflow", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "workflow stopped successfully", map[string]any{
		"workflowType": workflowType,
		"jobID":        jobID,
	})
}

func (wh *WorkflowHandler) ListCheckpoint(c *fiber.Ctx) error {
	jobStatusString := c.Query("job_status")
	workflowStatus, err := models.ParseJobStatus(jobStatusString)
	if err != nil {
		workflowStatus = models.JobStatusRunning
	}

	workflowTypeString := c.Query("workflow_type")
	workflowType, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		workflowType = models.ScreenshotWorkflowType
	}

	jobs, err := wh.workflowService.List(context.Background(), workflowType, workflowStatus)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "failed to list workflows", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "workflows listed successfully", map[string]any{
		"workflowStatus": workflowStatus,
		"workflowType":   workflowType,
		"jobs":           jobs,
	})
}

func (wh *WorkflowHandler) ListHistory(c *fiber.Ctx) error {
	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	workflowTypeString := c.Query("workflow_type", "screenshot")
	wf, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workflow type", err.Error())
	}

	schedules, err := wh.workflowService.History(c.Context(), &limits, &offsets, &wf)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Failed to list schedules", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully listed schedules", schedules)
}
