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
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/err"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
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
// Respects the status of the user and excludes soft-deleted users
func (r *userRepo) GetUserByUserID(ctx context.Context, userID uuid.UUID) (models.User, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
        SELECT id, clerk_id, email, name, status, created_at, updated_at
        FROM users
        WHERE id = $1
        AND status != $2
    `

	var user models.User
	err := runner.QueryRowContext(ctx, query,
		userID,
		models.AccountStatusInactive, // Add status check to exclude soft-deleted users
	).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Propagate any other error
	if err != nil {
		// Remap the error to ErrUserNotFound
		if err == sql.ErrNoRows {
			userErr.Add(repo.ErrUserNotFoundByID, map[string]any{
				"userID": userID,
			})
		} else {
			userErr.Add(err, map[string]any{
				"userID": userID,
			})
		}
		return models.User{}, userErr
	}

	return user, nil
}

// GetUserByEmail gets a user by email
// Respects the status of the user and excludes soft-deleted users
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (models.User, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	email = utils.NormalizeEmail(email)

	query := `
        SELECT id, clerk_id, email, name, status, created_at, updated_at
        FROM users
        WHERE email = $1
        AND status != $2
    `

	var user models.User
	err := runner.QueryRowContext(ctx, query,
		email,
		models.AccountStatusInactive, // Add status check to exclude soft-deleted users
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
		// Remap the error to ErrUserNotFound
		if err == sql.ErrNoRows {
			userErr.Add(repo.ErrUserNotFoundByEmail, map[string]any{
				"email": email,
			})
		} else {
			userErr.Add(err, map[string]any{
				"email": email,
			})
		}
		return models.User{}, userErr
	}

	return user, nil
}

// GetClerkUser gets a user by clerkID or clerkEmail
// Respects the status of the user and excludes soft-deleted users
func (r *userRepo) GetClerkUser(ctx context.Context, clerkID string, clerkEmail string) (models.User, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	clerkEmail = utils.NormalizeEmail(clerkEmail)

	query := `
        SELECT id, clerk_id, email, name, status, created_at, updated_at
        FROM users
        WHERE (clerk_id = $1 OR email = $2)
        AND status != $3
    `

	var user models.User
	err := runner.QueryRowContext(ctx, query, clerkID, clerkEmail, models.AccountStatusInactive).Scan(
		&user.ID,
		&user.ClerkID,
		&user.Email,
		&user.Name,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Propagate any other error
	if err != nil {
		// Remap the error to ErrUserNotFound
		if err == sql.ErrNoRows {
			userErr.Add(repo.ErrUserNotFoundByIDOrEmail, map[string]any{
				"clerkID":    clerkID,
				"clerkEmail": clerkEmail,
			})
		} else {
			userErr.Add(err, map[string]any{
				"clerkID":    clerkID,
				"clerkEmail": clerkEmail,
			})
		}
		return models.User{}, userErr
	}

	return user, nil
}

// GetOrCreateUser creates a user if it does not exist
// Respects the status of the user and excludes soft-deleted users during lookup
// And if the user is soft-deleted, it will be restored with the new data
// If the user does not exist, it will be created
func (r *userRepo) GetOrCreateUser(ctx context.Context, partialUser models.User) (models.User, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	// Validate email
	if partialUser.Email == nil {
		userErr.Add(repo.ErrInvalidUserEmail, map[string]any{
			"email":   partialUser.Email,
			"clerkID": partialUser.ClerkID,
		})
		return models.User{}, userErr
	}

	partialUser.Email = utils.ToPtr(utils.NormalizeEmail(*partialUser.Email))

	// First, check if the user exists (including soft-deleted)
	checkQuery := `
        SELECT id, clerk_id, email, name, status, created_at, updated_at
        FROM users
        WHERE (clerk_id = $1 OR email = $2)
    `

	var existingUser models.User
	err := runner.QueryRowContext(ctx, checkQuery,
		partialUser.ClerkID,
		partialUser.Email,
	).Scan(
		&existingUser.ID,
		&existingUser.ClerkID,
		&existingUser.Email,
		&existingUser.Name,
		&existingUser.Status,
		&existingUser.CreatedAt,
		&existingUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// User doesn't exist at all, create new user
		if partialUser.Status == "" {
			partialUser.Status = models.AccountStatusPending
		}

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
			userErr.Add(err, map[string]any{
				"clerkID": partialUser.ClerkID,
				"email":   partialUser.Email,
			})
			return models.User{}, userErr
		}

		err = row.Scan(
			&existingUser.ID,
			&existingUser.ClerkID,
			&existingUser.Email,
			&existingUser.Name,
			&existingUser.Status,
			&existingUser.CreatedAt,
			&existingUser.UpdatedAt,
		)

		if err != nil {
			userErr.Add(repo.ErrCouldnotScanUser, map[string]any{
				"clerkID": partialUser.ClerkID,
				"email":   partialUser.Email,
			})
			return models.User{}, userErr
		}
	} else if err != nil {
		// Handle other errors
		userErr.Add(err, map[string]any{
			"clerkID": partialUser.ClerkID,
			"email":   partialUser.Email,
		})
		return models.User{}, userErr
	} else if existingUser.Status == models.AccountStatusInactive {
		// User exists but is soft-deleted, restore it
		updateQuery := `
            UPDATE users
            SET clerk_id = $2,
                email = $3,
                name = $4,
                status = $5,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = $1
            RETURNING id, clerk_id, email, name, status, created_at, updated_at
        `

		// Set status for restoration
		if partialUser.Status == "" {
			partialUser.Status = models.AccountStatusPending
		}

		err = runner.QueryRowContext(ctx, updateQuery,
			existingUser.ID,
			partialUser.ClerkID,
			partialUser.Email,
			partialUser.Name,
			partialUser.Status,
		).Scan(
			&existingUser.ID,
			&existingUser.ClerkID,
			&existingUser.Email,
			&existingUser.Name,
			&existingUser.Status,
			&existingUser.CreatedAt,
			&existingUser.UpdatedAt,
		)

		if err != nil {
			userErr.Add(err, map[string]any{
				"userID":  existingUser.ID,
				"clerkID": partialUser.ClerkID,
				"email":   partialUser.Email,
			})
			return models.User{}, userErr
		}
	}

	return existingUser, nil
}

// GetOrCreateUserByEmail creates users if they do not exist
// Respects the status of the user and excludes soft-deleted users during lookup
// And if the user is soft-deleted, it will be restored with the new data
// If the user does not exist, it will be created
func (r *userRepo) GetOrCreateUserByEmail(ctx context.Context, emails []string) ([]models.User, err.Error) {
	userErr := err.New()
	if len(emails) == 0 {
		userErr.Add(repo.ErrUserEmailsNotSpecified, map[string]any{
			"emails": emails,
		})
		return nil, userErr
	}

	runner := r.tm.GetRunner(ctx)
	users := make([]models.User, 0, len(emails))

	for _, email := range emails {
		var user models.User
		email = utils.NormalizeEmail(email)
		user.Email = utils.ToPtr(email)

		// First check if user exists (including soft-deleted)
		checkQuery := `
            SELECT id, clerk_id, email, name, status, created_at, updated_at
            FROM users
            WHERE email = $1
        `
		err := runner.QueryRowContext(ctx, checkQuery, email).Scan(
			&user.ID,
			&user.ClerkID,
			&user.Email,
			&user.Name,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			// User doesn't exist, create new user
			insertQuery := `
                INSERT INTO users (email, status)
                VALUES ($1, $2)
                RETURNING id, clerk_id, email, name, status, created_at, updated_at
            `
			row := runner.QueryRowContext(ctx, insertQuery,
				email,
				models.AccountStatusPending,
			)

			if err := row.Err(); err != nil {
				userErr.Add(err, map[string]any{
					"email": email,
				})
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
				userErr.Add(repo.ErrCouldnotScanUser, map[string]any{
					"email": email,
				})
				continue
			}
		} else if err != nil {
			// Handle other errors
			userErr.Add(err, map[string]any{
				"email": email,
			})
			continue
		} else if user.Status == models.AccountStatusInactive {
			// User exists but is soft-deleted, restore it
			updateQuery := `
                UPDATE users
                SET status = $2,
                    email = $3,
                    updated_at = CURRENT_TIMESTAMP
                WHERE id = $1
                RETURNING id, clerk_id, email, name, status, created_at, updated_at
            `

			err = runner.QueryRowContext(ctx, updateQuery,
				user.ID,
				models.AccountStatusPending,
				email,
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
				userErr.Add(err, map[string]any{
					"email":  email,
					"userID": user.ID,
				})
				continue
			}
		}

		users = append(users, user)
	}

	// Return both users and any errors that occurred
	if userErr.HasErrors() {
		return users, userErr
	}

	return users, nil
}

// -----------------------------------------------------------
// --- Create & Update operations for workspace_user table ---
// -----------------------------------------------------------

// AddUsersToWorkspace adds users to a workspace
// Gets the user by email and creates a workspace_user entry
// Underlying dependents respect the soft-deleted status of user
func (r *userRepo) AddUsersToWorkspace(ctx context.Context, workspaceUserProps []models.WorkspaceUserProps, workspaceID uuid.UUID) ([]models.WorkspaceUser, err.Error) {
	userErr := err.New()
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)
	emailToUserID := make(map[string]uuid.UUID)

	// Get or create users for each email
	emails := make([]string, 0)
	for i, props := range workspaceUserProps {
		if err := utils.SetDefaultsAndValidate(&workspaceUserProps[i]); err != nil {
			userErr.Add(err, map[string]any{
				"workspaceUserProps": workspaceUserProps[i],
			})
			continue
		}
		emails = append(emails, props.Email)
	}

	users, errs := r.GetOrCreateUserByEmail(ctx, emails)
	if errs != nil && errs.HasErrors() {
		userErr.Merge(errs)
	}

	if len(users) == 0 {
		return nil, userErr
	}

	for _, user := range users {
		emailToUserID[*user.Email] = user.ID
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
		userErr.Add(err, map[string]any{
			"query": insertQuery,
		})
		return nil, userErr
	}
	defer rows.Close()

	var workspaceUsers []models.WorkspaceUser
	for rows.Next() {
		var wu models.WorkspaceUser

		if err := rows.Scan(
			&wu.WorkspaceID,
			&wu.ID, // user_id
			&wu.Role,
			&wu.WorkspaceUserStatus,
		); err != nil {
			userErr.Add(err, map[string]any{
				"query": insertQuery,
			})
			continue
		}

		// Populate the user fields
		user, err := r.GetUserByUserID(ctx, wu.ID)
		if err != nil {
			userErr.Merge(err)
			continue
		}

		wu.ClerkID = user.ClerkID
		wu.Email = user.Email
		wu.Name = user.Name
		wu.Status = user.Status

		workspaceUsers = append(workspaceUsers, wu)
	}

	if userErr.HasErrors() {
		return nil, userErr
	}

	if err = rows.Err(); err != nil {
		userErr.Add(err, map[string]any{
			"query": insertQuery,
		})
		return nil, userErr
	}

	return workspaceUsers, nil
}

// RemoveUsersFromWorkspace removes users from a workspace
// When userIDs is nil, all users are removed from the workspace
// Only removes users that are not already soft-deleted (either at user or workspace level)
func (r *userRepo) RemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) err.Error {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	var query string
	var args []interface{}

	if userIDs == nil {
		// Remove all active users from workspace
		query = `
            WITH updated AS (
                UPDATE workspace_users wu
                SET status = $2,
                    updated_at = CURRENT_TIMESTAMP
                FROM users u
                WHERE wu.workspace_id = $1
                AND wu.user_id = u.id
                AND wu.status != $2           -- Don't update already inactive workspace users
                AND u.status != $3            -- Don't update soft-deleted users
                RETURNING wu.user_id
            )
            SELECT COUNT(*) FROM updated
        `
		args = []interface{}{
			workspaceID,
			models.UserWorkspaceStatusInactive,
			models.AccountStatusInactive,
		}
	} else {
		// Remove specific users from workspace
		placeholders := make([]string, len(userIDs))
		args = make([]interface{}, 0, len(userIDs)+3)
		args = append(args,
			workspaceID,
			models.UserWorkspaceStatusInactive,
			models.AccountStatusInactive,
		)

		for i := range userIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+4) // Start from $4 now
			args = append(args, userIDs[i])
		}

		query = fmt.Sprintf(`
            WITH updated AS (
                UPDATE workspace_users wu
                SET status = $2,
                    updated_at = CURRENT_TIMESTAMP
                FROM users u
                WHERE wu.workspace_id = $1
                AND wu.user_id = u.id
                AND wu.user_id IN (%s)
                AND wu.status != $2           -- Don't update already inactive workspace users
                AND u.status != $3            -- Don't update soft-deleted users
                RETURNING wu.user_id
            )
            SELECT COUNT(*) FROM updated
        `, strings.Join(placeholders, ","))
	}

	var count int64
	err := runner.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"userIDs":     userIDs,
			"query":       query,
		})
		return userErr
	}

	// Only return error if no users were affected and specific users were requested
	if count == 0 && userIDs != nil {
		userErr.Add(repo.ErrWorkspaceUsersNotFound, map[string]any{
			"workspaceID": workspaceID,
			"userIDs":     userIDs,
		})
		return userErr
	}

	return nil
}

// GetWorkspaceUser gets a user from the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) GetWorkspaceUser(ctx context.Context, workspaceID, userID uuid.UUID) (models.WorkspaceUser, err.Error) {
	userErr := err.New()
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
        AND u.status != $3
        AND wu.status != $4
    `

	var wu models.WorkspaceUser
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		userID,
		models.AccountStatusInactive,
		models.UserWorkspaceStatusInactive,
	).Scan(
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
		userErr.Add(repo.ErrWorkspaceUsersNotFound, map[string]any{
			"workspaceID": workspaceID,
			"userID":      userID,
		})
		return models.WorkspaceUser{}, userErr
	}
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"userID":      userID,
		})
		return models.WorkspaceUser{}, userErr
	}

	return wu, nil
}

