package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	// Key prefixes for Redis
	workflowStatePrefix  = "workflow:state:"
	workflowStatusPrefix = "workflow:status:"

	// Default expiration time for workflow states
	defaultStateExpiration = 7 * 24 * time.Hour // 7 days
)

type workflowRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewWorkflowRepository(client *redis.Client, logger *logger.Logger) (interfaces.WorkflowRepository, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if client == nil {
		return nil, errors.New("client is required")
	}

	repo := workflowRepository{
		client: client,
		logger: logger,
	}

	if err := repo.checkConnection(context.Background()); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *workflowRepository) GetState(ctx context.Context, wi models.WorkflowIdentifier) (models.WorkflowState, error) {
	key := r.getStateKey(wi)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.WorkflowState{}, fmt.Errorf("workflow state not found: %v", wi)
		}
		return models.WorkflowState{}, fmt.Errorf("failed to get workflow state: %w", err)
	}

	var state models.WorkflowState
	if err := json.Unmarshal(data, &state); err != nil {
		return models.WorkflowState{}, fmt.Errorf("failed to unmarshal workflow state: %w", err)
	}

	return state, nil
}

func (r *workflowRepository) SetCheckpoint(ctx context.Context, wi models.WorkflowIdentifier, ws models.WorkflowStatus, wc models.WorkflowCheckpoint) error {
	// First get the existing state
	state, err := r.GetState(ctx, wi)
	if err != nil && err.Error() != "workflow state not found" {
		return err
	}

	// Update the state with new checkpoint and status
	state.Status = ws
	state.Checkpoint = wc

	// Store the updated state
	return r.SetState(ctx, wi, state)
}

func (r *workflowRepository) SetState(ctx context.Context, wi models.WorkflowIdentifier, ws models.WorkflowState) error {
	// Marshal the state
	data, err := json.Marshal(ws)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow state: %w", err)
	}

	// Store the state
	stateKey := r.getStateKey(wi)
	pipe := r.client.Pipeline()

	// Set the state with expiration
	pipe.Set(ctx, stateKey, data, defaultStateExpiration)

	// Update the status index
	statusKey := r.getStatusKey(ws.Status, *wi.Type)
	pipe.SAdd(ctx, statusKey, stateKey)
	pipe.Expire(ctx, statusKey, defaultStateExpiration)

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set workflow state: %w", err)
	}

	r.logger.Debug("stored workflow state",
		zap.Any("workflow_id", wi),
		zap.Any("status", ws.Status),
	)

	return nil
}

func (r *workflowRepository) List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]models.WorkflowState, error) {
	// Get all state keys for the given status and type
	statusKey := r.getStatusKey(ws, wt)
	stateKeys, err := r.client.SMembers(ctx, statusKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow keys: %w", err)
	}

	if len(stateKeys) == 0 {
		return []models.WorkflowState{}, nil
	}

	// Get all states in parallel using pipeline
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range stateKeys {
		cmds[key] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to get workflow states: %w", err)
	}

	// Process results
	var states []models.WorkflowState
	for key, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				// Key expired, remove from status set
				r.client.SRem(ctx, statusKey, key)
				continue
			}
			r.logger.Error("failed to get workflow state",
				zap.Any("key", key),
				zap.Any("error", err),
			)
			continue
		}

		var state models.WorkflowState
		if err := json.Unmarshal(data, &state); err != nil {
			r.logger.Error("failed to unmarshal workflow state",
				zap.Any("key", key),
				zap.Any("error", err),
			)
			continue
		}

		states = append(states, state)
	}

	return states, nil
}

// Helper methods for key management
func (r *workflowRepository) getStateKey(wi models.WorkflowIdentifier) string {
	return fmt.Sprintf("%s%s:%d:%d:%d",
		workflowStatePrefix,
		*wi.Type,
		*wi.Year,
		*wi.WeekNumber,
		*wi.WeekDay)
}

func (r *workflowRepository) getStatusKey(status models.WorkflowStatus, workflowType models.WorkflowType) string {
	return fmt.Sprintf("%s%s:%s",
		workflowStatusPrefix,
		workflowType,
		status)
}

// Cleanup method (optional, can be called periodically)
func (r *workflowRepository) Cleanup(ctx context.Context) error {
	// Scan for all status keys
	pattern := workflowStatusPrefix + "*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		statusKey := iter.Val()

		// Get all state keys in this status set
		stateKeys, err := r.client.SMembers(ctx, statusKey).Result()
		if err != nil {
			r.logger.Error("failed to get state keys",
				zap.Any("status_key", statusKey),
				zap.Any("error", err),
			)
			continue
		}

		// Check each state key
		for _, stateKey := range stateKeys {
			exists, err := r.client.Exists(ctx, stateKey).Result()
			if err != nil {
				r.logger.Error("failed to check state key",
					zap.Any("state_key", stateKey),
					zap.Any("error", err),
				)
				continue
			}

			if exists == 0 {
				// State doesn't exist, remove from status set
				r.client.SRem(ctx, statusKey, stateKey)
			}
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to cleanup workflow states: %w", err)
	}

	return nil
}

func (r *workflowRepository) checkConnection(ctx context.Context) error {
	// Set timeout for health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r.logger.Debug("checking redis client connection", zap.Any("addr", r.client.Options().Addr))

	// Try to ping Redis
	if err := r.client.Ping(timeoutCtx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	} else {
		r.logger.Debug("redis client connection successful")
	}

	return nil
}
