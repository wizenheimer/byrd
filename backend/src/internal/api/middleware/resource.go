// ./src/internal/api/middleware/resource.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type ResourceMiddleware struct {
	ws     workspace.WorkspaceService
	logger *logger.Logger
}

func NewResourceMiddleware(workspaceService workspace.WorkspaceService, logger *logger.Logger) *ResourceMiddleware {
	return &ResourceMiddleware{
		ws: workspaceService,
		logger: logger.WithFields(map[string]interface{}{
			"module": "resource_middleware",
		}),
	}
}

// ValidateWorkspaceResource checks if the workspace exists
// NOTE: Workspace checks are implicitly done by other middlewares
// Added here for completeness
func (m *ResourceMiddleware) ValidateWorkspaceResource(c *fiber.Ctx) error {
	// Verify workspace exists
	workspaceUUID, err := getWorkspaceIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceUUID})
	}
	exists, err := m.ws.WorkspaceExists(c.Context(), workspaceUUID)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusInternalServerError, "Could not verify workspace", err.Error())
	}

	if !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Workspace not found", map[string]interface{}{"workspaceID": workspaceUUID, "exists": exists})
	}

	return c.Next()
}

// ValidateCompetitorResource checks if the competitor exists
func (m *ResourceMiddleware) ValidateCompetitorResource(c *fiber.Ctx) error {
	workspaceUUID, err := getWorkspaceIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Couldn't get workspace ID from context", map[string]interface{}{"error": err.Error()})
	}

	competitorUUID, err := getCompetitorIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusUnauthorized, "Couldn't get competitor ID from context", map[string]interface{}{"error": err.Error()})
	}

	// Verify competitor exists
	if exists, err := m.ws.WorkspaceCompetitorExists(c.Context(), workspaceUUID, competitorUUID); err != nil || !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Competitor not found", map[string]interface{}{"competitorID": competitorUUID, "exists": exists})
	}

	return c.Next()
}

// ValidatePageResource checks if the page exists
func (m *ResourceMiddleware) ValidatePageResource(c *fiber.Ctx) error {
	workspaceUUID, err := getWorkspaceIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid workspace ID", map[string]interface{}{"workspaceID": workspaceUUID})
	}

	competitorUUID, err := getCompetitorIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid competitor ID", map[string]interface{}{"competitorID": competitorUUID})
	}

	pageUUID, err := getPageIDFromContext(c)
	if err != nil {
		return sendErrorResponse(c, m.logger, fiber.StatusBadRequest, "Invalid page ID", map[string]interface{}{"pageID": pageUUID})
	}

	// Verify page exists
	if exists, err := m.ws.WorkspaceCompetitorPageExists(c.Context(), workspaceUUID, competitorUUID, pageUUID); err != nil || !exists {
		return sendErrorResponse(c, m.logger, fiber.StatusNotFound, "Page not found", map[string]interface{}{"pageID": pageUUID, "exists": exists})
	}

	return c.Next()
}
