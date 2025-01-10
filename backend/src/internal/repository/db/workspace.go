package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewWorkspaceRepository(tm *transaction.TxManager, logger *logger.Logger) repo.WorkspaceRepository {
	return &workspaceRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "workspace_repository"}),
	}
}

// generateSlug creates a URL-friendly slug from workspace name
func generateSlug(name string) string {
	name += "-" + uuid.NewString()
	return slug.Make(name)
}

// CreateWorkspace creates a new workspace
func (r *workspaceRepo) CreateWorkspace(ctx context.Context, workspaceName, billingEmail string) (models.Workspace, errs.Error) {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	// Generate initial slug
	slug := generateSlug(workspaceName)

	query := `
		INSERT INTO workspaces (name, slug, billing_email, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, billing_email, status, created_at, updated_at
	`

	var workspace models.Workspace
	err := runner.QueryRowContext(ctx, query,
		workspaceName,
		slug,
		billingEmail,
		models.WorkspaceStatusActive,
	).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.BillingEmail,
		&workspace.Status,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err != nil {
		wErr.Add(err, map[string]any{
			"query": query,
		})
		return models.Workspace{}, wErr
	}

	return workspace, nil
}

// GetWorkspaces gets multiple workspaces by their IDs
func (r *workspaceRepo) GetWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, errs.Error) {
	wErr := errs.New()
	if len(workspaceIDs) == 0 {
		wErr.Add(repo.ErrNoWorkspaceSpecified, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return []models.Workspace{}, wErr
	}

	runner := r.tm.GetRunner(ctx)

	// Create placeholders for the IN clause
	placeholders := make([]string, len(workspaceIDs))
	args := make([]interface{}, len(workspaceIDs))
	for i, id := range workspaceIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT id, name, slug, billing_email, status, created_at, updated_at
        FROM workspaces
        WHERE id IN (%s)
        AND status = $%d
    `, strings.Join(placeholders, ","), len(workspaceIDs)+1)

	args = append(args, models.WorkspaceStatusActive)

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		wErr.Add(err, map[string]any{
			"query": query,
		})
		return nil, wErr
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
			&workspace.Status,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			wErr.Add(err, map[string]any{
				"query": query,
			})
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	if err = rows.Err(); err != nil {
		wErr.Add(err, map[string]any{
			"query": query,
		})
	}

	if wErr.HasErrors() {
		return nil, wErr
	}

	if len(workspaces) == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return nil, wErr
	}

	return workspaces, nil
}

// WorkspaceExists checks if a workspace exists
// And is active
func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, errs.Error) {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	var exists bool
	err := runner.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1 AND status = 'active')",
		workspaceID,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return false, wErr
	}

	return exists, nil
}

// UpdateWorkspaceBillingEmail updates the billing email
func (r *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) errs.Error {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET billing_email = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, billingEmail, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	if rowsAffected == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	return nil
}

// UpdateWorkspaceName updates the workspace name
func (r *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) errs.Error {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	// Generate and verify unique slug
	slug := generateSlug(workspaceName)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, workspaceName, slug, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	if rowsAffected == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	return nil
}

// UpdateWorkspaceStatus updates the workspace status
func (r *workspaceRepo) UpdateWorkspaceStatus(ctx context.Context, workspaceID uuid.UUID, status models.WorkspaceStatus) errs.Error {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	if rowsAffected == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	return nil
}

// UpdateWorkspace updates the workspace details
func (r *workspaceRepo) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceReq models.WorkspaceProps) errs.Error {
	wErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	// Generate and verify unique slug
	slug := generateSlug(workspaceReq.Name)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET name = $1,
			slug = $2,
			billing_email = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, workspaceReq.Name, slug, workspaceReq.BillingEmail, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	if rowsAffected == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr
	}

	return nil
}

func (r *workspaceRepo) RemoveWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) errs.Error {
	wErr := errs.New()
	if len(workspaceIDs) == 0 {
		wErr.Add(repo.ErrNoWorkspaceSpecified, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return wErr
	}

	runner := r.tm.GetRunner(ctx)

	placeholders := make([]string, len(workspaceIDs))
	args := make([]interface{}, len(workspaceIDs))
	for i, id := range workspaceIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        UPDATE workspaces
        SET status = '%s',
            updated_at = CURRENT_TIMESTAMP
        WHERE id IN (%s)
        AND status = '%s'
    `, models.WorkspaceStatusInactive, strings.Join(placeholders, ","), models.WorkspaceStatusActive)

	_, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return wErr
	}

	return nil
}
