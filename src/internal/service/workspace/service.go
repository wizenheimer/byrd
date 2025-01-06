package workspace

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

func NewWorkspaceService(workspaceRepo repo.WorkspaceRepository, competitorService svc.CompetitorService, userService svc.UserService, logger *logger.Logger) svc.WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		logger:            logger,
	}
}

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, err.Error) {
	wErr := err.New()
	// Step 1: Create a workspace
	// Generate a workspace name
	workspaceName := utils.GenerateWorkspaceName(workspaceOwner)
	billingEmail, err := utils.GetClerkUserEmail(workspaceOwner)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetClerkUserEmail, map[string]any{"error": err.Error()})
		return nil, wErr
	}

	// Create a workspace using the workspace name and billing email
	workspace, rErr := ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, billingEmail)
	if rErr != nil && rErr.HasErrors() {
		wErr.Merge(rErr)
		return nil, wErr
	}

	// Step 2: Associate the workspace with the clerkUser
	if _, err := ws.userService.CreateWorkspaceOwner(ctx, workspaceOwner, workspace.ID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	// Step 3: Invite users to the workspace
	emailMap := make(map[string]bool)
	emailMap[billingEmail] = true

	var invitedUsers []api.InviteUserToWorkspaceRequest
	for _, user := range workspaceReq.WorkspaceUserCreationRequest {
		// Ensure the owner is not invited again
		if _, ok := emailMap[user.Email]; ok {
			continue
		}

		invitedUsers = append(invitedUsers, api.InviteUserToWorkspaceRequest{
			Email: user.Email,
			Role:  models.UserRoleUser,
		})

		emailMap[user.Email] = true
	}

	// Batch invite users to the workspace
	_, rErr = ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
	if rErr != nil && rErr.HasErrors() {
		wErr.Merge(rErr)
		return nil, wErr
	}

	// Step 4: Create a competitor for the workspace
	// Flatten the competitor request to create a competitor
	for _, pages := range workspaceReq.CompetitorCreationRequest.Pages {
		competitorReq := api.CreateCompetitorRequest{
			Pages: []api.CreatePageRequest{pages},
		}

		_, err := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq)
		if err != nil && err.HasErrors() {
			// TODO: non-fatal error handling
			wErr.Merge(err)
		}
	}

	return &workspace, wErr
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, err.Error) {
	wErr := err.New()
	// List workspaces for a user
	workspaceIDs, err := ws.userService.ListUserWorkspaces(ctx, workspaceMember)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	workspaces, err := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	return workspaces, nil
}

func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, err.Error) {
	wErr := err.New()

	workspaceIDs := []uuid.UUID{workspaceID}
	workspaces, err := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	if len(workspaces) == 0 {
		wErr.Add(svc.ErrWorkspaceDoesntExist, map[string]any{"workspace_id": workspaceID})
		return nil, wErr
	}

	if workspaces[0].Status == models.WorkspaceStatusInactive {
		wErr.Add(svc.ErrWorkspaceDoesntExist, map[string]any{"workspace_id": workspaceID})
		return nil, wErr
	}

	return &workspaces[0], nil
}

func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) err.Error {
	wErr := err.New()
	// Get existing workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		return wErr
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
		return wErr
	}

	return nil
}

