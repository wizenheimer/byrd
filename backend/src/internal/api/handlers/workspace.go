// ./src/internal/api/handlers/workspace.go
package handlers

import (
	"net/url"
	"strings"

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
func (wh *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
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
		captureProfile := screenshot.GetDefaultScreenshotRequestOptions(competitorURL)
		pages = append(pages, models.PageProps{
			URL:            competitorURL,
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
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
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
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
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
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
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
func (wh *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
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

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
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

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
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

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	pageNumber := max(1, c.QueryInt("pageNumber", 1))
	pageSize := max(10, c.QueryInt("pageSize", 10))

	pagination := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limits := pagination.GetLimit()
	offsets := pagination.GetOffset()

	roleFilterString := strings.ToLower(c.Query("role", ""))

	var roleFilter *models.WorkspaceRole
	switch roleFilterString {
	case "admin":
		adminRole := models.RoleAdmin
		roleFilter = &adminRole
	case "user":
		userRole := models.RoleUser
		roleFilter = &userRole
	default:
		wh.logger.Debug("Invalid role filter, defaulting to user")
		roleFilter = nil
	}

	ctx := c.Context()
	users, hasMore, err := wh.workspaceService.ListWorkspaceMembers(ctx, workspaceID, &limits, &offsets, roleFilter)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list workspace users", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace users successfully", map[string]any{
		"users":   users,
		"hasMore": hasMore,
	})
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var users api.AddUsersToWorkspaceRequest
	if err := c.BodyParser(&users); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Unprocessable request body", err.Error())
	}

	err = utils.SetDefaultsAndValidate(&users)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := getClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()

	responses, err := wh.workspaceService.AddUsersToWorkspace(ctx, clerkUser, workspaceID, users.Emails)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not invite users to workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Invited users to workspace successfully", responses)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	ctx := c.Context()
	if err := wh.workspaceService.RemoveUserFromWorkspace(ctx, workspaceID, userID); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not remove user from workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Removed user from workspace successfully", map[string]any{
		"userId":      userID,
		"workspaceId": workspaceID,
	})
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	var req api.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()

	if err := wh.workspaceService.UpdateWorkspaceMemberRole(ctx, workspaceID, userID, req.Role); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update user role in workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated user role in workspace successfully", map[string]any{
		"role": req.Role,
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

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var pages []api.CreatePageRequest
	if err := c.BodyParser(&pages); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidateArray(&pages); err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	createdPages, err := wh.workspaceService.AddPageToCompetitor(ctx, competitorID, pages)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not add page to competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Added page to competitor successfully", createdPages)
}

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

	pageNumber := max(1, c.QueryInt("pageNumber", 1))
	pageSize := max(10, c.QueryInt("pageSize", 10))

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

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	pageNumber := max(1, c.QueryInt("pageNumber", 1))
	pageSize := max(10, c.QueryInt("pageSize", 10))

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	limit := params.GetLimit()
	offset := params.GetOffset()

	ctx := c.Context()
	competitors, hasMore, err := wh.workspaceService.ListCompetitorsForWorkspace(ctx, workspaceID, &limit, &offset)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not list workspace competitors", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace competitors successfully", map[string]any{
		"competitors": competitors,
		"hasMore":     hasMore,
	})
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

	pageNumber := max(1, c.QueryInt("pageNumber", 1))
	pageSize := max(10, c.QueryInt("pageSize", 10))

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

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
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

	ctx := c.Context()
	page, err := wh.workspaceService.UpdateCompetitorPage(ctx, competitorID, pageID, req)
	if err != nil {
		return sendErrorResponse(c, wh.logger, fiber.StatusInternalServerError, "Could not update page in competitor", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Updated page in competitor successfully", page)
}
