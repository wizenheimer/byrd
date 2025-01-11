// ./src/internal/interfaces/repository/user.go
package interfaces

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// UserRepository interacts with the user table and the workspace_user table
// This is used to interact with the user repository

type UserRepository interface {
	// ----- CRUD operations for user and workspace_user table -----

	// --------------------------------------------------
	// --- Create & Update operations for user table ---
	// --------------------------------------------------
	// GetUserByUserID gets a user by UserID
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, errs.Error)

	// GetUserByEmail gets a user by email
	GetUserByEmail(ctx context.Context, email string) (models.User, errs.Error)

	// GetClerkUser gets a user by clerkID or clerkEmail
	// It attempts to get a user by clerkID first
	// If the user does not exist, it attempts to get the user by clerkEmail
	GetClerkUser(ctx context.Context, clerkID string, clerkEmail string) (models.User, errs.Error)

	// GetOrCreateUser creates a user if it does not exist
	// This is triggered when a user signs up
	GetOrCreateUser(ctx context.Context, partialUser models.User) (models.User, errs.Error)

	// GetOrCreateUserByEmail creates a user if it does not exist
	// This is triggered when a user is invited
	GetOrCreateUserByEmail(ctx context.Context, emails []string) ([]models.User, errs.Error)

	// -----------------------------------------------------------
	// --- Create & Update operations for workspace_user table ---
	// -----------------------------------------------------------

	// AddUsersToWorkspace adds a batch of users to workspace
	AddUsersToWorkspace(ctx context.Context, workspaceUserProps []models.WorkspaceUserProps, workspaceID uuid.UUID) ([]models.WorkspaceUser, errs.Error)

	// RemoveUsersFromWorkspace removes a batch of users from a workspace
	RemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) errs.Error

	// GetWorkspaceUser gets a user from the workspace if it exists
	GetWorkspaceUser(ctx context.Context, workspaceID, userID uuid.UUID) (models.WorkspaceUser, errs.Error)

	// GetWorkspaceUserByClerkID gets a user from the workspace by clerkID
	GetWorkspaceClerkUser(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (models.WorkspaceUser, errs.Error)

	// ListWorkspaceUser lists all users from the workspace
	ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, errs.Error)

	// ListUserWorkspaces lists all workspaces of a user
	ListUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, errs.Error)

	// UpdateWorkspaceUserRole updates the role of a batch of users in the workspace
	UpdateWorkspaceUserRole(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, role models.UserWorkspaceRole) ([]models.UserWorkspaceRole, errs.Error)

	// UpdateWorkspaceUserStatus  updates the status of a batch of users in the workspace
	UpdateWorkspaceUserStatus(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, status models.UserWorkspaceStatus) ([]models.UserWorkspaceStatus, errs.Error)

	//  GetWorkspaceUserCountByRole gets the count of users by role in the workspace
	GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, errs.Error)

	// ------------------------------------------
	// ----- Sync operations for user table -----
	// -------------------------------------------

	// SyncUser syncs the user data with Clerk
	// This is triggered when a user updates their profile
	SyncUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) errs.Error

	// DeleteUser deletes a user from the user table
	// and removes the user from all workspaces
	DeleteUser(ctx context.Context, userID uuid.UUID) errs.Error

	// ------------------------------------------
	// -----  Optimized Lookup Operations -----
	// ------------------------------------------

	// UserExists checks if a user exists
	// Optimized for quick lookups over the user table
	UserExists(ctx context.Context, userID uuid.UUID) (bool, errs.Error)

	// ClerkUserExists checks if a user exists by clerkID or clerkEmail
	// Optimized for quick lookups over the user table
	ClerkUserExists(ctx context.Context, clerkID, clerkEmail string) (bool, errs.Error)

	// WorkspaceUserExists checks if a user exists in the workspace
	// Optimized for quick lookups over the workspace_user table
	WorkspaceUserExists(ctx context.Context, workspaceID, userID uuid.UUID) (bool, errs.Error)

	// WorkspaceClerkUserExists checks if a user exists in the workspace by clerkID or clerkEmail
	// Optimized for quick lookups over the workspace_user table
	WorkspaceClerkUserExists(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (bool, errs.Error)

	// WorkspaceUserIsAdmin checks if a user exists in the workspace and is an admin
	// Optimized for quick lookups over the workspace_user table
	ClerkUserIsAdmin(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, errs.Error)

	// WorkspaceUserIsMember checks if a user exists in the workspace
	// Optimized for quick lookups over the workspace_user table
	ClerkUserIsMember(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, errs.Error)
}