// GetWorkspaceClerkUser gets a user from the workspace by clerk credentials
// It respects the status of the user and workspace_user entries
func (r *userRepo) GetWorkspaceClerkUser(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (models.WorkspaceUser, err.Error) {
	userErr := err.New()

	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	clerkEmail = utils.NormalizeEmail(clerkEmail)

	query := `
        SELECT u.id, u.clerk_id, u.email, u.name, u.status,
               wu.role, wu.status as workspace_status,
               u.created_at, u.updated_at
        FROM users u
        JOIN workspace_users wu ON u.id = wu.user_id
        WHERE wu.workspace_id = $1
        AND (u.clerk_id = $2 OR u.email = $3)
        AND u.status != $4
        AND wu.status != $5
    `

	var wu models.WorkspaceUser
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		clerkID,
		clerkEmail,
		models.AccountStatusInactive,
		models.UserWorkspaceStatusInactive,
	).Scan(
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
		userErr.Add(repo.ErrWorkspaceUsersNotFound, map[string]any{
			"workspaceID": workspaceID,
			"clerkID":     clerkID,
			"clerkEmail":  clerkEmail,
		})
		return models.WorkspaceUser{}, userErr
	}
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"clerkID":     clerkID,
			"clerkEmail":  clerkEmail,
		})
		return models.WorkspaceUser{}, userErr
	}

	return wu, nil
}

