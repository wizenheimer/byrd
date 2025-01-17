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

func NewWorkspaceService(workspaceRepo workspace.WorkspaceRepository, competitorService competitor.CompetitorService, userService user.UserService, tm *transaction.TxManager, logger *logger.Logger) WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		logger:            logger.WithFields(map[string]interface{}{"module": "workspace_service"}),
		tm:                tm,
	}
}

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, pages []models.PageProps, userEmails []string) (*models.Workspace, error) {
	ws.logger.Debug("creating workspace", zap.Any("workspaceOwner", workspaceOwner), zap.Any("pages", pages), zap.Any("userEmails", userEmails))
	// Step 0: Create workspace along with the owner
	var workspace *models.Workspace
	if err := ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
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

	// STEP 1: Validate and Create Competitors for the workspace
	err := utils.SetDefaultsAndValidateArray(&pages)
	if err != nil {
		ws.logger.Debug("failed to validate pages", zap.Error(err))
	} else {
		// Add competitors to the workspace
		_, err = ws.competitorService.AddCompetitorsToWorkspace(ctx, workspace.ID, pages)
		if err != nil {
			ws.logger.Debug("failed to add competitors to workspace", zap.Error(err))
		}
	}

	// Step 2: Validate and invite users to the workspace
	ownerEmail, err := utils.GetClerkUserEmail(workspaceOwner)
	if err != nil {
		return nil, err
	}
	ownerEmail = utils.NormalizeEmail(ownerEmail)

	emailMap := make(map[string]bool)
	emailMap[ownerEmail] = true

	memberEmails := make([]string, 0)
	for _, memberEmail := range userEmails {
		// Normalize email address
		memberEmail = utils.NormalizeEmail(memberEmail)

		// Check if the email is already in the map
		if _, ok := emailMap[memberEmail]; ok {
			continue
		}
		emailMap[memberEmail] = true

		// Add the member to the list of invited users
		memberEmails = append(memberEmails, memberEmail)
	}

	members, err := ws.userService.BatchGetOrCreateUsers(ctx, memberEmails)
	if err != nil {
		ws.logger.Debug("failed to get or create users", zap.Error(err))
	}

	if len(members) == 0 {
		ws.logger.Debug("no users to invite")
	} else {
		// Add the users to the workspace
		memberIDs := make([]uuid.UUID, 0)
		for _, member := range members {
			memberIDs = append(memberIDs, member.ID)
		}

		_, err := ws.workspaceRepo.BatchAddUsersToWorkspace(ctx, workspace.ID, memberIDs)
		if err != nil {
			return nil, err
		}
	}

	return workspace, nil
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error) {
	ws.logger.Debug("listing user workspaces", zap.Any("workspaceMember", workspaceMember))
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
	ws.logger.Debug("getting workspace", zap.Any("workspaceID", workspaceID))
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

	ws.logger.Debug("updating workspace", zap.Any("workspaceID", workspaceID), zap.Any("workspaceProps", workspaceProps), zap.Bool("updatedNameRequiresUpdate", updatedNameRequiresUpdate), zap.Bool("billingEmailRequiresUpdate", billingEmailRequiresUpdate))

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
	ws.logger.Debug("deleting workspace", zap.Any("workspaceID", workspaceID))
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

func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, roleFilter *models.WorkspaceRole) ([]models.WorkspaceUser, bool, error) {
	ws.logger.Debug("listing workspace members", zap.Any("workspaceID", workspaceID), zap.Any("limit", limit), zap.Any("offset", offset), zap.Any("roleFilter", roleFilter))
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
			ws.logger.Debug("skipping, user not found", zap.Any("member", member))
			continue
		} else {
			ws.logger.Debug("adding user to workspace", zap.Any("member", member))
		}
		workspaceUser := models.WorkspaceUser{
			ID:               member.ID,
			WorkspaceID:      workspaceID,
			Role:             member.Role,
			Name:             utils.FromPtr(user.Name, ""),
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, hasMore, nil
}

