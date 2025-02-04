// ./src/internal/repository/workspace/workspace.go
package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewWorkspaceRepository(tm *transaction.TxManager, logger *logger.Logger) WorkspaceRepository {
	return &workspaceRepo{
		tm: tm,
		logger: logger.WithFields(map[string]interface{}{
			"module": "workspace_repository",
		}),
	}
}

func (r *workspaceRepo) CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID, workspacePlan models.WorkspacePlan) (*models.Workspace, error) {
	workspaceSlug := getSlug(workspaceName)
	workspace := &models.Workspace{}

	// Create workspace
	err := r.getQuerier(ctx).QueryRow(ctx, `
        INSERT INTO workspaces (name, slug, billing_email, workspace_status, workspace_plan)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, name, slug, billing_email, workspace_status, workspace_plan, created_at, updated_at`,
		workspaceName, workspaceSlug, billingEmail, models.WorkspaceActive, workspacePlan,
	).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.BillingEmail,
		&workspace.WorkspaceStatus,
		&workspace.WorkspacePlan,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Add creator as admin
	_, err = r.getQuerier(ctx).Exec(ctx, `
		INSERT INTO workspace_users (workspace_id, user_id, workspace_role, membership_status)
		VALUES ($1, $2, $3, $4)`,
		workspace.ID, workspaceCreatorUserID, models.RoleAdmin, models.ActiveMember,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add creator to workspace: %w", err)
	}

	return workspace, nil
}

func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	var exists bool
	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM workspaces
			WHERE id = $1 AND workspace_status != $2
		)`,
		workspaceID, models.WorkspaceInactive,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
	}

	return exists, nil
}

func (r *workspaceRepo) GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	workspace := &models.Workspace{}
	err := r.getQuerier(ctx).QueryRow(ctx, `
        SELECT id, name, slug, billing_email, workspace_status, workspace_plan, created_at, updated_at
        FROM workspaces
        WHERE id = $1 AND workspace_status != $2`,
		workspaceID, models.WorkspaceInactive,
	).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.BillingEmail,
		&workspace.WorkspaceStatus,
		&workspace.WorkspacePlan,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("workspace not found")
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	return workspace, nil
}

func (r *workspaceRepo) BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error) {
	if len(workspaceIDs) == 0 {
		return []models.Workspace{}, nil
	}

	rows, err := r.getQuerier(ctx).Query(ctx, `
    SELECT id, name, slug, billing_email, workspace_status, workspace_plan, created_at, updated_at
    FROM workspaces
    WHERE id = ANY($1) AND workspace_status != $2`,
		workspaceIDs, models.WorkspaceInactive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get workspaces: %w", err)
	}
	defer rows.Close()

	workspaces := make([]models.Workspace, 0)
	for rows.Next() {
		var workspace models.Workspace
		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Slug,
			&workspace.BillingEmail,
			&workspace.WorkspaceStatus,
			&workspace.WorkspacePlan,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, rows.Err()
}

func (r *workspaceRepo) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, workspaceRole *models.WorkspaceRole) ([]models.PartialWorkspaceUser, bool, error) {
	args := []interface{}{workspaceID, models.InactiveMember}
	argCount := 2

	query := `
		SELECT wu.user_id, wu.workspace_role, wu.membership_status
		FROM workspace_users wu
		WHERE wu.workspace_id = $1
		AND wu.membership_status != $2`

	if workspaceRole != nil {
		argCount++
		query += fmt.Sprintf(" AND wu.workspace_role = $%d", argCount)
		args = append(args, *workspaceRole)
	}

	query += " ORDER BY wu.created_at DESC"

	if limit != nil {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *limit)
	}

	if offset != nil {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, *offset)
	}

	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to list workspace members: %w", err)
	}
	defer rows.Close()

	var members []models.PartialWorkspaceUser
	for rows.Next() {
		var member models.PartialWorkspaceUser
		err := rows.Scan(
			&member.ID,
			&member.Role,
			&member.MembershipStatus,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan workspace member: %w", err)
		}
		members = append(members, member)
	}

	var hasMore bool
	if limit != nil && len(members) == *limit {
		hasMore = true
	} else {
		hasMore = false
	}

	return members, hasMore, rows.Err()
}

func (r *workspaceRepo) GetWorkspaceMemberByUserID(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error) {
	member := &models.PartialWorkspaceUser{}
	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT user_id, workspace_role, membership_status
		FROM workspace_users
		WHERE workspace_id = $1 AND user_id = $2 AND membership_status != $3`,
		workspaceID, userID, models.InactiveMember,
	).Scan(
		&member.ID,
		&member.Role,
		&member.MembershipStatus,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("workspace member not found")
		}
		return nil, fmt.Errorf("failed to get workspace member: %w", err)
	}

	return member, nil
}

