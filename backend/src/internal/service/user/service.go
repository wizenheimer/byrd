// ./src/internal/service/user/service.go
package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

func NewUserService(userRepository repo.UserRepository, logger *logger.Logger) svc.UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (us *userService) CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, errs.Error) {
	wErr := errs.New()
	email, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspaceOwner)
	}

	name := utils.GetClerkUserFullName(clerk)

	userID, err := uuid.NewUUID()
	if err != nil {
		wErr.Add(err, map[string]any{"userID": userID})
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspaceOwner)
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
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspaceOwner)
	}

	addUsers := []models.WorkspaceUserProps{
		{
			Email:  email,
			Role:   models.UserRoleAdmin,
			Status: models.UserWorkspaceStatusActive,
		},
	}

	if _, uErr := us.userRepository.AddUsersToWorkspace(ctx, addUsers, workspaceID); uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToCreateWorkspaceOwner)
	}

	return &user, nil
}

func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, errs.Error) {
	wErr := errs.New()
	if len(invitedUsers) == 0 {
		wErr.Add(svc.ErrInvitedUserCountShouldBeNonZero, map[string]any{"error": "no users to add"})
		return nil, wErr.Propagate(svc.ErrFailedToAddUserToWorkspace)
	}

	err := utils.SetDefaultsAndValidateArray(&invitedUsers)
	if err != nil {
		wErr.Add(err, map[string]any{"users": invitedUsers})
		return nil, wErr.Propagate(svc.ErrFailedToAddUserToWorkspace)
	}

	workspaceUsers, uErr := us.userRepository.AddUsersToWorkspace(ctx, invitedUsers, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToAddUserToWorkspace)
	}

	return workspaceUsers, wErr
}

func (us *userService) GetWorkspaceUser(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (models.WorkspaceUser, errs.Error) {
	wErr := errs.New()
	userEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return models.WorkspaceUser{}, wErr.Propagate(svc.ErrWorkspaceUserNotFound)
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, userEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr.Propagate(svc.ErrWorkspaceUserNotFound)
	}

	return workspaceUser, nil
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, errs.Error) {
	wErr := errs.New()
	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr.Propagate(svc.ErrWorkspaceUserNotFound)
	}

	return workspaceUser, nil
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, errs.Error) {
	wErr := errs.New()
	workspaceUsers, uErr := us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToListWorkspaceUsers)
	}

	return workspaceUsers, nil
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, errs.Error) {
	wErr := errs.New()

	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return nil, wErr.Propagate(svc.ErrFailedToListUserWorkspaces)
	}

	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToListUserWorkspaces)
	}

	workspaceIDs, uErr := us.userRepository.ListUserWorkspaces(ctx, user.ID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr.Propagate(svc.ErrFailedToListUserWorkspaces)
	}

	return workspaceIDs, nil
}

func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, errs.Error) {
	wErr := errs.New()

	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserRole(ctx, workspaceID, userIDs, role)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr.Propagate(svc.ErrFailedToUpdateWorkspaceUserRole)
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr.Propagate(svc.ErrFailedToUpdateWorkspaceUserRole)
	}

	return workspaceUser, nil
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) errs.Error {
	wErr := errs.New()
	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserStatus(ctx, workspaceID, userIDs, status)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToUpdateWorkspaceUserStatus)
	}

	return nil
}

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) errs.Error {
	wErr := errs.New()
	uErr := us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToRemoveWorkspaceUsers)
	}

	return nil
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, errs.Error) {
	wErr := errs.New()
	admin, member, uErr := us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return 0, 0, wErr.Propagate(svc.ErrFailedToGetWorkspaceUserCountByRole)
	}

	return admin, member, nil
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) errs.Error {
	wErr := errs.New()
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return wErr.Propagate(svc.ErrFailedToSyncUser)
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToSyncUser)
	}

	if uErr := us.userRepository.SyncUser(ctx, user.ID, clerk); uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToSyncUser)
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) errs.Error {
	wErr := errs.New()

	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return wErr.Propagate(svc.ErrFailedToDeleteUser)
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToDeleteUser)
	}

	if uErr := us.userRepository.DeleteUser(ctx, user.ID); uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr.Propagate(svc.ErrFailedToDeleteUser)
	}

	return nil
}
