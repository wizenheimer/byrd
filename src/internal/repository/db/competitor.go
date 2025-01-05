package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	// ---- validation errors ----
	ErrCompetitorNamesListEmpty = errors.New("competitor names list is empty")
	ErrInvalidLimit             = errors.New("invalid limit value for list operation")
	ErrInvalidOffset            = errors.New("invalid offset value for list operation")

	// ---- non fatal errors -----
	ErrFailedToScanCompetitors          = errors.New("failed to scan competitors")
	ErrFailedToConfirmCompetitorRemoval = errors.New("failed to confirm competitor removal")

	// --- remapped errors ----
	// remapping an existing error
	// remapping non error scenario to an error
	ErrCompetitorNotFound          = errors.New("competitor not found")
	ErrNoWorkspaceCompetitorsFound = errors.New("no competitors found for the workspace")
)

type competitorRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewCompetitorRepository(tm *transaction.TxManager, logger *logger.Logger) repo.CompetitorRepository {
	return &competitorRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "competitor_repository"}),
	}
}

// CreateCompetitors creates multiple competitors in a workspace
func (r *competitorRepo) CreateCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, []error) {
	if competitorNames == nil {
		return nil, []error{
			ErrCompetitorNamesListEmpty,
		}
	}

	// Get the runner (either tx from context or db)
	runner := r.tm.GetRunner(ctx)

	// Prepare the batch insert query
	valueStrings := make([]string, len(competitorNames))
	valueArgs := make([]interface{}, 0, len(competitorNames)*2)
	now := time.Now()
	for i, name := range competitorNames {
		valueStrings[i] = fmt.Sprintf("($%d, $%d, 'active', $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		valueArgs = append(valueArgs, workspaceID, name, now, now)
	}

	query := fmt.Sprintf(`
		INSERT INTO competitors (workspace_id, name, status, created_at, updated_at)
		VALUES %s
		RETURNING id, workspace_id, name, status, created_at, updated_at
	`, strings.Join(valueStrings, ","))

	// Execute the batch insert
	rows, err := runner.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, []error{
			err,
		}
	}
	defer rows.Close()

	// Collect the created competitors
	// Or collect the errors if any
	errs := make([]error, 0)
	competitors := make([]models.Competitor, 0, len(competitorNames))
	for rows.Next() {
		var competitor models.Competitor
		err := rows.Scan(
			&competitor.ID,
			&competitor.WorkspaceID,
			&competitor.Name,
			&competitor.Status,
			&competitor.CreatedAt,
			&competitor.UpdatedAt,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		competitors = append(competitors, competitor)
	}

	// In case of any errors, return them
	// along with the competitors created so far
	if len(errs) > 0 {
		return competitors, append([]error{ErrFailedToScanCompetitors}, errs...)
	}

	if err = rows.Err(); err != nil {
		return nil, []error{
			err,
		}
	}

	return competitors, nil
}

// GetCompetitor gets a competitor by its ID
func (r *competitorRepo) GetCompetitor(ctx context.Context, competitorID uuid.UUID) (models.Competitor, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE id = $1
	`

	var competitor models.Competitor
	err := runner.QueryRowContext(ctx, query, competitorID).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	// Remap sql.ErrNoRows to ErrCompetitorNotFound
	if err == sql.ErrNoRows {
		return models.Competitor{}, ErrCompetitorNotFound
	}

	// For any other error, return it as is
	if err != nil {
		return models.Competitor{}, err
	}

	return competitor, nil
}

// ListWorkspaceCompetitors lists all competitors in a workspace with pagination
func (r *competitorRepo) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]models.Competitor, []error) {
	errs := make([]error, 0)
	if limit < 1 {
		errs = append(errs, ErrInvalidLimit)
	}
	if offset < 0 {
		errs = append(errs, ErrInvalidOffset)
	}

	if len(errs) > 0 {
		// Return the errors if any
		return nil, errs
	}

	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := runner.QueryContext(ctx, query, workspaceID, limit, offset)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	competitors := make([]models.Competitor, 0)
	for rows.Next() {
		var competitor models.Competitor
		err := rows.Scan(
			&competitor.ID,
			&competitor.WorkspaceID,
			&competitor.Name,
			&competitor.Status,
			&competitor.CreatedAt,
			&competitor.UpdatedAt,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		competitors = append(competitors, competitor)
	}

	// In case of any errors, return them
	// along with the competitors created so far
	if len(errs) > 0 {
		return competitors, append([]error{ErrFailedToScanCompetitors}, errs...)
	}

	// Check if there are any errors while iterating over the rows
	// and return them if any along with the competitors created so far
	if err = rows.Err(); err != nil {
		return competitors, []error{err}
	}

	// If no competitors found in the workspace, return an error
	if len(competitors) == 0 {
		return nil, []error{ErrNoWorkspaceCompetitorsFound}
	}

	return competitors, nil
}

// RemoveWorkspaceCompetitors removes competitors from a workspace
func (r *competitorRepo) RemoveWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) []error {
	runner := r.tm.GetRunner(ctx)

	var query string
	var args []interface{}
	if competitorIDs == nil {
		// Remove all competitors from workspace
		query = `
			DELETE FROM competitors
			WHERE workspace_id = $1
		`
		args = []interface{}{workspaceID}
	} else {
		// Remove specific competitors
		placeholders := make([]string, len(competitorIDs))
		args = make([]interface{}, len(competitorIDs)+1)
		args[0] = workspaceID

		for i, id := range competitorIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args[i+1] = id
		}

		query = fmt.Sprintf(`
			DELETE FROM competitors
			WHERE workspace_id = $1
			AND id IN (%s)
		`, strings.Join(placeholders, ","))
	}

	result, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		return []error{ErrFailedToConfirmCompetitorRemoval, err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return []error{err}
	}

	if rowsAffected == 0 {
		return []error{ErrNoWorkspaceCompetitorsFound}
	}

	return nil
}

// WorkspaceCompetitorExists checks if a competitor exists in a workspace
func (r *competitorRepo) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	runner := r.tm.GetRunner(ctx)
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM competitors
			WHERE workspace_id = $1 AND id = $2
		)
	`

	var exists bool
	err := runner.QueryRowContext(ctx, query, workspaceID, competitorID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
