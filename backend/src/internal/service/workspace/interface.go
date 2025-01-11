package workspace

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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
	CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, workspaceReq api.WorkspaceCreationRequest) (*models.Workspace, error)

	// ListUserWorkspaces lists the workspaces of a clerk user
	// It returns the workspaces of the clerk user
	// It returns an error if the user does not exist
	ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error)

	// GetWorkspace gets a workspace by ID
	// It returns the workspace if it exists, otherwise it returns an error
	GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error)

	// UpdateWorkspace updates a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the workspace is inactive
	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req api.WorkspaceUpdateRequest) error

	// DeleteWorkspace deletes a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the workspace is inactive
	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error)

	// ListWorkspaceMembers gets the members of a workspace
	// It returns the members of the workspace
	// It returns an error if the workspace does not exist
	// If includeMembers is true, it includes members in the response
	// If includeAdmins is true, it includes admins in the response
	// If both includeMembers and includeAdmins are false, it returns an empty list
	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, params api.WorkspaceMembersListingParams) ([]models.WorkspaceUser, error)

	// InviteUserToWorkspace adds a user to a workspace
	// If the user does not exist, it creates a new user and adds it to the workspace
	// If the user exists, it adds it to the workspace
	InviteUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, error)

	// LeaveWorkspace exits clerk user from a workspace
	// It returns an error if the user does not exist in the workspace
	// It returns an error if the user is the last admin in the workspace
	LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error

	// UpdateWorkspaceMemberRole updates member role in a workspace
	// It returns an error if the user does not exist in the workspace
	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.UserWorkspaceRole) error

	// RemoveUserFromWorkspace removes a user from a workspace
	// It returns an error if the user does not exist in the workspace
	RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error

	// JoinWorkspace when a user accepts an invite to join a workspace
	// It returns an error if the user does not exist in the workspace
	// It returns an error if the user is not invited to the workspace
	JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error

	// <--------- Workspace x Middleware --------->

	// WorkspaceExists checks if a workspace exists and is active
	// It returns true if the workspace exists and is active
	// It returns false if the workspace does not exist or is not active
	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	// ClerkUserIsWorkspaceAdmin checks if a user is an admin in a workspace
	// It returns true if the user is an admin
	// It returns false if the user is not an admin
	ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

	// ClerkUserIsMember checks if a user is a member in a workspace
	// It returns true if the user is a member
	// It returns false if the user is not a member
	ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

	// WorkspaceCompetitorExists checks if a competitor exists and is active
	// It returns true if the competitor exists and is active
	// It returns false if the competitor does not exist or is not active
	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)

	// PageExists checks if a page exists and is active
	// It returns true if the page exists and is active
	// It returns false if the page does not exist or is not active
	WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error)

	// <--------- Workspace Competitor Management --------->
	// CRUD operations for competitors
	// CRUD operations for pages
	// Read operations for page history

	// CreateWorkspaceCompetitor adds a new competitor to a workspace
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	CreateWorkspaceCompetitor(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, page api.CreatePageRequest) error

	// AddPageToCompetitor adds a new page to an existing competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	AddPageToCompetitor(ctx context.Context, clerkUser *clerk.User, competitorID string, pages []api.CreatePageRequest) ([]models.Page, error)

	// ListPages lists the competitors of a workspace with tracked pages for each competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	// pagination params applies to competitors (not pages)
	ListWorkspaceCompetitors(ctx context.Context, clerkUser *clerk.User, workspaceID uuid.UUID, params api.PaginationParams) ([]api.GetWorkspaceCompetitorResponse, error)

	// ListWorkspacePageHistory lists the history of a page
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	// It returns an error if the page does not exist in the competitor of the workspace
	// pagination params applies to page history
	ListWorkspacePageHistory(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID, pageID uuid.UUID, param api.PaginationParams) ([]models.PageHistory, error)

	// RemovePage removes a page from a competitor
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	RemovePageFromWorkspace(ctx context.Context, clerkUser *clerk.User, competitorID, pageID uuid.UUID) error

	// RemoveCompetitor removes a competitor from a workspace, including all its pages
	// It returns an error if the workspace does not exist
	// It returns an error if the user is not an member of the workspace
	RemoveCompetitorFromWorkspace(ctx context.Context, clerkUser *clerk.User, workspaceID, competitorID uuid.UUID) error

	// UpdatePage updates a page
	// It returns an error if the workspace does not exist or page is not active
	// It returns an error if the user is not an member of the workspace
	// It returns an error if the page does not exist in the competitor of the workspace
	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, req api.UpdatePageRequest) error
}