func (ws *workspaceService) AddUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error) {
	ws.logger.Debug("adding users to workspace", zap.Any("workspaceID", workspaceID), zap.Any("workspaceMember", workspaceMember), zap.Any("emails", emails))

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

	partialWorkspaceUsers, err := ws.workspaceRepo.BatchAddUsersToWorkspace(ctx, workspaceID, userIDs)
	if err != nil {
		return nil, err
	}

	workspaceUsers := make([]models.WorkspaceUser, 0)
	for _, member := range partialWorkspaceUsers {
		user, ok := userIDsToUserMap[member.ID]
		if !ok {
			ws.logger.Debug("skipping, user not found", zap.Any("member", member))
			continue
		} else {
			ws.logger.Debug("adding user to workspace", zap.Any("member", member))
		}
		workspaceUser := models.WorkspaceUser{
			ID:               member.ID,
			WorkspaceID:      workspaceID,
			Role:             member.Role,
			Name:             utils.FromPtr(user.Name, ""),
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error {
	ws.logger.Debug("leaving workspace", zap.Any("workspaceID", workspaceID), zap.Any("workspaceMember", workspaceMember))
	// Get the workspace member
	workspaceUser, err := ws.userService.GetUserByClerkCredentials(ctx, workspaceMember)
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
			ws.logger.Debug("promoting a random active user to admin", zap.Any("activeUsers", activeUsers), zap.Any("pendingUsers", pendingUsers), zap.Any("activeAdmins", activeAdmins), zap.Any("pendingAdmins", pendingAdmins))
			return ws.workspaceRepo.PromoteRandomUserToAdmin(ctx, workspaceID)
		}
		// If there are no active users or admins, delete the workspace
		ws.logger.Debug("attempting to delete the workspace", zap.Any("activeUsers", activeUsers), zap.Any("pendingUsers", pendingUsers), zap.Any("activeAdmins", activeAdmins), zap.Any("pendingAdmins", pendingAdmins))
		_, err := ws.DeleteWorkspace(ctx, workspaceID)
		return err
	}

	// Remove the user from the workspace
	err = ws.workspaceRepo.RemoveUserFromWorkspace(ctx, workspaceID, workspaceUser.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.WorkspaceRole) error {
	ws.logger.Debug("updating workspace member role", zap.Any("workspaceID", workspaceID), zap.Any("workspaceMemberID", workspaceMemberID), zap.Any("role", role))
	err := ws.workspaceRepo.UpdateUserRoleForWorkspace(ctx, workspaceID, workspaceMemberID, role)
	if err != nil {
		return err
	}

	return nil
}

func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error {
	ws.logger.Debug("removing user from workspace", zap.Any("workspaceID", workspaceID), zap.Any("workspaceMemberID", workspaceMemberID))
	err := ws.workspaceRepo.RemoveUserFromWorkspace(ctx, workspaceID, workspaceMemberID)

	if err != nil {
		return err
	}
	return nil
}

func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error {
	ws.logger.Debug("joining workspace", zap.Any("workspaceID", workspaceID), zap.Any("invitedMember", invitedMember))
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
	ws.logger.Debug("checking if workspace exists", zap.Any("workspaceID", workspaceID))
	return ws.workspaceRepo.WorkspaceExists(ctx, workspaceID)
}

func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	ws.logger.Debug("checking if user is workspace admin", zap.Any("workspaceID", workspaceID), zap.Any("clerkUser", clerkUser))
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	ws.logger.Debug("got workspace user role", zap.Any("role", workspaceUser.Role))
	return workspaceUser.Role == models.RoleAdmin, nil
}

func (ws *workspaceService) ClerkUserIsActiveWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	ws.logger.Debug("checking if user is active workspace member", zap.Any("workspaceID", workspaceID), zap.Any("clerkUser", clerkUser))
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	ws.logger.Debug("got workspace user membership status", zap.Any("membershipStatus", workspaceUser.MembershipStatus))
	return workspaceUser.MembershipStatus == models.ActiveMember, nil
}