func (r *workspaceRepo) GetWorkspaceUserCountsByRoleAndStatus(ctx context.Context, workspaceID uuid.UUID) (activeUsers, pendingUsers, activeAdmins, pendingAdmins int, err error) {
	rows, err := r.getQuerier(ctx).Query(ctx, `
        SELECT
            workspace_role,
            membership_status,
            COUNT(*)
        FROM workspace_users
        WHERE workspace_id = $1
            AND membership_status IN ($2, $3)
        GROUP BY workspace_role, membership_status`,
		workspaceID, models.ActiveMember, models.PendingMember,
	)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get workspace user count: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var role models.WorkspaceRole
		var status models.MembershipStatus
		var count int

		if err := rows.Scan(&role, &status, &count); err != nil {
			return 0, 0, 0, 0, fmt.Errorf("failed to scan role count: %w", err)
		}

		switch {
		case role == models.RoleUser && status == models.ActiveMember:
			activeUsers = count
		case role == models.RoleUser && status == models.PendingMember:
			pendingUsers = count
		case role == models.RoleAdmin && status == models.ActiveMember:
			activeAdmins = count
		case role == models.RoleAdmin && status == models.PendingMember:
			pendingAdmins = count
		}
	}

	if err = rows.Err(); err != nil {
		return 0, 0, 0, 0, err
	}

	return activeUsers, pendingUsers, activeAdmins, pendingAdmins, nil
}

func (r *workspaceRepo) PromoteRandomUserToAdmin(ctx context.Context, workspaceID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
        UPDATE workspace_users
        SET workspace_role = $1
        WHERE user_id = (
            SELECT user_id
            FROM workspace_users
            WHERE workspace_id = $2
                AND workspace_role = $3
                AND membership_status = $4
            ORDER BY RANDOM()
            LIMIT 1
        )`,
		models.RoleAdmin,    // $1: new role
		workspaceID,         // $2: workspace ID
		models.RoleUser,     // $3: current role
		models.ActiveMember, // $4: membership status
	)
	if err != nil {
		return fmt.Errorf("failed to promote random user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("no eligible users found to promote")
	}

	return nil
}

func (r *workspaceRepo) GetWorkspacesForUserID(ctx context.Context, userID uuid.UUID, membershipStatus models.MembershipStatus) ([]models.Workspace, error) {
	rows, err := r.getQuerier(ctx).Query(ctx, `
        SELECT w.id, w.name, w.slug, w.billing_email, w.workspace_status, w.workspace_plan, w.created_at, w.updated_at
        FROM workspaces w
        INNER JOIN workspace_users wu ON w.id = wu.workspace_id
        WHERE wu.user_id = $1
        AND wu.membership_status = $2
        AND w.workspace_status != $3`,
		userID, membershipStatus, models.WorkspaceInactive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces for user: %w", err)
	}
	defer rows.Close()

	workspaces := make([]models.Workspace, 0)
	for rows.Next() {
		var workspace models.Workspace
		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Slug,
			&workspace.BillingEmail,
			&workspace.WorkspaceStatus,
			&workspace.WorkspacePlan,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, rows.Err()
}

func (r *workspaceRepo) BatchAddUsersToWorkspace(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID) ([]models.PartialWorkspaceUser, error) {
	if len(userIDs) == 0 {
		return []models.PartialWorkspaceUser{}, nil
	}

	// Create values string for bulk insert
	valueStrings := make([]string, 0, len(userIDs))
	valueArgs := make([]interface{}, 0, len(userIDs)*4)
	for i, userID := range userIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, workspaceID, userID, models.RoleUser, models.PendingMember)
	}

	query := fmt.Sprintf(`
		INSERT INTO workspace_users (workspace_id, user_id, workspace_role, membership_status)
		VALUES %s
		RETURNING user_id, workspace_role, membership_status`,
		strings.Join(valueStrings, ","))

	rows, err := r.getQuerier(ctx).Query(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to batch add users: %w", err)
	}
	defer rows.Close()

	var members []models.PartialWorkspaceUser
	for rows.Next() {
		var member models.PartialWorkspaceUser
		err := rows.Scan(
			&member.ID,
			&member.Role,
			&member.MembershipStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace member: %w", err)
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func (r *workspaceRepo) AddUserToWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error) {
	member := &models.PartialWorkspaceUser{}
	err := r.getQuerier(ctx).QueryRow(ctx, `
		INSERT INTO workspace_users (workspace_id, user_id, workspace_role, membership_status)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, workspace_role, membership_status`,
		workspaceID, userID, models.RoleUser, models.PendingMember,
	).Scan(
		&member.ID,
		&member.Role,
		&member.MembershipStatus,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add user to workspace: %w", err)
	}

	return member, nil
}

func (r *workspaceRepo) BatchRemoveUsersFromWorkspace(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	_, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET membership_status = $1
		WHERE workspace_id = $2 AND user_id = ANY($3)`,
		models.InactiveMember, workspaceID, userIDs,
	)

	if err != nil {
		return fmt.Errorf("failed to batch remove users: %w", err)
	}

	return nil
}

