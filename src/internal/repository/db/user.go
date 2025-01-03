package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrWorkspaceUserNotFound  = errors.New("workspace user not found")
	ErrDuplicateUser          = errors.New("user already exists")
	ErrDuplicateWorkspaceUser = errors.New("user already exists in workspace")
	ErrInvalidUserData        = errors.New("invalid user data")
	ErrNoUsersFound           = errors.New("no users found")
	ErrInvalidRole            = errors.New("invalid role")
	ErrInvalidStatus          = errors.New("invalid status")
)

type userRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewUserRepository(tm *transaction.TxManager, logger *logger.Logger) repo.UserRepository {
	return &userRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "user_repository"}),
	}
}

// --------------------------------------------------
// --- Create & Update operations for user table ---
// --------------------------------------------------

// GetUserByUserID gets a user by UserID
func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := runner.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := runner.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetClerkUser gets a user by clerkID or clerkEmail
func (r *userRepo) GetClerkUser(ctx context.Context, clerkID string, clerkEmail string) (models.User, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE clerk_id = $1 OR email = $2
	`

	var user models.User
	err := runner.QueryRowContext(ctx, query, clerkID, clerkEmail).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get clerk user: %w", err)
	}

	return user, nil
}

// GetOrCreateUser creates a user if it does not exist
func (r *userRepo) GetOrCreateUser(ctx context.Context, partialUser models.User) (models.User, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	// Try to get existing user
	query := `
		SELECT id, clerk_id, email, name, status, created_at, updated_at
		FROM users
		WHERE clerk_id = $1 OR email = $2
	`

	var user models.User
	err := runner.QueryRowContext(ctx, query, partialUser.ClerkID, partialUser.Email).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new user
		insertQuery := `
			INSERT INTO users (clerk_id, email, name, status)
			VALUES ($1, $2, $3, $4)
			RETURNING id, clerk_id, email, name, status, created_at, updated_at
		`

		err = runner.QueryRowContext(ctx, insertQuery,
			partialUser.ClerkID,
			partialUser.Email,
			partialUser.Name,
			partialUser.Status,
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
			return models.User{}, fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return models.User{}, fmt.Errorf("failed to get existing user: %w", err)
	}

	return user, nil
}

// GetOrCreateUserByEmail creates users if they do not exist
func (r *userRepo) GetOrCreateUserByEmail(ctx context.Context, emails []string) ([]models.User, []error) {
	if len(emails) == 0 {
		return nil, []error{errors.New("no emails provided")}
	}

	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	users := make([]models.User, 0, len(emails))
	errs := make([]error, 0)

	for _, email := range emails {
		// Try to get existing user
		var user models.User
		err := runner.QueryRowContext(ctx, `
			SELECT id, clerk_id, email, name, status, created_at, updated_at
			FROM users
			WHERE email = $1
		`, email).Scan(
			&user.ID,
			&user.ClerkID,
			&user.Email,
			&user.Name,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			// Create new user
			err = runner.QueryRowContext(ctx, `
				INSERT INTO users (email, status)
				VALUES ($1, $2)
				RETURNING id, clerk_id, email, name, status, created_at, updated_at
			`, email, models.AccountStatusPending).Scan(
				&user.ID,
				&user.ClerkID,
				&user.Email,
				&user.Name,
				&user.Status,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
		}

		if err != nil {
			errs = append(errs, fmt.Errorf("failed to process email %s: %w", email, err))
			continue
		}

		users = append(users, user)
	}

	if len(errs) > 0 {
		return users, errs
	}

	return users, nil
}

// -----------------------------------------------------------
// --- Create & Update operations for workspace_user table ---
// -----------------------------------------------------------

// AddUsersToWorkspace adds users to a workspace
func (r *userRepo) AddUsersToWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) ([]models.WorkspaceUser, []error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	// Prepare batch insert
	valueStrings := make([]string, len(userIDs))
	valueArgs := make([]interface{}, 0, len(userIDs)*2)
	for i, userID := range userIDs {
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		valueArgs = append(valueArgs, workspaceID, userID)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO workspace_users (workspace_id, user_id)
		VALUES %s
		ON CONFLICT (workspace_id, user_id) DO NOTHING
		RETURNING workspace_id, user_id, role, status
	`, strings.Join(valueStrings, ","))

	rows, err := runner.QueryContext(ctx, insertQuery, valueArgs...)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to add users to workspace: %w", err)}
	}
	defer rows.Close()

	var workspaceUsers []models.WorkspaceUser
	for rows.Next() {
		var wu models.WorkspaceUser
		err := rows.Scan(
			&wu.WorkspaceID,
			&wu.ID, // user_id
			&wu.Role,
			&wu.Status,
		)
		if err != nil {
			return nil, []error{fmt.Errorf("failed to scan workspace user: %w", err)}
		}
		workspaceUsers = append(workspaceUsers, wu)
	}

	if err = rows.Err(); err != nil {
		return nil, []error{fmt.Errorf("error iterating over rows: %w", err)}
	}

	return workspaceUsers, nil
}

