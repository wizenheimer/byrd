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
	"github.com/wizenheimer/iris/src/pkg/utils"
)

var (
	// ---- validation errors -----
	ErrInvalidUserEmail          = errors.New("user email invalid")
	ErrUserEmailsNotSpecified    = errors.New("user emails unspecified")
	ErrUserIDsNotSpecified       = errors.New("user ids unspecified")
	ErrCouldnotGetClerkUserEmail = errors.New("couldn't get user email from profile")

	// ---- non fatal errors ----
	ErrCouldnotScanUser            = errors.New("couldn't scan user")
	ErrNoWorkspaceFoundForUser     = errors.New("no workspaces found for the user")
	ErrCouldnotConfirmSyncStatus   = errors.New("couldn't confirm sync status")
	ErrCouldnotConfirmDeleteStatus = errors.New("couldn't confirm delete status")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrUserNotFoundByID        = errors.New("user not found by ID")
	ErrUserNotFoundByEmail     = errors.New("user not found by email")
	ErrUserNotFoundByIDOrEmail = errors.New("user not found by ID or Email")
	ErrWorkspaceUsersNotFound  = errors.New("users not found in workspace")
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

	// Remap the error to ErrUserNotFound
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFoundByID
	}

	// Propagate any other error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	runner := r.tm.GetRunner(ctx)

	email = utils.NormalizeEmail(email)

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

	// Remap the error to ErrUserNotFound
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFoundByEmail
	}

	// Propagate any other error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetClerkUser gets a user by clerkID or clerkEmail
func (r *userRepo) GetClerkUser(ctx context.Context, clerkID string, clerkEmail string) (models.User, error) {
	runner := r.tm.GetRunner(ctx)

	clerkEmail = utils.NormalizeEmail(clerkEmail)

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

	// Remap the error to ErrUserNotFound
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFoundByIDOrEmail
	}

	// Propagate any other error
	if err != nil {
		return models.User{}, err
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

	if partialUser.Email == nil {
		return models.User{}, ErrInvalidUserEmail
	}

	partialUser.Email = utils.ToPtr(utils.NormalizeEmail(*partialUser.Email))

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
		// User does not exist
		// Create new user
		insertQuery := `
			INSERT INTO users (clerk_id, email, name, status)
			VALUES ($1, $2, $3, $4)
			RETURNING id, clerk_id, email, name, status, created_at, updated_at
		`

		row := runner.QueryRowContext(ctx, insertQuery,
			partialUser.ClerkID,
			partialUser.Email,
			partialUser.Name,
			partialUser.Status,
		)
		if err := row.Err(); err != nil {
			return models.User{}, err
		}

		err = row.Scan(
			&user.ID,
			&user.ClerkID,
			&user.Email,
			&user.Name,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return models.User{}, ErrCouldnotScanUser
		}
	}

	// Propagate any other error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetOrCreateUserByEmail creates users if they do not exist
func (r *userRepo) GetOrCreateUserByEmail(ctx context.Context, emails []string) ([]models.User, []error) {
	if len(emails) == 0 {
		return nil, []error{ErrUserEmailsNotSpecified}
	}

	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	users := make([]models.User, 0, len(emails))
	errs := make([]error, 0)

	for _, email := range emails {
		// Try to get existing user
		var user models.User

		// Normalize email
		user.Email = utils.ToPtr(utils.NormalizeEmail(email))

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
			row := runner.QueryRowContext(ctx, `
				INSERT INTO users (email, status)
				VALUES ($1, $2)
				RETURNING id, clerk_id, email, name, status, created_at, updated_at
			`, email, models.AccountStatusPending)

			if err := row.Err(); err != nil {
				errs = append(errs, err)
				continue
			}

			err = row.Scan(
				&user.ID,
				&user.ClerkID,
				&user.Email,
				&user.Name,
				&user.Status,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
			if err != nil {
				err = ErrCouldnotScanUser
			}
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}

		users = append(users, user)
	}

	// If there are any errors, return them along with the users created so far
	// Otherwise, return the users
	if len(errs) > 0 {
		return users, errs
	}

	return users, nil
}

// -----------------------------------------------------------
// --- Create & Update operations for workspace_user table ---
// -----------------------------------------------------------

