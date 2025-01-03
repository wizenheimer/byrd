package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	ErrCompetitorNotFound     = errors.New("competitor not found")
	ErrNoWorkspaceCompetitors = errors.New("no competitors found in workspace")
	ErrInvalidLimit           = errors.New("invalid limit value")
	ErrInvalidOffset          = errors.New("invalid offset value")
)

type competitorRepo struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCompetitorRepository(db *sql.DB, logger *logger.Logger) repo.CompetitorRepository {
	return &competitorRepo{
		db:     db,
		logger: logger.WithFields(map[string]interface{}{"module": "competitor_repository"}),
	}
}

// CreateCompetitors creates multiple competitors in a workspace
func (r *competitorRepo) CreateCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, []error) {
	if len(competitorNames) == 0 {
		return nil, []error{errors.New("competitor names list is empty")}
	}

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	// Prepare the batch insert query
	valueStrings := make([]string, len(competitorNames))
	valueArgs := make([]interface{}, 0, len(competitorNames)*2)
	for i, name := range competitorNames {
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		valueArgs = append(valueArgs, workspaceID, name)
	}

	query := fmt.Sprintf(`
		INSERT INTO competitors (workspace_id, name, status)
		VALUES %s
		RETURNING id, workspace_id, name, status, created_at, updated_at
	`, strings.Join(valueStrings, ","))

	// Execute the batch insert
	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to create competitors: %w", err)}
	}
	defer rows.Close()

	// Collect the created competitors
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
			return nil, []error{fmt.Errorf("failed to scan competitor: %w", err)}
		}
		competitors = append(competitors, competitor)
	}

	if err = rows.Err(); err != nil {
		return nil, []error{fmt.Errorf("error iterating over rows: %w", err)}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, []error{fmt.Errorf("failed to commit transaction: %w", err)}
	}

	return competitors, nil
}

// GetCompetitor gets a competitor by its ID
func (r *competitorRepo) GetCompetitor(ctx context.Context, competitorID uuid.UUID) (models.Competitor, error) {
	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE id = $1
	`

	var competitor models.Competitor
	err := r.db.QueryRowContext(ctx, query, competitorID).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.Competitor{}, ErrCompetitorNotFound
	}
	if err != nil {
		return models.Competitor{}, fmt.Errorf("failed to get competitor: %w", err)
	}

	return competitor, nil
}

// ListWorkspaceCompetitors lists all competitors in a workspace with pagination
func (r *competitorRepo) ListWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]models.Competitor, error) {
	if limit < 1 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}

	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, workspaceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspace competitors: %w", err)
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
			return nil, fmt.Errorf("failed to scan competitor: %w", err)
		}
		competitors = append(competitors, competitor)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(competitors) == 0 {
		return nil, ErrNoWorkspaceCompetitors
	}

	return competitors, nil
}

// RemoveWorkspaceCompetitors removes competitors from a workspace
func (r *competitorRepo) RemoveWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) []error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer tx.Rollback()

	var query string
	var args []interface{}

	if len(competitorIDs) == 0 {
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

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return []error{fmt.Errorf("failed to remove competitors: %w", err)}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return []error{fmt.Errorf("failed to get affected rows: %w", err)}
	}

	if rowsAffected == 0 {
		return []error{ErrNoWorkspaceCompetitors}
	}

	if err := tx.Commit(); err != nil {
		return []error{fmt.Errorf("failed to commit transaction: %w", err)}
	}

	return nil
}

// WorkspaceCompetitorExists checks if a competitor exists in a workspace
func (r *competitorRepo) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM competitors
			WHERE workspace_id = $1 AND id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, workspaceID, competitorID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check competitor existence: %w", err)
	}

	return exists, nil
}
