// ./src/internal/service/workspace/service.go
package workspace

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

func NewWorkspaceService(workspaceRepo repo.WorkspaceRepository, competitorService svc.CompetitorService, userService svc.UserService, logger *logger.Logger) svc.WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		logger:            logger,
	}
}

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, errs.Error) {
	wErr := errs.New()
	// Step 1: Create a workspace
	// Generate a workspace name
	workspaceName := utils.GenerateWorkspaceName(workspaceOwner)
	billingEmail, err := utils.GetClerkUserEmail(workspaceOwner)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetClerkUserEmail, map[string]any{"error": err.Error()})
		wErr.Log("Failed to get clerk user email", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspace)
	}

	// Create a workspace using the workspace name and billing email
	workspace, rErr := ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, billingEmail)
	if rErr != nil && rErr.HasErrors() {
		wErr.Merge(rErr)
		wErr.Log("Failed to create workspace", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspace)
	}

	// Step 2: Invite users to the workspace
	emailMap := make(map[string]bool)
	emailMap[billingEmail] = true

	var invitedUsers []api.InviteUserToWorkspaceRequest
	for _, user := range workspaceReq.WorkspaceUserCreationRequest {
		// Ensure that the user is not already invited
		if _, ok := emailMap[user.Email]; ok {
			continue
		}

		invitedUsers = append(invitedUsers, api.InviteUserToWorkspaceRequest{
			Email:  user.Email,
			Role:   models.UserRoleUser,
			Status: models.UserWorkspaceStatusPending,
		})

		emailMap[user.Email] = true
	}

	// Add the workspace owner to the list of invited users
	invitedUsers = append(invitedUsers, api.InviteUserToWorkspaceRequest{
		Email:  billingEmail,
		Role:   models.UserRoleAdmin,
		Status: models.UserWorkspaceStatusActive,
	})

	ws.logger.Debug("Inviting users to workspace", zap.Any("invited_users", invitedUsers))

	// Batch invite users to the workspace
	_, rErr = ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
	if rErr != nil && rErr.HasErrors() {
		wErr.Merge(rErr)
		wErr.Log("Failed to invite users to workspace", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspace)
	}

	// Step 3: Create a competitor for the workspace
	// Flatten the competitor request to create a competitor
	for _, pages := range workspaceReq.CompetitorCreationRequest.Pages {
		competitorReq := api.CreateCompetitorRequest{
			Pages: []api.CreatePageRequest{pages},
		}

		_, err := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq)
		if err != nil && err.HasErrors() {
			err.Log("Failed to create competitor", ws.logger)
			wErr.Merge(err)
		}
	}

	return &workspace, wErr.Propagate(svc.ErrFailedToCreateWorkspace, nil)
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, errs.Error) {
	wErr := errs.New()
	// List workspaces for a user
	workspaceIDs, err := ws.userService.ListUserWorkspaces(ctx, workspaceMember)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to list user workspaces", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListUserWorkspaces)
	}

	workspaces, err := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspaces", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListUserWorkspaces)
	}

	return workspaces, nil
}

func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, errs.Error) {
	wErr := errs.New()

	workspaceIDs := []uuid.UUID{workspaceID}
	workspaces, err := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspaces", ws.logger)
		return nil, wErr.Propagate(svc.ErrWorkspaceNotFound)
	}

	if len(workspaces) == 0 {
		wErr.Add(svc.ErrWorkspaceLookupEmpty, map[string]any{"workspace_id": workspaceID})
		return nil, wErr.Propagate(svc.ErrWorkspaceNotFound)
	}

	if workspaces[0].Status == models.WorkspaceStatusInactive {
		wErr.Add(svc.ErrWorkspaceInactive, map[string]any{"workspace_id": workspaceID})
		return nil, wErr.Propagate(svc.ErrWorkspaceNotFound)
	}

	return &workspaces[0], nil
}

func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) errs.Error {
	wErr := errs.New()
	// Get existing workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		return wErr.Propagate(svc.ErrFailedToUpdateWorkspace)
	}

	// Check if workspace requires an update
	if req.BillingEmail == workspace.BillingEmail && req.Name == workspace.Name {
		return nil
	}

	if req.Name == "" {
		req.Name = workspace.Name
	}

	if req.BillingEmail == "" {
		req.BillingEmail = workspace.BillingEmail
	}

	workspaceReq := api.WorkspaceUpdateRequest{
		Name:         req.Name,
		BillingEmail: req.BillingEmail,
	}

	if err := ws.workspaceRepo.UpdateWorkspace(ctx, workspaceID, workspaceReq); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to update workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToUpdateWorkspace)
	}

	return nil
}