// AddUsersToWorkspace adds users to a workspace
func (r *userRepo) AddUsersToWorkspace(ctx context.Context, workspaceUserProps []models.WorkspaceUserProps, workspaceID uuid.UUID) ([]models.WorkspaceUser, []error) {
	runner := r.tm.GetRunner(ctx)
	var errors []error
	emailToUserID := make(map[string]uuid.UUID)

	// Get or create users for each email
	emails := make([]string, len(workspaceUserProps))
	for i, props := range workspaceUserProps {
		emails[i] = props.Email
	}

	users, errs := r.GetOrCreateUserByEmail(ctx, emails)
	if len(errs) > 0 {
		return nil, errs
	}

	for _, user := range users {
		emailToUserID[*user.Email] = user.ID
	}

	if len(errors) > 0 {
		return nil, errors
	}

	valueStrings := make([]string, len(workspaceUserProps))
	valueArgs := make([]interface{}, 0, len(workspaceUserProps)*4)

	for i, props := range workspaceUserProps {
		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		valueArgs = append(valueArgs,
			workspaceID,
			emailToUserID[props.Email],
			props.Role,
			props.Status, // Now using the status from props
		)
	}

	insertQuery := fmt.Sprintf(`
        INSERT INTO workspace_users (workspace_id, user_id, role, status)
        VALUES %s
        ON CONFLICT (workspace_id, user_id) DO UPDATE
        SET role = EXCLUDED.role,
            status = EXCLUDED.status
        RETURNING workspace_id, user_id, role, status
    `, strings.Join(valueStrings, ","))

	rows, err := runner.QueryContext(ctx, insertQuery, valueArgs...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	var workspaceUsers []models.WorkspaceUser
	errs = make([]error, 0)
	for rows.Next() {
		var wu models.WorkspaceUser
		err := rows.Scan(
			&wu.WorkspaceID,
			&wu.ID, // user_id
			&wu.Role,
			&wu.WorkspaceUserStatus,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// Populate the user fields
		user, err := r.GetUserByUserID(ctx, wu.ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		wu.ClerkID = user.ClerkID
		wu.Email = user.Email
		wu.Name = user.Name
		wu.Status = user.Status

		workspaceUsers = append(workspaceUsers, wu)
	}

	if len(errs) > 0 {
		return workspaceUsers, errs
	}

	if err = rows.Err(); err != nil {
		return nil, []error{err}
	}

	return workspaceUsers, nil
}

// RemoveUsersFromWorkspace removes users from a workspace
// When userIDs is nil, all users are removed from the workspace
func (r *userRepo) RemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) []error {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	var query string
	var args []interface{}

	if userIDs == nil {
		// Remove all users from workspace
		query = `
            WITH updated AS (
                UPDATE workspace_users
                SET status = 'inactive',
                    updated_at = CURRENT_TIMESTAMP
                WHERE workspace_id = $1
                RETURNING user_id
            )
            SELECT COUNT(*) FROM updated
        `
		args = []interface{}{workspaceID}
	} else {
		// Remove specific users from workspace
		placeholders := make([]string, len(userIDs))
		args = make([]interface{}, 0, len(userIDs)+1)
		args = append(args, workspaceID)

		for i := range userIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, userIDs[i])
		}

		query = fmt.Sprintf(`
            WITH updated AS (
                UPDATE workspace_users
                SET status = 'inactive',
                    updated_at = CURRENT_TIMESTAMP
                WHERE workspace_id = $1
                AND user_id IN (%s)
                RETURNING user_id
            )
            SELECT COUNT(*) FROM updated
        `, strings.Join(placeholders, ","))
	}

	var count int64
	err := runner.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return []error{err}
	}

	if count == 0 {
		return []error{ErrWorkspaceUsersNotFound}
	}

	return nil
}

