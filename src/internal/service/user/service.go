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
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

func NewUserService(userRepository repo.UserRepository, logger *logger.Logger) svc.UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (us *userService) CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, err.Error) {
	wErr := err.New()
	email, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToCreateWorkspaceOwner, map[string]any{"error": err.Error()})
		return nil, wErr
	}

	name := utils.GetClerkUserFullName(clerk)

	userID, err := uuid.NewUUID()
	if err != nil {
		wErr.Add(svc.ErrFailedToCreateWorkspaceOwner, map[string]any{"error": err.Error()})
		return nil, wErr
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
		return nil, wErr
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
		return nil, wErr
	}

	return &user, nil
}

func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) ([]api.CreateWorkspaceUserResponse, err.Error) {
	wErr := err.New()
	if len(invitedUsers) == 0 {
		wErr.Add(svc.ErrFailedToAddUserToWorkspace, map[string]any{"error": errors.New("no users to add")})
		return nil, wErr
	}

	var responses []api.CreateWorkspaceUserResponse
	workspaceUsers, uErr := us.userRepository.AddUsersToWorkspace(ctx, invitedUsers, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr
	}

	responses = append(responses, workspaceUsers...)

	return responses, wErr
}

func (us *userService) GetWorkspaceUser(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (models.WorkspaceUser, err.Error) {
	wErr := err.New()
	userEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return models.WorkspaceUser{}, wErr
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, userEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr
	}

	return workspaceUser, nil
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, err.Error) {
	wErr := err.New()
	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr
	}

	return workspaceUser, nil
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, err.Error) {
	wErr := err.New()
	workspaceUsers, uErr := us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr
	}

	return workspaceUsers, nil
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, err.Error) {
	wErr := err.New()

	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return nil, wErr
	}

	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr
	}

	workspaceIDs, uErr := us.userRepository.ListUserWorkspaces(ctx, user.ID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return nil, wErr
	}

	return workspaceIDs, nil
}

func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, err.Error) {
	wErr := err.New()

	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserRole(ctx, workspaceID, userIDs, role)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr
	}

	workspaceUser, uErr := us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return models.WorkspaceUser{}, wErr
	}

	return workspaceUser, nil
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) err.Error {
	wErr := err.New()
	userIDs := []uuid.UUID{userID}
	_, uErr := us.userRepository.UpdateWorkspaceUserStatus(ctx, workspaceID, userIDs, status)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	return nil
}

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) err.Error {
	wErr := err.New()
	uErr := us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	return nil
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, err.Error) {
	wErr := err.New()
	admin, member, uErr := us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return 0, 0, wErr
	}

	return admin, member, nil
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) err.Error {
	wErr := err.New()
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return wErr
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	if uErr := us.userRepository.SyncUser(ctx, user.ID, clerk); uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) err.Error {
	wErr := err.New()

	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		wErr.Add(svc.ErrFailedToGetUserEmail, map[string]any{"error": err.Error()})
		return wErr
	}
	user, uErr := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	if uErr := us.userRepository.DeleteUser(ctx, user.ID); uErr != nil && uErr.HasErrors() {
		wErr.Merge(uErr)
		return wErr
	}

	return nil
}
