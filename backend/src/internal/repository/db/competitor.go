package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
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
func (r *competitorRepo) CreateCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, errs.Error) {
	competitorErr := errs.New()
	if competitorNames == nil {
		competitorErr.Add(repo.ErrCompetitorNamesListEmpty, map[string]any{
			"competitorNames": competitorNames,
		})
		return nil, competitorErr
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
		competitorErr.Add(err, map[string]any{
			"query": query,
		})
		return nil, competitorErr
	}
	defer rows.Close()

	// Collect the created competitors
	competitors := make([]models.Competitor, 0, len(competitorNames))
	index := 0
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
			competitorErr.Add(err, map[string]any{
				"competitor": competitorNames[index],
			})
			continue
		}
		competitors = append(competitors, competitor)
		index++
	}

	if err = rows.Err(); err != nil {
		competitorErr.Add(err, map[string]any{
			"query": query,
		})
	}

	return competitors, competitorErr
}

// GetCompetitor gets a competitor by its ID
func (r *competitorRepo) GetCompetitor(ctx context.Context, competitorID uuid.UUID) (models.Competitor, errs.Error) {
	competitorErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
        FROM competitors
        WHERE id = $1
        AND status = $2
	`

	var competitor models.Competitor
	err := runner.QueryRowContext(ctx, query, competitorID, models.CompetitorStatusActive).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Remap sql.ErrNoRows to ErrCompetitorNotFound
			competitorErr.Add(repo.ErrCompetitorNotFound, map[string]any{
				"competitorID": competitorID,
			})
		} else {
			// For any other error, return it as is
			competitorErr.Add(err, map[string]any{
				"competitorID": competitorID,
			})
		}
	}

	return competitor, competitorErr
}

// ListWorkspaceCompetitors lists all competitors in a workspace with pagination
func (r *competitorRepo) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]models.Competitor, errs.Error) {
	listErr := errs.New()
	if limit < 1 {
		listErr.Add(repo.ErrInvalidLimit, map[string]any{
			"limit": limit,
		})
	}
	if offset < 0 {
		listErr.Add(repo.ErrInvalidOffset, map[string]any{
			"offset": offset,
		})
	}

	if listErr.HasErrors() {
		return nil, listErr
	}

	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
        FROM competitors
        WHERE workspace_id = $1
        AND status = $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4
	`

	rows, err := runner.QueryContext(ctx, query, workspaceID, models.CompetitorStatusActive, limit, offset)
	if err != nil {
		listErr.Add(err, map[string]any{
			"query": query,
		})
		return nil, listErr
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
			listErr.Add(err, nil)
			continue
		}
		competitors = append(competitors, competitor)
	}

	// Check if there are any errors while iterating over the rows
	// and return them if any along with the competitors created so far
	if err = rows.Err(); err != nil {
		listErr.Add(err, nil)
	}

	if listErr.HasErrors() {
		return nil, listErr
	}

	// If no competitors found in the workspace, return an error
	if len(competitors) == 0 {
		listErr.Add(repo.ErrNoWorkspaceCompetitorsFound, map[string]any{
			"workspaceID": workspaceID,
		})
		return nil, listErr
	}

	return competitors, nil
}

// RemoveWorkspaceCompetitors removes competitors from a workspace
func (r *competitorRepo) RemoveWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) errs.Error {
	remErr := errs.New()
	runner := r.tm.GetRunner(ctx)

	var query string
	var args []interface{}
	now := time.Now()

	if competitorIDs == nil {
		// Remove all competitors from workspace
		query = `
			UPDATE competitors
            SET status = $2,
                updated_at = $3
            WHERE workspace_id = $1
            AND status = $4
		`
		args = []interface{}{workspaceID, models.CompetitorStatusInactive, now, models.CompetitorStatusActive}
	} else {
		// Deactivate specific competitors
		placeholders := make([]string, len(competitorIDs))
		args = make([]interface{}, len(competitorIDs)+4)
		args[0] = workspaceID
		args[1] = models.CompetitorStatusInactive
		args[2] = now
		args[3] = models.CompetitorStatusActive

		for i, id := range competitorIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+5)
			args[i+4] = id
		}

		query = fmt.Sprintf(`
            UPDATE competitors
            SET status = $2,
                updated_at = $3
            WHERE workspace_id = $1
            AND id IN (%s)
            AND status = $4
        `, strings.Join(placeholders, ","))
	}

	_, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		remErr.Add(err, map[string]any{
			"query": query,
		})
		return remErr
	}

	return nil
}

// WorkspaceCompetitorExists checks if a competitor exists in a workspace
func (r *competitorRepo) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, errs.Error) {
	competitorErr := errs.New()
	runner := r.tm.GetRunner(ctx)
	query := `
		SELECT EXISTS(
            SELECT 1
            FROM competitors
            WHERE workspace_id = $1
            AND id = $2
            AND status = $3
        )
	`

	var exists bool
	err := runner.QueryRowContext(ctx, query, workspaceID, competitorID, models.CompetitorStatusActive).Scan(&exists)
	if err != nil {
		competitorErr.Add(err, map[string]any{
			"workspaceID":  workspaceID,
			"competitorID": competitorID,
		})
		return false, competitorErr
	}

	return exists, nil
}
