// ./src/internal/service/user/service.go
package user

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/repository/user"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type userService struct {
	userRepository  user.UserRepository
	templateLibrary template.TemplateLibrary
	logger          *logger.Logger
	errorRecord     *recorder.ErrorRecorder
}

func NewUserService(
	userRepository user.UserRepository,
	templateLibrary template.TemplateLibrary,
	logger *logger.Logger,
	errorRecord *recorder.ErrorRecorder,
) (UserService, error) {

	us := userService{
		userRepository:  userRepository,
		templateLibrary: templateLibrary,
		errorRecord:     errorRecord,
		logger:          logger.WithFields(map[string]interface{}{"module": "user_service"}),
	}

	return &us, nil
}

// GetOrCreateWorkspaceOwner gets or creates a single user.
func (us *userService) GetOrCreateUser(ctx context.Context, clerk *clerk.User) (*models.User, error) {
	if clerk == nil {
		return nil, errors.New("non-fatal: clerk user is nil")
	}

	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, err
	}

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
	if len(emails) > maxUserBatchSize {
		return nil, errors.New("non-fatal: user batch size exceeds the maximum limit")
	}

	emails = utils.CleanEmailList(emails, nil)

	var users []models.User
	var err error
	if len(emails) == 1 {
		// If there is only one email, get or create a single user
		user, err := us.userRepository.GetOrCreatePartialUsers(ctx, emails[0])
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		users = append(users, *user)
	} else {
		// If there are multiple emails, get or create a batch of users
		users, err = us.userRepository.BatchGetOrCreatePartialUsers(ctx, emails)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// ListUsersByUserIDs lists users by userIDs.
// This is used to get the user details.
func (us *userService) ListUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("non-fatal: no userIDs provided")
	}

	var users []models.User
	var err error
	if len(userIDs) == 1 {
		// If there is only one user, get the user
		user, err := us.userRepository.GetUserByUserID(ctx, userIDs[0])
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		users = append(users, *user)
	} else {
		// If there are multiple users, get a batch of users
		users, err = us.userRepository.BatchGetUsersByUserIDs(ctx, userIDs)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GetUserByClerk gets a clerk user by clerk credentials.
// This is used to get the clerk user details.
func (us *userService) GetUserByClerkCredentials(ctx context.Context, clerk *clerk.User) (*models.User, error) {
	if clerk == nil {
		return nil, errors.New("non-fatal: clerk user is nil")
	}

	userEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return nil, err
	}

	user, err := us.userRepository.GetUserByClerkCredentials(ctx, clerk.ID, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SyncUser syncs a user with Clerk.
// It returns an error if the user could not be synced.
// When sync is triggered, it marks the account status as active.
// And updates the user's email and name if they have changed.
func (us *userService) SyncUser(ctx context.Context, clerk *clerk.User) error {
	if clerk == nil {
		return errors.New("non-fatal: clerk user is nil")
	}

	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}

	err = us.userRepository.SyncUser(ctx, clerk.ID, clerkEmail)
	if err != nil {
		return err
	}

	return nil
}

// ActivateUser activates a user in Clerk.
// It returns an error if the user could not be activated.
func (us *userService) ActivateUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) error {
	if clerkUser == nil {
		return errors.New("non-fatal: clerk user is nil")
	}

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return err
	}

	clerkID := clerkUser.ID
	if err := us.userRepository.ActivateUser(ctx, userID, clerkID, userEmail); err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from Clerk.
// It returns an error if the user could not be deleted.
// It returns nil if the user was deleted successfully.
// This is the only user-facing and handler-owned method.
func (us *userService) DeleteUser(ctx context.Context, clerk *clerk.User) error {
	if clerk == nil {
		return errors.New("non-fatal: clerk user is nil")
	}

	clerkEmail, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return err
	}

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
	if clerk == nil {
		return false, errors.New("non-fatal: clerk user is nil")
	}

	email, err := utils.GetClerkUserEmail(clerk)
	if err != nil {
		return false, err
	}

	return us.userRepository.ClerkUserExists(ctx, clerk.ID, email)
}
