// ./src/internal/interfaces/service/workspace.go
package interfaces

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// Note: WorkspaceService is user-facing and handler-owned service
// It's an entry point for all workspace-related operations

// It embeds other services and repositories
// Repository:
// - WorkspaceRepository
// Service:
// - UserService
// - CompetitorService (embeds PageService)
// - BillingService

type WorkspaceService interface {
	// <---------  Workspace Management --------->
	// CRUD operations for workspaces
	// CRUD operations for workspace members

	// CreateWorkspace creates a new workspace for a user
	// If the user exists, it creates a new workspace and returns it
	// If the user does not exist, it creates a new user and a new workspace and returns the workspace
	CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, errs.Error)

	// ListUserWorkspaces lists the workspaces of a clerk user
	// It returns the workspaces of the clerk user
	// It returns an error if the user does not exist
	ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, errs.Error)

	// GetWorkspace gets a workspace by ID
	// It returns the workspace if it exists, otherwise it returns an error
	GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, errs.Error)

	// UpdateWorkspace updates a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the workspace is inactive
	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) errs.Error

	// DeleteWorkspace deletes a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the workspace is inactive
	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, errs.Error)

	// ListWorkspaceMembers gets the members of a workspace
	// It returns the members of the workspace
	// It returns an error if the workspace does not exist
	// If includeMembers is true, it includes members in the response
	// If includeAdmins is true, it includes admins in the response
	// If both includeMembers and includeAdmins are false, it returns an empty list
	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, errs.Error)

	// InviteUserToWorkspace adds a user to a workspace
	// If the user does not exist, it creates a new user and adds it to the workspace
	// If the user exists, it adds it to the workspace
	InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, errs.Error)

	// LeaveWorkspace exits clerk user from a workspace
	// It returns an error if the user does not exist in the workspace
	// It returns an error if the user is the last admin in the workspace
	LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) errs.Error

	// UpdateWorkspaceMemberRole updates member role in a workspace
	// It returns an error if the user does not exist in the workspace
	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) errs.Error

	// RemoveUserFromWorkspace removes a user from a workspace
	// It returns an error if the user does not exist in the workspace
	RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) errs.Error

	// JoinWorkspace when a user accepts an invite to join a workspace
	// It returns an error if the user does not exist in the workspace
	// It returns an error if the user is not invited to the workspace
	JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) errs.Error

	// <--------- Workspace x Middleware --------->

	// WorkspaceExists checks if a workspace exists and is active
	// It returns true if the workspace exists and is active
	// It returns false if the workspace does not exist or is not active
	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, errs.Error)

	// ClerkUserIsWorkspaceAdmin checks if a user is an admin in a workspace
	// It returns true if the user is an admin
	// It returns false if the user is not an admin
	ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, errs.Error)

	// ClerkUserIsMember checks if a user is a member in a workspace
	// It returns true if the user is a member
	// It returns false if the user is not a member
	ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, errs.Error)

	// WorkspaceCompetitorExists checks if a competitor exists and is active
	// It returns true if the competitor exists and is active
	// It returns false if the competitor does not exist or is not active
	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, errs.Error)

	// PageExists checks if a page exists and is active
	// It returns true if the page exists and is active
	// It returns false if the page does not exist or is not active
	WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, errs.Error)

	// <--------- Workspace Competitor Management --------->
	// CRUD operations for competitors
	// CRUD operations for pages
	// Read operations for page history

	// CreateWorkspaceCompetitor adds a new competitor to a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) errs.Error

	// AddPageToCompetitor adds a new page to an existing competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pages []api.CreatePageRequest) ([]models.Page, errs.Error)

	// ListPages lists the competitors of a workspace with tracked pages for each competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	// pagination params applies to competitors (not pages)
	ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, errs.Error)

	// ListWorkspacePageHistory lists the history of a page
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	// It returns an error if the page does not exist in the competitor of the workspace
	// pagination params applies to page history
	ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, errs.Error)

	// RemovePage removes a page from a competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) errs.Error

	// RemoveCompetitor removes a competitor from a workspace, including all its pages
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) errs.Error

	// UpdatePage updates a page
	// It returns an error if the workspace does not exist or page is not active
	// It returns an error if the user is not an member of the workspace
	// It returns an error if the page does not exist in the competitor of the workspace
	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) errs.Error
}

