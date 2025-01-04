package db

import (
	"context"
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
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrInvalidWorkspace  = errors.New("invalid workspace data")
	ErrDuplicateSlug     = errors.New("workspace slug already exists")
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
	baseSlug := generateSlug(workspaceName)
	slug := baseSlug

	// Handle potential slug collisions
	for i := 1; ; i++ {
		var exists bool
		err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM workspaces WHERE slug = $1)", slug).Scan(&exists)
		if err != nil {
			return models.Workspace{}, fmt.Errorf("failed to check slug existence: %w", err)
		}
		if !exists {
			break
		}
		// If slug exists, append a number and try again
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

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
		return models.Workspace{}, fmt.Errorf("failed to create workspace: %w", err)
	}

	return workspace, nil
}

// GetWorkspaces gets multiple workspaces by their IDs
func (r *workspaceRepo) GetWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, []error) {
	if len(workspaceIDs) == 0 {
		return nil, []error{errors.New("no workspace IDs provided")}
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
		return nil, []error{fmt.Errorf("failed to get workspaces: %w", err)}
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
			errs = append(errs, fmt.Errorf("failed to scan workspace: %w", err))
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	if err = rows.Err(); err != nil {
		errs = append(errs, fmt.Errorf("error iterating over rows: %w", err))
	}

	if len(workspaces) == 0 {
		return nil, []error{ErrWorkspaceNotFound}
	}

	if len(errs) > 0 {
		return workspaces, errs
	}

	return workspaces, nil
}

// WorkspaceExists checks if a workspace exists
func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	runner := r.tm.GetRunner(ctx)

	var exists bool
	err := runner.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1)",
		workspaceID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
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
		return fmt.Errorf("failed to update workspace billing email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
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
	baseSlug := generateSlug(workspaceName)
	slug := baseSlug

	for i := 1; ; i++ {
		var exists bool
		err := runner.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM workspaces WHERE slug = $1 AND id != $2)",
			slug, workspaceID,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check slug existence: %w", err)
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, workspaceName, slug, workspaceID)

	if err != nil {
		return fmt.Errorf("failed to update workspace name: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
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
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
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
	baseSlug := generateSlug(workspaceReq.Name)
	slug := baseSlug

	for i := 1; ; i++ {
		var exists bool
		err := runner.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM workspaces WHERE slug = $1 AND id != $2)",
			slug, workspaceID,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check slug existence: %w", err)
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	result, err := runner.ExecContext(ctx, `
		UPDATE workspaces
		SET name = $1,
			slug = $2,
			billing_email = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, workspaceReq.Name, slug, workspaceReq.BillingEmail, workspaceID)

	if err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}
