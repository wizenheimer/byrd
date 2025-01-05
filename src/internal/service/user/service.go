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

	user, err := us.userRepository.GetOrCreateUser(ctx, partialUser)
	if err != nil {
		return nil, ErrFailedToGetUser
	}

	addUsers := []models.WorkspaceUserProps{
		{
			Email:  email,
			Role:   models.UserRoleAdmin,
			Status: models.UserWorkspaceStatusActive,
		},
	}

	_, errs := us.userRepository.AddUsersToWorkspace(ctx, addUsers, workspaceID)
	if len(errs) > 0 {
		return nil, ErrFailedToAddUserToWorkspace
	}

	return &user, nil
}

func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) []api.CreateWorkspaceUserResponse {
	if len(invitedUsers) == 0 {
		return []api.CreateWorkspaceUserResponse{}
	}

	var responses []api.CreateWorkspaceUserResponse
	workspaceUsers, errs := us.userRepository.AddUsersToWorkspace(ctx, invitedUsers, workspaceID)
	if len(errs) > 0 {
		for range errs {
			responses = append(responses, api.CreateWorkspaceUserResponse{
				Error: ErrFailedToAddUserToWorkspace,
			})
		}
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

	workspaceUser, err := us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, userEmail)
	if err != nil {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	workspaceUser, err := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if err != nil {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	workspaceUsers, errs := us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
	if len(errs) > 0 {
		return nil, ErrFailedToListWorkspaceUsers
	}

	return workspaceUsers, nil
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error) {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, ErrFailedToGetUserEmail
	}

	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return nil, ErrFailedToGetUser
	}

	workspaceIDs, errs := us.userRepository.ListUserWorkspaces(ctx, user.ID)
	if len(errs) > 0 {
		return nil, ErrFailedToListWorkspaceForUser
	}

	return workspaceIDs, nil
}

func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, error) {
	userIDs := []uuid.UUID{userID}
	_, errs := us.userRepository.UpdateWorkspaceUserRole(ctx, workspaceID, userIDs, role)
	if len(errs) > 0 {
		return models.WorkspaceUser{}, ErrFailedToUpdateWorkspaceUserRole
	}

	workspaceUser, err := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if err != nil {
		return models.WorkspaceUser{}, ErrFailedToGetUser
	}

	return workspaceUser, nil
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) error {
	userIDs := []uuid.UUID{userID}
	_, errs := us.userRepository.UpdateWorkspaceUserStatus(ctx, workspaceID, userIDs, status)
	if len(errs) > 0 {
		return ErrFailedToUpdateWorkspaceUserStatus
	}

	return nil
}

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	errs := us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
	if len(errs) > 0 {
		return []error{ErrFailedToRemoveWorkspaceUsers}
	}

	return nil
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error) {
	admin, member, err := us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
	if err != nil {
		return 0, 0, ErrFailedToGetWorkspaceUserRoleCount
	}

	return admin, member, nil
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return ErrFailedToGetUserEmail
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return ErrFailedToGetUser
	}

	if err := us.userRepository.SyncUser(ctx, user.ID, clerk); err != nil {
		return ErrFailedToSyncUser
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return ErrFailedToGetUserEmail
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return ErrFailedToGetUser
	}

	if err := us.userRepository.DeleteUser(ctx, user.ID); err != nil {
		return ErrFailedToDeleteUser
	}

	return nil
}
