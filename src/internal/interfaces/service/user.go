package interfaces

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// UserService is a service that manages users
// It holds the business logic for user management
// It embeds UserRespository to interact with the database
// UserRepository inturn interacts primarily with 2 tables: users and workspace_users
type UserService interface {
	// <------ Workspace User Management ------>

	// CreateWorkspaceOwner creates a owner in a workspace if it does not exist
	// Once the owner is created or found, it returns the owner's user model
	// It returns an error if the user could not be created
	// This is triggered when workspace owner creates a workspace
	CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, error)

	// AddUserToWorkspace adds a user to a workspace
	// It returns an error if the user could not be added to the workspace
	// It returns nil if the user was added successfully
	// This is triggered when workspace owner invites a user to a workspace
	AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) []api.CreateWorkspaceUserResponse

	// GetWorkspaceUser gets a user from a workspace
	// Once the user is found, it returns the user
	// It returns an error if the user could not be found
	// This is triggered whenever signs in to the application
	GetWorkspaceUser(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (models.WorkspaceUser, error)

	// GetWorkspaceUserByID gets a user from a workspace by ID
	// Once the user is found, it returns the user
	// This is triggered when workspace owner or member wants to get a user by ID
	GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error)

	// ListWorkspaceUsers lists the users of a workspace
	// It returns the users of the workspace
	// It returns an error if the workspace does not exist
	ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error)

	// ListUserWorkspaces lists the workspaces of a user
	// It returns the workspaces of the user
	// It returns an error if the user does not exist
	// This is triggered when the user signs in to the application
	ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error)

	// UpdateWorkspaceUserRole updates a user role in a workspace
	// It returns an error if the user could not be updated
	// It returns nil if the user was updated successfully
	// It can update the user's role in the workspace
	UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, error)

	UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) error

	// RemoveUserFromWorkspace removes users from a workspace
	// It returns an error if the users could not be removed
	// It returns nil if the users were removed successfully
	// When the userIDs are nil, it removes all users from the workspace
	RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error

	// GetWorkspaceUserCount gets the count of users in a workspace
	// It returns the count of users in the workspace
	// It returns an error if the workspace does not exist
	// This is triggered when workspace owner or member wants to get the count of users in a workspace
	GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error)

	// <------ User Management ------>

	// SyncUser syncs a user with Clerk
	// It returns an error if the user could not be synced
	// When sync is triggered, it marks the account status as active
	// And updates the user's email and name if they have changed
	SyncUser(ctx context.Context, clerk *clerk.User) error

	// DeleteUser deletes a user from Clerk
	// It returns an error if the user could not be deleted
	// It returns nil if the user was deleted successfully
	// This is the only user-facing and handler-owned method
	DeleteUser(ctx context.Context, clerk *clerk.User) error
}
