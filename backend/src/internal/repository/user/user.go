package user

import (
	"context"

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

func (r *userRepo) GetOrCreateClerkUser(ctx context.Context, clerkID, normalizedClerkEmail, userName string) (*models.User, error) {
	return nil, nil
}

func (r *userRepo) GetOrCreatePartialUsers(ctx context.Context, normalizedUserEmail string) (*models.User, error) {
	return nil, nil
}

func (r *userRepo) BatchGetOrCreatePartialUsers(ctx context.Context, normalizedUserEmail []string) ([]models.User, error) {
	return nil, nil
}

func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return nil, nil
}

func (r *userRepo) BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	return nil, nil
}

func (r *userRepo) GetUserByClerkCredentials(ctx context.Context, clerkID, normalizedClerkEmail string) (*models.User, error) {
	return nil, nil
}

func (r *userRepo) SyncUser(ctx context.Context, clerkID, normalizedUserEmail string) error {
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (r *userRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	return false, nil
}

func (r *userRepo) ClerkUserExists(ctx context.Context, clerkID, normalizedClerkEmail string) (bool, error) {
	return false, nil
}
