// ./src/internal/api/handlers/workspace.go
package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/api/auth"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type WorkspaceHandler struct {
	workspaceService svc.WorkspaceService
	logger           *logger.Logger
	tx               *transaction.TxManager
}

func NewWorkspaceHandler(
	workspaceService svc.WorkspaceService,
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
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidRequest", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidRequest", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	var workspace *models.Workspace
	var e errs.Error
	ctx := c.Context()
	err = wh.tx.RunInTx(ctx, nil, func(ctx context.Context) error {
		workspace, e = wh.workspaceService.CreateWorkspace(ctx, clerkUser, req)
		if e != nil && e.HasErrors() {
			return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create workspace", e)
		}
		return nil
	})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create workspace", err.Error())
	}

	return sendDataResponse(c, fiber.StatusCreated, "Created workspace successfully", workspace)

}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var workspaces []models.Workspace
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	workspaces, e = wh.workspaceService.ListUserWorkspaces(ctx, clerkUser)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspaces", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspaces", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Listed workspaces successfully", workspaces)
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	ctx := c.Context()
	var workspace *models.Workspace
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	workspace, e = wh.workspaceService.GetWorkspace(ctx, workspaceID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get workspace", e)
	}
	// return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not get workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Fetched workspace successfully", workspace)
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req api.WorkspaceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.UpdateWorkspace(ctx, workspaceID, req); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update workspace", err.Error())
	// }

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
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	ctx := c.Context()
	var status models.WorkspaceStatus
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	status, e = wh.workspaceService.DeleteWorkspace(ctx, workspaceID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not delete workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not delete workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Deleted workspace successfully", map[string]any{
		"workspaceId": workspaceID,
		"status":      status,
	})
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.LeaveWorkspace(ctx, clerkUser, workspaceID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not exit workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not exit workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Exited workspace successfully", map[string]any{
		"workspaceId": workspaceID,
	})
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.JoinWorkspace(ctx, clerkUser, workspaceID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not join workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not join workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Joined workspace successfully", map[string]any{
		"workspaceId": workspaceID,
	})
}

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	params := api.WorkspaceMembersListingParams{
		IncludeMembers: c.QueryBool("includeMembers", true),
		IncludeAdmins:  c.QueryBool("includeAdmins", true),
	}

	ctx := c.Context()
	var users []models.WorkspaceUser
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	users, e = wh.workspaceService.ListWorkspaceMembers(ctx, workspaceID, params)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace users", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace users", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace users successfully", users)
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var reqs []api.InviteUserToWorkspaceRequest
	if err := c.BodyParser(&reqs); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	err = utils.SetDefaultsAndValidateArray(&reqs)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var responses []api.CreateWorkspaceUserResponse
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	responses, e = wh.workspaceService.InviteUsersToWorkspace(ctx, clerkUser, workspaceID, reqs)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not invite users to workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not invite users to workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusCreated, "Invited users to workspace successfully", responses)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	ctx := c.Context()
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	e := wh.workspaceService.RemoveUserFromWorkspace(ctx, workspaceID, userID)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove user from workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove user from workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Removed user from workspace successfully", map[string]any{
		"userId":      userID,
		"workspaceId": workspaceID,
	})
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
	}

	var req api.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	e := wh.workspaceService.UpdateWorkspaceMemberRole(ctx, workspaceID, userID, req.Role)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update user role in workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update user role in workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Updated user role in workspace successfully", map[string]any{
		"role": req.Role,
	})
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	var req api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	var e errs.Error
	ctx := c.Context()
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	e = wh.workspaceService.CreateWorkspaceCompetitor(ctx, clerkUser, workspaceID, req)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create competitor", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not create competitor", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusCreated, "Created competitor successfully", map[string]any{
		"url": req.URL,
	})
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	var pages []api.CreatePageRequest
	if err := c.BodyParser(&pages); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidateArray(&pages); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	var createdPages []models.Page
	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	createdPages, e = wh.workspaceService.AddPageToCompetitor(ctx, clerkUser, competitorID.String(), pages)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not add page to competitor", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not add page to competitor", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusCreated, "Added page to competitor successfully", createdPages)
}

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	pageNumber := offset/limit + 1
	pageSize := limit

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	ctx := c.Context()
	var response []api.GetWorkspaceCompetitorResponse
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	response, e = wh.workspaceService.ListWorkspaceCompetitors(ctx, clerkUser, workspaceID, params)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace competitors", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list workspace competitors", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Listed workspace competitors successfully", response)
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid workspace ID format", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid competitor ID format", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid page ID format", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	pageNumber := c.QueryInt("pageNumber", 10)
	pageSize := c.QueryInt("pageSize", 0)

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	ctx := c.Context()
	var history []models.PageHistory
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	history, e = wh.workspaceService.ListWorkspacePageHistory(ctx, clerkUser, workspaceID, competitorID, pageID, params)
	if e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list page history", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not list page history", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Listed page history successfully", history)
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.RemovePageFromWorkspace(ctx, clerkUser, competitorID, pageID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove page from competitor", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove page from competitor", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Removed page from competitor successfully", nil)
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidWorkspaceID", err.Error())
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.RemoveCompetitorFromWorkspace(ctx, clerkUser, workspaceID, competitorID); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove competitor from workspace", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not remove competitor from workspace", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Removed competitor from workspace successfully", nil)
}

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidCompetitorID", err.Error())
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "InvalidPageID", err.Error())
	}

	var req api.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx := c.Context()
	var e errs.Error
	// err = wh.tx.RunInTx(c.Context(), nil, func(ctx context.Context) error {
	if e = wh.workspaceService.UpdateCompetitorPage(ctx, competitorID, pageID, req); e != nil && e.HasErrors() {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update page in competitor", e)
	}
	// 	return nil
	// })
	// if err != nil {
	// 	return sendErrorResponse(c, fiber.StatusInternalServerError, "Could not update page in competitor", err.Error())
	// }

	return sendDataResponse(c, fiber.StatusOK, "Updated page in competitor successfully", nil)
}
