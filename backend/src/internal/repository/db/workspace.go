// ./src/internal/repository/db/workspace.go
package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
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
	querier := r.tm.GetQuerier(ctx)

	// Generate initial slug
	slug := generateSlug(workspaceName)

	query := `
		INSERT INTO workspaces (name, slug, billing_email, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, billing_email, status, created_at, updated_at
	`

	var workspace models.Workspace
	err := querier.QueryRow(ctx, query,
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
		return models.Workspace{}, wErr.Propagate(repo.ErrFailedToCreateWorkspaceInWorkspaceRepository)
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
		return []models.Workspace{}, wErr.Propagate(repo.ErrFailedToGetWorkspaceFromWorkspaceRepository)
	}

	querier := r.tm.GetQuerier(ctx)

	query := `
        SELECT id, name, slug, billing_email, status, created_at, updated_at
        FROM workspaces
        WHERE id = ANY($1)
        AND status = $2
    `

	rows, err := querier.Query(ctx, query, workspaceIDs, models.WorkspaceStatusActive)
	if err != nil {
		wErr.Add(err, map[string]any{
			"query": query,
		})
		return nil, wErr.Propagate(repo.ErrFailedToGetWorkspaceFromWorkspaceRepository)
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
		return nil, wErr.Propagate(repo.ErrFailedToGetWorkspaceFromWorkspaceRepository)
	}

	if len(workspaces) == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return nil, wErr.Propagate(repo.ErrFailedToGetWorkspaceFromWorkspaceRepository)
	}

	return workspaces, nil
}

// WorkspaceExists checks if a workspace exists and is active
func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, errs.Error) {
	wErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	var exists bool
	err := querier.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1 AND status = 'active')",
		workspaceID,
	).Scan(&exists)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return false, wErr.Propagate(repo.ErrFailedToCheckIfWorkspaceExistsInWorkspaceRepository)
	}

	return exists, nil
}

// UpdateWorkspaceBillingEmail updates the billing email
func (r *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) errs.Error {
	wErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	commandTag, err := querier.Exec(ctx, `
		UPDATE workspaces
		SET billing_email = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, billingEmail, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceBillingEmailInWorkspaceRepository)
	}

	if commandTag.RowsAffected() == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceBillingEmailInWorkspaceRepository)
	}

	return nil
}

// UpdateWorkspaceName updates the workspace name
func (r *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) errs.Error {
	wErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	// Generate and verify unique slug
	slug := generateSlug(workspaceName)

	commandTag, err := querier.Exec(ctx, `
		UPDATE workspaces
		SET name = $1, slug = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, workspaceName, slug, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceNameInWorkspaceRepository)
	}

	if commandTag.RowsAffected() == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceNameInWorkspaceRepository)
	}

	return nil
}

// UpdateWorkspaceStatus updates the workspace status
func (r *workspaceRepo) UpdateWorkspaceStatus(ctx context.Context, workspaceID uuid.UUID, status models.WorkspaceStatus) errs.Error {
	wErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	commandTag, err := querier.Exec(ctx, `
		UPDATE workspaces
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, workspaceID)

	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceStatusInWorkspaceRepository)
	}

	if commandTag.RowsAffected() == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceStatusInWorkspaceRepository)
	}

	return nil
}

// UpdateWorkspace updates the workspace details
func (r *workspaceRepo) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceReq models.WorkspaceProps) errs.Error {
	wErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	// Generate and verify unique slug
	slug := generateSlug(workspaceReq.Name)

	commandTag, err := querier.Exec(ctx, `
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
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceInWorkspaceRepository)
	}

	if commandTag.RowsAffected() == 0 {
		wErr.Add(repo.ErrWorkspaceNotFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return wErr.Propagate(repo.ErrFailedToUpdateWorkspaceInWorkspaceRepository)
	}

	return nil
}

// RemoveWorkspaces removes workspaces by setting their status to inactive
func (r *workspaceRepo) RemoveWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) errs.Error {
	wErr := errs.New()
	if len(workspaceIDs) == 0 {
		wErr.Add(repo.ErrNoWorkspaceSpecified, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return wErr.Propagate(repo.ErrFailedToRemoveWorkspacesInWorkspaceRepository)
	}

	querier := r.tm.GetQuerier(ctx)

	query := `
        UPDATE workspaces
        SET status = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ANY($2)
        AND status = $3
    `

	_, err := querier.Exec(ctx, query,
		models.WorkspaceStatusInactive,
		workspaceIDs,
		models.WorkspaceStatusActive,
	)
	if err != nil {
		wErr.Add(err, map[string]any{
			"workspaceIDs": workspaceIDs,
		})
		return wErr.Propagate(repo.ErrFailedToRemoveWorkspacesInWorkspaceRepository)
	}

	return nil
}
