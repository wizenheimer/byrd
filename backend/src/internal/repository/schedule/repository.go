package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type scheduleRepo struct {
	logger *logger.Logger
	tm     *transaction.TxManager
}

func NewScheduleRepo(tm *transaction.TxManager, logger *logger.Logger) ScheduleRepository {
	return &scheduleRepo{
		logger: logger,
		tm:     tm,
	}
}

func (r *scheduleRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

// CreateSchedule schedules a new workflow in the repository
func (r *scheduleRepo) CreateScheduleWithID(ctx context.Context, id models.ScheduleID, workflowProps models.WorkflowScheduleProps) (models.ScheduleID, error) {
	q := r.getQuerier(ctx)

	sql := `
        INSERT INTO workflow_schedules (
            id, workflow_type, about, spec, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, NOW(), NOW()
        )
        RETURNING id`

	err := q.QueryRow(ctx, sql,
		id,
		workflowProps.WorkflowType,
		workflowProps.About,
		workflowProps.Spec,
	).Scan(&id)

	if err != nil {
		return models.NilScheduleID(), fmt.Errorf("failed to create schedule: %w", err)
	}

	return id, nil
}

// CreateSchedule schedules a new workflow in the repository
func (r *scheduleRepo) CreateSchedule(ctx context.Context, workflowProps models.WorkflowScheduleProps) (models.ScheduleID, error) {
	id := models.NewScheduleID()
	return r.CreateScheduleWithID(ctx, id, workflowProps)
}

// GetSchedule returns the schedule of a workflow
func (r *scheduleRepo) GetSchedule(ctx context.Context, scheduleID models.ScheduleID) (models.WorkflowSchedule, error) {
	q := r.getQuerier(ctx)
	var schedule models.WorkflowSchedule

	sql := `
        SELECT
            id, workflow_type, about, spec,
            last_run, next_run, created_at, updated_at
        FROM workflow_schedules
        WHERE id = $1 AND deleted_at IS NULL`

	err := q.QueryRow(ctx, sql, scheduleID).Scan(
		&schedule.ID,
		&schedule.WorkflowType,
		&schedule.About,
		&schedule.Spec,
		&schedule.LastRun, // sql.NullTime will handle NULL properly
		&schedule.NextRun, // sql.NullTime will handle NULL properly
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.WorkflowSchedule{}, errors.New("schedule not found")
		}
		return models.WorkflowSchedule{}, fmt.Errorf("failed to get schedule: %w", err)
	}

	return schedule, nil
}

// UpdateSchedule updates the schedule of a workflow in the repository
func (r *scheduleRepo) UpdateSchedule(ctx context.Context, scheduleID models.ScheduleID, workflowProps models.WorkflowScheduleProps) error {
	q := r.getQuerier(ctx)

	sql := `
        UPDATE workflow_schedules
        SET
            workflow_type = $1,
            about = $2,
            spec = $3,
            updated_at = NOW()
        WHERE id = $4 AND deleted_at IS NULL`

	result, err := q.Exec(ctx, sql,
		workflowProps.WorkflowType,
		workflowProps.About,
		workflowProps.Spec,
		scheduleID,
	)

	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

// DeleteSchedule deletes the scheduled workflows
func (r *scheduleRepo) DeleteSchedule(ctx context.Context, scheduleID models.ScheduleID) error {
	q := r.getQuerier(ctx)

	sql := `
        UPDATE workflow_schedules
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL`

	result, err := q.Exec(ctx, sql, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

// ListScheduledWorkflows returns the list of scheduled workflows
func (r *scheduleRepo) ListScheduledWorkflows(ctx context.Context, limit, offset *int, workflowType *models.WorkflowType) ([]models.WorkflowSchedule, error) {
	q := r.getQuerier(ctx)
	var args []interface{}
	argPosition := 1

	sql := `
        SELECT
            id, workflow_type, about, spec,
            last_run, next_run, created_at, updated_at
        FROM workflow_schedules
        WHERE deleted_at IS NULL`

	if workflowType != nil {
		sql += fmt.Sprintf(" AND workflow_type = $%d", argPosition)
		args = append(args, *workflowType)
		argPosition++
	}

	sql += " ORDER BY created_at DESC"

	if limit != nil {
		sql += fmt.Sprintf(" LIMIT $%d", argPosition)
		args = append(args, *limit)
		argPosition++
	}

	if offset != nil {
		sql += fmt.Sprintf(" OFFSET $%d", argPosition)
		args = append(args, *offset)
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.WorkflowSchedule
	for rows.Next() {
		var schedule models.WorkflowSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.WorkflowType,
			&schedule.About,
			&schedule.Spec,
			&schedule.LastRun, // sql.NullTime will handle NULL properly
			&schedule.NextRun, // sql.NullTime will handle NULL properly
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating schedules: %w", err)
	}

	if len(schedules) == 0 {
		schedules = make([]models.WorkflowSchedule, 0)
	}

	return schedules, nil
}

// Sync
func (r *scheduleRepo) Sync(ctx context.Context, scheduleID models.ScheduleID, lastRun, nextRun time.Time) error {
	q := r.getQuerier(ctx)

	query := `
        UPDATE workflow_schedules
        SET
            last_run = $1,
            next_run = $2,
            updated_at = NOW()
        WHERE id = $3 AND deleted_at IS NULL`

	// Convert time.Time to sql.NullTime
	lastRunNull := sql.NullTime{Time: lastRun, Valid: !lastRun.IsZero()}
	nextRunNull := sql.NullTime{Time: nextRun, Valid: !nextRun.IsZero()}

	result, err := q.Exec(ctx, query, lastRunNull, nextRunNull, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to sync schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}
