package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	// Single key prefix for all workflow data
	workflowKeyPrefix = "workflow:"
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
		return fmt.Errorf("client is required")
	}
	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to ping client: %w", err)
	}
	return nil
}

func (r *workflowRepository) SetState(ctx context.Context, wi core_models.WorkflowIdentifier, ws api_models.WorkflowState) error {
	key := r.generateKey(wi)

	data, err := json.Marshal(ws)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := r.client.Set(ctx, key, data, workflowStateTTL).Err(); err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	return nil
}

func (r *workflowRepository) GetState(ctx context.Context, wi core_models.WorkflowIdentifier) (api_models.WorkflowState, error) {
	key := r.generateKey(wi)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return api_models.WorkflowState{}, fmt.Errorf("workflow not found: %s", key)
		}
		return api_models.WorkflowState{}, fmt.Errorf("failed to get state: %w", err)
	}

	var state api_models.WorkflowState
	if err := json.Unmarshal(data, &state); err != nil {
		return api_models.WorkflowState{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return state, nil
}

func (r *workflowRepository) SetCheckpoint(ctx context.Context, wi core_models.WorkflowIdentifier, status core_models.WorkflowStatus, checkpoint core_models.WorkflowCheckpoint) error {
	state := api_models.WorkflowState{
		Status:     status,
		Checkpoint: checkpoint,
	}
	return r.SetState(ctx, wi, state)
}

func (r *workflowRepository) List(ctx context.Context, status core_models.WorkflowStatus, wfType core_models.WorkflowType) ([]api_models.WorkflowResponse, error) {
	r.logger.Debug("listing workflows",
		zap.String("status", string(status)),
		zap.String("type", string(wfType)))

	var responses []api_models.WorkflowResponse
	var cursor uint64
	pattern := fmt.Sprintf("%s%s:*", workflowKeyPrefix, string(wfType))

	// Debug log the pattern
	r.logger.Debug("scanning with pattern", zap.String("pattern", pattern))

	for {
		// Scan for keys
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		r.logger.Debug("found keys", zap.Strings("keys", keys))

		// Process found keys
		for _, key := range keys {
			// Parse identifier from key
			identifier, err := r.parseKey(key)
			if err != nil {
				r.logger.Error("failed to parse key",
					zap.String("key", key),
					zap.Error(err))
				continue
			}

			// Get workflow state
			data, err := r.client.Get(ctx, key).Bytes()
			if err != nil {
				if err == redis.Nil {
					continue // Key expired
				}
				r.logger.Error("failed to get workflow data",
					zap.String("key", key),
					zap.Error(err))
				continue
			}

			var state api_models.WorkflowState
			if err := json.Unmarshal(data, &state); err != nil {
				r.logger.Error("failed to unmarshal state",
					zap.String("key", key),
					zap.Error(err))
				continue
			}

			// Filter by status if specified
			if status != "" && state.Status != status {
				continue
			}

			// Debug log the parsed data
			r.logger.Debug("parsed workflow data",
				zap.Any("identifier", identifier),
				zap.Any("state", state))

			responses = append(responses, api_models.WorkflowResponse{
				WorkflowID: core_models.WorkflowIdentifier{
					Type:       &wfType, // Use the input type
					Year:       identifier.Year,
					WeekNumber: identifier.WeekNumber,
					WeekDay:    identifier.WeekDay,
				},
				WorkflowState: state,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return responses, nil
}

func (r *workflowRepository) parseKey(key string) (core_models.WorkflowIdentifier, error) {
	// Remove prefix
	key = strings.TrimPrefix(key, workflowKeyPrefix)

	// Split the key parts
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return core_models.WorkflowIdentifier{}, fmt.Errorf("invalid key format: %s", key)
	}

	// Parse values
	workflowType := core_models.WorkflowType(parts[0])

	var year, weekNum, weekDay int
	if _, err := fmt.Sscanf(parts[1], "%d", &year); err != nil {
		return core_models.WorkflowIdentifier{}, fmt.Errorf("invalid year: %s", parts[1])
	}
	if _, err := fmt.Sscanf(parts[2], "%d", &weekNum); err != nil {
		return core_models.WorkflowIdentifier{}, fmt.Errorf("invalid week number: %s", parts[2])
	}
	if _, err := fmt.Sscanf(parts[3], "%d", &weekDay); err != nil {
		return core_models.WorkflowIdentifier{}, fmt.Errorf("invalid week day: %s", parts[3])
	}

	// Create new variables for the pointers
	wfTypeCopy := workflowType
	yearCopy := year
	weekNumCopy := weekNum
	weekDayCopy := weekDay

	// Return identifier with proper pointer values
	return core_models.WorkflowIdentifier{
		Type:       &wfTypeCopy,
		Year:       &yearCopy,
		WeekNumber: &weekNumCopy,
		WeekDay:    &weekDayCopy,
	}, nil
}

func (r *workflowRepository) generateKey(wi core_models.WorkflowIdentifier) string {
	if wi.Type == nil || wi.Year == nil || wi.WeekNumber == nil || wi.WeekDay == nil {
		r.logger.Error("invalid workflow identifier - nil values",
			zap.Any("identifier", wi))
		return ""
	}

	key := fmt.Sprintf("%s%s:%d:%d:%d",
		workflowKeyPrefix,
		*wi.Type,
		*wi.Year,
		*wi.WeekNumber,
		*wi.WeekDay)

	r.logger.Debug("generated key", zap.String("key", key))
	return key
}
