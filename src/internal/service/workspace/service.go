package workspace

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, error) {

	// Create a new workspace
	workspaceName := *workspaceOwner.FirstName + "'s Workspace"
	billingEmail := *workspaceOwner.PrimaryEmailAddressID
	workspace, err := ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, billingEmail)
	if err != nil {
		return nil, err
	}

	// Create workspace owner
	_, err = ws.userService.CreateWorkspaceOwner(ctx, workspaceOwner, workspace.ID)
	if err != nil {
		return nil, err
	}

	// Add other users to the workspace
	emailMap := make(map[string]bool)
	emailMap[*workspaceOwner.PrimaryEmailAddressID] = true

	var invitedUsers []api.InviteUserToWorkspaceRequest
	for _, user := range workspaceReq.WorkspaceUserCreationRequest {
		if _, ok := emailMap[user.Email]; ok {
			continue
		}

		invitedUsers = append(invitedUsers, api.InviteUserToWorkspaceRequest{
			Email: user.Email,
			Role:  models.UserRoleUser,
		})

		emailMap[user.Email] = true
	}

	batchResponse := ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
	for _, response := range batchResponse {
		// TODO: non-fatal error handling
		if response.Error != nil {
			return nil, response.Error
		}
	}

	// Create a competitor for the workspace
	// Flatten the competitor request to create a competitor
	for _, pages := range workspaceReq.CompetitorCreationRequest.Pages {
		competitorReq := api.CreateCompetitorRequest{
			Pages: []api.CreatePageRequest{pages},
		}

		_, errs := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq)
		if len(errs) > 0 {
			// TODO: non-fatal error handling
			return nil, errs[0]
		}
	}

	return &workspace, nil
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error) {
	workspaceIDs, err := ws.userService.ListUserWorkspaces(ctx, workspaceMember)
	if err != nil {
		return nil, err
	}

	workspaces, errs := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if len(errs) > 0 {
		// TODO: non-fatal error handling
		return nil, errs[0]
	}

	return workspaces, nil
}

func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	workspaceIDs := []uuid.UUID{workspaceID}
	workspaces, errs := ws.workspaceRepo.GetWorkspaces(ctx, workspaceIDs)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if len(workspaces) == 0 {
		// return nil, domain.ErrWorkspaceNotFound
		return nil, errors.New("workspace not found")
	}

	return &workspaces[0], nil
}

func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) error {
	// Get existing workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
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

	return ws.workspaceRepo.UpdateWorkspace(ctx, workspaceReq)
}

func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error) {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return models.WorkspaceStatusInactive, err
	}

	// Handle workspace user deletion
	errs := ws.userService.RemoveWorkspaceUsers(ctx, nil, workspaceID)
	if len(errs) > 0 {
		// TODO: non-fatal error handling
		return models.WorkspaceStatusInactive, errs[0]
	}

	// Handle workspace competitor deletion
	errs = ws.competitorService.RemoveCompetitors(ctx, workspace.ID, nil)
	if len(errs) > 0 {
		// TODO: non-fatal error handling
		return models.WorkspaceStatusInactive, errs[0]
	}

	// Handle workspace deletion
	err = ws.workspaceRepo.UpdateWorkspaceStatus(ctx, workspace.ID, models.WorkspaceStatusInactive)
	if err != nil {
		return models.WorkspaceStatusInactive, err
	}

	return models.WorkspaceStatusInactive, nil
}

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, error) {
	_, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	return ws.userService.ListWorkspaceUsers(ctx, workspaceID)
}

