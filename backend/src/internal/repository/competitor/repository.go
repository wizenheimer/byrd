package competitor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type competitorRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewCompetitorRepository(tm *transaction.TxManager, logger *logger.Logger) CompetitorRepository {
	return &competitorRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "competitor_repository"}),
	}
}

func (r *competitorRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func (r *competitorRepo) CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorName string) (*models.Competitor, error) {
	competitor := &models.Competitor{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		INSERT INTO competitors (workspace_id, name, status)
		VALUES ($1, $2, $3)
		RETURNING id, workspace_id, name, status, created_at, updated_at`,
		workspaceID, competitorName, models.CompetitorStatusActive,
	).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create competitor: %w", err)
	}

	return competitor, nil
}

func (r *competitorRepo) BatchCreateCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorNames []string) ([]models.Competitor, error) {
	if len(competitorNames) == 0 {
		return []models.Competitor{}, nil
	}

	// Create values string for bulk insert
	valueStrings := make([]string, 0, len(competitorNames))
	valueArgs := make([]interface{}, 0, len(competitorNames)*3)
	for i, name := range competitorNames {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)",
			i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, workspaceID, name, models.CompetitorStatusActive)
	}

	query := fmt.Sprintf(`
		INSERT INTO competitors (workspace_id, name, status)
		VALUES %s
		RETURNING id, workspace_id, name, status, created_at, updated_at`,
		strings.Join(valueStrings, ","))

	rows, err := r.getQuerier(ctx).Query(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create competitors: %w", err)
	}
	defer rows.Close()

	var competitors []models.Competitor
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

	return competitors, rows.Err()
}

func (r *competitorRepo) GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error) {
	competitor := &models.Competitor{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE workspace_id = $1 AND id = $2 AND status != $3`,
		workspaceID, competitorID, models.CompetitorStatusInactive,
	).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("competitor not found")
		}
		return nil, fmt.Errorf("failed to get competitor: %w", err)
	}

	return competitor, nil
}

func (r *competitorRepo) BatchGetCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error) {
	if len(competitorIDs) == 0 {
		return []models.Competitor{}, nil
	}

	rows, err := r.getQuerier(ctx).Query(ctx, `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE workspace_id = $1 AND id = ANY($2) AND status != $3`,
		workspaceID, competitorIDs, models.CompetitorStatusInactive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get competitors: %w", err)
	}
	defer rows.Close()

	var competitors []models.Competitor
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

	return competitors, rows.Err()
}

func (r *competitorRepo) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error) {
	query := `
		SELECT id, workspace_id, name, status, created_at, updated_at
		FROM competitors
		WHERE workspace_id = $1 AND status != $2
		ORDER BY created_at DESC`

	args := []interface{}{workspaceID, models.CompetitorStatusInactive}

	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *limit)
	}

	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *offset)
	}

	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to list competitors: %w", err)
	}
	defer rows.Close()

	var competitors []models.Competitor
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
			return nil, false, fmt.Errorf("failed to scan competitor: %w", err)
		}
		competitors = append(competitors, competitor)
	}

	hasMore := limit != nil && len(competitors) == *limit

	return competitors, hasMore, rows.Err()
}

func (r *competitorRepo) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	competitor := &models.Competitor{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
		UPDATE competitors
		SET name = $1
		WHERE workspace_id = $2 AND id = $3 AND status != $4
		RETURNING id, workspace_id, name, status, created_at, updated_at`,
		competitorName, workspaceID, competitorID, models.CompetitorStatusInactive,
	).Scan(
		&competitor.ID,
		&competitor.WorkspaceID,
		&competitor.Name,
		&competitor.Status,
		&competitor.CreatedAt,
		&competitor.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("competitor not found")
		}
		return nil, fmt.Errorf("failed to update competitor: %w", err)
	}

	return competitor, nil
}

func (r *competitorRepo) RemoveCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE competitors
		SET status = $1
		WHERE workspace_id = $2 AND id = $3 AND status != $1`,
		models.CompetitorStatusInactive, workspaceID, competitorID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove competitor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("competitor not found")
	}

	return nil
}

func (r *competitorRepo) BatchRemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error {
	if len(competitorIDs) == 0 {
		return nil
	}

	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE competitors
		SET status = $1
		WHERE workspace_id = $2 AND id = ANY($3) AND status != $1`,
		models.CompetitorStatusInactive, workspaceID, competitorIDs,
	)

	if err != nil {
		return fmt.Errorf("failed to batch remove competitors: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("competitors not found")
	}

	return nil
}

func (r *competitorRepo) RemoveAllCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE competitors
		SET status = $1
		WHERE workspace_id = $2 AND status != $1`,
		models.CompetitorStatusInactive, workspaceID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove all competitors: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("no competitors found")
	}

	return nil
}

func (r *competitorRepo) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	var exists bool
	err := r.getQuerier(ctx).QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM competitors
			WHERE workspace_id = $1 AND id = $2 AND status != $3
		)`,
		workspaceID, competitorID, models.CompetitorStatusInactive,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check competitor existence: %w", err)
	}

	return exists, nil
}
