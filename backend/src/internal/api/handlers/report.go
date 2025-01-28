package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	"github.com/wizenheimer/byrd/src/internal/service/report"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type ReportHandler struct {
	logger        *logger.Logger
	reportService report.ReportService
}

func NewReportHandler(logger *logger.Logger, reportService report.ReportService) *ReportHandler {
	return &ReportHandler{
		logger:        logger.WithFields(map[string]interface{}{"handler": "report"}),
		reportService: reportService,
	}
}

// GetReport gets a specific report by ID
func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("reportID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid report ID format", err.Error())
	}

	ctx := c.Context()
	report, err := h.reportService.Get(ctx, reportID)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not get report", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Retrieved report successfully", report)
}

// GetLatestReport gets the latest report for a workspace and competitor
func (h *ReportHandler) GetLatestReport(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	ctx := c.Context()
	report, err := h.reportService.GetLatest(ctx, workspaceID, competitorID)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not get latest report", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Retrieved latest report successfully", report)
}

// ListReports lists reports for a workspace and competitor with pagination
func (h *ReportHandler) ListReports(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	// Parse pagination parameters
	limit := utils.QueryIntPtr(c, "limit", 10)
	offset := utils.QueryIntPtr(c, "offset", 0)

	ctx := c.Context()
	reports, err := h.reportService.List(ctx, workspaceID, competitorID, limit, offset)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not list reports", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Retrieved reports successfully", reports)
}

// CreateReport creates a new report
func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var req api.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	report, err := h.reportService.Create(ctx, workspaceID, competitorID, req.History)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not create report", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created report successfully", report)
}

// DispatchReport dispatches a report to subscribers
func (h *ReportHandler) DispatchReport(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var req api.DispatchReportRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	err = h.reportService.Dispatch(ctx, workspaceID, competitorID, req.CompetitorName, req.SubscriberEmails)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not dispatch report", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Dispatched report successfully", nil)
}
