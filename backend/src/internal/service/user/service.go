package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/user"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

var _ UserService = (*userService)(nil)

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

// GetOrCreateWorkspaceOwner gets or creates a single user.
func (us *userService) GetOrCreateUser(ctx context.Context, clerk *clerk.User) (*models.User, error) {
	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, err
	}
	clerkEmail = utils.NormalizeEmail(clerkEmail)

	clerkUserName := utils.GetClerkUserFullName(clerk)

	user, err := us.userRepository.GetOrCreateClerkUser(ctx, clerk.ID, clerkEmail, clerkUserName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// BatchGetOrCreateUsers creates a batch of users if they do not exist.
// It returns the created or found users.
// It returns an error if the users could not be created.
func (us *userService) BatchGetOrCreateUsers(ctx context.Context, emails []string) ([]models.User, error) {
	for i, email := range emails {
		emails[i] = utils.NormalizeEmail(email)
	}

	if len(emails) == 1 {
		user, err := us.userRepository.GetOrCreatePartialUsers(ctx, emails[0])
		if err != nil {
			return nil, err
		}
		return []models.User{*user}, err
	}

	users, err := us.userRepository.BatchGetOrCreatePartialUsers(ctx, emails)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// ListUsersByUserIDs lists users by userIDs.
// This is used to get the user details.
func (us *userService) ListUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	if len(userIDs) == 1 {
		user, err := us.userRepository.GetUserByUserID(ctx, userIDs[0])
		if err != nil {
			return nil, err
		}

		return []models.User{*user}, nil
	}

	users, err := us.userRepository.BatchGetUsersByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByClerk gets a clerk user by clerk credentials.
// This is used to get the clerk user details.
func (us *userService) GetUserByClerkCredentials(ctx context.Context, clerk *clerk.User) (*models.User, error) {
	userEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, err
	}
	userEmail = utils.NormalizeEmail(userEmail)

	user, err := us.userRepository.GetUserByClerkCredentials(ctx, clerk.ID, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SyncUser syncs a user with Clerk
// It returns an error if the user could not be synced
// When sync is triggered, it marks the account status as active
// And updates the user's email and name if they have changed
func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}
	clerkEmail = utils.NormalizeEmail(clerkEmail)

	err = us.userRepository.SyncUser(ctx, clerk.ID, clerkEmail)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from Clerk
// It returns an error if the user could not be deleted
// It returns nil if the user was deleted successfully
// This is the only user-facing and handler-owned method
func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}
	clerkEmail = utils.NormalizeEmail(clerkEmail)

	user, err := us.userRepository.GetUserByClerkCredentials(ctx, clerk.ID, clerkEmail)
	if err != nil {
		return err
	}

	return us.userRepository.DeleteUser(ctx, user.ID)
}

// UserExistsByUserID checks if a user exists by UserID.
// It returns true if the user exists, otherwise it returns false.
func (us *userService) UserExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	return us.userRepository.UserExists(ctx, userID)
}

// UserExistsByClerkID checks if a user exists by ClerkID.
// It returns true if the user exists, otherwise it returns false.
func (us *userService) ClerkUserExists(ctx context.Context, clerk *clerk.User) (bool, error) {
	email, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return false, err
	}

	email = utils.NormalizeEmail(email)

	return us.userRepository.ClerkUserExists(ctx, clerk.ID, email)
}
