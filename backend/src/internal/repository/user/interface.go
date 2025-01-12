package user

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// UserRepository interacts with the user table and the workspace_user table
// This is used to interact with the user repository

type UserRepository interface {
	GetOrCreateClerkUser(ctx context.Context, clerkID, normalizedClerkEmail, userName string) (*models.User, error)

	GetOrCreatePartialUsers(ctx context.Context, normalizedUserEmail string) (*models.User, error)

	BatchGetOrCreatePartialUsers(ctx context.Context, normalizedUserEmail []string) ([]models.User, error)

	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error)

	BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	GetUserByClerkCredentials(ctx context.Context, clerkID, normalizedClerkEmail string) (*models.User, error)

	SyncUser(ctx context.Context, clerkID, normalizedUserEmail string) error

	DeleteUser(ctx context.Context, userID uuid.UUID) error

	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)

	ClerkUserExists(ctx context.Context, clerkID, normalizedClerkEmail string) (bool, error)
}