func (ws *workspaceService) ClerkUserIsPendingWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	ws.logger.Debug("checking if user is active workspace member", zap.Any("workspaceID", workspaceID), zap.Any("clerkUser", clerkUser))
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	ws.logger.Debug("got workspace user membership status", zap.Any("membershipStatus", workspaceUser.MembershipStatus))
	return workspaceUser.MembershipStatus == models.PendingMember, nil
}

func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	ws.logger.Debug("checking if user is workspace member", zap.Any("workspaceID", workspaceID), zap.Any("clerkUser", clerkUser))
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return false, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return false, err
	}

	ws.logger.Debug("got workspace user role", zap.Any("role", workspaceUser.Role))
	return workspaceUser.Role == models.RoleAdmin || workspaceUser.Role == models.RoleUser, nil
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	ws.logger.Debug("checking if competitor exists", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID))
	exists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error) {
	ws.logger.Debug("checking if competitor page exists", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID), zap.Any("pageID", pageID))
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
	ws.logger.Debug("adding competitors to workspace", zap.Any("workspaceID", workspaceID), zap.Any("pages", pages))
	competitors, err := ws.competitorService.AddCompetitorsToWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	if len(competitors) == 0 {
		ws.logger.Debug("failed to create competitor", zap.Int("numPages", len(pages)), zap.Int("numCreatedCompetitors", len(competitors)))
		return nil, errors.New("failed to create competitor")
	}

	return &competitors[0], nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pageProps []models.PageProps) ([]models.Page, error) {
	ws.logger.Debug("adding pages to competitor", zap.Any("competitorID", competitorID), zap.Any("pages", pageProps))
	createdPages, err := ws.competitorService.AddPagesToCompetitor(ctx, competitorID, pageProps)
	if err != nil {
		return nil, err
	}
	if len(pageProps) != len(createdPages) {
		if len(createdPages) == 0 {
			ws.logger.Debug("failed to create pages", zap.Int("numPages", len(pageProps)), zap.Int("numCreatedPages", len(createdPages)))
			return nil, errors.New("failed to create pages")
		}
		ws.logger.Debug("failed to create some pages", zap.Int("numPages", len(pageProps)), zap.Int("numCreatedPages", len(createdPages)))
		return createdPages, errors.New("failed to create some pages")
	}

	return createdPages, nil
}

func (ws *workspaceService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error) {
	ws.logger.Debug("listing competitors for workspace", zap.Any("workspaceID", workspaceID), zap.Any("limit", limit), zap.Any("offset", offset))
	competitors, hasMore, err := ws.competitorService.ListCompetitorsForWorkspace(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, false, err
	}

	return competitors, hasMore, nil
}

// ListPagesForCompetitor lists the pages for a competitor
func (ws *workspaceService) ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	ws.logger.Debug("listing pages for competitor", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID), zap.Any("limit", limit), zap.Any("offset", offset))
	pages, hasMore, err := ws.competitorService.ListCompetitorPages(ctx, competitorID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pages, hasMore, nil
}

// ListHistoryForPage lists the history of a page
func (ws *workspaceService) ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	ws.logger.Debug("listing history for page", zap.Any("pageID", pageID), zap.Any("limit", limit), zap.Any("offset", offset))
	pageHistory, hasMore, err := ws.competitorService.ListPageHistory(ctx, pageID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pageHistory, hasMore, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error {
	ws.logger.Debug("removing page from workspace", zap.Any("competitorID", competitorID), zap.Any("pageID", pageID))
	return ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, []uuid.UUID{pageID})
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	ws.logger.Debug("removing competitor from workspace", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID))
	return ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, pageProps models.PageProps) (*models.Page, error) {
	ws.logger.Debug("updating competitor page", zap.Any("pageProps", pageProps), zap.Any("competitorID", competitorID), zap.Any("pageID", pageID))
	return ws.competitorService.UpdatePage(ctx, competitorID, pageID, pageProps)
}
