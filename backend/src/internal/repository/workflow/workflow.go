package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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

type workflowRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewWorkflowRepository(client *redis.Client, logger *logger.Logger) (WorkflowRepository, error) {
	workflowRepo := workflowRepository{
		client: client,
		logger: logger.WithFields(
            map[string]interface{}{
                "module": "workflow_repository",
            },
        ),
	}
	if err := validateClient(context.Background(), client); err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	return &workflowRepo, nil
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

func (r *workflowRepository) getKey(jobID uuid.UUID, workflowType models.WorkflowType, status models.JobStatus) string {
	return fmt.Sprintf(keyFormat, workflowType, status, jobID.String())
}

func (r *workflowRepository) GetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) (models.JobState, error) {
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

func (r *workflowRepository) SetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobState models.JobState) error {
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

func (r *workflowRepository) SetCheckpoint(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobCheckpoint models.JobCheckpoint) error {
	r.logger.Debug("setting checkpoint for job",
		zap.Any("jobID", jobID),
		zap.Any("workflowType", workflowType))

	// Get current state
	state, err := r.GetState(ctx, jobID, workflowType)
	if err != nil {
		return fmt.Errorf("failed to get job state: %w", err)
	}

	// Update checkpoint
	state.Checkpoint = jobCheckpoint
	return r.SetState(ctx, jobID, workflowType, state)
}

func (r *workflowRepository) List(ctx context.Context, workflowType models.WorkflowType, jobStatus models.JobStatus) ([]models.Job, error) {
	r.logger.Debug("listing jobs",
		zap.Any("workflowType", workflowType),
		zap.Any("jobStatus", jobStatus))

	pattern := fmt.Sprintf(keyPattern, workflowType, jobStatus)
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

func (r *workflowRepository) Initialize(ctx context.Context, jobProps models.Job, workflowType models.WorkflowType) error {
	r.logger.Debug("initializing workflow", zap.Any("jobProps", jobProps), zap.Any("workflowType", workflowType))
	// TODO: Persist the triggered state in db
	return nil
}
