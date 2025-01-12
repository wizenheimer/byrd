package user

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// UserService is a service that manages users
// It holds the business logic for user management
// It embeds UserRespository to interact with the database
// UserRepository inturn interacts primarily with 2 tables: users and workspace_users
type UserService interface {
	// ---- User Management ----

	GetOrCreateUser(ctx context.Context, clerk *clerk.User) (*models.User, error)

	BatchGetOrCreateUsers(ctx context.Context, emails []string) ([]models.User, error)

	ListUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	GetUserByClerkCredentials(ctx context.Context, clerk *clerk.User) (*models.User, error)

	SyncUser(ctx context.Context, clerk *clerk.User) error

	DeleteUser(ctx context.Context, clerk *clerk.User) error

	UserExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)

	ClerkUserExists(ctx context.Context, clerk *clerk.User) (bool, error)
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