// ListWorkspaceUsers lists all users from the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) ListWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceUser, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	// Modified query to:
	// 1. Only exclude inactive accounts (allow pending and active)
	// 2. Exclude inactive workspace_user entries
	query := `
        SELECT u.id, u.clerk_id, u.email, u.name, u.status, wu.workspace_id,
               wu.role, wu.status as workspace_status,
               u.created_at, u.updated_at
        FROM users u
        JOIN workspace_users wu ON u.id = wu.user_id
        WHERE wu.workspace_id = $1
        AND u.status != $2         -- Exclude inactive users
        AND wu.status != $3        -- Exclude inactive workspace_user entries
        ORDER BY u.created_at DESC
    `

	rows, err := runner.QueryContext(ctx, query,
		workspaceID,
		models.AccountStatusInactive,       // Only exclude inactive users
		models.UserWorkspaceStatusInactive, // Exclude inactive workspace users
	)
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return nil, userErr
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
			&wu.WorkspaceID,
			&wu.Role,
			&wu.WorkspaceUserStatus,
			&wu.CreatedAt,
			&wu.UpdatedAt,
		)
		if err != nil {
			userErr.Add(repo.ErrCouldnotScanUser, map[string]any{
				"workspaceID": workspaceID,
			})
			continue
		}
		users = append(users, wu)
	}

	if err = rows.Err(); err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return nil, userErr
	}

	// Only return ErrWorkspaceUsersNotFound if we got no results
	if len(users) == 0 {
		userErr.Add(repo.ErrWorkspaceUsersNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return nil, userErr
	}

	return users, nil
}