var (
	// ErrFailedToCreateWorkspace is returned when the workspace creation fails
	ErrFailedToCreateWorkspace = errors.New("failed to create workspace")

	// ErrFailedToListUserWorkspaces is returned when the user workspace listing fails
	ErrFailedToListUserWorkspaces = errors.New("failed to list user workspaces")

	// ErrFailedToGetWorkspace is returned when the workspace retrieval fails
	ErrFailedToGetWorkspace = errors.New("failed to get workspace")

	// ErrFailedToUpdateWorkspace is returned when the workspace update fails
	ErrFailedToUpdateWorkspace = errors.New("failed to update workspace")

	// ErrFailedToUpdateWorkspaceMemberStatus is returned when the workspace member status update fails
	ErrFailedToUpdateWorkspaceMemberStatus = errors.New("failed to update workspace member status")

	// ErrFailedToUpdateWorkspaceMemberRole is returned when the workspace member role update fails
	ErrFailedToUpdateWorkspaceMemberRole = errors.New("failed to update workspace member role")

	// ErrFailedToDeleteWorkspace is returned when the workspace deletion fails
	ErrFailedToDeleteWorkspace = errors.New("failed to delete workspace")

	// ErrFailedToListWorkspaceMembers is returned when the workspace member listing fails
	ErrFailedToListWorkspaceMembers = errors.New("failed to list workspace members")

	// ErrFailedToJoinWorkspace is returned when the user joining the workspace fails
	ErrFailedToJoinWorkspace = errors.New("failed to join workspace")

	// ErrFailedToInviteUserToWorkspace is returned when the user invitation to workspace fails
	ErrFailedToInviteUserToWorkspace = errors.New("failed to invite user to workspace")

	// ErrUserNotInvitedToWorkspace is returned when the user is not invited to the workspace
	ErrUserNotInvitedToWorkspace = errors.New("user not invited to workspace")

	// ErrFailedToLeaveWorkspace is returned when the user leaving the workspace fails
	ErrFailedToLeaveWorkspace = errors.New("failed to leave workspace")

	// ErrFailedToLeaveWorkspaceAsSoleAdmin is returned when the user leaving the workspace as the sole admin fails
	ErrFailedToLeaveWorkspaceAsSoleAdmin = errors.New("failed to leave workspace as the sole admin")

	// ErrFailedToRemoveUserFromWorkspace is returned when the user removal from workspace fails
	ErrFailedToRemoveUserFromWorkspace = errors.New("failed to remove user from workspace")

	// ErrWorkspaceNotFound is returned when the workspace is not found
	ErrWorkspaceNotFound = errors.New("workspace not found")

	// ErrWorkspaceUserNotFound is returned when the workspace user is not found
	ErrWorkspaceUserNotFound = errors.New("workspace user not found")

	// ErrWorkspaceMemberNotFound is returned when the workspace member is not found
	ErrWorkspaceMemberNotFound = errors.New("workspace member not found")

	// ErrWorkspaceAdminNotFound is returned when the workspace admin is not found
	ErrWorkspaceAdminNotFound = errors.New("workspace admin not found")

	// ErrWorkspaceCompetitorNotFound is returned when the workspace competitor is not found
	ErrWorkspaceCompetitorNotFound = errors.New("workspace competitor not found")

	// ErrCompetitorPageNotFound is returned when the competitor page is not found
	ErrCompetitorPageNotFound = errors.New("competitor page not found")

	// ErrFailedToCreateCompetitor is returned when the competitor creation fails
	ErrFailedToCreateWorkspaceCompetitor = errors.New("failed to create workspace competitor")

	// ErrWorkspacceInactive is returned when the workspace is inactive
	ErrWorkspaceInactive = errors.New("workspace is inactive")

	// ErrWorkspaceLookupEmpty is returned when the workspace lookup is empty
	ErrWorkspaceLookupEmpty = errors.New("workspace lookup is empty")

	// ErrUserLeftWorkspaceNeedReinvite is returned when the user left the workspace and needs to be reinvited
	ErrUserLeftWorkspaceNeedReinvite = errors.New("user left workspace and needs to be reinvited")

    // ErrFailedToAddPageToCompetitor is returned when the page addition to competitor fails
    ErrFailedToAddPageToCompetitor = errors.New("failed to add page to competitor")

    // ErrFailedToListWorkspacePageHistory is returned when the page history listing fails
    ErrFailedToListWorkspacePageHistory = errors.New("failed to list workspace page history")

    // ErrFailedToRemovePageFromWorkspace is returned when the page removal from workspace fails
    ErrFailedToRemovePageFromWorkspace = errors.New("failed to remove page from workspace")

    // ErrFailedToRemoveCompetitorFromWorkspace is returned when the competitor removal from workspace fails
    ErrFailedToRemoveCompetitorFromWorkspace = errors.New("failed to remove competitor from workspace")

    // ErrFailedToUpdateCompetitorPage is returned when the competitor page update fails
    ErrFailedToUpdateCompetitorPage = errors.New("failed to update competitor page")

    // ErrFailedToGetWorkspaceUserCountByRole is returned when the workspace user count by role retrieval fails
    ErrFailedToGetWorkspaceUserCountByRole = errors.New("failed to get workspace user count by role")
)
