package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// UserRepository interacts with the user table and the workspace_user table
// This is used to interact with the user repository

type UserRepository interface {
	// ----- CRUD operations for user table -----

	// ---- Create operations for user table ----

	// CreateUser creates a new user.
	// This is triggered when a user signs up.
	CreateUser(ctx context.Context, clerkID, normalizedEmail string) (uuid.UUID, error)

	// CreatePartialUser creates a new user with partial details.
	// This is triggered when a user is invited.
	CreatePartialUser(ctx context.Context, normalizedUserEmail string) (uuid.UUID, error)

	// BatchCreatePartialUsers creates a new user with partial details.
	// This is triggered when a users are invited.
	BatchCreatePartialUsers(ctx context.Context, normalizedUserEmail []string) ([]uuid.UUID, error)

	// ---- Read operations for user table ----

	// GetUserByClerkID gets a user by ClerkID.
	// This is used to get users own details by ClerkID or ClerkEmail.
	// Internally, it attempts to get the user by ClerkID first.
	// If the user does not exist, it attempts to get the user by ClerkEmail.
	GetUserByClerkID(ctx context.Context, clerkID, normalizedClerkEmail string) (models.User, error)

	// GetUserByUserID gets a user by UserID.
	// This is used to get the user details.
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, error)

	// BatchGetUsersByUserIDs gets users by userIDs.
	// This is used to get other users details.
	BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	// getUserByUserEmail gets a user by email.
	// This is used to get other users details.
	GetUserByUserEmail(ctx context.Context, normalizedUserEmail string) (models.User, error)

	// BatchGetUsersByUserEmails gets users by emails.
	// This is used to get other users details.
	BatchGetUsersByUserEmails(ctx context.Context, normalizedUserEmails []string) ([]models.User, error)

	// ---- Get Or Create operations for user table ----

	// GetOrCreateUserByEmail checks if a user exists by email.
	// If the user exists, it returns the user.
	// If the user does not exist, it creates a new user and returns the user.
	GetOrCreateUserByEmail(ctx context.Context, normalizedUserEmail string) (models.User, error)

	// BatchGetOrCreateUsersByEmails checks if users exist by emails.
	// If the user exists, it returns the user.
	// If the user does not exist, it creates a new user and returns the user.
	BatchGetOrCreateUsersByEmails(ctx context.Context, normalizedUserEmails []string) ([]models.User, error)

	// ---- Sync operations for user table ----

	// SyncUser syncs the user data with Clerk
	// This is triggered when a user updates their profile
	SyncUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) error

	// DeleteUser deletes a user from the user table
	// and removes the user from all workspaces
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// -----  Optimized Lookup Operations -----

	// UserExists checks if a user exists
	// Optimized for quick lookups over the user table
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)

	// ClerkUserExists checks if a user exists by clerkID or clerkEmail
	// Optimized for quick lookups over the user table
	ClerkUserExists(ctx context.Context, clerkID, normalizedClerkEmail string) (bool, error)
}