// ListUserWorkspaces lists all workspaces of a user
// It respects the status of the user and workspace_user entries
func (r *userRepo) ListUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
        WITH workspace_list AS (
            SELECT wu.workspace_id
            FROM workspace_users wu
            JOIN users u ON u.id = wu.user_id
            WHERE wu.user_id = $1
            AND wu.status = $2          -- Only active workspace memberships
            AND u.status != $3          -- Exclude soft-deleted users
            ORDER BY wu.created_at DESC
        )
        SELECT wl.workspace_id
        FROM workspace_list wl
    `

	rows, err := runner.QueryContext(ctx, query,
		userID,
		models.UserWorkspaceStatusActive, // Only include active workspace memberships
		models.AccountStatusInactive,     // Exclude soft-deleted users
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userErr.Add(repo.ErrNoWorkspaceFoundForUser, map[string]any{
				"userID": userID,
			})
		} else {
			userErr.Add(err, map[string]any{
				"userID": userID,
				"query":  query,
			})
		}
		return nil, userErr
	}
	defer rows.Close()

	workspaces := make([]uuid.UUID, 0) // Initialize with empty slice instead of nil
	for rows.Next() {
		var workspaceID uuid.UUID
		if err := rows.Scan(&workspaceID); err != nil {
			userErr.Add(err, map[string]any{
				"userID": userID,
			})
			continue
		}
		workspaces = append(workspaces, workspaceID)
	}

	if err = rows.Err(); err != nil {
		userErr.Add(err, map[string]any{
			"userID": userID,
			"query":  query,
		})
	}

	if userErr.HasErrors() {
		return nil, userErr
	}

	return workspaces, nil // Will return empty slice if no rows found
}

// UpdateWorkspaceUserRole updates the role of users in the workspace
// Skips soft-deleted users and inactive workspace users
func (r *userRepo) UpdateWorkspaceUserRole(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, role models.UserWorkspaceRole) ([]models.UserWorkspaceRole, err.Error) {
	userErr := err.New()
	if userIDs == nil {
		userErr.Add(repo.ErrUserIDsNotSpecified, map[string]any{
			"workspaceID": workspaceID,
			"role":        role,
		})
		return nil, userErr
	}

	runner := r.tm.GetRunner(ctx)

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, 0, len(userIDs)+4)
	args = append(args,
		workspaceID,
		role,
		models.AccountStatusInactive,
		models.UserWorkspaceStatusInactive,
	)

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+5) // Start from $5
		args = append(args, id)
	}

	query := fmt.Sprintf(`
        UPDATE workspace_users wu
        SET role = $2,
            updated_at = CURRENT_TIMESTAMP
        FROM users u
        WHERE wu.workspace_id = $1
        AND wu.user_id = u.id
        AND u.status != $3                    -- Skip soft-deleted users
        AND wu.status != $4                   -- Skip inactive workspace users
        AND wu.user_id IN (%s)
        RETURNING wu.user_id, wu.role
    `, strings.Join(placeholders, ","))

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"role":        role,
			"userIDs":     userIDs,
			"query":       query,
		})
		return nil, userErr
	}
	defer rows.Close()

	updatedRoles := make([]models.UserWorkspaceRole, 0, len(userIDs))
	for rows.Next() {
		var userID uuid.UUID
		var updatedRole models.UserWorkspaceRole
		if err := rows.Scan(&userID, &updatedRole); err != nil {
			userErr.Add(err, map[string]any{
				"workspaceID": workspaceID,
				"userID":      userID,
			})
			continue
		}
		updatedRoles = append(updatedRoles, updatedRole)
	}

	if err = rows.Err(); err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"role":        role,
			"userIDs":     userIDs,
			"query":       query,
		})
	}

	if userErr.HasErrors() {
		return nil, userErr
	}

	return updatedRoles, nil
}

// UpdateWorkspaceUserStatus updates the status of users in the workspace
// Skip soft-deleted user accounts as they are no longer considered active in the app
func (r *userRepo) UpdateWorkspaceUserStatus(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID, status models.UserWorkspaceStatus) ([]models.UserWorkspaceStatus, err.Error) {
	userErr := err.New()
	if userIDs == nil {
		userErr.Add(repo.ErrUserIDsNotSpecified, map[string]any{
			"workspaceID": workspaceID,
			"status":      status,
		})
		return nil, userErr
	}
	runner := r.tm.GetRunner(ctx)

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, 0, len(userIDs)+3)
	args = append(args,
		workspaceID,
		status,
		models.AccountStatusInactive,
	)

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+4)
		args = append(args, id)
	}

	// Only allow updates for non-soft-deleted users
	query := fmt.Sprintf(`
       WITH valid_users AS (
           SELECT id
           FROM users
           WHERE status != $3
           AND id IN (%s)
       )
       UPDATE workspace_users wu
       SET status = $2,
           updated_at = CURRENT_TIMESTAMP
       FROM valid_users vu
       WHERE wu.workspace_id = $1
       AND wu.user_id = vu.id
       RETURNING wu.user_id, wu.status
   `, strings.Join(placeholders, ","))

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"status":      status,
			"userIDs":     userIDs,
			"query":       query,
		})
		return nil, userErr
	}
	defer rows.Close()

	updatedStatuses := make([]models.UserWorkspaceStatus, 0, len(userIDs))
	for rows.Next() {
		var userID uuid.UUID
		var updatedStatus models.UserWorkspaceStatus
		if err := rows.Scan(&userID, &updatedStatus); err != nil {
			userErr.Add(err, map[string]any{
				"workspaceID": workspaceID,
				"userID":      userID,
			})
			continue
		}
		updatedStatuses = append(updatedStatuses, updatedStatus)
	}

	if err = rows.Err(); err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"status":      status,
			"userIDs":     userIDs,
			"query":       query,
		})
	}

	// If no users were updated but specific IDs were provided, return an error
	if len(updatedStatuses) == 0 && len(userIDs) > 0 {
		userErr.Add(repo.ErrWorkspaceUsersNotFound, map[string]any{
			"workspaceID": workspaceID,
			"userIDs":     userIDs,
		})
		return nil, userErr
	}

	return updatedStatuses, nil
}

// GetWorkspaceUserCountByRole gets the count of users by role in the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (int, int, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
        SELECT
            COUNT(CASE WHEN wu.role = 'admin' AND wu.status = $2 THEN 1 END) as admin_count,
            COUNT(CASE WHEN wu.role = 'user' AND wu.status = $2 THEN 1 END) as user_count
        FROM workspace_users wu
        JOIN users u ON u.id = wu.user_id
        WHERE wu.workspace_id = $1
        AND u.status != $3
    `

	var adminCount, userCount int
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		models.UserWorkspaceStatusActive,
		models.AccountStatusInactive,
	).Scan(&adminCount, &userCount)

	if err != nil {
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"query":       query,
		})
		return 0, 0, userErr
	}

	return adminCount, userCount, nil
}