func (ws *workspaceService) InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) []api.CreateWorkspaceUserResponse {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return []api.CreateWorkspaceUserResponse{
			{
				User:  nil,
				Error: err,
			},
		}
	}

	return ws.userService.AddUserToWorkspace(ctx, workspace.ID, invitedUsers)
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	adminUsers, _, err := ws.userService.GetWorkspaceUserCountByRole(ctx, workspace.ID)
	if err != nil {
		return err
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, workspaceMember, workspace.ID)
	if err != nil {
		return err
	}

	if adminUsers == 1 && workspaceUser.Role == models.UserRoleAdmin {
		return errors.New("cannot leave workspace as the only admin")
	}

	users := []uuid.UUID{workspace.ID}
	errs := ws.userService.RemoveWorkspaceUsers(ctx, users, workspace.ID)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) error {
	if _, err := ws.GetWorkspace(ctx, workspaceID); err != nil {
		return err
	}

	if _, err := ws.userService.UpdateWorkspaceUserRole(ctx, workspaceMemberID, workspaceID, role); err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	workspaceUser, err := ws.userService.GetWorkspaceUserByID(ctx, workspaceMemberID, workspaceMemberID)
	if err != nil {
		return err
	}

	errs := ws.userService.RemoveWorkspaceUsers(ctx, []uuid.UUID{workspaceUser.ID}, workspace.ID)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// JoinWorkspace joins a user using the invite link
func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	if workspace.Status == models.WorkspaceStatusInactive {
		return errors.New("workspace is inactive")
	}

	// Sync the user to the workspace
	if err := ws.userService.SyncUser(ctx, invitedMember); err != nil {
		return err
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, invitedMember, workspaceID)

	if err != nil {
		return err
	}

	errs := ws.userService.AddWorkspaceUsers(ctx, []uuid.UUID{workspaceUser.ID}, workspace.ID)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	// TODO: optimize this for quick lookiups
	_, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil {
		return false, err
	}

	return workspaceUser.Role == models.UserRoleAdmin, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.userService.GetWorkspaceUser(ctx, clerkUser, workspace.ID)
	if err != nil {
		return false, err
	}

	return workspaceUser.Role == models.UserRoleAdmin || workspaceUser.Role == models.UserRoleUser, nil
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	workspaceExists, err := ws.WorkspaceExists(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	if !workspaceExists {
		return false, nil
	}

	return ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error) {
	workspaceAndCompetitorExists, err := ws.WorkspaceCompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	if !workspaceAndCompetitorExists {
		return false, nil
	}

	return ws.competitorService.PageExists(ctx, competitorID, pageID)
}

func (ws *workspaceService) CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) []error {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return []error{err}
	}

	competitorReq := api.CreateCompetitorRequest{
		Pages: []api.CreatePageRequest{page},
	}

	_, errs := ws.competitorService.CreateCompetitor(ctx, workspace.ID, competitorReq)
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pageRequest []api.CreatePageRequest) ([]models.Page, []error) {
	competitorUUID, err := uuid.Parse(competitorID)
	if err != nil {
		return nil, []error{errors.New("couldn't parse competitor uuid")}
	}

	return ws.competitorService.AddPagesToCompetitor(ctx, competitorUUID, pageRequest)
}

func (ws *workspaceService) ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) (*api.PaginatedResponse, error) {
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	competitorsWithPages, err := ws.competitorService.ListWorkspaceCompetitors(ctx, workspace.ID, params)
	if err != nil {
		return nil, err
	}

	return &api.PaginatedResponse{
		Data: competitorsWithPages,
	}, nil
}

func (ws *workspaceService) ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error) {
	pageHistory, err := ws.competitorService.GetCompetitorPage(ctx, competitorID, pageID, param)
	if err != nil {
		return nil, err
	}

	return pageHistory.History, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) error {
	pageIDs := []uuid.UUID{pageID}
	errs := ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, pageIDs)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) error {
	errs := ws.competitorService.RemoveCompetitors(ctx, workspaceID, []uuid.UUID{competitorID})
	if len(errs) > 0 {
		return errs[0]
	}

	errs = ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, nil)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) error {
	_, err := ws.competitorService.UpdatePage(ctx, competitorID, pageID, req)

	return err
}
