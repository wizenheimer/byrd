// ./src/internal/repository/user/interface.go
package user

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// UserRepository interacts with the user table and the workspace_user table
// This is used to interact with the user repository

type UserRepository interface {
	GetOrCreateUser(ctx context.Context, userEmail string) (*models.User, error)

	BatchGetOrCreateUsers(ctx context.Context, userEmail []string) ([]models.User, error)

	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error)

	BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)

	ActivateUser(ctx context.Context, userEmail string) (*models.User, error)

	DeleteUserByID(ctx context.Context, userID uuid.UUID) error

	DeleteUserByEmail(ctx context.Context, userEmail string) error

	UserIDExists(ctx context.Context, userID uuid.UUID) (bool, error)

	UserEmailExists(ctx context.Context, userEmail string) (bool, error)
}