func (r *workspaceRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET membership_status = $1
		WHERE workspace_id = $2 AND user_id = $3`,
		models.InactiveMember, workspaceID, userID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove user from workspace: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found in workspace")
	}

	return nil
}

func (r *workspaceRepo) UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.WorkspaceRole) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET workspace_role = $1
		WHERE workspace_id = $2 AND user_id = $3 AND membership_status != $4`,
		role, workspaceID, userID, models.InactiveMember,
	)

	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found in workspace")
	}

	return nil
}

func (r *workspaceRepo) UpdateUserMembershipStatusForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, status models.MembershipStatus) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET membership_status = $1
		WHERE workspace_id = $2 AND user_id = $3 AND membership_status != $4`,
		status, workspaceID, userID, models.InactiveMember,
	)

	if err != nil {
		return fmt.Errorf("failed to update membership status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found in workspace")
	}

	return nil
}

func (r *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspaces
		SET billing_email = $1
		WHERE id = $2 AND workspace_status != $3`,
		billingEmail, workspaceID, models.WorkspaceInactive,
	)

	if err != nil {
		return fmt.Errorf("failed to update workspace billing email: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error {
	workspaceSlug := getSlug(workspaceName)
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2
		WHERE id = $3 AND workspace_status != $4`,
		workspaceName, workspaceSlug, workspaceID, models.WorkspaceInactive,
	)

	if err != nil {
		return fmt.Errorf("failed to update workspace name: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *workspaceRepo) UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error {
	workspaceSlug := getSlug(workspaceName)
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2, billing_email = $3
		WHERE id = $4 AND workspace_status != $5`,
		workspaceName, workspaceSlug, billingEmail, workspaceID, models.WorkspaceInactive,
	)

	if err != nil {
		return fmt.Errorf("failed to update workspace details: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *workspaceRepo) UpdateWorkspacePlan(ctx context.Context, workspaceID uuid.UUID, plan models.WorkspacePlan) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
        UPDATE workspaces
        SET workspace_plan = $1
        WHERE id = $2 AND workspace_status != $3`,
		plan, workspaceID, models.WorkspaceInactive,
	)

	if err != nil {
		return fmt.Errorf("failed to update workspace plan: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *workspaceRepo) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	// Start by soft deleting all workspace memberships
	_, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET membership_status = $1
		WHERE workspace_id = $2 AND membership_status != $1`,
		models.InactiveMember, workspaceID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete workspace memberships: %w", err)
	}

	// Then soft delete the workspace itself
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspaces
		SET workspace_status = $1
		WHERE id = $2 AND workspace_status != $1`,
		models.WorkspaceInactive, workspaceID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *workspaceRepo) BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error {
	if len(workspaceIDs) == 0 {
		return nil
	}

	// Start by soft deleting all workspace memberships
	_, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspace_users
		SET membership_status = $1
		WHERE workspace_id = ANY($2) AND membership_status != $1`,
		models.InactiveMember, workspaceIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to batch delete workspace memberships: %w", err)
	}

	// Then soft delete the workspaces
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE workspaces
		SET workspace_status = $1
		WHERE id = ANY($2) AND workspace_status != $1`,
		models.WorkspaceInactive, workspaceIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to batch delete workspaces: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspaces not found")
	}

	return nil
}

func (r *workspaceRepo) ListActiveWorkspaces(ctx context.Context, batchSize int, lastWorkspaceID *uuid.UUID) (models.ActiveWorkspaceBatch, error) {
	if batchSize <= 0 {
		return models.ActiveWorkspaceBatch{}, errors.New("invalid batch size")
	}

	// Build base query
	query := `
		SELECT id
		FROM workspaces
		WHERE workspace_status = $1`
	args := []interface{}{models.WorkspaceActive}

	// Add cursor-based pagination using lastWorkspaceID
	if lastWorkspaceID != nil {
		// Ensure UUID ordering is appropriate for pagination
		query += ` AND id > $2`
		args = append(args, *lastWorkspaceID)
	}

	// Order by ID for consistent pagination
	query += ` ORDER BY id ASC`

	// Add limit
	query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
	args = append(args, batchSize+1) // Request one extra to determine if there are more workspaces

	// Execute query
	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return models.ActiveWorkspaceBatch{}, fmt.Errorf("failed to query active workspaces: %w", err)
	}
	if rows == nil {
		return models.ActiveWorkspaceBatch{}, errors.New("query returned nil rows")
	}
	defer rows.Close()

	// Collect results
	workspaceIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		workspaceID := uuid.UUID{}
		if err := rows.Scan(&workspaceID); err != nil {
			return models.ActiveWorkspaceBatch{}, fmt.Errorf("failed to scan workspace ID: %w", err)
		}
		workspaceIDs = append(workspaceIDs, workspaceID)
	}

	if err = rows.Err(); err != nil {
		return models.ActiveWorkspaceBatch{}, fmt.Errorf("error iterating workspaces: %w", err)
	}

	// Determine if there are more workspaces
	hasMore := len(workspaceIDs) > batchSize
	if hasMore {
		workspaceIDs = workspaceIDs[:batchSize] // Remove the extra item we requested
	}

	// Set the last seen ID
	var lastSeen *uuid.UUID
	if len(workspaceIDs) > 0 {
		lastSeen = &workspaceIDs[len(workspaceIDs)-1]
	}

	return models.ActiveWorkspaceBatch{
		WorkspaceIDs: workspaceIDs,
		HasMore:      hasMore,
		LastSeen:     lastSeen,
	}, nil
}

// Get total workspace count for a user
func (r *workspaceRepo) GetWorkspaceCountForUser(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int

	err := r.getQuerier(ctx).QueryRow(ctx, `
        SELECT COUNT(DISTINCT w.id)
        FROM workspaces w
        INNER JOIN workspace_users wu ON w.id = wu.workspace_id
        WHERE wu.user_id = $1
        AND wu.membership_status != $2
        AND w.workspace_status != $3`,
		userID, models.InactiveMember, models.WorkspaceInactive,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get workspace count: %w", err)
	}

	return count, nil
}

// Get total active + pending members for a workspace
func (r *workspaceRepo) GetActivePendingMemberCounts(ctx context.Context, workspaceID uuid.UUID) (activeCount, pendingCount int, err error) {
	rows, err := r.getQuerier(ctx).Query(ctx, `
        SELECT membership_status, COUNT(*)
        FROM workspace_users
        WHERE workspace_id = $1
        AND membership_status IN ($2, $3)
        GROUP BY membership_status`,
		workspaceID, models.ActiveMember, models.PendingMember,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get member counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status models.MembershipStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return 0, 0, fmt.Errorf("failed to scan member count: %w", err)
		}

		switch status {
		case models.ActiveMember:
			activeCount = count
		case models.PendingMember:
			pendingCount = count
		}
	}

	if err = rows.Err(); err != nil {
		return 0, 0, err
	}

	return activeCount, pendingCount, nil
}

func (r *workspaceRepo) getQuerier(ctx context.Context) interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func getSlug(name string) string {
	// remove non alphanumeric characters
	return slug.Make(name) + "-" + uuid.New().String()
}