var (
	// ErrFailedToGetUserFromUserRepository is an error that indicates that the user could not be fetched
	ErrFailedToGetUserFromUserRepository = errors.New("failed to get user from user repository")

	// ErrFailedToGetClerkUserFromUserRepository is an error that indicates that the clerk user could not be fetched
	ErrFailedToGetClerkUserFromUserRepository = errors.New("failed to get clerk user from user repository")

	// ErrFailedToGetOrCreateUserFromUserRepository is an error that indicates that the user could not be created or fetched
	ErrFailedToGetOrCreateUserFromUserRepository = errors.New("failed to get or create user from user repository")

	// ErrFailedToAddUsersToWorkspaceInUserRepository is an error that indicates that the users could not be added
	ErrFailedToAddUsersToWorkspaceInUserRepository = errors.New("failed to add users to workspace in user repository")

	// ErrFailedToRemoveUsersFromWorkspaceInUserRepository is an error that indicates that the users could not be removed
	ErrFailedToRemoveUsersFromWorkspaceInUserRepository = errors.New("failed to remove users from workspace in user repository")

	// ErrFailedToListWorkspaceUsersFromUserRepository is an error that indicates that the workspace users could not be listed
	ErrFailedToListWorkspaceUsersFromUserRepository = errors.New("failed to list workspace users from user repository")

	// ErrFailedToListUserWorkspacesFromUserRepository is an error that indicates that the user workspaces could not be listed
	ErrFailedToListUserWorkspacesFromUserRepository = errors.New("failed to list user workspaces from user repository")

	// ErrFailedToUpdateWorkspaceUserRoleInUserRepository is an error that indicates that the workspace user role could not be updated
	ErrFailedToUpdateWorkspaceUserRoleInUserRepository = errors.New("failed to update workspace user role in user repository")

	// ErrFailedToUpdateWorkspaceUserStatusInUserRepository is an error that indicates that the workspace user status could not be updated
	ErrFailedToUpdateWorkspaceUserStatusInUserRepository = errors.New("failed to update workspace user status in user repository")

	// ErrFailedToGetWorkspaceUserCountByRoleFromUserRepository is an error that indicates that the workspace user count by role could not be fetched
	ErrFailedToGetWorkspaceUserCountByRoleFromUserRepository = errors.New("failed to get workspace user count by role from user repository")

    // ErrFailedToSyncUserInUserRepository is an error that indicates that the user could not be synced
    ErrFailedToSyncUserInUserRepository = errors.New("failed to sync user in user repository")

    // ErrFailedToDeleteUserInUserRepository is an error that indicates that the user could not be deleted
    ErrFailedToDeleteUserInUserRepository = errors.New("failed to delete user in user repository")

    // ErrFailedToCheckIfUserExistsInUserRepository is an error that indicates that the user existence could not be checked
    ErrFailedToCheckIfUserExistsInUserRepository = errors.New("failed to check if user exists in user repository")

    // ErrFailedToCheckIfClerkUserExistsInUserRepository is an error that indicates that the clerk user existence could not be checked
    ErrFailedToCheckIfClerkUserExistsInUserRepository = errors.New("failed to check if clerk user exists in user repository")

    // ErrFailedToCheckUserRoleInUserRepository is an error that indicates that the user role could not be checked
    ErrFailedToCheckUserRoleInUserRepository = errors.New("failed to check user role in user repository")
)
