// ./src/internal/api/handlers/schedule.go
package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/scheduler"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type ScheduleHandler struct {
	schedulerService scheduler.SchedulerService
	logger           *logger.Logger
}

func NewScheduleHandler(schedulerService scheduler.SchedulerService, logger *logger.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		schedulerService: schedulerService,
		logger: logger.WithFields(map[string]interface{}{
			"module": "schedule_handler",
		}),
	}
}

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req models.WorkflowScheduleProps
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	scheduleID, err := h.schedulerService.Schedule(c.Context(), req)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to create schedule", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully created schedule", map[string]any{
		"scheduleID": scheduleID.String(),
	})
}

func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	scheduleIDString := c.Params("scheduleID")
	scheduleIDUUID, err := uuid.Parse(scheduleIDString)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid schedule ID", "Invalid schedule ID")
	}
	scheduleID := models.ScheduleID(scheduleIDUUID)
	schedule, err := h.schedulerService.Get(c.Context(), scheduleID)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to get schedule", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully retrieved schedule", schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	scheduleIDString := c.Params("scheduleID")
	scheduleIDUUID, err := uuid.Parse(scheduleIDString)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid schedule ID", "Invalid schedule ID")
	}
	scheduleID := models.ScheduleID(scheduleIDUUID)

	var req models.WorkflowScheduleProps
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	updatedScheduleID, err := h.schedulerService.Reschedule(c.Context(), scheduleID, req)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to update schedule", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully updated schedule", map[string]any{
		"scheduleID": updatedScheduleID.String(),
	})
}

func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	scheduleIDString := c.Params("scheduleID")
	scheduleIDUUID, err := uuid.Parse(scheduleIDString)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid schedule ID", "Invalid schedule ID")
	}
	scheduleID := models.ScheduleID(scheduleIDUUID)

	err = h.schedulerService.Unschedule(c.Context(), scheduleID)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to delete schedule", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully deleted schedule", nil)
}

func (h *ScheduleHandler) ListSchedules(c *fiber.Ctx) error {
	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	workflowTypeString := c.Query("workflowType", "screenshot")
	wf, err := models.ParseWorkflowType(workflowTypeString)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid workflow type", err.Error())
	}

	schedules, err := h.schedulerService.List(c.Context(), &limits, &offsets, &wf)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to list schedules", err.Error())
	}

	return sendDataResponse(c, http.StatusOK, "Successfully listed schedules", schedules)
}