// ------------------------------------------
// ----- Sync operations for user table -----
// -------------------------------------------

// SyncUser syncs user data with Clerk
// Whenever sync is triggered, it will activate the user from inactive state
func (r *userRepo) SyncUser(ctx context.Context, userID uuid.UUID, clerkUser *clerk.User) err.Error {
	userErr := err.New()
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
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return userErr
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
		userErr.Add(err, map[string]any{
			"userID": userID,
			"query":  query,
		})
		return userErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return userErr
	}

	if rowsAffected == 0 {
		userErr.Add(repo.ErrUserNotFoundByID, map[string]any{
			"userID": userID,
		})
		return userErr
	}

	return nil
}

// DeleteUser deletes a user and removes them from all workspaces
// It respects the status of the user and workspace_user entries
func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) err.Error {
	userErr := err.New()
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	// First, verify the user exists
	var exists bool
	err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			userErr.Add(repo.ErrUserNotFoundByID, map[string]any{
				"userID": userID,
			})
		} else {
			userErr.Add(err, map[string]any{
				"userID": userID,
			})
		}
		return userErr
	}
	if !exists {
		userErr.Add(repo.ErrUserNotFoundByID, map[string]any{
			"userID": userID,
		})
		return userErr
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
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return userErr
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
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return userErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return userErr
	}

	if rowsAffected == 0 {
		userErr.Add(repo.ErrUserNotFoundByID, map[string]any{
			"userID": userID,
		})
		return userErr
	}

	return nil
}