// GetWorkspaceUser gets a user from the workspace
func (r *userRepo) GetWorkspaceUser(ctx context.Context, workspaceID, userID uuid.UUID) (models.WorkspaceUser, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT
			u.id,
			u.clerk_id,
			u.email,
			u.name,
			u.status,
			wu.role,
			wu.workspace_id,
			wu.status AS workspace_status,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN workspace_users wu ON u.id = wu.user_id
		WHERE wu.workspace_id = $1
		AND wu.user_id = $2
	`

	var wu models.WorkspaceUser
	err := runner.QueryRowContext(ctx, query, workspaceID, userID).Scan(
		&wu.ID,
		&wu.ClerkID,
		&wu.Email,
		&wu.Name,
		&wu.Status,
		&wu.Role,
		&wu.WorkspaceID,
		&wu.WorkspaceUserStatus, // matches the AS workspace_status in query
		&wu.CreatedAt,
		&wu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.WorkspaceUser{}, ErrWorkspaceUsersNotFound
	}
	if err != nil {
		return models.WorkspaceUser{}, err
	}

	return wu, nil
}

// GetWorkspaceClerkUser gets a user from the workspace by clerk credentials
func (r *userRepo) GetWorkspaceClerkUser(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (models.WorkspaceUser, error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	clerkEmail = utils.NormalizeEmail(clerkEmail)

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
		return models.WorkspaceUser{}, ErrWorkspaceUsersNotFound
	}
	if err != nil {
		return models.WorkspaceUser{}, err
	}

	return wu, nil
}

// ListWorkspaceUsers lists all users from the workspace
func (r *userRepo) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, []error) {
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT u.id, u.clerk_id, u.email, u.name, u.status, wu.workspace_id,
			   wu.role, wu.status as workspace_status,
			   u.created_at, u.updated_at
		FROM users u
		JOIN workspace_users wu ON u.id = wu.user_id
		WHERE wu.workspace_id = $1
		ORDER BY u.created_at DESC
	`

	rows, err := runner.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	var users []models.WorkspaceUser
	errs := make([]error, 0)
	for rows.Next() {
		var wu models.WorkspaceUser
		err := rows.Scan(
			&wu.ID,
			&wu.ClerkID,
			&wu.Email,
			&wu.Name,
			&wu.Status,
			&wu.WorkspaceID,
			&wu.Role,
			&wu.WorkspaceUserStatus,
			&wu.CreatedAt,
			&wu.UpdatedAt,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		users = append(users, wu)
	}

	if errs != nil {
		return users, errs
	}

	if err = rows.Err(); err != nil {
		return nil, []error{err}
	}

	if len(users) == 0 {
		return nil, []error{ErrWorkspaceUsersNotFound}
	}

	return users, nil
}

// ListUserWorkspaces lists all workspaces of a user
func (r *userRepo) ListUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, []error) {
	runner := r.tm.GetRunner(ctx)

	query := `
        WITH workspace_list AS (
            SELECT wu.workspace_id
            FROM workspace_users wu
            WHERE wu.user_id = $1 AND wu.status = 'active'
            ORDER BY wu.created_at DESC
        )
        SELECT wl.workspace_id
        FROM workspace_list wl
    `

	rows, err := runner.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []uuid.UUID{}, []error{ErrNoWorkspaceFoundForUser}
		}
		return nil, []error{err}
	}
	defer rows.Close()

	workspaces := make([]uuid.UUID, 0) // Initialize with empty slice instead of nil
	errs := make([]error, 0)
	for rows.Next() {
		var workspaceID uuid.UUID
		if err := rows.Scan(&workspaceID); err != nil {
			errs = append(errs, err)
			continue
		}
		workspaces = append(workspaces, workspaceID)
	}

	if len(errs) > 0 {
		return workspaces, errs
	}

	if err = rows.Err(); err != nil {
		return nil, []error{err}
	}

	return workspaces, nil // Will return empty slice if no rows found
}

// UpdateWorkspaceUserRole updates the role of users in the workspace
func (r *userRepo) UpdateWorkspaceUserRole(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, role models.UserWorkspaceRole) ([]models.UserWorkspaceRole, []error) {
	if userIDs == nil {
		return nil, []error{ErrUserIDsNotSpecified}
	}
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
		return nil, []error{err}
	}
	defer rows.Close()

	updatedRoles := make([]models.UserWorkspaceRole, 0, len(userIDs))
	errs := make([]error, 0)
	for rows.Next() {
		var userID uuid.UUID
		var updatedRole models.UserWorkspaceRole
		if err := rows.Scan(&userID, &updatedRole); err != nil {
			errs = append(errs, err)
			continue
		}
		updatedRoles = append(updatedRoles, updatedRole)
	}

	if len(errs) > 0 {
		return updatedRoles, errs
	}

	if err = rows.Err(); err != nil {
		return nil, []error{err}
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
		return nil, []error{err}
	}
	defer rows.Close()

	updatedStatuses := make([]models.UserWorkspaceStatus, 0, len(userIDs))
	errs := make([]error, 0)
	for rows.Next() {
		var userID uuid.UUID
		var updatedStatus models.UserWorkspaceStatus
		if err := rows.Scan(&userID, &updatedStatus); err != nil {
			errs = append(errs, err)
			continue
		}
		updatedStatuses = append(updatedStatuses, updatedStatus)
	}

	if len(errs) > 0 {
		return updatedStatuses, errs
	}

	if err = rows.Err(); err != nil {
		return nil, []error{err}
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

	userEmail, err := utils.GetClerkUserEmail(clerkUser)
	if err != nil {
		return ErrCouldnotGetClerkUserEmail
	}

	userFullName := utils.GetClerkUserFullName(clerkUser)

	result, err := runner.ExecContext(ctx, query,
		clerkUser.ID,
		userEmail,
		userFullName,
		models.AccountStatusActive,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrCouldnotConfirmSyncStatus
	}

	if rowsAffected == 0 {
		return ErrUserNotFoundByIDOrEmail
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
		return err
	}
	if !exists {
		return ErrUserNotFoundByIDOrEmail
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
		return err
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
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrCouldnotConfirmDeleteStatus
	}

	if rowsAffected == 0 {
		return ErrUserNotFoundByIDOrEmail
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return isMember, nil
}
