// ./src/internal/repository/user/user.go
package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserInactive = errors.New("user is inactive")
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

func (r *userRepo) GetOrCreateUser(ctx context.Context, userEmail string) (*models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		WITH existing_user AS (
			SELECT id, email, status, created_at, updated_at
			FROM users
			WHERE email = $1 AND status IN ($2, $3)
		), new_user AS (
			INSERT INTO users (id, email, status, created_at, updated_at)
			SELECT $4, $1, $2, $5, $5
			WHERE NOT EXISTS (SELECT 1 FROM existing_user)
			RETURNING id, email, status, created_at, updated_at
		)
		SELECT id, email, status, created_at, updated_at
		FROM existing_user
		UNION ALL
		SELECT id, email, status, created_at, updated_at
		FROM new_user`

	user := &models.User{}
	err := q.QueryRow(
		ctx,
		query,
		userEmail,
		models.AccountStatusPending,
		models.AccountStatusActive,
		uuid.New(),
		time.Now(),
	).Scan(
		&user.ID,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) BatchGetOrCreateUsers(ctx context.Context, userEmails []string) ([]models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		WITH existing_users AS (
			SELECT id, email, status, created_at, updated_at
			FROM users
			WHERE email = ANY($1) AND status IN ($2, $3)
		), new_users AS (
			INSERT INTO users (id, email, status, created_at, updated_at)
			SELECT
				gen_random_uuid(),
				e.email,
				$2,
				$4,
				$4
			FROM unnest($1::text[]) AS e(email)
			WHERE NOT EXISTS (
				SELECT 1 FROM existing_users WHERE existing_users.email = e.email
			)
			RETURNING id, email, status, created_at, updated_at
		)
		SELECT id, email, status, created_at, updated_at
		FROM existing_users
		UNION ALL
		SELECT id, email, status, created_at, updated_at
		FROM new_users`

	rows, err := q.Query(
		ctx,
		query,
		userEmails,
		models.AccountStatusPending,
		models.AccountStatusActive,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		SELECT id, email, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status != $2`

	user := &models.User{}
	err := q.QueryRow(ctx, query, userID, models.AccountStatusInactive).Scan(
		&user.ID,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) BatchGetUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		SELECT id, email, status, created_at, updated_at
		FROM users
		WHERE id = ANY($1) AND status != $2`

	rows, err := q.Query(ctx, query, userIDs, models.AccountStatusInactive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		SELECT id, email, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status != $2`

	user := &models.User{}
	err := q.QueryRow(ctx, query, userEmail, models.AccountStatusInactive).Scan(
		&user.ID,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) ActivateUser(ctx context.Context, userEmail string) (*models.User, error) {
	q := r.getQuerier(ctx)

	const query = `
		UPDATE users
		SET status = $1, updated_at = $3
		WHERE email = $2 AND status != $4
		RETURNING id, email, status, created_at, updated_at`

	user := &models.User{}
	err := q.QueryRow(
		ctx,
		query,
		models.AccountStatusActive,
		userEmail,
		time.Now(),
		models.AccountStatusInactive,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	q := r.getQuerier(ctx)

	const query = `
		UPDATE users
		SET status = $1, updated_at = $3
		WHERE id = $2 AND status != $1`

	result, err := q.Exec(
		ctx,
		query,
		models.AccountStatusInactive,
		userID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepo) DeleteUserByEmail(ctx context.Context, userEmail string) error {
	q := r.getQuerier(ctx)

	const query = `
		UPDATE users
		SET status = $1, updated_at = $3
		WHERE email = $2 AND status != $1`

	result, err := q.Exec(
		ctx,
		query,
		models.AccountStatusInactive,
		userEmail,
		time.Now(),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UserIDExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	q := r.getQuerier(ctx)

	const query = `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE id = $1 AND status != $2
		)`

	var exists bool
	err := q.QueryRow(ctx, query, userID, models.AccountStatusInactive).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepo) UserEmailExists(ctx context.Context, userEmail string) (bool, error) {
	q := r.getQuerier(ctx)

	const query = `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE email = $1 AND status != $2
		)`

	var exists bool
	err := q.QueryRow(ctx, query, userEmail, models.AccountStatusInactive).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