// ------------------------------------------
// -----  Optimized Lookup Operations -----
// ------------------------------------------

// UserExists checks if a user exists by userID
// It respects the status of the user
func (r *userRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, err.Error) {
	userErr := err.New()
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM users
            WHERE id = $1
            AND status != $2
        )
    `
	err := runner.QueryRowContext(ctx, query, userID, models.AccountStatusInactive).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"userID": userID,
		})
		return false, userErr
	}
	return exists, nil
}

// ClerkUserExists checks if a user exists by clerk credentials
// It respects the status of the user
func (r *userRepo) ClerkUserExists(ctx context.Context, clerkID, clerkEmail string) (bool, err.Error) {
	// Normalize email
	clerkEmail = utils.NormalizeEmail(clerkEmail)

	userErr := err.New()

	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)
	var exists bool

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM users
            WHERE (clerk_id = $1 OR email = $2)
            AND status != $3
        )`

	if err := runner.QueryRowContext(ctx,
		query,
		clerkID,
		clerkEmail,
		models.AccountStatusInactive, // Add status check to exclude soft-deleted users
	).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"clerkID":    clerkID,
			"clerkEmail": clerkEmail,
		})
		return false, userErr
	}
	return exists, nil
}

// WorkspaceUserExists checks if a user exists in the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) WorkspaceUserExists(ctx context.Context, workspaceID, userID uuid.UUID) (bool, err.Error) {
	userErr := err.New()
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM workspace_users wu
            JOIN users u ON u.id = wu.user_id
            WHERE wu.workspace_id = $1
            AND wu.user_id = $2
            AND u.status != $3
            AND wu.status != $4
        )
    `
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		userID,
		models.AccountStatusInactive,
		models.UserWorkspaceStatusInactive,
	).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"userID":      userID,
		})
		return false, userErr
	}
	return exists, nil
}

// WorkspaceClerkUserExists checks if a user exists in the workspace by clerk credentials
// It respects the status of the user and workspace_user entries
func (r *userRepo) WorkspaceClerkUserExists(ctx context.Context, workspaceID uuid.UUID, clerkID, clerkEmail string) (bool, err.Error) {
	userErr := err.New()
	// Get a transaction runner
	runner := r.tm.GetRunner(ctx)

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM workspace_users wu
            JOIN users u ON u.id = wu.user_id
            WHERE wu.workspace_id = $1
            AND (u.clerk_id = $2 OR u.email = $3)
            AND u.status != $4
            AND wu.status != $5
        )
    `

	var exists bool
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		clerkID,
		clerkEmail,
		models.AccountStatusInactive,
		models.UserWorkspaceStatusInactive,
	).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"clerkID":     clerkID,
			"clerkEmail":  clerkEmail,
		})
		return false, userErr
	}

	return exists, nil
}