func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, errs.Error) {
	wErr := errs.New()
	// Get existing workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return models.WorkspaceStatusInactive, wErr.Propagate(svc.ErrFailedToDeleteWorkspace)
	}

	// Handle workspace user deletion
	if err := ws.userService.RemoveWorkspaceUsers(ctx, nil, workspaceID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove workspace users", ws.logger)
		return models.WorkspaceStatusInactive, wErr.Propagate(svc.ErrFailedToDeleteWorkspace)
	}

	// Handle workspace competitor deletion

	if err = ws.competitorService.RemoveCompetitors(ctx, workspace.ID, nil); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove workspace competitors", ws.logger)
		return models.WorkspaceStatusInactive, wErr.Propagate(svc.ErrFailedToDeleteWorkspace)
	}

	// Handle workspace deletion
	if err := ws.workspaceRepo.RemoveWorkspaces(ctx, []uuid.UUID{workspace.ID}); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove workspace", ws.logger)
		return models.WorkspaceStatusInactive, wErr.Propagate(svc.ErrFailedToDeleteWorkspace)
	}

	return models.WorkspaceStatusInactive, nil
}

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, errs.Error) {
	wErr := errs.New()
	_, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspaceMembers)
	}

	wu, err := ws.userService.ListWorkspaceUsers(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to list workspace users", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspaceMembers)
	}

	return wu, nil
}

func (ws *workspaceService) InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, errs.Error) {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return []api.CreateWorkspaceUserResponse{}, wErr.Propagate(svc.ErrFailedToInviteUserToWorkspace)
	}

	resp, err := ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to invite users to workspace", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToInviteUserToWorkspace)
	}

	return resp, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) errs.Error {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToLeaveWorkspace)
	}

	adminUsers, _, err := ws.userService.GetWorkspaceUserCountByRole(ctx, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user count by role", ws.logger)
		return wErr.Propagate(svc.ErrFailedToLeaveWorkspace)
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, workspaceMember, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user", ws.logger)
		return wErr.Propagate(svc.ErrFailedToLeaveWorkspace)
	}

	if adminUsers == 1 && workspaceUser.Role == models.UserRoleAdmin {
		wErr.Add(svc.ErrFailedToLeaveWorkspaceAsSoleAdmin, map[string]any{"workspace_id": workspace.ID})
		wErr.Log("Failed to leave workspace as sole admin", ws.logger)
		return wErr.Propagate(svc.ErrFailedToLeaveWorkspace)
	}

	users := []uuid.UUID{workspaceUser.ID}
	if err := ws.userService.RemoveWorkspaceUsers(ctx, users, workspace.ID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove workspace users", ws.logger)
		return wErr.Propagate(svc.ErrFailedToLeaveWorkspace)
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) errs.Error {
	wErr := errs.New()
	if _, err := ws.GetWorkspace(ctx, workspaceID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToUpdateWorkspaceMemberRole)
	}

	if _, err := ws.userService.UpdateWorkspaceUserRole(ctx, workspaceMemberID, workspaceID, role); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to update workspace user role", ws.logger)
		return wErr.Propagate(svc.ErrFailedToUpdateWorkspaceMemberRole)
	}

	return nil
}

func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) errs.Error {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveUserFromWorkspace)
	}

	workspaceUser, err := ws.userService.GetWorkspaceUserByID(ctx, workspaceMemberID, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveUserFromWorkspace)
	}

	if workspaceUser.WorkspaceUserStatus == models.UserWorkspaceStatusInactive {
		wErr.Add(svc.ErrUserLeftWorkspaceNeedReinvite, map[string]any{"workspace_id": workspace.ID})
		wErr.Log("User left workspace need reinvite", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveUserFromWorkspace)
	}

	if err := ws.userService.RemoveWorkspaceUsers(ctx, []uuid.UUID{workspaceMemberID}, workspace.ID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove workspace users", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveUserFromWorkspace)
	}

	return nil
}

