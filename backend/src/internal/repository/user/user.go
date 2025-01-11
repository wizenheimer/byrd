package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type userRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewUserRepository(tm *transaction.TxManager, logger *logger.Logger) UserRepository {
	return &userRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "user_repository"}),
	}
}

// --------------------------------------------------
// --- Create & Update operations for user table ---
// --------------------------------------------------

// GetUserByUserID gets a user by UserID
// Respects the status of the user and excludes soft-deleted users

// CreateUser creates a new user.
// This is triggered when a user signs up.
func (r *userRepo) CreateUser(ctx context.Context, clerkID, normalizedEmail string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

// CreatePartialUser creates a new user with partial details.
// This is triggered when a user is invited.
func (r *userRepo) CreatePartialUser(ctx context.Context, normalizedUserEmail string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

// BatchCreatePartialUsers creates a new user with partial details.
// This is triggered when a users are invited.
func (r *userRepo) BatchCreatePartialUsers(ctx context.Context, normalizedUserEmail []string) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

// ---- Read operations for user table ----

// GetUserByClerkID gets a user by ClerkID.
// This is used to get users own details by ClerkID or ClerkEmail.
// Internally, it attempts to get the user by ClerkID first.
// If the user does not exist, it attempts to get the user by ClerkEmail.
func (r *userRepo) GetUserByClerkID(ctx context.Context, clerkID, normalizedClerkEmail string) (models.User, error) {
	return models.User{}, nil
}

// GetUserByUserID gets a user by UserID.
// This is used to get the user details.
func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	return models.User{}, nil
}

// BatchGetUsersByUserIDs gets users by userIDs.
// This is used to get other users details.
func (r *userRepo) BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	return []models.User{}, nil
}

// getUserByUserEmail gets a user by email.
// This is used to get other users details.
func (r *userRepo) GetUserByUserEmail(ctx context.Context, normalizedUserEmail string) (models.User, error) {
	return models.User{}, nil
}

// BatchGetUsersByUserEmails gets users by emails.
// This is used to get other users details.
func (r *userRepo) BatchGetUsersByUserEmails(ctx context.Context, normalizedUserEmails []string) ([]models.User, error) {
	return []models.User{}, nil
}

// ---- Get Or Create operations for user table ----

// GetOrCreateUserByEmail checks if a user exists by email.
// If the user exists, it returns the user.
// If the user does not exist, it creates a new user and returns the user.
func (r *userRepo) GetOrCreateUserByEmail(ctx context.Context, normalizedUserEmail string) (models.User, error) {
	return models.User{}, nil
}

// BatchGetOrCreateUsersByEmails checks if users exist by emails.
// If the user exists, it returns the user.
// If the user does not exist, it creates a new user and returns the user.
func (r *userRepo) BatchGetOrCreateUsersByEmails(ctx context.Context, normalizedUserEmails []string) ([]models.User, error) {
	return []models.User{}, nil
}

// ---- Sync operations for user table ----

// SyncUser syncs the user data with Clerk
// This is triggered when a user updates their profile
func (r *userRepo) SyncUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) error {
	return nil
}

// DeleteUser deletes a user from the user table
// and removes the user from all workspaces
func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}

// -----  Optimized Lookup Operations -----

// UserExists checks if a user exists
// Optimized for quick lookups over the user table
func (r *userRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	return false, nil
}

// ClerkUserExists checks if a user exists by clerkID or clerkEmail
// Optimized for quick lookups over the user table
func (r *userRepo) ClerkUserExists(ctx context.Context, clerkID, normalizedClerkEmail string) (bool, error) {
	return false, nil
}
