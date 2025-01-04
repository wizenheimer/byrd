package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
	"go.uber.org/zap"
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
		return nil, err
	}

	name := utils.GetClerkUserFullName(clerk)

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, errs[0]
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
		for _, err := range errs {
			responses = append(responses, api.CreateWorkspaceUserResponse{
				Error: err,
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
		return models.WorkspaceUser{}, err
	}

	return us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, userEmail)
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	return us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	return us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error) {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, err
	}

	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return nil, err
	}

	return us.userRepository.ListUserWorkspaces(ctx, user.ID)
}

func (us *userService) UpdateWorkspaceUserRole(ctx context.Context, userID, workspaceID uuid.UUID, role models.UserWorkspaceRole) (models.WorkspaceUser, error) {
	userIDs := []uuid.UUID{userID}
	_, errs := us.userRepository.UpdateWorkspaceUserRole(ctx, workspaceID, userIDs, role)
	if len(errs) > 0 {
		return models.WorkspaceUser{}, errs[0]
	}

	return us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
}

func (us *userService) UpdateWorkspaceUserStatus(ctx context.Context, userID, workspaceID uuid.UUID, status models.UserWorkspaceStatus) error {
	userIDs := []uuid.UUID{userID}
	_, errs := us.userRepository.UpdateWorkspaceUserStatus(ctx, workspaceID, userIDs, status)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {

	us.logger.Debug("Removing users from workspace", zap.Any("userIDs", userIDs), zap.Any("workspaceID", workspaceID))

	return us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error) {
	return us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return err
	}

	if err := us.userRepository.SyncUser(ctx, user.ID, clerk); err != nil {
		return err
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	primaryEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, primaryEmail)
	if err != nil {
		return err
	}

	return us.userRepository.DeleteUser(ctx, user.ID)
}
