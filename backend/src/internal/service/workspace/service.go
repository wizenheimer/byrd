// ./src/internal/service/workspace/service.go
package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
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
	library           template.TemplateLibrary
	emailClient       email.EmailClient
	userService       user.UserService
	logger            *logger.Logger
	errorRecord       *recorder.ErrorRecorder
	tm                *transaction.TxManager
}

func NewWorkspaceService(
	workspaceRepo workspace.WorkspaceRepository,
	competitorService competitor.CompetitorService,
	userService user.UserService,
	library template.TemplateLibrary,
	tm *transaction.TxManager,
	emailClient email.EmailClient,
	logger *logger.Logger,
	errorRecord *recorder.ErrorRecorder,
) (WorkspaceService, error) {

	ws := workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		library:           library,
		logger: logger.WithFields(map[string]any{
			"module": "workspace_service",
		}),
		emailClient: emailClient,
		errorRecord: errorRecord,
		tm:          tm,
	}

	return &ws, nil
}

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, pages []models.PageProps, userEmails []string) (*models.Workspace, error) {
	// Step 0: Create workspace along with the owner
	var workspace *models.Workspace
	var workspaceCreator *models.User
	var err error

	// Step 1: Get or create the workspace owner
	workspaceCreator, err = ws.userService.GetOrCreateUser(ctx, workspaceOwner)
	if err != nil {
		return nil, err
	}

	// Step 2: Check if the user can create a workspace
	canCreate, err := ws.CanCreateWorkspace(ctx, workspaceCreator.ID)
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, errors.New("user cannot create workspace")
	}

	// Step 3: Create the workspace
	workspaceName := utils.GenerateWorkspaceName(workspaceOwner)

	err = ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		workspace, err = ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, *workspaceCreator.Email, workspaceCreator.ID, models.WorkspaceTrial)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("failed to create workspace")
	}
	// Step 4: Check if the user creation request is valid
	// If not, truncate the list to the maximum allowed limit
	pageLimit, err := workspace.GetMaxPages()
	if err != nil {
		return nil, err
	}
	if len(pages) > pageLimit {
		pages = pages[:pageLimit]
	}

	// Step 5: Add competitors to the workspace
	if competitors, err := ws.competitorService.BatchCreateCompetitorForWorkspace(ctx, workspace.ID, pages); err != nil {
		ws.logger.Error("failed to create competitors for workspace", zap.Error(err), zap.Any("workspaceID", workspace.ID), zap.Any("pages", len(pages)), zap.Any("competitors", len(competitors)))
	} else if len(competitors) != len(pages) {
		ws.logger.Error("failed to create all competitors for workspace", zap.Any("workspaceID", workspace.ID), zap.Any("pages", len(pages)), zap.Any("competitors", len(competitors)))
	}

	// Step 6: Check if the user creation request is valid
	// If not, truncate the list to the maximum allowed limit
	userLimit, err := workspace.GetMaxUsers()
	if err != nil {
		return nil, err
	}
	if len(userEmails) > userLimit {
		userEmails = userEmails[:userLimit]
	}

	// Step 7: Add users to the workspace
	// If the user is the owner, add them as an admin
	workspaceUsers, err := ws.AddUsersToWorkspace(ctx, workspaceOwner, workspace.ID, userEmails)
	if err != nil {
		ws.logger.Error("failed to add users to workspace", zap.Error(err), zap.Any("workspaceID", workspace.ID), zap.Any("userEmails", len(userEmails)))
	} else if len(workspaceUsers) != len(userEmails) {
		ws.logger.Error("failed to add all users to workspace", zap.Any("workspaceID", workspace.ID), zap.Any("userEmails", len(userEmails)), zap.Any("workspaceUsers", len(workspaceUsers)))
	}

	// Step 8: Return the workspace
	return workspace, nil
}

