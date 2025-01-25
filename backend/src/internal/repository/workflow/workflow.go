// ./src/internal/repository/workflow/workflow.go
package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	// Key format: workflow:{type}:{status}:{jobId}
	keyFormat = "workflow:%s:%s:%s"
	// Key pattern for listing: workflow:{type}:{status}:*
	keyPattern = "workflow:%s:%s:*"
	// Default TTL for workflow keys (3 days)
	defaultTTL = 72 * time.Hour
)

// WorkflowRepository represents the repository for managing workflows
type workflowRepo struct {
	// client is the Redis client
	client *redis.Client
	// logger is the logger
	logger *logger.Logger
	// tm is the transaction manager
	tm *transaction.TxManager
}

func (r *workflowRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func NewWorkflowRepository(client *redis.Client, tm *transaction.TxManager, logger *logger.Logger) (WorkflowRepository, error) {
	repo := workflowRepo{
		client: client,
		logger: logger.WithFields(
			map[string]interface{}{
				"module": "workflow_repository",
			},
		),
		tm: tm,
	}
	if err := validateClient(context.Background(), client); err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	return &repo, nil
}

func validateClient(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("client is required")
	}
	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to ping client: %w", err)
	}
	return nil
}

func (r *workflowRepo) getKey(jobID uuid.UUID, workflowType models.WorkflowType, status models.JobStatus) string {
	return fmt.Sprintf(keyFormat, workflowType, status, jobID.String())
}

func (r *workflowRepo) GetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) (models.JobState, error) {
	r.logger.Debug("getting state for job", zap.Any("jobID", jobID), zap.Any("workflowType", workflowType))

	// We need to check all possible statuses since we don't know the current status
	for _, status := range []models.JobStatus{
		models.JobStatusRunning,
		models.JobStatusCompleted,
		models.JobStatusFailed,
		models.JobStatusAborted,
		models.JobStatusUnknown,
	} {
		key := r.getKey(jobID, workflowType, status)
		data, err := r.client.Get(ctx, key).Bytes()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return models.JobState{}, fmt.Errorf("failed to get job state: %w", err)
		}

		var state models.JobState
		if err := json.Unmarshal(data, &state); err != nil {
			return models.JobState{}, fmt.Errorf("failed to unmarshal job state: %w", err)
		}

		return state, nil
	}

	return models.JobState{}, fmt.Errorf("job not found")
}

func (r *workflowRepo) SetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobState models.JobState) error {
	r.logger.Debug("setting state for job",
		zap.Any("jobID", jobID),
		zap.Any("workflowType", workflowType),
		zap.Any("status", jobState.Status))

	// Delete old status key if it exists
	oldState, err := r.GetState(ctx, jobID, workflowType)
	if err == nil && oldState.Status != jobState.Status {
		oldKey := r.getKey(jobID, workflowType, oldState.Status)
		if err := r.client.Del(ctx, oldKey).Err(); err != nil {
			return fmt.Errorf("failed to delete old state: %w", err)
		}
	}

	// Set new state
	data, err := json.Marshal(jobState)
	if err != nil {
		return fmt.Errorf("failed to marshal job state: %w", err)
	}

	key := r.getKey(jobID, workflowType, jobState.Status)
	if err := r.client.Set(ctx, key, data, defaultTTL).Err(); err != nil {
		return fmt.Errorf("failed to set job state: %w", err)
	}

	return nil
}