// JoinWorkspace joins a user using the invite link
func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) errs.Error {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	if workspace.Status == models.WorkspaceStatusInactive {
		wErr.Add(svc.ErrWorkspaceInactive, map[string]any{"workspace_id": workspace.ID})
		wErr.Log("Workspace is inactive", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	// Sync the user to the workspace
	if err := ws.userService.SyncUser(ctx, invitedMember); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to sync user", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, invitedMember, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	// The user was removed from the workspace, requires re-invite
	if workspaceUser.WorkspaceUserStatus == models.UserWorkspaceStatusInactive {
		wErr.Add(svc.ErrUserLeftWorkspaceNeedReinvite, map[string]any{"workspace_id": workspace.ID})
		wErr.Log("User left workspace need reinvite", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	// Activate the user in the workspace
	if err := ws.userService.UpdateWorkspaceUserStatus(ctx, workspaceUser.ID, workspaceID, models.UserWorkspaceStatusActive); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to update workspace user status", ws.logger)
		return wErr.Propagate(svc.ErrFailedToJoinWorkspace)
	}

	return nil
}

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, errs.Error) {
	wErr := errs.New()
	if _, err := ws.GetWorkspace(ctx, workspaceID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceNotFound)
	}

	return true, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, errs.Error) {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceAdminNotFound)
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceAdminNotFound)
	}

	return workspaceUser.Role == models.UserRoleAdmin, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, errs.Error) {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceMemberNotFound)
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace user", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceMemberNotFound)
	}

	return workspaceUser.Role == models.UserRoleAdmin || workspaceUser.Role == models.UserRoleUser, nil
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, errs.Error) {
	wErr := errs.New()
	workspaceExists, err := ws.WorkspaceExists(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to check if workspace exists", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceCompetitorNotFound)
	}

	if !workspaceExists {
		return false, nil
	}

	exists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to check if competitor exists", ws.logger)
		return false, wErr.Propagate(svc.ErrWorkspaceCompetitorNotFound)
	}

	return exists, nil
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, errs.Error) {
	wErr := errs.New()

	workspaceAndCompetitorExists, err := ws.WorkspaceCompetitorExists(ctx, workspaceID, competitorID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to check if workspace competitor exists", ws.logger)
		return false, wErr.Propagate(svc.ErrCompetitorPageNotFound)
	}

	if !workspaceAndCompetitorExists {
		return false, nil
	}

	pageExists, err := ws.competitorService.PageExists(ctx, competitorID, pageID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to check if page exists", ws.logger)
		return false, wErr.Propagate(svc.ErrCompetitorPageNotFound)
	}

	return pageExists, nil
}

func (ws *workspaceService) CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) errs.Error {
	wErr := errs.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return wErr.Propagate(svc.ErrFailedToCreateWorkspaceCompetitor)
	}

	competitorReq := api.CreateCompetitorRequest{
		Pages: []api.CreatePageRequest{page},
	}

	if _, err := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq); err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to create competitor", ws.logger)
		return wErr.Propagate(svc.ErrFailedToCreateWorkspaceCompetitor)
	}

	return nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pageRequest []api.CreatePageRequest) ([]models.Page, errs.Error) {
	wErr := errs.New()
	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		wErr.Add(err, map[string]any{"uuid": competitorID})
		wErr.Log("Failed to parse competitor id", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToAddPageToCompetitor)
	}

	pages, pErr := ws.competitorService.AddPagesToCompetitor(ctx, competitorUUID, pageRequest)
	if pErr != nil && pErr.HasErrors() {
		wErr.Merge(pErr)
		wErr.Log("Failed to add pages to competitor", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToAddPageToCompetitor)
	}

	return pages, nil
}

func (ws *workspaceService) ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, errs.Error) {
	wErr := errs.New()

	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		wErr.Log("Failed to get workspace", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspaceCompetitors)
	}

	competitorsWithPages, err := ws.competitorService.ListWorkspaceCompetitors(ctx, workspace.ID, params)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to list workspace competitors", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspaceCompetitors)
	}

	return competitorsWithPages, nil
}

func (ws *workspaceService) ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, errs.Error) {
	wErr := errs.New()
	pageHistory, err := ws.competitorService.GetCompetitorPage(ctx, competitorID, pageID, param)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to get competitor page", ws.logger)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspacePageHistory)
	}

	return pageHistory.History, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) errs.Error {
	wErr := errs.New()
	pageIDs := []uuid.UUID{pageID}
	err := ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, pageIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove pages from competitor", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemovePageFromWorkspace)
	}

	return nil
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) errs.Error {
	wErr := errs.New()

	err := ws.competitorService.RemoveCompetitors(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove competitors", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveCompetitorFromWorkspace)
	}

	err = ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, nil)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to remove pages from competitor", ws.logger)
		return wErr.Propagate(svc.ErrFailedToRemoveCompetitorFromWorkspace)
	}

	return nil
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) errs.Error {
	wErr := errs.New()
	_, err := ws.competitorService.UpdatePage(ctx, competitorID, pageID, req)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		wErr.Log("Failed to update competitor page", ws.logger)
		return wErr.Propagate(svc.ErrFailedToUpdateCompetitorPage)
	}

	return nil
}
