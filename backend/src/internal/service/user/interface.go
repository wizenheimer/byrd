// ./src/internal/service/user/interface.go
package user

import (
	"context"

	// "github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

var (
	maxUserBatchSize int = 25
)

// UserService is a service that manages users
// It holds the business logic for user management
// It embeds UserRespository to interact with the database
type UserService interface {
	GetOrCreateUser(ctx context.Context, userEmail string) (*models.User, error)

	BatchGetOrCreateUsers(ctx context.Context, userEmails []string) ([]models.User, error)

	ListUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)

	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)

	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error)

	ActivateUser(ctx context.Context, userEmail string) (*models.User, error)

	DeleteUserByEmail(ctx context.Context, userEmail string) error

	DeleteUserByID(ctx context.Context, userID uuid.UUID) error

	UserExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)

	UserExistsByUserEmail(ctx context.Context, userEmail string) (bool, error)
}
