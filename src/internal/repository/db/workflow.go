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
	client         *redis.Client
	logger         *logger.Logger
	workflowTTL    time.Duration
	workflowPrefix string
}

func NewRedisWorkflowRepository(client *redis.Client, ttl time.Duration, prefix string, logger *logger.Logger) (interfaces.WorkflowRepository, error) {
	repo := RedisWorkflowRepository{
		client:         client,
		logger:         logger,
		workflowTTL:    ttl,
		workflowPrefix: prefix,
	}

	// Verify connection on startup
	if err := repo.CheckConnection(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &repo, nil
}

// SetStatus implements WorkflowRepository.SetStatus
func (r *RedisWorkflowRepository) SetStatus(ctx context.Context, id *models.WorkflowIdentifier, status models.WorkflowStatus, batchID *string) error {
	// Create workflow response object
	workflow := &models.WorkflowResponse{
		Type:         id.Type,
		Year:         id.Year,
		WeekNumber:   id.WeekNumber,
		BucketNumber: id.BucketNumber,
		Status:       status,
		BatchID:      batchID,
	}

	// Marshal workflow to JSON
	workflowJSON, err := json.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	// Set workflow status with TTL
	key := r.workflowPrefix + id.Serialize()
	if err := r.client.Set(ctx, key, workflowJSON, r.workflowTTL).Err(); err != nil {
		return fmt.Errorf("failed to set workflow status: %w", err)
	}

	return nil
}

// GetStatus implements WorkflowRepository.GetStatus
func (r *RedisWorkflowRepository) GetStatus(ctx context.Context, id *models.WorkflowIdentifier) (*models.WorkflowResponse, error) {
	key := r.workflowPrefix + id.Serialize()

	// Get workflow data
	workflowJSON, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist or has expired
			return &models.WorkflowResponse{
				Type:         id.Type,
				Year:         id.Year,
				WeekNumber:   id.WeekNumber,
				BucketNumber: id.BucketNumber,
				Status:       models.WorkflowStatusExpired,
				BatchID:      nil,
			}, nil
		}
		return nil, fmt.Errorf("failed to get workflow status: %w", err)
	}

	// Unmarshal workflow data
	workflow := &models.WorkflowResponse{}
	if err := json.Unmarshal(workflowJSON, workflow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	return workflow, nil
}

// ListWorkflows implements WorkflowRepository.ListWorkflows
func (r *RedisWorkflowRepository) ListWorkflows(ctx context.Context, limit int) ([]models.WorkflowResponse, int, error) {
	var cursor uint64
	var total int
	var workflows []models.WorkflowResponse
	pattern := r.workflowPrefix + "*"

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
			for _, cmd := range cmds {
				workflowJSON, err := cmd.Bytes()
				if err == redis.Nil {
					// Skip expired workflows
					continue
				}
				if err != nil {
					continue // Skip workflows that can't be retrieved
				}

				workflow := models.WorkflowResponse{}
				if err := json.Unmarshal(workflowJSON, &workflow); err != nil {
					continue // Skip workflows that can't be unmarshaled
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

// CheckConnection verifies the Redis connection is healthy
func (r *RedisWorkflowRepository) CheckConnection(ctx context.Context) error {
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