func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User, membershipStatus models.MembershipStatus) ([]models.Workspace, error) {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, workspaceMember)
	if err != nil {
		return nil, err
	}

	workspaces, err := ws.workspaceRepo.GetWorkspacesForUserID(ctx, user.ID, membershipStatus)
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
	// Check if the workspace name or billing email is being updated
	updatedNameRequiresUpdate := workspaceProps.Name != ""
	billingEmailRequiresUpdate := workspaceProps.BillingEmail != ""

	// Update the workspace
	if updatedNameRequiresUpdate && billingEmailRequiresUpdate {
		// Update both the workspace name and billing email
		return ws.workspaceRepo.UpdateWorkspaceDetails(ctx, workspaceID, workspaceProps.Name, workspaceProps.BillingEmail)
	} else if updatedNameRequiresUpdate {
		// Update the workspace name
		return ws.workspaceRepo.UpdateWorkspaceName(ctx, workspaceID, workspaceProps.Name)
	} else if billingEmailRequiresUpdate {
		// Update the billing email
		return ws.workspaceRepo.UpdateWorkspaceBillingEmail(ctx, workspaceID, workspaceProps.BillingEmail)
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspacePlan(ctx context.Context, workspaceID uuid.UUID, plan models.WorkspacePlan) error {
	return ws.workspaceRepo.UpdateWorkspacePlan(ctx, workspaceID, plan)
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
			Name:             utils.FromPtr(user.Name, ""),
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	return workspaceUsers, hasMore, nil
}

func (ws *workspaceService) AddUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error) {
	canCreate, err := ws.CanAddUsers(ctx, workspaceID, len(emails))
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, errors.New("member cannot add users to workspace")
	}

	inviterEmail, err := utils.GetClerkUserEmail(workspaceMember)
	if err != nil {
		return nil, err
	}

	// Step 1: Normalize the emails
	normalizedEmails := utils.CleanEmailList(emails, []string{inviterEmail})

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
			Name:             utils.FromPtr(user.Name, ""),
			Email:            utils.FromPtr(user.Email, ""),
			MembershipStatus: member.MembershipStatus,
		}
		workspaceUsers = append(workspaceUsers, workspaceUser)
	}

	go func() {
		// TODO: Send Email to the users acknowledging the invitation to the workspace
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
		// TODO: fix the cancelation of the email sending
		ws.sendEmail(email)
	}()

	return workspaceUsers, nil
}

func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error {
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

func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error {
	// Get the user by clerk credentials
	user, err := ws.userService.GetUserByClerkCredentials(ctx, invitedMember)
	if err != nil {
		return err
	}
	userEmail, err := utils.GetClerkUserEmail(invitedMember)
	if err != nil {
		return err
	}

	firstTimeUser := user.ClerkID == nil || user.Status != models.AccountStatusActive
	// Synchronize the user with the database if the user is first time user
	if firstTimeUser {
		if err := ws.userService.ActivateUser(ctx, user.ID, invitedMember); err != nil {
			return err
		}
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
			To:           []string{userEmail},
			EmailFormat:  models.EmailFormatHTML,
			EmailContent: emailHTML,
			EmailSubject: "And You're In! Own Your Competitor's Next Move As They Make It",
		}
		ws.sendEmail(email)
	}()

	return nil
}

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	return ws.workspaceRepo.WorkspaceExists(ctx, workspaceID)
}

func (ws *workspaceService) GetClerkWorkspaceUser(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (*models.PartialWorkspaceUser, error) {
	user, err := ws.userService.GetUserByClerkCredentials(ctx, clerkUser)
	if err != nil {
		return nil, err
	}

	workspaceUser, err := ws.workspaceRepo.GetWorkspaceMemberByUserID(ctx, workspaceID, user.ID)
	if err != nil {
		return nil, err
	}

	return workspaceUser, nil
}

func (ws *workspaceService) GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error) {
	competitors, err := ws.competitorService.GetCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil {
		return nil, err
	}
	if len(competitors) == 0 {
		return nil, errors.New("competitor not found")
	}

	return &competitors[0], nil
}

func (ws *workspaceService) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return ws.competitorService.UpdateCompetitorForWorkspace(ctx, workspaceID, competitorID, competitorName)
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
	canCreateCompetitor, err := ws.CanCreateCompetitor(ctx, workspaceID, 1, len(pages))
	if err != nil {
		return nil, err
	}
	if !canCreateCompetitor {
		return nil, errors.New("user cannot create competitor")
	}

	competitor, err := ws.competitorService.CreateCompetitorForWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	return &competitor, nil
}

func (ws *workspaceService) BatchAddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error) {
	canCreateCompetitor, err := ws.CanCreateCompetitor(ctx, workspaceID, len(pages), len(pages))
	if err != nil {
		return nil, err
	}
	if !canCreateCompetitor {
		return nil, errors.New("user cannot create competitor")
	}

	competitors, err := ws.competitorService.BatchCreateCompetitorForWorkspace(ctx, workspaceID, pages)
	if err != nil {
		return nil, err
	}

	if len(competitors) == 0 {
		ws.logger.Error("failed to create competitors", zap.Any("workspaceID", workspaceID), zap.Any("pageProps", pages), zap.Any("numCompetitors", len(competitors)))
		return nil, errors.New("failed to create competitors")
	}

	return competitors, nil
}

