// ./src/internal/api/handlers/pages.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/commons"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// ListPagesForCompetitor lists pages for a competitor
func (wh *WorkspaceHandler) ListPagesForCompetitor(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	ctx := c.Context()
	pages, hasMore, err := wh.workspaceService.ListPagesForCompetitor(ctx, workspaceID, competitorID, &limits, &offsets)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not add page to competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Listed page for competitor successfully", map[string]any{
		"pages":   pages,
		"hasMore": hasMore,
	})
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPagesToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var req []api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidateArray(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	var pages []models.PageProps
	for _, r := range req {
		r.DiffProfile, err = ai.Sanitize(r.DiffProfile)
		if err != nil {
			continue
		}
		page, err := r.ToProps()
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}

	ctx := c.Context()
	createdPages, err := wh.workspaceService.AddPageToCompetitor(ctx, competitorID, pages)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not add page to competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Added page to competitor successfully", createdPages)
}

func (wh *WorkspaceHandler) GetPageForCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	ctx := c.Context()
	page, err := wh.workspaceService.GetPageForCompetitor(ctx, competitorID, pageID)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not get page for competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Fetched page for competitor successfully", page)
}

func (wh *WorkspaceHandler) UpdatePageForCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	var req api.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	prop, err := req.ToProps()
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	page, err := wh.workspaceService.UpdateCompetitorPage(ctx, competitorID, pageID, prop)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update page in competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated page in competitor successfully", page)
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.RemovePageFromWorkspace(ctx, competitorID, pageID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not remove page from competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed page from competitor successfully", nil)
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	_, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	_, err = uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid page ID format", err.Error())
	}

	pageNumber := max(1, c.QueryInt("_page", commons.DefaultPageNumber))
	pageSize := max(10, c.QueryInt("_limit", commons.DefaultPageSize))

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limit := params.GetLimit()
	offset := params.GetOffset()

	ctx := c.Context()
	history, hasMore, err := wh.workspaceService.ListHistoryForPage(ctx, pageID, &limit, &offset)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list page history", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed page history successfully", map[string]any{
		"history": history,
		"hasMore": hasMore,
	})
}