// RemoveUsersFromWorkspace removes users from a workspace
func (r *userRepo) RemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	// Create placeholders for the IN clause
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, 0, len(userIDs)+1)
	args = append(args, workspaceID)

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE workspace_users
		SET status = 'inactive', updated_at = CURRENT_TIMESTAMP
		WHERE workspace_id = $1 AND user_id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		return []error{fmt.Errorf("failed to remove users from workspace: %w", err)}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return []error{fmt.Errorf("failed to get affected rows: %w", err)}
	}

	if rowsAffected == 0 {
		return []error{ErrWorkspaceUserNotFound}
	}

	return nil
}

// GetWorkspaceUser gets a user from the workspace
func (r *userRepo) GetWorkspaceUser(ctx context.Context, workspaceID, userID uuid.UUID) (models.WorkspaceUser, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT u.id, u.clerk_id, u.email, u.name, u.status,
			   wu.role, wu.status as workspace_status,
			   u.created_at, u.updated_at
		FROM users u
		JOIN workspace_users wu ON u.id = wu.user_id
		WHERE wu.workspace_id = $1 AND wu.user_id = $2
	`

	var wu models.WorkspaceUser
	err := runner.QueryRowContext(ctx, query, workspaceID, userID).Scan(
		&wu.ID,
		&wu.ClerkID,
		&wu.Email,
		&wu.Name,
		&wu.Status,
		&wu.Role,
		&wu.WorkspaceUserStatus,
		&wu.CreatedAt,
		&wu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.WorkspaceUser{}, ErrWorkspaceUserNotFound
	}
	if err != nil {
		return models.WorkspaceUser{}, fmt.Errorf("failed to get workspace user: %w", err)
	}

	return wu, nil
}

// GetWorkspaceClerkUser gets a user from the workspace by clerk credentials
func (r *userRepo) GetWorkspaceClerkUser(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (models.WorkspaceUser, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT u.id, u.clerk_id, u.email, u.name, u.status,
			   wu.role, wu.status as workspace_status,
			   u.created_at, u.updated_at
		FROM users u
		JOIN workspace_users wu ON u.id = wu.user_id
		WHERE wu.workspace_id = $1 AND (u.clerk_id = $2 OR u.email = $3)
	`

	var wu models.WorkspaceUser
	err := runner.QueryRowContext(ctx, query, workspaceID, clerkID, clerkEmail).Scan(
		&wu.ID,
		&wu.ClerkID,
		&wu.Email,
		&wu.Name,
		&wu.Status,
		&wu.Role,
		&wu.WorkspaceUserStatus,
		&wu.CreatedAt,
		&wu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.WorkspaceUser{}, ErrWorkspaceUserNotFound
	}
	if err != nil {
		return models.WorkspaceUser{}, fmt.Errorf("failed to get workspace clerk user: %w", err)
	}

	return wu, nil
}