func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pageProps []models.PageProps) ([]models.Page, error) {
	canCreatePage, err := ws.CanCreatePage(ctx, workspaceID, len(pageProps))
	if err != nil {
		return nil, err
	}
	if !canCreatePage {
		return nil, errors.New("user cannot create page")
	}

	createdPages, err := ws.competitorService.AddPagesToCompetitor(ctx, competitorID, pageProps)
	if err != nil {
		return nil, err
	}

	if len(pageProps) != len(createdPages) {
		ws.logger.Error("failed to create some pages", zap.Int("numPages", len(pageProps)), zap.Int("numCreatedPages", len(createdPages)), zap.Any("pages", pageProps), zap.Any("createdPages", createdPages))
	}

	return createdPages, nil
}

func (ws *workspaceService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error) {
	competitors, hasMore, err := ws.competitorService.ListCompetitorsForWorkspace(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, false, err
	}

	return competitors, hasMore, nil
}

// ListPagesForCompetitor lists the pages for a competitor
func (ws *workspaceService) ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	pages, hasMore, err := ws.competitorService.ListCompetitorPages(ctx, competitorID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pages, hasMore, nil
}

// ListHistoryForPage lists the history of a page
func (ws *workspaceService) ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	pageHistory, hasMore, err := ws.competitorService.ListPageHistory(ctx, pageID, limit, offset)
	if err != nil {
		return nil, hasMore, err
	}

	return pageHistory, hasMore, nil
}

func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error {
	return ws.competitorService.RemovePagesFromCompetitor(ctx, competitorID, []uuid.UUID{pageID})
}

func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	return ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
}

func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, pageProps models.PageProps) (*models.Page, error) {
	return ws.competitorService.UpdatePage(ctx, competitorID, pageID, pageProps)
}

func (ws *workspaceService) GetPageForCompetitor(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	return ws.competitorService.GetCompetitorPage(ctx, competitorID, pageID)
}

func (ws *workspaceService) ListReports(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	return ws.competitorService.ListReports(ctx, workspaceID, competitorID, limit, offset)
}

func (ws *workspaceService) CreateReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) (*models.Report, error) {
	return ws.competitorService.CreateReport(ctx, workspaceID, competitorID)
}

func (ws *workspaceService) DispatchReportToWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) error {
	members, hasMore, err := ws.workspaceRepo.ListWorkspaceMembers(ctx, workspaceID, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return errors.New("no members found")
	}
	if hasMore {
		ws.logger.Warn("more members found", zap.Any("workspaceID", workspaceID))
	}

	userIDs := make([]uuid.UUID, 0)
	for _, member := range members {
		userIDs = append(userIDs, member.ID)
	}

	users, err := ws.userService.ListUsersByUserIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	subscriberEmails := make([]string, 0)
	for _, user := range users {
		if user.Email != nil {
			subscriberEmails = append(subscriberEmails, *user.Email)
		}
	}

	return ws.competitorService.DispatchReport(ctx, workspaceID, competitorID, subscriberEmails)
}

func (ws *workspaceService) DispatchReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID, subscriberEmails []string) error {
	// Clean up the email list for duplicates and nil values
	subscriberEmails = utils.CleanEmailList(subscriberEmails, nil)

	// Dispatch the report
	return ws.competitorService.DispatchReport(ctx, workspaceID, competitorID, subscriberEmails)
}

func (ws *workspaceService) ListActiveWorkspaces(ctx context.Context, batchSize int, lastWorkspaceID *uuid.UUID) (<-chan []uuid.UUID, <-chan error) {
	workspaceChan := make(chan []uuid.UUID)
	errorChan := make(chan error)

	go func() {
		defer close(workspaceChan)
		defer close(errorChan)

		hasMore := true
		for hasMore {
			activeWorkspaces, err := ws.workspaceRepo.ListActiveWorkspaces(ctx, batchSize, lastWorkspaceID)
			if err != nil {
				errorChan <- err
				return
			}

			hasMore = activeWorkspaces.HasMore
			workspaceChan <- activeWorkspaces.WorkspaceIDs
		}
	}()

	return workspaceChan, errorChan
}

