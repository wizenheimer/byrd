package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListCompetitorsForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}
	var competitorResponse []api.CompetitorResponse
	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limit := params.GetLimit()
	offset := params.GetOffset()

	ctx := c.Context()

	// List out the competitors for the workspace
	competitors, hasMore, err := wh.workspaceService.ListCompetitorsForWorkspace(ctx, workspaceID, &limit, &offset)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list workspace competitors", err.Error())
	}

	includePages := c.Query("includePages", "false") == "true"

	// Iterate through the competitors and list out the pages for each competitor if includePages is true
	for _, competitor := range competitors {
		// If includePages is true, list out the pages for the competitor
		p := []models.Page{} // Set to empty array
		if includePages {
			p, _, err = wh.workspaceService.ListPagesForCompetitor(ctx, workspaceID, competitor.ID, nil, nil)
			if err != nil {
				wh.logger.Debug("Could not list pages for competitor", zap.Error(err))
			}
		}

		// Prepare the response
		competitorResponse = append(competitorResponse, api.NewCompetitorResponse(&competitor, p))
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace competitors successfully", map[string]any{
		"competitors": competitorResponse,
		"hasMore":     hasMore,
	})
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req []api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidateArray(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	competitor, err := wh.workspaceService.AddCompetitorToWorkspace(ctx, workspaceID, req)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not create competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created competitor successfully", competitor)
}

func (wh *WorkspaceHandler) GetCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	ctx := c.Context()
	competitor, err := wh.workspaceService.GetCompetitorForWorkspace(ctx, workspaceID, competitorID)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not retrieve competitor", err.Error())
	}

	includePages := c.Query("includePages", "false") == "true"
	var pages []models.Page
	if includePages {
		pages, _, err = wh.workspaceService.ListPagesForCompetitor(ctx, workspaceID, competitorID, nil, nil)
		if err != nil {
			return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not retrieve pages for competitor", err.Error())
		}
	}

	return sendDataResponse(c, fiber.StatusOK, "Retrieved competitor successfully", api.NewCompetitorResponse(competitor, pages))
}

func (wh *WorkspaceHandler) UpdateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var req api.UpdateCompetitorRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	competitor, err := wh.workspaceService.UpdateCompetitorForWorkspace(ctx, workspaceID, competitorID, req.CompetitorName)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated competitor successfully", competitor)
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidWorkspaceID", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.RemoveCompetitorFromWorkspace(ctx, workspaceID, competitorID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not remove competitor from workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed competitor from workspace successfully", nil)
}
