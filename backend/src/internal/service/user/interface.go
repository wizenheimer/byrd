package user

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// UserService is a service that manages users
// It holds the business logic for user management
// It embeds UserRespository to interact with the database
type UserService interface {
	GetOrCreateUser(ctx context.Context, clerk *clerk.User) (*models.User, error)

	BatchGetOrCreateUsers(ctx context.Context, emails []string) ([]models.User, error)

	ListUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	GetUserByClerkCredentials(ctx context.Context, clerk *clerk.User) (*models.User, error)

	SyncUser(ctx context.Context, clerk *clerk.User) error

	DeleteUser(ctx context.Context, clerk *clerk.User) error

	UserExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)

	ClerkUserExists(ctx context.Context, clerk *clerk.User) (bool, error)
}
