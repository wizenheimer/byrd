// ./src/internal/service/user/service.go
package user

import (
	"context"
	"errors"

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
func (us *userService) GetOrCreateUser(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := us.userRepository.GetOrCreateUser(ctx, userEmail)
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
		return nil, errors.New("user batch size exceeds the maximum limit")
	}

	emails = utils.CleanEmailList(emails, nil)

	var users []models.User
	var err error
	if len(emails) == 1 {
		// If there is only one email, get or create a single user
		user, err := us.userRepository.GetOrCreateUser(ctx, emails[0])
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		users = append(users, *user)
	} else {
		// If there are multiple emails, get or create a batch of users
		users, err = us.userRepository.BatchGetOrCreateUsers(ctx, emails)
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
		return nil, errors.New("no userIDs provided")
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
func (us *userService) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := us.userRepository.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ActivateUser activates a user in Clerk.
// It returns an error if the user could not be activated.
func (us *userService) ActivateUser(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := us.userRepository.ActivateUser(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user from Clerk.
// It returns an error if the user could not be deleted.
// It returns nil if the user was deleted successfully.
// This is the only user-facing and handler-owned method.
func (us *userService) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	return us.userRepository.DeleteUserByID(ctx, userID)
}

func (us *userService) DeleteUserByEmail(ctx context.Context, userEmail string) error {
	return us.userRepository.DeleteUserByEmail(ctx, userEmail)
}

// UserExistsByUserID checks if a user exists by UserID.
// It returns true if the user exists, otherwise it returns false.
func (us *userService) UserExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	return us.userRepository.UserIDExists(ctx, userID)
}

func (us *userService) UserExistsByUserEmail(ctx context.Context, userEmail string) (bool, error) {
	return us.userRepository.UserEmailExists(ctx, userEmail)
}