// ListWorkspaceUsers lists all users from the workspace
func (r *userRepo) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT u.id, u.clerk_id, u.email, u.name, u.status,
			   wu.role, wu.status as workspace_status,
			   u.created_at, u.updated_at
		FROM users u
		JOIN workspace_users wu ON u.id = wu.user_id
		WHERE wu.workspace_id = $1
		ORDER BY u.created_at DESC
	`

	rows, err := runner.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspace users: %w", err)
	}
	defer rows.Close()

	var users []models.WorkspaceUser
	for rows.Next() {
		var wu models.WorkspaceUser
		err := rows.Scan(
			&wu.ID,
			&wu.ClerkID,
			&wu.Email,
			&wu.Name,
			&wu.Status,
			&wu.Role,
			&wu.WorkspaceUserStatus,
			&wu.CreatedAt,
			&wu.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace user: %w", err)
		}
		users = append(users, wu)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrNoUsersFound
	}

	return users, nil
}

// ListUserWorkspaces lists all workspaces of a user
func (r *userRepo) ListUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT workspace_id
		FROM workspace_users
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
	`

	rows, err := runner.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []uuid.UUID
	for rows.Next() {
		var workspaceID uuid.UUID
		if err := rows.Scan(&workspaceID); err != nil {
			return nil, fmt.Errorf("failed to scan workspace ID: %w", err)
		}
		workspaces = append(workspaces, workspaceID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return workspaces, nil
}

// UpdateWorkspaceUserRole updates the role of users in the workspace
func (r *userRepo) UpdateWorkspaceUserRole(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, role models.UserWorkspaceRole) ([]models.UserWorkspaceRole, []error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, 0, len(userIDs)+2)
	args = append(args, workspaceID, role)

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+3)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE workspace_users
		SET role = $2, updated_at = CURRENT_TIMESTAMP
		WHERE workspace_id = $1 AND user_id IN (%s)
		RETURNING user_id, role
	`, strings.Join(placeholders, ","))

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to update workspace user roles: %w", err)}
	}
	defer rows.Close()

	updatedRoles := make([]models.UserWorkspaceRole, 0, len(userIDs))
	for rows.Next() {
		var userID uuid.UUID
		var updatedRole models.UserWorkspaceRole
		if err := rows.Scan(&userID, &updatedRole); err != nil {
			return nil, []error{fmt.Errorf("failed to scan updated role: %w", err)}
		}
		updatedRoles = append(updatedRoles, updatedRole)
	}

	if err = rows.Err(); err != nil {
		return nil, []error{fmt.Errorf("error iterating over rows: %w", err)}
	}

	return updatedRoles, nil
}

// UpdateWorkspaceUserStatus updates the status of users in the workspace
func (r *userRepo) UpdateWorkspaceUserStatus(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, status models.UserWorkspaceStatus) ([]models.UserWorkspaceStatus, []error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, 0, len(userIDs)+2)
	args = append(args, workspaceID, status)

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+3)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE workspace_users
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE workspace_id = $1 AND user_id IN (%s)
		RETURNING user_id, status
	`, strings.Join(placeholders, ","))

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to update workspace user statuses: %w", err)}
	}
	defer rows.Close()

	updatedStatuses := make([]models.UserWorkspaceStatus, 0, len(userIDs))
	for rows.Next() {
		var userID uuid.UUID
		var updatedStatus models.UserWorkspaceStatus
		if err := rows.Scan(&userID, &updatedStatus); err != nil {
			return nil, []error{fmt.Errorf("failed to scan updated status: %w", err)}
		}
		updatedStatuses = append(updatedStatuses, updatedStatus)
	}

	if err = rows.Err(); err != nil {
		return nil, []error{fmt.Errorf("error iterating over rows: %w", err)}
	}

	return updatedStatuses, nil
}

