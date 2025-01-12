// ./src/internal/service/workspace/service.go
package workspace

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workspace"
	"github.com/wizenheimer/byrd/src/internal/service/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type workspaceService struct {
	workspaceRepo     workspace.WorkspaceRepository
	competitorService competitor.CompetitorService
	userService       user.UserService
	logger            *logger.Logger
	tm                *transaction.TxManager
}

func NewWorkspaceService(workspaceRepo workspace.WorkspaceRepository, competitorService competitor.CompetitorService, userService user.UserService, logger *logger.Logger) WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		logger:            logger,
	}
}

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, pages []models.PageProps, users []models.UserProps) (*models.Workspace, error) {
	// STEP 0: Validate the workspace creation request
	err := utils.SetDefaultsAndValidateArray(&pages)
	if err != nil {
		return nil, err
	}

	err = utils.SetDefaultsAndValidateArray(&users)
	if err != nil {
		return nil, err
	}

	ownerEmail, err := utils.GetClerkUserEmail(workspaceOwner)
	if err != nil {
		return nil, err
	}

	ownerEmail = utils.NormalizeEmail(ownerEmail)

	emailMap := make(map[string]bool)
	emailMap[ownerEmail] = true

	memberEmails := make([]string, 0)
	for _, member := range users {
		// Normalize email address
		member.Email = utils.NormalizeEmail(member.Email)

		// Check if the email is already in the map
		if _, ok := emailMap[member.Email]; ok {
			continue
		}
		emailMap[member.Email] = true

		// Add the member to the list of invited users
		memberEmails = append(memberEmails, member.Email)
	}

	// Step 1: Create workspace along with the owner
	var workspace *models.Workspace
	if err = ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		createdUser, err := ws.userService.GetOrCreateUser(ctx, workspaceOwner)
		if err != nil {
			return err
		}

		// Step 2: Create the workspace
		workspaceName := utils.GenerateWorkspaceName(workspaceOwner)
		workspace, err = ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, *createdUser.Email, createdUser.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Step 2: Invite users to the workspace
	members, err := ws.userService.BatchGetOrCreateUsers(ctx, memberEmails)
	if err != nil {
		ws.logger.Debug("Failed to get or create users", zap.Error(err))
	}
	if len(members) == 0 {
		ws.logger.Debug("No users to invite")
	} else {
		// Add the users to the workspace
		memberIDs := make([]uuid.UUID, 0)
		for _, member := range members {
			memberIDs = append(memberIDs, member.ID)
		}

		ws.workspaceRepo.BatchAddUsersToWorkspace(ctx, memberIDs, workspace.ID)
		if err != nil {
			return nil, err
		}
	}

	// Step 3: Create competitors for the workspace
	// Add competitors to the workspace
	_, err = ws.competitorService.AddCompetitorsToWorkspace(ctx, workspace.ID, pages)
	if err != nil {
		ws.logger.Debug("Failed to add competitors to workspace", zap.Error(err))
	}

	return workspace, nil
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error) {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, workspaceMember)
	if err != nil {
		return nil, err
	}

	workspaces, err := ws.workspaceRepo.GetWorkspacesForUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	workspace, err := ws.workspaceRepo.GetWorkspaceByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceProps models.WorkspaceProps) error {
	// Check which fields are being updated, and update them
	updatedBillingEmail := utils.NormalizeEmail(workspaceProps.BillingEmail)
	// TODO: add normalization for workspace name
	updateWorkspaceName := workspaceProps.Name

	updatedNameRequiresUpdate := updateWorkspaceName != ""
	billingEmailRequiresUpdate := updatedBillingEmail != ""

	if updatedNameRequiresUpdate && billingEmailRequiresUpdate {
		return ws.workspaceRepo.UpdateWorkspaceDetails(ctx, workspaceID, updateWorkspaceName, updatedBillingEmail)
	}

	if updatedNameRequiresUpdate {
		return ws.workspaceRepo.UpdateWorkspaceName(ctx, workspaceID, updateWorkspaceName)
	}

	if billingEmailRequiresUpdate {
		return ws.workspaceRepo.UpdateWorkspaceBillingEmail(ctx, workspaceID, updatedBillingEmail)
	}

	return nil
}

