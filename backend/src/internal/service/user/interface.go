package user

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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
	AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, error)

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
	RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error

	// GetWorkspaceUserCount gets the count of users in a workspace
	// It returns the count of users in the workspace
	// It returns an error if the workspace does not exist
	// This is triggered when workspace owner or member wants to get the count of users in a workspace
	GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.UserWorkspaceRole]int, error)

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

var (
	// ErrFailedToGetClerkUserEmail is an error that indicates that the clerk user email could not be fetched
	ErrFailedToGetClerkUserEmail = errors.New("failed to get clerk user email")

	// ErrFailedToGetUserEmail is an error that indicates that the user email could not be fetched
	ErrFailedToGetUserEmail = errors.New("failed to get user email")

	// ErrFailedToCreateWorkspaceOwner is an error that indicates that the workspace owner could not be created
	ErrFailedToCreateWorkspaceOwner = errors.New("failed to create workspace owner")

	ErrFailedToGetUser = errors.New("failed to get user")

	ErrFailedToListWorkspaceUsers = errors.New("failed to list workspace users")

	ErrFailedToListWorkspaceForUser = errors.New("failed to list workspace for user")

	ErrFailedToCreateUser = errors.New("failed to create user")

	ErrFailedToAddUserToWorkspace = errors.New("failed to add user to workspace")

	ErrFailedToUpdateWorkspaceUserRole = errors.New("failed to update workspace user role")

	ErrFailedToUpdateWorkspaceUserStatus = errors.New("failed to update workspace user status")

	ErrFailedToRemoveWorkspaceUsers = errors.New("failed to remove workspace users")

	ErrFailedToGetWorkspaceUserRoleCount = errors.New("failed to get workspace user role count")

	ErrFailedToSyncUser = errors.New("failed to sync user")

	ErrFailedToDeleteUser = errors.New("failed to delete user")

	ErrInvitedUserCountShouldBeNonZero = errors.New("invited user count should be non-zero")
)