// GetWorkspaceUserCountByRole gets the count of users by role in the workspace
func (r *userRepo) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT
			COUNT(CASE WHEN role = 'admin' AND status = 'active' THEN 1 END) as admin_count,
			COUNT(CASE WHEN role = 'user' AND status = 'active' THEN 1 END) as user_count
		FROM workspace_users
		WHERE workspace_id = $1
	`

	var adminCount, userCount int
	err := runner.QueryRowContext(ctx, query, workspaceID).Scan(&adminCount, &userCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get workspace user counts: %w", err)
	}

	return adminCount, userCount, nil
}

// ------------------------------------------
// ----- Sync operations for user table -----
// -------------------------------------------

// SyncUser syncs user data with Clerk
func (r *userRepo) SyncUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) error {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		UPDATE users
		SET clerk_id = $1,
			email = $2,
			name = $3,
			status = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	result, err := runner.ExecContext(ctx, query,
		clerkUser.ID,
		clerkUser.PrimaryEmailAddressID,
		fmt.Sprintf("%s %s", *clerkUser.FirstName, *clerkUser.LastName),
		models.AccountStatusActive,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to sync user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes a user and removes them from all workspaces
func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	// First, verify the user exists
	var exists bool
	err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return ErrUserNotFound
	}

	// First update all workspace_users entries to inactive
	// This maintains referential history while effectively removing the user
	_, err = runner.ExecContext(ctx, `
		UPDATE workspace_users
		SET status = 'inactive',
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to update workspace users: %w", err)
	}

	// Then mark the user as inactive
	// We use UPDATE instead of DELETE to maintain referential integrity and history
	result, err := runner.ExecContext(ctx, `
		UPDATE users
		SET status = 'inactive',
			email = NULL,
			clerk_id = NULL,
			name = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ------------------------------------------
// -----  Optimized Lookup Operations -----
// ------------------------------------------

// UserExists checks if a user exists by userID
func (r *userRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	var exists bool
	err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}

// ClerkUserExists checks if a user exists by clerk credentials
func (r *userRepo) ClerkUserExists(ctx context.Context, clerkID, clerkEmail string) (bool, error) {
	runner := r.tm.GetRunner(ctx)
	var exists bool
	err := runner.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE clerk_id = $1 OR email = $2)",
		clerkID, clerkEmail,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check clerk user existence: %w", err)
	}
	return exists, nil
}

// WorkspaceUserExists checks if a user exists in the workspace
func (r *userRepo) WorkspaceUserExists(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	var exists bool
	err := runner.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM workspace_users WHERE workspace_id = $1 AND user_id = $2)",
		workspaceID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check workspace user existence: %w", err)
	}
	return exists, nil
}

// WorkspaceClerkUserExists checks if a user exists in the workspace by clerk credentials
func (r *userRepo) WorkspaceClerkUserExists(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (bool, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM workspace_users wu
			JOIN users u ON u.id = wu.user_id
			WHERE wu.workspace_id = $1 AND (u.clerk_id = $2 OR u.email = $3)
		)
	`

	var exists bool
	err := runner.QueryRowContext(ctx, query, workspaceID, clerkID, clerkEmail).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check workspace clerk user existence: %w", err)
	}

	return exists, nil
}

// ClerkUserIsAdmin checks if a clerk user is an admin in the workspace
func (r *userRepo) ClerkUserIsAdmin(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM workspace_users wu
			JOIN users u ON u.id = wu.user_id
			WHERE wu.workspace_id = $1
			AND u.clerk_id = $2
			AND wu.role = 'admin'
			AND wu.status = 'active'
		)
	`

	var isAdmin bool
	err := runner.QueryRowContext(ctx, query, workspaceID, clerkID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("failed to check if clerk user is admin: %w", err)
	}

	return isAdmin, nil
}

// ClerkUserIsMember checks if a clerk user is a member of the workspace
func (r *userRepo) ClerkUserIsMember(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, error) {
	// Get a transaction runner
    runner := r.tm.GetRunner(ctx)

    query := `
		SELECT EXISTS(
			SELECT 1
			FROM workspace_users wu
			JOIN users u ON u.id = wu.user_id
			WHERE wu.workspace_id = $1
			AND u.clerk_id = $2
			AND wu.status = 'active'
		)
	`

	var isMember bool
	err := runner.QueryRowContext(ctx, query, workspaceID, clerkID).Scan(&isMember)
	if err != nil {
		return false, fmt.Errorf("failed to check if clerk user is member: %w", err)
	}

	return isMember, nil
}
