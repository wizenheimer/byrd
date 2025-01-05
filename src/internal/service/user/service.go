package user

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

var (
	ErrFailedToGetUserEmail              = errors.New("failed to get user email")
	ErrFailedToCreateWorkspaceOwner      = errors.New("failed to create workspace owner")
	ErrFailedToGetUser                   = errors.New("failed to get user")
	ErrFailedToListWorkspaceUsers        = errors.New("failed to list workspace users")
	ErrFailedToListWorkspaceForUser      = errors.New("failed to list workspace for user")
	ErrFailedToCreateUser                = errors.New("failed to create user")
	ErrFailedToAddUserToWorkspace        = errors.New("failed to add user to workspace")
	ErrFailedToUpdateWorkspaceUserRole   = errors.New("failed to update workspace user role")
	ErrFailedToUpdateWorkspaceUserStatus = errors.New("failed to update workspace user status")
	ErrFailedToRemoveWorkspaceUsers      = errors.New("failed to remove workspace users")
	ErrFailedToGetWorkspaceUserRoleCount = errors.New("failed to get workspace user role count")
	ErrFailedToSyncUser                  = errors.New("failed to sync user")
	ErrFailedToDeleteUser                = errors.New("failed to delete user")
)

func NewUserService(userRepository repo.UserRepository, logger *logger.Logger) svc.UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (us *userService) CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, error) {
	email, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, ErrFailedToCreateWorkspaceOwner
	}

	name := utils.GetClerkUserFullName(clerk)

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, ErrFailedToCreateWorkspaceOwner
	}

	partialUser := models.User{
		ID:      userID,
		ClerkID: &clerk.ID,
		Name:    &name,
		Email:   &email,
		Status:  models.AccountStatusActive,
	}

	user, uErr := us.userRepository.GetOrCreateUser(ctx, partialUser)
	if uErr != nil && uErr.HasErrors() {
		return nil, ErrFailedToGetUser
	}

	addUsers := []models.WorkspaceUserProps{
		{
			Email:  email,
			Role:   models.UserRoleAdmin,
			Status: models.UserWorkspaceStatusActive,
		},
	}

	if _, uErr := us.userRepository.AddUsersToWorkspace(ctx, addUsers, workspaceID); uErr != nil && uErr.HasErrors() {
		return nil, ErrFailedToAddUserToWorkspace
	}

	return &user, nil
}

func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) []api.CreateWorkspaceUserResponse {
	if len(invitedUsers) == 0 {
		return []api.CreateWorkspaceUserResponse{}
	}

	var responses []api.CreateWorkspaceUserResponse
	workspaceUsers, uErr := us.userRepository.AddUsersToWorkspace(ctx, invitedUsers, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		// TODO: This is a bug. Fix this during service refactor
		responses = append(responses, api.CreateWorkspaceUserResponse{
			Error: ErrFailedToAddUserToWorkspace,
		})
	}

	for _, wu := range workspaceUsers {
		responses = append(responses, api.CreateWorkspaceUserResponse{
			Error: nil,
			User:  &wu,
		})
	}

	return responses
}

func (us *userService) GetWorkspaceUser(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	userEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, userEmail)
	if uErr != nil && uErr.HasErrors() {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	workspaceUsers, uErr := us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		return nil, ErrFailedToListWorkspaceUsers
	}

	return workspaceUsers, nil
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error) {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, ErrFailedToGetUserEmail
	}

	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		return nil, ErrFailedToGetUser
	}

	workspaceIDs, uErr := us.userRepository.ListUserWorkspaces(ctx, user.ID)
	if uErr != nil && uErr.HasErrors() {
		return nil, ErrFailedToListWorkspaceForUser
	}

	return workspaceIDs, nil
}

func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, error) {
	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserRole(ctx, workspaceID, userIDs, role)
	if uErr != nil && uErr.HasErrors() {
		return models.WorkspaceUser{}, ErrFailedToUpdateWorkspaceUserRole
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) error {
	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserStatus(ctx, workspaceID, userIDs, status)
	if uErr != nil && uErr.HasErrors() {
		return ErrFailedToUpdateWorkspaceUserStatus
	}

	return nil
}

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	uErr := us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		return []error{ErrFailedToRemoveWorkspaceUsers}
	}

	return nil
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error) {
	admin, member, uErr := us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		return 0, 0, ErrFailedToGetWorkspaceUserRoleCount
	}

	return admin, member, nil
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return ErrFailedToGetUserEmail
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		return ErrFailedToGetUser
	}

	if uErr := us.userRepository.SyncUser(ctx, user.ID, clerk); uErr != nil && uErr.HasErrors() {
		return ErrFailedToSyncUser
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return ErrFailedToGetUserEmail
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		return ErrFailedToGetUser
	}

	if uErr := us.userRepository.DeleteUser(ctx, user.ID); uErr != nil && uErr.HasErrors() {
		return ErrFailedToDeleteUser
	}

	return nil
}
