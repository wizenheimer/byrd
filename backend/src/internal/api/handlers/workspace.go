// ./src/internal/api/handlers/workspace.go
package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type WorkspaceHandler struct {
	workspaceService workspace.WorkspaceService
	logger           *logger.Logger
	tx               *transaction.TxManager
}

func NewWorkspaceHandler(
	workspaceService workspace.WorkspaceService,
	tx *transaction.TxManager,
	logger *logger.Logger,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: workspaceService,
		logger:           logger,
		tx:               tx,
	}
}

// CreateWorkspace creates a new workspace
func (wh *WorkspaceHandler) CreateWorkspaceForUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "User not found in request context", err.Error())
	}

	var req api.WorkspaceCreationRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Couldn't parse the workspace creation request", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Validation failed for workspace creation request", err.Error())
	}

	req.Profiles, err = ai.Sanitize(req.Profiles)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Couldn't sanitize the profiles", err.Error())
	}

	var pages []models.PageProps
	for _, competitorURL := range req.Competitors {
		page, err := models.NewPageProps(competitorURL, req.Profiles)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}
	if len(pages) == 0 {
		pages = make([]models.PageProps, 0)
	}

	ctx := c.Context()
	workspace, err := wh.workspaceService.CreateWorkspace(ctx, clerkUser, pages, req.Team)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not create workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created workspace successfully", workspace)

}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspacesForUser(c *fiber.Ctx) error {
	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "User is not authorized to list the workspace", err.Error())
	}

	membershipStatusString := strings.ToLower(c.Query("membership_status", "active"))
	var membershipStatus models.MembershipStatus
	switch membershipStatusString {
	case "active":
		membershipStatus = models.ActiveMember
	case "pending":
		membershipStatus = models.PendingMember
	default:
		membershipStatus = models.ActiveMember
	}

	ctx := c.Context()
	workspaces, err := wh.workspaceService.ListUserWorkspaces(ctx, clerkUser, membershipStatus)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Workspace couldn't be listed for the user", err.Error())
	}
	return sendDataResponse(c, fiber.StatusOK, "Listed workspaces successfully", map[string]any{
		"workspaces":        workspaces,
		"membership_status": membershipStatus,
	})
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspaceByID(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	ctx := c.Context()
	workspace, err := wh.workspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not get workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Fetched workspace successfully", workspace)
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspaceByID(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req api.WorkspaceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.UpdateWorkspace(ctx, workspaceID, req.ToProps()); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated workspace successfully",
		map[string]any{
			"workspace_id": workspaceID,
		})
}

// DeleteWorkspace deletes a workspace by ID
func (wh *WorkspaceHandler) DeleteWorkspaceByID(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	ctx := c.Context()
	status, err := wh.workspaceService.DeleteWorkspace(ctx, workspaceID)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not delete workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Deleted workspace successfully", map[string]any{
		"workspace_id":     workspaceID,
		"workspace_status": status,
	})
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspaceByID(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.JoinWorkspace(ctx, clerkUser, workspaceID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not join workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Joined workspace successfully", map[string]any{
		"workspace_id": workspaceID,
	})
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspaceByID(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.LeaveWorkspace(ctx, clerkUser, workspaceID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not exit workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Exited workspace successfully", map[string]any{
		"workspace_id": workspaceID,
	})
}