func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, err.Error) {
	wErr := err.New()
	// Get existing workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return models.WorkspaceStatusInactive, wErr
	}

	// Handle workspace user deletion
	if err := ws.userService.RemoveWorkspaceUsers(ctx, nil, workspaceID); err != nil && err.HasErrors() {
		// TODO: non-fatal error handling
		wErr.Merge(err)
		return models.WorkspaceStatusInactive, wErr
	}

	// Handle workspace competitor deletion

	if err = ws.competitorService.RemoveCompetitors(ctx, workspace.ID, nil); err != nil && err.HasErrors() {
		// TODO: non-fatal error handling
		wErr.Merge(err)
		return models.WorkspaceStatusInactive, wErr
	}

	// Handle workspace deletion
	if err := ws.workspaceRepo.UpdateWorkspaceStatus(ctx, workspace.ID, models.WorkspaceStatusInactive); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return models.WorkspaceStatusInactive, wErr
	}

	return models.WorkspaceStatusInactive, nil
}

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, err.Error) {
	wErr := err.New()
	_, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	wu, err := ws.userService.ListWorkspaceUsers(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	return wu, nil
}

func (ws *workspaceService) InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, err.Error) {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return []api.CreateWorkspaceUserResponse{}, wErr
	}

	resp, err := ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	return resp, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) err.Error {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	adminUsers, _, err := ws.userService.GetWorkspaceUserCountByRole(ctx, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, workspaceMember, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	if adminUsers == 1 && workspaceUser.Role == models.UserRoleAdmin {
		wErr.Add(svc.ErrCannotLeaveWorkspaceAsOnlyAdmin, map[string]any{"workspace_id": workspace.ID})
		return wErr
	}

	users := []uuid.UUID{workspaceUser.ID}
	if err := ws.userService.RemoveWorkspaceUsers(ctx, users, workspace.ID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) err.Error {
	wErr := err.New()
	if _, err := ws.GetWorkspace(ctx, workspaceID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	if _, err := ws.userService.UpdateWorkspaceUserRole(ctx, workspaceMemberID, workspaceID, role); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) err.Error {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	workspaceUser, err := ws.userService.GetWorkspaceUserByID(ctx, workspaceMemberID, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	if workspaceUser.WorkspaceUserStatus == models.UserWorkspaceStatusInactive {
		wErr.Add(svc.ErrUserNotInvitedToWorkspace, map[string]any{"workspace_id": workspace.ID})
		return wErr
	}

	if err := ws.userService.RemoveWorkspaceUsers(ctx, []uuid.UUID{workspaceMemberID}, workspace.ID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

// JoinWorkspace joins a user using the invite link
func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) err.Error {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	if workspace.Status == models.WorkspaceStatusInactive {
		wErr.Add(svc.ErrWorkspaceDoesNotExist, map[string]any{"workspace_id": workspace.ID})
		return wErr
	}

	// Sync the user to the workspace
	if err := ws.userService.SyncUser(ctx, invitedMember); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, invitedMember, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	// The user was removed from the workspace, requires re-invite
	if workspaceUser.WorkspaceUserStatus == models.UserWorkspaceStatusInactive {

		wErr.Add(svc.ErrUserNotInvitedToWorkspace, map[string]any{"workspace_id": workspace.ID})
		return wErr
	}

	// Activate the user in the workspace
	if err := ws.userService.UpdateWorkspaceUserStatus(ctx, workspaceUser.ID, workspaceID, models.UserWorkspaceStatusActive); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, err.Error) {
	wErr := err.New()
	// TODO: optimize this for quick lookiups
	if _, err := ws.GetWorkspace(ctx, workspaceID); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	return true, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, err.Error) {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	return workspaceUser.Role == models.UserRoleAdmin, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, err.Error) {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, err
	}

	return workspaceUser.Role == models.UserRoleAdmin || workspaceUser.Role == models.UserRoleUser, nil
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, err.Error) {
	wErr := err.New()
	workspaceExists, err := ws.WorkspaceExists(ctx, workspaceID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	if !workspaceExists {
		return false, nil
	}

	return ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, err.Error) {
	wErr := err.New()

	workspaceAndCompetitorExists, err := ws.WorkspaceCompetitorExists(ctx, workspaceID, competitorID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	if !workspaceAndCompetitorExists {
		return false, nil
	}

	pageExists, err := ws.competitorService.PageExists(ctx, competitorID, pageID)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return false, wErr
	}

	return pageExists, nil
}

func (ws *workspaceService) CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) err.Error {
	wErr := err.New()
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		return wErr
	}

	competitorReq := api.CreateCompetitorRequest{
		Pages: []api.CreatePageRequest{page},
	}

	if _, err := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq); err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pageRequest []api.CreatePageRequest) ([]models.Page, err.Error) {
	wErr := err.New()
	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		wErr.Add(err, map[string]any{"uuid": competitorID})
		return nil, wErr
	}

	pages, pErr := ws.competitorService.AddPagesToCompetitor(ctx, competitorUUID, pageRequest)
	if pErr != nil && pErr.HasErrors() {
		wErr.Merge(pErr)
		return nil, wErr
	}

	return pages, nil
}

func (ws *workspaceService) ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) (*api.PaginatedResponse, err.Error) {
	wErr := err.New()

	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		wErr.Merge(err)
		return nil, wErr
	}

	competitorsWithPages, err := ws.competitorService.ListWorkspaceCompetitors(ctx, workspace.ID, params)
	if err != nil && err.HasErrors() {
		return nil, wErr
	}

	return &api.PaginatedResponse{
		Data: competitorsWithPages,
	}, nil
}

func (ws *workspaceService) ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, err.Error) {
	wErr := err.New()
	pageHistory, err := ws.competitorService.GetCompetitorPage(ctx, competitorID, pageID, param)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return nil, wErr
	}

	return pageHistory.History, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) err.Error {
	wErr := err.New()
	pageIDs := []uuid.UUID{pageID}
	err := ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, pageIDs)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) err.Error {
	wErr := err.New()

	err := ws.competitorService.RemoveCompetitors(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	err = ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, nil)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) err.Error {
	wErr := err.New()
	_, err := ws.competitorService.UpdatePage(ctx, competitorID, pageID, req)
	if err != nil && err.HasErrors() {
		wErr.Merge(err)
		return wErr
	}

	return nil
}