func (ws *workspaceService) sendEmail(email models.Email) {
	// Create a context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Important to avoid context leak

	if err := ws.emailClient.Send(ctx, email); err != nil {
		ws.errorRecord.RecordError(ctx, err, zap.Any("subscriberEmails", email.To), zap.Any("emailSubject", email.EmailSubject))
	}
}

// CanCreateWorkspace checks if the user can create a workspace
// based on the user's current workspace count and the maximum workspace limit
func (ws *workspaceService) CanCreateWorkspace(ctx context.Context, userID uuid.UUID) (bool, error) {
	totalIncomingWorkspaces := 1

	currentCount, err := ws.CountUserWorkspaces(ctx, userID)
	if err != nil {
		return false, err
	}

	maxCount, err := models.GetWorkspaceCreationLimit()
	if err != nil {
		return false, err
	}

	// Check if the user can create a workspace
	return currentCount+totalIncomingWorkspaces <= maxCount, nil
}

// CanCreateCompetitor checks if the user can create a competitor
// based on the user's current competitor count and the maximum competitor limit
// it also checks if the user can create a page based on the user's current page count and the maximum page limit
func (ws *workspaceService) CanCreateCompetitor(ctx context.Context, workspaceID uuid.UUID, totalIncomingCompetitors int, totalIncomingPages int) (bool, error) {
	canCreatePage, err := ws.CanCreatePage(ctx, workspaceID, totalIncomingPages)
	if err != nil {
		return false, err
	}
	if !canCreatePage {
		return false, errors.New("user cannot create page")
	}

	// Get the workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Get the max competitors for the workspace
	competitorLimit, err := workspace.GetMaxCompetitors()
	if err != nil {
		return false, err
	}

	// Get the current competitor count
	currentCount, err := ws.CountWorkspaceCompetitors(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Check if the user can create a competitor
	return currentCount+totalIncomingCompetitors <= competitorLimit, nil
}

// CanCreatePage checks if the user can create a page
// based on the user's current page count and the maximum page limit
func (ws *workspaceService) CanCreatePage(ctx context.Context, workspaceID uuid.UUID, totalIncomingPages int) (bool, error) {
	// Get the workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Get the max pages for the competitor
	limit, err := workspace.GetMaxPages()
	if err != nil {
		return false, err
	}

	// Get the current page count
	pageCount, err := ws.CountWorkspacePages(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Check if the user can create a page
	return pageCount+totalIncomingPages <= limit, nil
}

// CanAddUsers checks if the user can add users to the workspace
// based on the user's current user count and the maximum user limit
func (ws *workspaceService) CanAddUsers(ctx context.Context, workspaceID uuid.UUID, totalIncomingUsers int) (bool, error) {
	// Get the workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Get the max users for the workspace
	maxCount, err := workspace.GetMaxUsers()
	if err != nil {
		return false, err
	}

	// Get the current user count
	activeCount, pendingCount, err := ws.CountWorkspaceMembers(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	// Check if the user can add users
	return activeCount+pendingCount+totalIncomingUsers <= maxCount, nil
}

// CountUserWorkspaces counts the number of workspaces for a user
func (ws *workspaceService) CountUserWorkspaces(ctx context.Context, userID uuid.UUID) (totalCount int, err error) {
	totalCount, err = ws.workspaceRepo.GetWorkspaceCountForUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// CountWorkspaceMembers counts the number of active and pending members for a workspace
func (ws *workspaceService) CountWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) (activeCount, pendingCount int, err error) {
	activeCount, pendingCount, err = ws.workspaceRepo.GetActivePendingMemberCounts(ctx, workspaceID)
	if err != nil {
		return 0, 0, err
	}

	return activeCount, pendingCount, nil
}

// CountWorkspaceCompetitors counts the number of competitors for a workspace
func (ws *workspaceService) CountWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	count, err := ws.competitorService.CountCompetitorsForWorkspace(ctx, workspaceID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountWorkspacePages counts the number of pages for a workspace
func (ws *workspaceService) CountWorkspacePages(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	competitors, _, err := ws.ListCompetitorsForWorkspace(ctx, workspaceID, nil, nil)
	if err != nil {
		return 0, err
	}

	competitorsIDs := make([]uuid.UUID, 0)
	for _, competitor := range competitors {
		competitorsIDs = append(competitorsIDs, competitor.ID)
	}

	count, err := ws.competitorService.CountPagesForCompetitors(ctx, competitorsIDs)
	if err != nil {
		return 0, err
	}

	return count, nil
}
