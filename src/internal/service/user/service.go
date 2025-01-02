package user

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

func (us *userService) CreateWorkspaceOwner(ctx context.Context, clerk *clerk.User, workspaceID uuid.UUID) (*models.User, error) {
	if clerk.PrimaryEmailAddressID == nil {
		return nil, errors.New("PrimaryEmailAddressID is required")
	}

	email := *clerk.PrimaryEmailAddressID

	var name string
	if clerk.FirstName != nil {
		name = *clerk.FirstName
	}
	if clerk.LastName != nil {
		name += " " + *clerk.LastName
	}

	if name == "" {
		name = utils.GenerateNameFromEmail(email)
	}

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

	return &user, nil
}

func (us *userService) AddUserToWorkspace(ctx context.Context, workspaceID uuid.UUID, invitedUsers []api.InviteUserToWorkspaceRequest) []api.CreateWorkspaceUserResponse {
	var responses []api.CreateWorkspaceUserResponse
	var emails []string
	for _, user := range invitedUsers {
		emails = append(emails, user.Email)
	}

	users, errs := us.userRepository.GetOrCreateUserByEmail(ctx, emails)
	if len(errs) > 0 {
		for _, err := range errs {
			responses = append(responses, api.CreateWorkspaceUserResponse{
				Error: err,
			})
		}
		return responses
	}

	userIDs := make([]uuid.UUID, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	workspaceUsers, errs := us.userRepository.AddUsersToWorkspace(ctx, userIDs, workspaceID)
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
	if clerk.PrimaryEmailAddressID == nil {
		return models.WorkspaceUser{}, errors.New("PrimaryEmailAddressID is required")
	}

	return us.userRepository.GetWorkspaceClerkUser(ctx, workspaceID, clerk.ID, *clerk.PrimaryEmailAddressID)
}

func (us *userService) GetWorkspaceUserByID(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	return us.userRepository.GetWorkspaceUser(ctx, workspaceID, userID)
}

func (us *userService) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	return us.userRepository.ListWorkspaceUsers(ctx, workspaceID)
}

func (us *userService) ListUserWorkspaces(ctx context.Context, clerk *clerk.User) ([]uuid.UUID, error) {
	if clerk.PrimaryEmailAddressID == nil {
		return nil, errors.New("PrimaryEmailAddressID is required")
	}

	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, *clerk.PrimaryEmailAddressID)
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

func (us *userService) RemoveWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	return us.userRepository.RemoveUsersFromWorkspace(ctx, userIDs, workspaceID)
}

func (us *userService) AddWorkspaceUsers(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	_, errs := us.userRepository.AddUsersToWorkspace(ctx, userIDs, workspaceID)
	return errs
}

func (us *userService) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error) {
	return us.userRepository.GetWorkspaceUserCountByRole(ctx, workspaceID)
}

func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	if clerk.PrimaryEmailAddressID == nil {
		return errors.New("PrimaryEmailAddress can't be nil")
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, *clerk.PrimaryEmailAddressID)
	if err != nil {
		return err
	}

	return us.userRepository.SyncUser(ctx, user.ID, clerk)
}

func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	if clerk.PrimaryEmailAddressID == nil {
		return errors.New("PrimaryEmailAddress can't be nil")
	}
	user, err := us.userRepository.GetClerkUser(ctx, clerk.ID, *clerk.PrimaryEmailAddressID)
	if err != nil {
		return err
	}

	return us.userRepository.DeleteUser(ctx, user.ID)
}
