// ./src/internal/service/workspace/service.go
package workspace

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workspace"
	"github.com/wizenheimer/byrd/src/internal/service/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceService struct {
	workspaceRepo     workspace.WorkspaceRepository
	competitorService competitor.CompetitorService
	userService       user.UserService
	logger            *logger.Logger
}

// compile time check if the interface is implemented
var _ WorkspaceService = (*workspaceService)(nil)

func NewWorkspaceService(workspaceRepo workspace.WorkspaceRepository, competitorService competitor.CompetitorService, userService user.UserService, logger *logger.Logger) WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		logger:            logger,
	}
}

// CreateWorkspace creates a new workspace for a user
// If the user exists, it creates a new workspace and returns it
// If the user does not exist, it creates a new user and a new workspace and returns the workspace
func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, error) {
	return nil, nil
}

// ListUserWorkspaces lists the workspaces of a clerk user
// It returns the workspaces of the clerk user
// It returns an error if the user does not exist
func (ws *workspaceService) ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error) {
	return nil, nil
}

// GetWorkspace gets a workspace by ID
// It returns the workspace if it exists, otherwise it returns an error
func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	return nil, nil
}

// UpdateWorkspace updates a workspace
// It returns an error if the workspace does not exist
// It returns an error if the workspace is inactive
func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) error {
	return nil
}

// DeleteWorkspace deletes a workspace
// It returns an error if the workspace does not exist
// It returns an error if the workspace is inactive
func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error) {
	return models.WorkspaceStatusActive, nil
}

// ListWorkspaceMembers gets the members of a workspace
// It returns the members of the workspace
// It returns an error if the workspace does not exist
// If includeMembers is true, it includes members in the response
// If includeAdmins is true, it includes admins in the response
// If both includeMembers and includeAdmins are false, it returns an empty list
func (ws *workspaceService) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, error) {
	return nil, nil
}

// InviteUserToWorkspace adds a user to a workspace
// If the user does not exist, it creates a new user and adds it to the workspace
// If the user exists, it adds it to the workspace
func (ws *workspaceService) InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, error) {
	return nil, nil
}

// LeaveWorkspace exits clerk user from a workspace
// It returns an error if the user does not exist in the workspace
// It returns an error if the user is the last admin in the workspace
func (ws *workspaceService) LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error {
	return nil
}

// UpdateWorkspaceMemberRole updates member role in a workspace
// It returns an error if the user does not exist in the workspace
func (ws *workspaceService) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) error {
	return nil
}

// RemoveUserFromWorkspace removes a user from a workspace
// It returns an error if the user does not exist in the workspace
func (ws *workspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error {
	return nil
}

// JoinWorkspace when a user accepts an invite to join a workspace
// It returns an error if the user does not exist in the workspace
// It returns an error if the user is not invited to the workspace
func (ws *workspaceService) JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error {
	return nil
}

// <--------- Workspace x Middleware --------->

// WorkspaceExists checks if a workspace exists and is active
// It returns true if the workspace exists and is active
// It returns false if the workspace does not exist or is not active
func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	return false, nil
}

// ClerkUserIsWorkspaceAdmin checks if a user is an admin in a workspace
// It returns true if the user is an admin
// It returns false if the user is not an admin
func (ws *workspaceService) ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	return false, nil
}

// ClerkUserIsMember checks if a user is a member in a workspace
// It returns true if the user is a member
// It returns false if the user is not a member
func (ws *workspaceService) ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error) {
	return false, nil
}

// WorkspaceCompetitorExists checks if a competitor exists and is active
// It returns true if the competitor exists and is active
// It returns false if the competitor does not exist or is not active
func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	return false, nil
}

// PageExists checks if a page exists and is active
// It returns true if the page exists and is active
// It returns false if the page does not exist or is not active
func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error) {
	return false, nil
}

// <--------- Workspace Competitor Management --------->
// CRUD operations for competitors
// CRUD operations for pages
// Read operations for page history

// CreateWorkspaceCompetitor adds a new competitor to a workspace
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
func (ws *workspaceService) CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) error {
	return nil
}

// AddPageToCompetitor adds a new page to an existing competitor
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
func (ws *workspaceService) AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pages []api.CreatePageRequest) ([]models.Page, error) {
	return nil, nil
}

// ListPages lists the competitors of a workspace with tracked pages for each competitor
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
// pagination params applies to competitors (not pages)
func (ws *workspaceService) ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, error) {
	return nil, nil
}

// ListWorkspacePageHistory lists the history of a page
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
// It returns an error if the page does not exist in the competitor of the workspace
// pagination params applies to page history
func (ws *workspaceService) ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error) {
	return nil, nil
}

// RemovePage removes a page from a competitor
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
func (ws *workspaceService) RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) error {
	return nil
}

// RemoveCompetitor removes a competitor from a workspace, including all its pages
// It returns an error if the workspace does not exist
// It returns an error if the user is not an member of the workspace
func (ws *workspaceService) RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) error {
	return nil
}

// UpdatePage updates a page
// It returns an error if the workspace does not exist or page is not active
// It returns an error if the user is not an member of the workspace
// It returns an error if the page does not exist in the competitor of the workspace
func (ws *workspaceService) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) error {
	return nil
}
