package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/api/auth"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
	"go.uber.org/zap"
)

type WorkspaceHandler struct {
	workspaceService svc.WorkspaceService
	logger           *logger.Logger
}

func NewWorkspaceHandler(
	workspaceService svc.WorkspaceService,
	logger *logger.Logger,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: workspaceService,
		logger:           logger,
	}
}

// CreateWorkspace creates a new workspace
func (wh *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	var req api.WorkspaceCreationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	workspace, err := wh.workspaceService.CreateWorkspace(c.Context(), clerkUser, req)
	if err != nil {
		wh.logger.Error("Failed to create workspace", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create workspace",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workspace)

}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	workspaces, err := wh.workspaceService.ListUserWorkspaces(c.Context(), clerkUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list workspaces",
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(workspaces)
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	workspace, err := wh.workspaceService.GetWorkspace(c.Context(), workspaceID)
	if err != nil {
		// if errors.Is(err, svc.ErrWorkspaceNotFound) {
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		// }
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(workspace)
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var req api.WorkspaceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := utils.SetDefaultsAndValidate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := wh.workspaceService.UpdateWorkspace(c.Context(), workspaceID, req); err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		// case errors.Is(err, svc.ErrWorkspaceInactive):
		// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		// 		Error:   "WorkspaceInactive",
		// 		Code:    fiber.StatusBadRequest,
		// 		Message: "Workspace is inactive",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to update workspace",
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteWorkspace deletes a workspace by ID
func (wh *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	status, err := wh.workspaceService.DeleteWorkspace(c.Context(), workspaceID)
	if err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		// case errors.Is(err, svc.ErrWorkspaceInactive):
		// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		// 		Error:   "WorkspaceInactive",
		// 		Code:    fiber.StatusBadRequest,
		// 		Message: "Workspace is already inactive",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to delete workspace",
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := wh.workspaceService.LeaveWorkspace(c.Context(), clerkUser, workspaceID); err != nil {
		switch {
		// case errors.Is(err, svc.ErrLastAdmin):
		// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		// 		Error:   "LastAdmin",
		// 		Code:    fiber.StatusBadRequest,
		// 		Message: "Cannot leave workspace as last admin",
		// 	})
		// case errors.Is(err, svc.ErrWorkspaceUserNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "UserNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "User not found in workspace",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := wh.workspaceService.JoinWorkspace(c.Context(), clerkUser, workspaceID); err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		// case errors.Is(err, svc.ErrNotInvited):
		// 	return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
		// 		Error:   "NotInvited",
		// 		Code:    fiber.StatusForbidden,
		// 		Message: "User not invited to workspace",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to join workspace",
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	params := api.WorkspaceMembersListingParams{
		IncludeMembers: c.QueryBool("includeMembers", true),
		IncludeAdmins:  c.QueryBool("includeAdmins", true),
	}

	users, err := wh.workspaceService.ListWorkspaceMembers(c.Context(), workspaceID, params)
	if err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to list workspace users",
			})
		}
	}

	return c.JSON(users)
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var reqs []api.InviteUserToWorkspaceRequest
	if err := c.BodyParser(&reqs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	for index := range reqs {
		if err := utils.SetDefaultsAndValidate(&reqs[index]); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "InvalidRequest",
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	responses := wh.workspaceService.InviteUsersToWorkspace(c.Context(), clerkUser, workspaceID, reqs)
	return c.JSON(responses)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidUserID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	err = wh.workspaceService.RemoveUserFromWorkspace(c.Context(), workspaceID, userID)
	if err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		//     return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		//         Error:   "WorkspaceNotFound",
		//         Code:    fiber.StatusNotFound,
		//         Message: "Workspace not found",
		//     })
		// case errors.Is(err, svc.ErrWorkspaceUserNotFound):
		//     return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		//         Error:   "UserNotFound",
		//         Code:    fiber.StatusNotFound,
		//         Message: "User not found in workspace",
		//     })
		// case errors.Is(err, svc.ErrLastAdmin):
		//     return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		//         Error:   "LastAdmin",
		//         Code:    fiber.StatusBadRequest,
		//         Message: "Cannot remove last admin from workspace",
		//     })
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidUserID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	var req struct {
		Role models.UserWorkspaceRole `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate role
	if !isValidRole(req.Role) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRole",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid role specified",
		})
	}

	err = wh.workspaceService.UpdateWorkspaceMemberRole(c.Context(), workspaceID, userID, req.Role)
	if err != nil {
		switch {
		// case errors.Is(err, svc.ErrWorkspaceNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "WorkspaceNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "Workspace not found",
		// 	})
		// case errors.Is(err, svc.ErrWorkspaceUserNotFound):
		// 	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		// 		Error:   "UserNotFound",
		// 		Code:    fiber.StatusNotFound,
		// 		Message: "User not found in workspace",
		// 	})
		// case errors.Is(err, svc.ErrLastAdmin):
		// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		// 		Error:   "LastAdmin",
		// 		Code:    fiber.StatusBadRequest,
		// 		Message: "Cannot change role of last admin",
		// 	})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "InternalServerError",
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// Helper function to validate workspace roles
func isValidRole(role models.UserWorkspaceRole) bool {
	switch role {
	case models.UserRoleAdmin, models.UserRoleUser, models.UserRoleViewer:
		return true
	default:
		return false
	}
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	var req api.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	errs := wh.workspaceService.CreateWorkspaceCompetitor(c.Context(), clerkUser, workspaceID, req)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "CompetitorCreationFailed",
			Code:    fiber.StatusBadRequest,
			Message: "Failed to create competitor: " + errs[0].Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	var pages []api.CreatePageRequest
	if err := c.BodyParser(&pages); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	createdPages, errs := wh.workspaceService.AddPageToCompetitor(c.Context(), clerkUser, competitorID.String(), pages)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "PageCreationFailed",
			Code:    fiber.StatusBadRequest,
			Message: "Failed to add pages: " + errs[0].Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPages)
}

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	pageNumber := offset/limit + 1
	pageSize := limit

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	response, err := wh.workspaceService.ListWorkspaceCompetitors(c.Context(), clerkUser, workspaceID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list competitors",
		})
	}

	return c.JSON(response)
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	pageNumber := offset/limit + 1
	pageSize := limit

	params := api.PaginationParams{
		Page:     pageNumber,
		PageSize: pageSize,
	}

	history, err := wh.workspaceService.ListWorkspacePageHistory(c.Context(), clerkUser, workspaceID, competitorID, pageID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list page history",
		})
	}

	return c.JSON(history)
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := wh.workspaceService.RemovePageFromWorkspace(c.Context(), clerkUser, competitorID, pageID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove page",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	workspaceID, err := uuid.Parse(c.Params("workspaceID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidWorkspaceID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid workspace ID format",
		})
	}

	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	clerkUser, err := auth.GetClerkUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := wh.workspaceService.RemoveCompetitorFromWorkspace(c.Context(), clerkUser, workspaceID, competitorID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove competitor",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
	competitorID, err := uuid.Parse(c.Params("competitorID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidCompetitorID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid competitor ID format",
		})
	}

	pageID, err := uuid.Parse(c.Params("pageID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidPageID",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid page ID format",
		})
	}

	var req api.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "InvalidRequest",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := wh.workspaceService.UpdateCompetitorPage(c.Context(), competitorID, pageID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "InternalServerError",
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update page",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
