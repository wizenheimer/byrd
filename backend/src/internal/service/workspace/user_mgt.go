package workspace

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

func (ws *workspaceService) ListWorkspacesForUser(ctx context.Context, workspaceUserEmail string, membershipStatus *models.MembershipStatus, limit, offset *int) ([]models.WorkspaceWithMembership, bool, error) {
	user, err := ws.userService.GetUserByEmail(ctx, workspaceUserEmail)
	if err != nil {
		return nil, false, err
	}

	workspaces, hasMore, err := ws.workspaceRepo.ListWorkspacesForUser(ctx, user.ID, membershipStatus, limit, offset)
	if err != nil {
		return nil, false, err
	}

	return workspaces, hasMore, nil
}

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, roleFilter *models.WorkspaceRole) ([]models.WorkspaceUser, bool, error) {
	members, hasMore, err := ws.workspaceRepo.ListWorkspaceMembers(ctx, workspaceID, limit, offset, roleFilter)
	if err != nil {
		return nil, false, err
	}

	membersIDs := make([]uuid.UUID, 0)
	for _, member := range members {
		membersIDs = append(membersIDs, member.ID)
	}

	users, err := ws.userService.ListUsersByUserIDs(ctx, membersIDs)
	if err != nil {
		return nil, false, err
	}
	userIDToUserMap := make(map[uuid.UUID]models.User)
	for _, user := range users {
		userIDToUserMap[user.ID] = user
	}

	workspaceUsers := make([]models.WorkspaceUser, 0)
	for _, member := range members {
		user, ok := userIDToUserMap[member.ID]
		if !ok {
			continue
		}
		workspaceUser := models.WorkspaceUser{
			ID:               member.ID,
			WorkspaceID:      workspaceID,
			Role:             member.Role,
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, hasMore, nil
}

func (ws *workspaceService) AddUsersToWorkspace(ctx context.Context, workspaceInviterEmail string, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error) {
	// Step 1: Normalize the emails
	normalizedEmails := utils.CleanEmailList(emails, []string{workspaceInviterEmail})

	// Step 2: Get or create users
	users, err := ws.userService.BatchGetOrCreateUsers(ctx, normalizedEmails)
	if err != nil {
		return nil, err
	}

	// Step 3: Add users to the workspace
	userIDs := make([]uuid.UUID, 0)
	userIDsToUserMap := make(map[uuid.UUID]models.User)
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
		userIDsToUserMap[user.ID] = user
	}

	partialWorkspaceUsers, err := ws.workspaceRepo.BatchAddUsersToWorkspace(ctx, workspaceID, userIDs)
	if err != nil {
		return nil, err
	}

	workspaceUsers := make([]models.WorkspaceUser, 0)
	workspaceUsersEmail := make([]string, 0)
	for _, member := range partialWorkspaceUsers {
		user, ok := userIDsToUserMap[member.ID]
		if !ok {
			continue
		}
		if user.Email != nil {
			workspaceUsersEmail = append(workspaceUsersEmail, *user.Email)
		}
		workspaceUser := models.WorkspaceUser{
			ID:               member.ID,
			WorkspaceID:      workspaceID,
			Role:             member.Role,
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	go func() {
		emailTemplate, err := ws.library.GetTemplate(template.WorkspaceInvitePendingTemplate)
		if err != nil {
			ws.logger.Error("couldn't get email template", zap.Error(err), zap.Any("template", template.WorkspaceInvitePendingTemplate))
			return
		}
		emailHTML, err := emailTemplate.RenderHTML()
		if err != nil {
			ws.logger.Error("couldn't template convert to html", zap.Error(err), zap.Any("template", template.WorkspaceInvitePendingTemplate))
			return
		}
		email := models.Email{
			To:           workspaceUsersEmail,
			EmailFormat:  models.EmailFormatHTML,
			EmailContent: emailHTML,
			EmailSubject: "You've been invited to Byrd",
		}
		ws.sendEmail(email)
	}()

	return workspaceUsers, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMemberEmail string, workspaceID uuid.UUID) error {
	// Get the workspace member
	workspaceUser, err := ws.userService.GetUserByEmail(ctx, workspaceMemberEmail)
	if err != nil {
		return err
	}

	// Get workspace member count
	activeUsers, pendingUsers, activeAdmins, pendingAdmins, err := ws.workspaceRepo.GetWorkspaceUserCountsByRoleAndStatus(ctx, workspaceID)
	if err != nil {
		return errors.New("couldn't get workspace user count")
	}

	// If the user is the last admin in the workspace
	if activeAdmins == 1 {
		if activeUsers != 0 {
			// If there are active users, promote one to admin
			return ws.workspaceRepo.PromoteRandomUserToAdmin(ctx, workspaceID)
		}
		// If there are no active users or admins, delete the workspace
		wstatus, err := ws.DeleteWorkspace(ctx, workspaceID)
		if err != nil {
			ws.logger.Error("failed to delete workspace", zap.Error(err), zap.Any("workspaceID", workspaceID), zap.Any("workspaceStatus", wstatus), zap.Any("workspaceMemberID", workspaceUser.ID), zap.Any("totalPendingUsers", pendingUsers), zap.Any("totalPendingAdmins", pendingAdmins), zap.Any("totalActiveUsers", activeUsers), zap.Any("totalActiveAdmins", activeAdmins))
		}

		return nil
	}

	// Remove the user from the workspace
	err = ws.workspaceRepo.RemoveUserFromWorkspace(ctx, workspaceID, workspaceUser.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.WorkspaceRole) error {
	err := ws.workspaceRepo.UpdateUserRoleForWorkspace(ctx, workspaceID, workspaceMemberID, role)
	if err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error {
	if err := ws.workspaceRepo.RemoveUserFromWorkspace(ctx, workspaceID, workspaceMemberID); err != nil {
		return err
	}
	return nil
}

func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedUserEmail string, workspaceID uuid.UUID) error {
	user, err := ws.userService.ActivateUser(ctx, invitedUserEmail)
	if err != nil {
		return err
	}

	// Update the user's membership status
	err = ws.workspaceRepo.UpdateUserMembershipStatusForWorkspace(ctx, workspaceID, user.ID, models.ActiveMember)
	if err != nil {
		return err
	}

	go func() {
		emailTemplate, err := ws.library.GetTemplate(template.WorkspaceInviteAcceptedTemplate)
		if err != nil {
			ws.logger.Error("couldn't get email template", zap.Error(err), zap.Any("template", template.WorkspaceInviteAcceptedTemplate))
			return
		}
		emailHTML, err := emailTemplate.RenderHTML()
		if err != nil {
			ws.logger.Error("couldn't template convert to html", zap.Error(err), zap.Any("template", template.WorkspaceInviteAcceptedTemplate))
			return
		}
		email := models.Email{
			To:           []string{invitedUserEmail},
			EmailFormat:  models.EmailFormatHTML,
			EmailContent: emailHTML,
			EmailSubject: "And You're In! Own Your Competitor's Next Move As They Make It",
		}
		ws.sendEmail(email)
	}()

	return nil
}

func (ws *workspaceService) GetWorkspaceUser(ctx context.Context, workspaceID uuid.UUID, workspaceMemberEmail string) (*models.PartialWorkspaceUser, error) {
	user, err := ws.userService.GetUserByEmail(ctx, workspaceMemberEmail)
	if err != nil {
		return nil, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return nil, err
	}

	return workspaceUser, nil
}
