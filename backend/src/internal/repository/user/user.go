package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *userRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func (r *userRepo) GetOrCreateClerkUser(ctx context.Context, clerkID, normalizedClerkEmail, userName string) (*models.User, error) {
	user := &models.User{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		INSERT INTO users (clerk_id, email, name, status)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE
		SET clerk_id = EXCLUDED.clerk_id,
			name = EXCLUDED.name,
			status = CASE
				WHEN users.status = $5 THEN $4
				ELSE users.status
			END
		RETURNING id, clerk_id, email, name, status, created_at, updated_at`,
		clerkID, normalizedClerkEmail, userName, models.AccountStatusActive, models.AccountStatusInactive,
	).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get or create clerk user: %w", err)
	}

	return user, nil
}

func (r *userRepo) GetOrCreatePartialUsers(ctx context.Context, normalizedUserEmail string) (*models.User, error) {
	user := &models.User{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		INSERT INTO users (email, status)
		VALUES ($1, $2)
		ON CONFLICT (email) DO UPDATE
		SET status = CASE
			WHEN users.status = $3 THEN $2
			ELSE users.status
		END
		RETURNING id, clerk_id, email, name, status, created_at, updated_at`,
		normalizedUserEmail, models.AccountStatusPending, models.AccountStatusInactive,
	).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get or create partial user: %w", err)
	}

	return user, nil
}

func (r *userRepo) BatchGetOrCreatePartialUsers(ctx context.Context, normalizedUserEmails []string) ([]models.User, error) {
	if len(normalizedUserEmails) == 0 {
		return []models.User{}, nil
	}

	// Create a temporary table for bulk insertion
	_, err := r.getQuerier(ctx).Exec(ctx, `
		CREATE TEMPORARY TABLE temp_users (
			email VARCHAR(255) PRIMARY KEY
		) ON COMMIT DROP`)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary table: %w", err)
	}

	// Insert emails into temporary table
	for _, email := range normalizedUserEmails {
		_, err = r.getQuerier(ctx).Exec(ctx, `
			INSERT INTO temp_users (email) VALUES ($1)
			ON CONFLICT DO NOTHING`,
			email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert into temporary table: %w", err)
		}
	}

	// Perform bulk insert/update and return results
	rows, err := r.getQuerier(ctx).Query(ctx, `
		WITH inserted AS (
			INSERT INTO users (email, status)
			SELECT t.email, $1
			FROM temp_users t
			ON CONFLICT (email) DO UPDATE
			SET status = CASE
				WHEN users.status = $2 THEN $1
				ELSE users.status
			END
			RETURNING id, clerk_id, email, name, status, created_at, updated_at
		)
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM inserted`,
		models.AccountStatusPending, models.AccountStatusInactive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get or create users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.ClerkID,
			&user.Email,
			&user.Name,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status != $2`,
		userID, models.AccountStatusInactive,
	).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepo) BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	if len(userIDs) == 0 {
		return []models.User{}, nil
	}

	rows, err := r.getQuerier(ctx).Query(ctx, `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE id = ANY($1) AND status != $2`,
		userIDs, models.AccountStatusInactive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.ClerkID,
			&user.Email,
			&user.Name,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *userRepo) GetUserByClerkCredentials(ctx context.Context, clerkID, normalizedClerkEmail string) (*models.User, error) {
	user := &models.User{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE (clerk_id = $1 OR email = $2) AND status != $3`,
		clerkID, normalizedClerkEmail, models.AccountStatusInactive,
	).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by clerk credentials: %w", err)
	}

	return user, nil
}

func (r *userRepo) SyncUser(ctx context.Context, clerkID, normalizedUserEmail string) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE users
		SET email = $1,
			status = $2
		WHERE email = $3 OR clerk_id != $4`,
		normalizedUserEmail, models.AccountStatusActive, normalizedUserEmail, clerkID,
	)

	if err != nil {
		return fmt.Errorf("failed to sync user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE users
		SET status = $1
		WHERE id = $2 AND status != $1`,
		models.AccountStatusInactive, userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE id = $1 AND status != $2
		)`,
		userID, models.AccountStatusInactive,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *userRepo) ClerkUserExists(ctx context.Context, clerkID, normalizedClerkEmail string) (bool, error) {
	var exists bool
	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE (clerk_id = $1 OR email = $2) AND status != $3
		)`,
		clerkID, normalizedClerkEmail, models.AccountStatusInactive,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check clerk user existence: %w", err)
	}

	return exists, nil
}
