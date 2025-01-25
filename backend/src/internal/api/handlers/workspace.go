package handlers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
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
	var req api.WorkspaceCreationRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "InvalidRequest", err.Error())
	}

	// Prepare the diff profiles
	var diffProfiles []string
	addedProfile := make(map[string]bool)

	// Loop through the profiles
	for _, profile := range req.Profiles {
		// Check if the profile exists in the AI service
		// Or if the profile has already been added
		if _, err := ai.GetField(profile); addedProfile[profile] || err != nil {
			wh.logger.Debug("skipping profile", zap.Any("profile", profile))
			continue
		}
		// Add the profile to the list of diff profiles
		addedProfile[profile] = true
		// Append the profile to the list of diff profiles
		diffProfiles = append(diffProfiles, profile)
	}

	// If there are no diff profiles, add the default field
	if len(diffProfiles) == 0 {
		diffProfiles = append(diffProfiles, ai.DefaultField)
		wh.logger.Debug("No valid profiles found, defaulting to default profile")
	}

	var pages []models.PageProps
	for _, competitorURL := range req.Competitors {
		if _, err := url.Parse(competitorURL); err != nil {
			continue
		}
		pageTitle, err := utils.GetPageTitle(competitorURL)
		if err != nil {
			continue
		}
		captureProfile := screenshot.GetDefaultScreenshotRequestOptions(competitorURL)
		pages = append(pages, models.PageProps{
			URL:            competitorURL,
			Title:          pageTitle,
			CaptureProfile: &captureProfile,
			DiffProfile:    diffProfiles,
		})
	}
	if len(pages) == 0 {
		pages = make([]models.PageProps, 0)
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "User not found in request context", err.Error())
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
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	workspaces, err := wh.workspaceService.ListUserWorkspaces(ctx, clerkUser)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list workspaces", err.Error())
	}
	return sendDataResponse(c, fiber.StatusOK, "Listed workspaces successfully", workspaces)
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
	if err := wh.workspaceService.UpdateWorkspace(ctx, workspaceID, req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated workspace successfully",
		map[string]any{
			"workspaceId":  workspaceID,
			"billingEmail": req.BillingEmail,
			"name":         req.Name,
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
		"workspaceId": workspaceID,
		"status":      status,
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
		"workspaceId": workspaceID,
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
		"workspaceId": workspaceID,
	})
}