// ClerkUserIsAdmin checks if a clerk user is an admin in the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) ClerkUserIsAdmin(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM workspace_users wu
            JOIN users u ON u.id = wu.user_id
            WHERE wu.workspace_id = $1
            AND u.clerk_id = $2
            AND wu.role = 'admin'
            AND wu.status = $3
            AND u.status != $4
        )
    `

	var isAdmin bool
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		clerkID,
		models.UserWorkspaceStatusActive,
		models.AccountStatusInactive,
	).Scan(&isAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"clerkID":     clerkID,
		})
		return false, userErr
	}

	return isAdmin, nil
}

// ClerkUserIsMember checks if a clerk user is a member of the workspace
// It respects the status of the user and workspace_user entries
func (r *userRepo) ClerkUserIsMember(ctx context.Context, workspaceID uuid.UUID, clerkID string) (bool, err.Error) {
	userErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM workspace_users wu
            JOIN users u ON u.id = wu.user_id
            WHERE wu.workspace_id = $1
            AND u.clerk_id = $2
            AND wu.status = $3
            AND u.status != $4
        )
    `

	var isMember bool
	err := runner.QueryRowContext(ctx, query,
		workspaceID,
		clerkID,
		models.UserWorkspaceStatusActive,
		models.AccountStatusInactive,
	).Scan(&isMember)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		userErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
			"clerkID":     clerkID,
		})
		return false, userErr
	}

	return isMember, nil
}