func (r *workflowRepo) ListActiveJobs(ctx context.Context, workflowType models.WorkflowType) ([]models.Job, error) {
	r.logger.Debug("listing active jobs",
		zap.Any("workflowType", workflowType))

	pattern := fmt.Sprintf(keyPattern, workflowType, models.JobStatusRunning)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var jobs []models.Job
	for _, key := range keys {
		// Extract jobID from key
		parts := strings.Split(key, ":")
		if len(parts) != 4 {
			r.logger.Warn("invalid key format", zap.String("key", key))
			continue
		}

		jobID, err := uuid.Parse(parts[3])
		if err != nil {
			r.logger.Warn("invalid job ID", zap.String("jobID", parts[3]))
			continue
		}

		// Get job state
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			r.logger.Warn("failed to get job state",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		var state models.JobState
		if err := json.Unmarshal(data, &state); err != nil {
			r.logger.Warn("failed to unmarshal job state",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		jobs = append(jobs, models.Job{
			JobID:    jobID,
			JobState: state,
		})
	}

	return jobs, nil
}

func (r *workflowRepo) StartJob(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) error {
	r.logger.Debug("starting job",
		zap.Any("jobID", jobID),
		zap.Any("workflowType", workflowType))

	// Set the state of the job to running in checkpoint repository
	if err := r.SetState(ctx, jobID, workflowType, *models.NewJobState()); err != nil {
		return fmt.Errorf("failed to set job state: %w", err)
	}

	// Start the job in the state repository
	q := r.getQuerier(ctx)

	sql := `
        INSERT INTO job_records (
            job_id, workflow_type, start_time, created_at, updated_at
        ) VALUES (
            $1, $2, NOW(), NOW(), NOW()
        )`

	_, err := q.Exec(ctx, sql, jobID, workflowType)
	if err != nil {
		return fmt.Errorf("failed to start job: %w", err)
	}

	return nil
}

func (r *workflowRepo) CompleteJob(ctx context.Context, jobID uuid.UUID, jobContext *models.JobState, workflowType models.WorkflowType) error {
	r.logger.Debug("completing job",
		zap.Any("jobID", jobID),
		zap.Any("workflowType", workflowType))

	// Set the state of the job to completed in checkpoint repository
	if jobContext == nil {
		jobContext = models.NewJobState()
		jobContext.Status = models.JobStatusCompleted
	}
	if err := r.SetState(ctx, jobID, workflowType, *jobContext); err != nil {
		return fmt.Errorf("failed to set job state: %w", err)
	}

	// Set the state of the job to completed in the state repository
	q := r.getQuerier(ctx)

	sql := `
        UPDATE job_records
        SET
            end_time = NOW(),
            updated_at = NOW()
        WHERE job_id = $1
        AND workflow_type = $2
        AND deleted_at IS NULL
        AND end_time IS NULL
        AND cancel_time IS NULL`

	result, err := q.Exec(ctx, sql, jobID, workflowType)
	if err != nil {
		return fmt.Errorf("failed to complete job: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no active job found with id %s", jobID)
	}

	return nil
}

func (r *workflowRepo) CancelJob(ctx context.Context, jobID uuid.UUID, jobContext *models.JobState, workflowType models.WorkflowType) error {
	r.logger.Debug("cancelling job",
		zap.Any("jobID", jobID),
		zap.Any("workflowType", workflowType))

	// Set the state of the job to aborted in checkpoint repository
	if jobContext == nil {
		jobContext = models.NewJobState()
		jobContext.Status = models.JobStatusAborted
	}
	if err := r.SetState(ctx, jobID, workflowType, *jobContext); err != nil {
		return fmt.Errorf("failed to set job state: %w", err)
	}

	// Cancel the job in the state repository
	q := r.getQuerier(ctx)

	sql := `
        UPDATE job_records
        SET
            cancel_time = NOW(),
            updated_at = NOW()
        WHERE job_id = $1
        AND workflow_type = $2
        AND deleted_at IS NULL
        AND end_time IS NULL
        AND cancel_time IS NULL`

	result, err := q.Exec(ctx, sql, jobID, workflowType)
	if err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no active job found with id %s", jobID)
	}

	return nil
}

func (r *workflowRepo) ListRecords(ctx context.Context, workflowType *models.WorkflowType, limit, offset *int) ([]models.JobRecord, error) {
	q := r.getQuerier(ctx)
	var args []interface{}
	argPosition := 1

	sql := `
        SELECT
            id, job_id, workflow_type,
            start_time, end_time, cancel_time,
            preemptions
        FROM job_records
        WHERE deleted_at IS NULL`

	if workflowType != nil {
		sql += fmt.Sprintf(" AND workflow_type = $%d", argPosition)
		args = append(args, *workflowType)
		argPosition++
	}

	sql += " ORDER BY start_time DESC"

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
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.JobRecord
	for rows.Next() {
		var job models.JobRecord
		err := rows.Scan(
			&job.ID,
			&job.JobID,
			&job.WorkflowType,
			&job.StartTime,
			&job.EndTime,
			&job.CancelTime,
			&job.Preemptions,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jobs: %w", err)
	}

	return jobs, nil
}

// IncrementPreemptions increments the preemption count for a job
func (r *workflowRepo) IncrementPreemptions(ctx context.Context, jobID uuid.UUID) error {
	q := r.getQuerier(ctx)

	sql := `
        UPDATE job_records
        SET
            preemptions = preemptions + 1,
            updated_at = NOW()
        WHERE job_id = $1
        AND deleted_at IS NULL
        AND end_time IS NULL
        AND cancel_time IS NULL`

	result, err := q.Exec(ctx, sql, jobID)
	if err != nil {
		return fmt.Errorf("failed to increment preemptions: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no active job found with id %s", jobID)
	}

	return nil
}
