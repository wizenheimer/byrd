package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type RedisWorkflowRepository struct {
	client      *redis.Client
	logger      *logger.Logger
	workflowTTL time.Duration
}

func NewRedisWorkflowRepository(client *redis.Client, ttl time.Duration, logger *logger.Logger) (interfaces.WorkflowRepository, error) {
	repo := RedisWorkflowRepository{
		client:      client,
		logger:      logger,
		workflowTTL: ttl,
	}

	// Verify connection on startup
	if err := repo.checkConnection(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &repo, nil
}

// SetStatus implements WorkflowRepository.SetStatus
func (r *RedisWorkflowRepository) SetStatus(ctx context.Context, id *models.WorkflowIdentifier, status models.WorkflowStatus, batchID *string, stage *int) error {
	r.logger.Debug("setting workflow status", zap.String("status", string(status)), zap.Any("id", id), zap.Any("batchID", batchID), zap.Any("stage", stage))

	// Create workflow response object
	checkpoint := models.Checkpoint{
		BatchID: batchID,
		Stage:   stage,
	}

	// Marshal checkpoint to JSON
	checkpointJSON, err := json.Marshal(checkpoint)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	// Set workflow status with TTL
	prefix := id.Type.Prefix()
	key := id.Serialize(prefix, status)
	if err := r.client.Set(ctx, key, checkpointJSON, r.workflowTTL).Err(); err != nil {
		return fmt.Errorf("failed to set workflow status: %w", err)
	}

	r.logger.Debug("set workflow status", zap.String("key", key), zap.Any("value", checkpoint))
	return nil
}

// GetStatus implements WorkflowRepository.GetStatus
func (r *RedisWorkflowRepository) GetStatus(ctx context.Context, id *models.WorkflowIdentifier) (*models.WorkflowResponse, error) {
	possibleStatus := []models.WorkflowStatus{
		models.WorkflowStatusRunning,
		// models.WorkflowStatusPending, // Pending workflows are not stored in Redis
		models.WorkflowStatusCompleted,
		models.WorkflowStatusFailed,
		models.WorkflowStatusAborted,
		// models.WorkflowStatusExpired, // Expired workflows are not stored in Redis
	}

	// Iterate through possible statuses
	var checkpointJSON []byte
	var err error
	var workflowKey string
	prefix := id.Type.Prefix()
	for _, status := range possibleStatus {
		// Serialize the workflow identifier
		key := id.Serialize(prefix, status)

		// Get workflow data
		checkpointJSON, err = r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Key doesn't exist or has expired
				continue
			}
			return nil, fmt.Errorf("failed to get workflow: %w", err)
		} else {
			// Key exists
			workflowKey = key
			r.logger.Debug("found workflow", zap.String("key", key), zap.String("status", string(status)))
			break
		}
	}

	// Check if the workflow was found
	if workflowKey == "" {
		// Check if week number is greater than the current week number
		_, currentWeekNumber := time.Now().UTC().ISOWeek()
		status := models.WorkflowStatusExpired
		if id.WeekNumber > currentWeekNumber {
			status = models.WorkflowStatusPending
		}
		return &models.WorkflowResponse{
			Type:         id.Type,
			Year:         id.Year,
			WeekNumber:   id.WeekNumber,
			BucketNumber: id.BucketNumber,
			BatchID:      nil,
			Stage:        nil,
			Status:       status,
		}, nil
	}

	// Unmarshal workflow data
	checkpoint := &models.Checkpoint{}
	if err := json.Unmarshal(checkpointJSON, checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	// Create workflow response object
	id, _, status, err := models.ParseWorkflowID(workflowKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow id: %w", err)
	}

	// Prepare workflow response
	workflow := &models.WorkflowResponse{
		Type:         id.Type,
		Year:         id.Year,
		WeekNumber:   id.WeekNumber,
		BucketNumber: id.BucketNumber,
		BatchID:      checkpoint.BatchID,
		Stage:        checkpoint.Stage,
		Status:       status,
	}

	return workflow, nil
}

// ListWorkflows implements WorkflowRepository.ListWorkflows
func (r *RedisWorkflowRepository) ListWorkflows(ctx context.Context, status models.WorkflowStatus, workflowType models.WorkflowType, limit int) ([]models.WorkflowResponse, int, error) {
	var cursor uint64
	var total int
	var workflows []models.WorkflowResponse
	prefix := workflowType.Prefix()
	pattern := fmt.Sprintf("%s-%v-%s-*", prefix, status, workflowType)

	// If limit is negative, set it to 0
	// This will return all workflows
	scanLimit := int64(max(limit, 0))

	// Use scan to iterate through keys
	for {
		// Scan for next batch of keys
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, scanLimit).Result()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan workflows: %w", err)
		}

		if len(keys) > 0 {
			// Use pipeline to get multiple workflows at once
			pipe := r.client.Pipeline()
			cmds := make([]*redis.StringCmd, len(keys))
			for i, key := range keys {
				cmds[i] = pipe.Get(ctx, key)
			}

			_, err = pipe.Exec(ctx)
			if err != nil && err != redis.Nil {
				return nil, 0, fmt.Errorf("failed to get workflows: %w", err)
			}

			// Process results
			for index, cmd := range cmds {
				checkpointJSON, err := cmd.Bytes()
				if err == redis.Nil {
					// Skip expired workflows
					continue
				}
				if err != nil {
					continue // Skip workflows that can't be retrieved
				}

				checkpoint := models.Checkpoint{}
				if err := json.Unmarshal(checkpointJSON, &checkpoint); err != nil {
					continue // Skip workflows that can't be unmarshaled
				}

				// Parse workflow ID
				identifier, _, status, err := models.ParseWorkflowID(keys[index])
				if err != nil {
					r.logger.Error("failed to parse workflow ID", zap.Error(err), zap.Any("key", keys[index]), zap.Any("value", checkpointJSON))
					continue // Skip workflows that can't be parsed
				}

				// Create workflow response object
				workflow := models.WorkflowResponse{
					Type:         identifier.Type,
					Year:         identifier.Year,
					WeekNumber:   identifier.WeekNumber,
					BucketNumber: identifier.BucketNumber,
					BatchID:      checkpoint.BatchID,
					Stage:        checkpoint.Stage,
					Status:       status,
				}

				workflows = append(workflows, workflow)
				total++

				// Check if we've reached the limit
				if limit > 0 && len(workflows) >= limit {
					return workflows[:limit], total, nil
				}
			}
		}

		// Exit if we've scanned all keys
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return workflows, total, nil
}

// StopWorkflow implements WorkflowRepository.StopWorkflow
func (r *RedisWorkflowRepository) StopWorkflow(ctx context.Context, id *models.WorkflowIdentifier) error {
	r.logger.Debug("stopping workflow", zap.Any("id", id))

	// Removing the workflow key will effectively stop the workflow
	prefix := id.Type.Prefix()
	key := id.Serialize(prefix, models.WorkflowStatusRunning)
	err := r.client.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to stop workflow: %w", err)
	}

	// Now that the workflow is stopped, create a new key with status as aborted
	return r.SetStatus(ctx, id, models.WorkflowStatusAborted, nil, nil)
}

// CheckConnection verifies the Redis connection is healthy
func (r *RedisWorkflowRepository) checkConnection(ctx context.Context) error {
	// Set timeout for health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r.logger.Debug("checking Redis connection", zap.Any("addr", r.client.Options().Addr))

	// Try to ping Redis
	if err := r.client.Ping(timeoutCtx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}
