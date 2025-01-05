package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	// ---- validation errors -----
	ErrNoWorkspaceSpecified = errors.New("no workspace specified")

	// ---- non fatal errors ----
	ErrFailedToConfirmUpdatedBillingEmail   = errors.New("failed to confirm billing email")
	ErrFailedToConfirmWorkspaceNameUpdate   = errors.New("failed to confirm workspace name update")
	ErrFailedToConfirmWorkspaceStatusUpdate = errors.New("failed to confirm workspace status update")
	ErrFailedToConfirmWorkspaceUpdate       = errors.New("failed to confirm workspace update")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrWorkspaceNotFound = errors.New("workspace not found")
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
func (r *workspaceRepo) CreateWorkspace(ctx context.Context, workspaceName, billingEmail string) (models.Workspace, error) {
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
		return models.Workspace{}, err
	}

	return workspace, nil
}

// GetWorkspaces gets multiple workspaces by their IDs
func (r *workspaceRepo) GetWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, []error) {
	if len(workspaceIDs) == 0 {
		return []models.Workspace{}, []error{ErrNoWorkspaceSpecified}
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
	`, strings.Join(placeholders, ","))

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	workspaces := make([]models.Workspace, 0)
	errs := make([]error, 0)

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
			errs = append(errs, err)
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	if err = rows.Err(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return workspaces, errs
	}

	if len(workspaces) == 0 {
		return nil, []error{ErrWorkspaceNotFound}
	}

	return workspaces, nil
}

// WorkspaceExists checks if a workspace exists
// And is active
func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
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
		return false, err
	}

	return exists, nil
}

// UpdateWorkspaceBillingEmail updates the billing email
func (r *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error {
	runner := r.tm.GetRunner(ctx)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET billing_email = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, billingEmail, workspaceID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrFailedToConfirmUpdatedBillingEmail
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

// UpdateWorkspaceName updates the workspace name
func (r *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error {
	runner := r.tm.GetRunner(ctx)

	// Generate and verify unique slug
	slug := generateSlug(workspaceName)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, workspaceName, slug, workspaceID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrFailedToConfirmWorkspaceNameUpdate
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

// UpdateWorkspaceStatus updates the workspace status
func (r *workspaceRepo) UpdateWorkspaceStatus(ctx context.Context, workspaceID uuid.UUID, status models.WorkspaceStatus) error {
	runner := r.tm.GetRunner(ctx)

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, workspaceID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrFailedToConfirmWorkspaceStatusUpdate
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

// UpdateWorkspace updates the workspace details
func (r *workspaceRepo) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceReq models.WorkspaceProps) error {
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
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrFailedToConfirmWorkspaceUpdate
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}
