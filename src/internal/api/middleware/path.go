package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

const (
	CompetitorIDParamKey = "competitorId"
	PageIDParamKey       = "pageId"
	WorkspaceIDParamKey  = "workspaceId"
)

type WorkspacePathValidationMiddleware struct {
	ws     svc.WorkspaceService
	logger *logger.Logger
}

func NewWorkspacePathValidationMiddleware(ws svc.WorkspaceService, logger *logger.Logger) *WorkspacePathValidationMiddleware {
	return &WorkspacePathValidationMiddleware{
		ws:     ws,
		logger: logger,
	}
}

// ValidateWorkspacePath checks if the workspace exists
func (m *WorkspacePathValidationMiddleware) ValidateWorkspacePath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	// Verify workspace exists
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}
	if exists, err := m.ws.WorkspaceExists(c.Context(), workspaceUUID); err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace not found",
		})
	}

	return c.Next()
}

// ValidateCompetitorPath checks if the competitor exists
func (m *WorkspacePathValidationMiddleware) ValidateCompetitorPath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	competitorID := c.Params(CompetitorIDParamKey)
	if competitorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid competitor ID",
		})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid competitor ID",
		})
	}

	// Verify competitor exists
	if exists, err := m.ws.WorkspaceCompetitorExists(c.Context(), workspaceUUID, competitorUUID); err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Competitor not found",
		})
	}

	return c.Next()
}

// ValidatePagePath checks if the page exists
func (m *WorkspacePathValidationMiddleware) ValidatePagePath(c *fiber.Ctx) error {
	workspaceID := c.Params(WorkspaceIDParamKey)
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	competitorID := c.Params(CompetitorIDParamKey)
	if competitorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid competitor ID",
		})
	}

	pageID := c.Params(PageIDParamKey)
	if pageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page ID",
		})
	}

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workspace ID",
		})
	}

	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid competitor ID",
		})
	}

	pageUUID, err := uuid.Parse(pageID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page ID",
		})
	}

	// Verify page exists
	if exists, err := m.ws.WorkspaceCompetitorPageExists(c.Context(), workspaceUUID, competitorUUID, pageUUID); err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	return c.Next()
}
