package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/user"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

var _ UserService = (*userService)(nil)

// TODO: rethink retrieval methods
type userService struct {
	userRepository user.UserRepository
	logger         *logger.Logger
}

func NewUserService(userRepository user.UserRepository, logger *logger.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

// CreateWorkspaceOwner creates a owner in a workspace if it does not exist
// Once the owner is created or found, it returns the owner's user model
// It returns an error if the user could not be created
// This is triggered when workspace owner creates a workspace
func (us *userService) CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, error) {
	return nil, nil
}

// AddUserToWorkspace adds a user to a workspace
// It returns an error if the user could not be added to the workspace
// It returns nil if the user was added successfully
// This is triggered when workspace owner invites a user to a workspace
func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, error) {
	return nil, nil
}

// GetWorkspaceUser gets a user from a workspace
// Once the user is found, it returns the user
// It returns an error if the user could not be found
// This is triggered whenever signs in to the application
func (us *userService) GetWorkspaceUser(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	return models.WorkspaceUser{}, nil
}

// GetWorkspaceUserByID gets a user from a workspace by ID
// Once the user is found, it returns the user
// This is triggered when workspace owner or member wants to get a user by ID
func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	return models.WorkspaceUser{}, nil
}

// ListWorkspaceUsers lists the users of a workspace
// It returns the users of the workspace
// It returns an error if the workspace does not exist
func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	return nil, nil
}

// ListUserWorkspaces lists the workspaces of a user
// It returns the workspaces of the user
// It returns an error if the user does not exist
// This is triggered when the user signs in to the application
func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error) {
	return nil, nil
}

// UpdateWorkspaceUserRole updates a user role in a workspace
// It returns an error if the user could not be updated
// It returns nil if the user was updated successfully
// It can update the user's role in the workspace
func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, error) {
	return models.WorkspaceUser{}, nil
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) error {
	return nil
}

// RemoveUserFromWorkspace removes users from a workspace
// It returns an error if the users could not be removed
// It returns nil if the users were removed successfully
// When the userIDs are nil, it removes all users from the workspace
func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error {
	return nil
}

// GetWorkspaceUserCount gets the count of users in a workspace
// It returns the count of users in the workspace
// It returns an error if the workspace does not exist
// This is triggered when workspace owner or member wants to get the count of users in a workspace
func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.UserWorkspaceRole]int, error) {
	return nil, nil
}

// <------ User Management ------>

// SyncUser syncs a user with Clerk
// It returns an error if the user could not be synced
// When sync is triggered, it marks the account status as active
// And updates the user's email and name if they have changed
func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	return nil
}

// DeleteUser deletes a user from Clerk
// It returns an error if the user could not be deleted
// It returns nil if the user was deleted successfully
// This is the only user-facing and handler-owned method
func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	return nil
}
