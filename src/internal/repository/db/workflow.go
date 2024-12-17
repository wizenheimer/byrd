package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	// Key prefixes for Redis
	workflowKeyPrefix = "workflow:"
	statusKeyPrefix   = "status:"

	// TTL for workflow states (30 days)
	workflowStateTTL = 30 * 24 * time.Hour
)

type workflowRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewWorkflowRepository(client *redis.Client, logger *logger.Logger) (interfaces.WorkflowRepository, error) {
	workflowRepo := workflowRepository{
		client: client,
		logger: logger.WithFields(map[string]interface{}{"module": "workflow_repository"}),
	}
	if err := validateClient(context.Background(), client); err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	return &workflowRepo, nil
}

func validateClient(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("nil client")
	}

	// Set timeout for health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to ping Redis
	if err := client.Ping(timeoutCtx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	return nil
}

// GetState retrieves workflow state from Redis
func (r *workflowRepository) GetState(ctx context.Context, wi models.WorkflowIdentifier) (models.WorkflowState, error) {
	key := r.generateWorkflowKey(wi)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.WorkflowState{}, fmt.Errorf("workflow state not found: %s", key)
		}
		return models.WorkflowState{}, fmt.Errorf("failed to get workflow state: %w", err)
	}

	var state models.WorkflowState
	if err := json.Unmarshal(data, &state); err != nil {
		return models.WorkflowState{}, fmt.Errorf("failed to unmarshal workflow state: %w", err)
	}

	return state, nil
}

// SetCheckpoint updates workflow checkpoint in Redis
func (r *workflowRepository) SetCheckpoint(
	ctx context.Context,
	wi models.WorkflowIdentifier,
	ws models.WorkflowStatus,
	wc models.WorkflowCheckpoint,
) error {
	// Get existing state
	state, err := r.GetState(ctx, wi)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to get existing state: %w", err)
	}

	// Update state
	state.Status = ws
	state.Checkpoint = wc

	// Store updated state
	return r.SetState(ctx, wi, state)
}

// SetState stores complete workflow state in Redis
func (r *workflowRepository) SetState(
	ctx context.Context,
	wi models.WorkflowIdentifier,
	ws models.WorkflowState,
) error {
	key := r.generateWorkflowKey(wi)

	// Marshal state
	data, err := json.Marshal(ws)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow state: %w", err)
	}

	// Store in Redis with TTL
	pipe := r.client.Pipeline()

	// Store workflow state
	pipe.Set(ctx, key, data, workflowStateTTL)

	// Add to status set
	statusKey := r.generateStatusKey(ws.Status, *wi.Type)
	pipe.SAdd(ctx, statusKey, key)
	pipe.Expire(ctx, statusKey, workflowStateTTL)

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to store workflow state: %w", err)
	}

	r.logger.Debug("stored workflow state",
		zap.String("key", key),
		zap.String("status", string(ws.Status)),
		zap.Any("checkpoint", ws.Checkpoint))

	return nil
}

// List retrieves workflows by status and type
func (r *workflowRepository) List(
	ctx context.Context,
	ws models.WorkflowStatus,
	wt models.WorkflowType,
) ([]models.WorkflowResponse, error) {
	statusKey := r.generateStatusKey(ws, wt)

	// Get all workflow keys for the status
	keys, err := r.client.SMembers(ctx, statusKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow keys: %w", err)
	}

	if len(keys) == 0 {
		return []models.WorkflowResponse{}, nil
	}

	// Get all workflow states in a pipeline
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to get workflow states: %w", err)
	}

	// Process results
	var responses []models.WorkflowResponse

	for key, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				// Key expired between SMEMBERS and GET
				continue
			}
			r.logger.Error("failed to get workflow state",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		var state models.WorkflowState
		if err := json.Unmarshal(data, &state); err != nil {
			r.logger.Error("failed to unmarshal workflow state",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		identifier, err := r.parseWorkflowKey(key)
		if err != nil {
			r.logger.Error("failed to parse workflow key",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		responses = append(responses, models.WorkflowResponse{
			WorkflowID:    identifier,
			WorkflowState: state,
		})
	}

	return responses, nil
}

// Helper methods for key management
func (r *workflowRepository) generateWorkflowKey(wi models.WorkflowIdentifier) string {
	return fmt.Sprintf("%s%s:%d:%d:%d",
		workflowKeyPrefix,
		*wi.Type,
		*wi.Year,
		*wi.WeekNumber,
		*wi.WeekDay)
}

func (r *workflowRepository) generateStatusKey(status models.WorkflowStatus, workflowType models.WorkflowType) string {
	return fmt.Sprintf("%s%s:%s", statusKeyPrefix, workflowType, status)
}

func (r *workflowRepository) parseWorkflowKey(key string) (models.WorkflowIdentifier, error) {
	// Remove prefix
	key = strings.TrimPrefix(key, workflowKeyPrefix)

	// Split remaining parts
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return models.WorkflowIdentifier{}, fmt.Errorf("invalid workflow key format: %s", key)
	}

	// Parse values
	wfType := models.WorkflowType(parts[0])
	var year, weekNum, weekDay int

	_, err := fmt.Sscanf(parts[1], "%d", &year)
	if err != nil {
		return models.WorkflowIdentifier{}, fmt.Errorf("invalid year format: %s", parts[1])
	}

	_, err = fmt.Sscanf(parts[2], "%d", &weekNum)
	if err != nil {
		return models.WorkflowIdentifier{}, fmt.Errorf("invalid week number format: %s", parts[2])
	}

	_, err = fmt.Sscanf(parts[3], "%d", &weekDay)
	if err != nil {
		return models.WorkflowIdentifier{}, fmt.Errorf("invalid week day format: %s", parts[3])
	}

	return models.WorkflowIdentifier{
		Type:       &wfType,
		Year:       &year,
		WeekNumber: &weekNum,
		WeekDay:    &weekDay,
	}, nil
}
