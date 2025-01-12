// ./src/internal/api/middleware/path.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

const (
	CompetitorIDParamKey = "competitorID"
	PageIDParamKey       = "pageID"
	WorkspaceIDParamKey  = "workspaceID"
)

type WorkspacePathValidationMiddleware struct {
	ws     workspace.WorkspaceService
	logger *logger.Logger
}

func NewWorkspacePathValidationMiddleware(ws workspace.WorkspaceService, logger *logger.Logger) *WorkspacePathValidationMiddleware {
	return &WorkspacePathValidationMiddleware{
		ws:     ws,
		logger: logger.WithFields(map[string]interface{}{"module": "workspace_path_validation_middleware"}),
	}
}

// ValidateWorkspacePath checks if the workspace exists
func (m *WorkspacePathValidationMiddleware) ValidateWorkspacePath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	// Verify workspace exists
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}
	exists, err := m.ws.WorkspaceExists(c.Context(), workspaceUUID)

	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusInternalServerError, "Could not verify workspace", err)
	}

	if !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Workspace not found", map[string]interface{}{"workspaceID": workspaceID, "exists": exists})
	}

	return c.Next()
}

// ValidateCompetitorPath checks if the competitor exists
func (m *WorkspacePathValidationMiddleware) ValidateCompetitorPath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	competitorID := c.Params(CompetitorIDParamKey)
	if competitorID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid competitor ID", map[string]interface{}{"competitorID": competitorID})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid competitor ID", map[string]interface{}{"competitorID": competitorID})
	}

	// Verify competitor exists
	if exists, err := m.ws.WorkspaceCompetitorExists(c.Context(), workspaceUUID, competitorUUID); err != nil || !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Competitor not found", map[string]interface{}{"competitorID": competitorID, "exists": exists})
	}

	return c.Next()
}

// ValidatePagePath checks if the page exists
func (m *WorkspacePathValidationMiddleware) ValidatePagePath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	competitorID := c.Params(CompetitorIDParamKey)
	if competitorID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid competitor ID", map[string]interface{}{"competitorID": competitorID})
	}

	pageID := c.Params(PageIDParamKey)
	if pageID == "" {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid page ID", map[string]interface{}{"pageID": pageID})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceID})
	}

	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid competitor ID", map[string]interface{}{"competitorID": competitorID})
	}

	pageUUID, err := uuid.Parse(pageID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid page ID", map[string]interface{}{"pageID": pageID})
	}

	// Verify page exists
	if exists, err := m.ws.WorkspaceCompetitorPageExists(c.Context(), workspaceUUID, competitorUUID, pageUUID); err != nil || !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Page not found", map[string]interface{}{"pageID": pageID, "exists": exists})
	}

	return c.Next()
}