func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error) {
	err := ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		// Step 1: Remove all competitors from the workspace
		err := ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, nil)
		if err != nil {
			return err
		}

		// Step 2: Remove all users from the workspace
		err = ws.workspaceRepo.DeleteWorkspace(ctx, workspaceID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.WorkspaceActive, err
	}

	// Step 3: Delete the workspace
	return models.WorkspaceInactive, nil
}

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, roleFilter *models.WorkspaceRole) ([]models.WorkspaceUser, error) {
	// TODO: add filtering and pagination support
	members, err := ws.workspaceRepo.GetWorkspaceMembers(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	membersIDs := make([]uuid.UUID, 0)
	for _, member := range members {
		membersIDs = append(membersIDs, member.ID)
	}

	users, err := ws.userService.ListUsersByUserIDs(ctx, membersIDs)
	if err != nil {
		return nil, err
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
			Name:             *user.Name,
			Email:            *user.Email,
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, nil
}

func (ws *workspaceService) AddUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error) {
	// Step 1: Normalize the emails
	normalizedEmails := make([]string, 0)
	for _, email := range emails {
		normalizedEmails = append(normalizedEmails, utils.NormalizeEmail(email))
	}

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

	partialWorkspaceUsers, err := ws.workspaceRepo.BatchAddUsersToWorkspace(ctx, userIDs, workspaceID)
	if err != nil {
		return nil, err
	}

	workspaceUsers := make([]models.WorkspaceUser, 0)
	for _, member := range partialWorkspaceUsers {
		user, ok := userIDsToUserMap[member.ID]
		if !ok {
			continue
		}
		workspaceUser := models.WorkspaceUser{
			ID:               member.ID,
			WorkspaceID:      workspaceID,
			Role:             member.Role,
			Name:             *user.Name,
			Email:            *user.Email,
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error {
	// Get the workspace member
	workspaceUser, err := ws.userService.GetUserByClerkCredentials(ctx, workspaceMember)
	if err != nil {
		return nil
	}

	// Get workspace member count
	workspaceCount, err := ws.workspaceRepo.GetWorkspaceUserCountByRole(ctx, workspaceID)
	if err != nil {
		return nil
	}

	// Check if the user is the last admin in the workspace
	adminRoleCount := workspaceCount[models.RoleAdmin]
	userRoleCount := workspaceCount[models.RoleUser]

	if userRoleCount != 0 && adminRoleCount == 1 {
		return errors.New("cannot leave workspace as the last admin")
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
	err := ws.workspaceRepo.RemoveUserFromWorkspace(ctx, workspaceID, workspaceMemberID)

	if err != nil {
		return err
	}
	return nil
}

func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, invitedMember)
	if err != nil {
		return err
	}

	err = ws.workspaceRepo.UpdateUserMembershipStatusForWorkspace(ctx, workspaceID, user.ID, models.ActiveMember)
	if err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	return ws.workspaceRepo.WorkspaceExists(ctx, workspaceID)
}

func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	return workspaceUser.Role == models.RoleAdmin, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	return workspaceUser.Role == models.RoleAdmin || workspaceUser.Role == models.RoleUser, nil
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	exists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error) {
	competitorExists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	if !competitorExists {
		return false, nil
	}

	pageExists, err := ws.competitorService.PageExists(ctx, competitorID, pageID)
	if err != nil {
		return false, err
	}

	if !pageExists {
		return false, nil
	}

	return true, nil
}

func (ws *workspaceService) AddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (*models.Competitor, error) {
	competitors, err := ws.competitorService.AddCompetitorsToWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	if len(competitors) == 0 {
		return nil, errors.New("failed to create competitor")
	}

	return &competitors[0], nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	createdPages, err := ws.competitorService.AddPagesToCompetitor(ctx, competitorID, pages)
	if err != nil {
		return nil, err
	}
	if len(pages) != len(createdPages) {
		if len(createdPages) == 0 {
			return nil, errors.New("failed to create pages")
		}
		return createdPages, errors.New("failed to create some pages")
	}

	return createdPages, nil
}

func (ws *workspaceService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, error) {
	competitors, err := ws.competitorService.ListCompetitorsForWorkspace(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, err
	}

	return competitors, nil
}

// ListPagesForCompetitor lists the pages for a competitor
func (ws *workspaceService) ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error) {
	pages, err := ws.competitorService.ListCompetitorPages(ctx, competitorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// ListHistoryForPage lists the history of a page
func (ws *workspaceService) ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
	pageHistory, err := ws.competitorService.ListPageHistory(ctx, pageID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pageHistory, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error {
	return ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, []uuid.UUID{pageID})
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	return ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return ws.competitorService.UpdatePage(ctx, competitorID, pageID, page)
}
