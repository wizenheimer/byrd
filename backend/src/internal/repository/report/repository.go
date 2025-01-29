package report

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type reportRespository struct {
	logger *logger.Logger
	tm     *transaction.TxManager
}

func (r *reportRespository) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

// NewReportRepository creates a new report repository
func NewReportRepository(tm *transaction.TxManager, logger *logger.Logger) ReportRepository {
	return &reportRespository{
		logger: logger.WithFields(map[string]interface{}{"repository": "report"}),
		tm:     tm,
	}
}

// Set creates a new report
func (r *reportRespository) Set(ctx context.Context, report *models.Report) error {
	querier := r.getQuerier(ctx)

	const insertSQL = `
        INSERT INTO reports (id, workspace_id, competitor_id, changes, time)
        VALUES ($1, $2, $3, $4, $5)
    `

	changesJSON, err := json.Marshal(report.Changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}

	_, err = querier.Exec(ctx, insertSQL,
		report.ID,
		report.WorkspaceID,
		report.CompetitorID,
		changesJSON,
		report.Time,
	)
	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	return nil
}

// Get returns the report with the given ID
func (r *reportRespository) Get(ctx context.Context, reportID uuid.UUID) (*models.Report, error) {
	querier := r.getQuerier(ctx)

	const getSQL = `
        SELECT id, workspace_id, competitor_id, changes, time
        FROM reports
        WHERE id = $1
    `

	report := &models.Report{}
	var changesJSON []byte

	err := querier.QueryRow(ctx, getSQL, reportID).Scan(
		&report.ID,
		&report.WorkspaceID,
		&report.CompetitorID,
		&changesJSON,
		&report.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("report not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	err = json.Unmarshal(changesJSON, &report.Changes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal changes: %w", err)
	}

	return report, nil
}

// List returns a list of reports for the given workspace and competitor
// List returns a list of reports for the given workspace and competitor
func (r *reportRespository) List(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	querier := r.getQuerier(ctx)

	var args []interface{}
	args = append(args, workspaceID, competitorID)

	// Fetch one extra record to determine if there are more results
	limitValue := 0
	if limit != nil {
		limitValue = *limit + 1 // Fetch one extra record
	}

	query := `
        SELECT id, workspace_id, competitor_id, changes, time
        FROM reports
        WHERE workspace_id = $1 AND competitor_id = $2
        ORDER BY time DESC
    `

	if limitValue > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limitValue)
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *offset)
	}

	rows, err := querier.Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to list reports: %w", err)
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		var changesJSON []byte

		err := rows.Scan(
			&report.ID,
			&report.WorkspaceID,
			&report.CompetitorID,
			&changesJSON,
			&report.Time,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan report: %w", err)
		}

		err = json.Unmarshal(changesJSON, &report.Changes)
		if err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal changes: %w", err)
		}

		reports = append(reports, report)
	}

	hasMore := false
	if limit != nil && len(reports) > *limit {
		hasMore = true
		reports = reports[:*limit] // Remove the extra record
	}

	return reports, hasMore, nil
}

// GetLatest returns the latest report for the given workspace and competitor
func (r *reportRespository) GetLatest(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Report, error) {
	querier := r.getQuerier(ctx)

	const getLatestSQL = `
        SELECT id, workspace_id, competitor_id, changes, time
        FROM reports
        WHERE workspace_id = $1 AND competitor_id = $2
        ORDER BY time DESC
        LIMIT 1
    `

	report := &models.Report{}
	var changesJSON []byte

	err := querier.QueryRow(ctx, getLatestSQL, workspaceID, competitorID).Scan(
		&report.ID,
		&report.WorkspaceID,
		&report.CompetitorID,
		&changesJSON,
		&report.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("report not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get latest report: %w", err)
	}

	err = json.Unmarshal(changesJSON, &report.Changes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal changes: %w", err)
	}

	return report, nil
}

func (r *reportRespository) GetForPeriod(ctx context.Context, workspaceID, competitorID uuid.UUID, since time.Time) (*models.Report, error) {
	querier := r.getQuerier(ctx)

	const getSQL = `
        SELECT id, workspace_id, competitor_id, changes, time
        FROM reports
        WHERE workspace_id = $1
        AND competitor_id = $2
        AND time >= $3
        ORDER BY time DESC
        LIMIT 1
    `

	report := &models.Report{}
	var changesJSON []byte

	err := querier.QueryRow(ctx, getSQL, workspaceID, competitorID, since).Scan(
		&report.ID,
		&report.WorkspaceID,
		&report.CompetitorID,
		&changesJSON,
		&report.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get report for period: %w", err)
	}

	err = json.Unmarshal(changesJSON, &report.Changes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal changes: %w", err)
	}

	return report, nil
}
