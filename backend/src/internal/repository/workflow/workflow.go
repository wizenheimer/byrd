package workflow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewWorkflowRepository(client *redis.Client, logger *logger.Logger) (WorkflowRepository, error) {
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

func (r *workflowRepository) Initialize(ctx context.Context, jobProps models.Job, workflowType models.WorkflowType) error {
	r.logger.Debug("initializing workflow", zap.Any("jobProps", jobProps), zap.Any("workflowType", workflowType))
	// TODO: Persist the triggered state in db
	return nil
}

func (r *workflowRepository) GetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) (models.JobState, error) {
	r.logger.Debug("getting state for job", zap.Any("jobID", jobID), zap.Any("workflowType", workflowType))
	return models.JobState{}, nil
}

func (r *workflowRepository) SetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobState models.JobState) error {
	r.logger.Debug("setting state for job", zap.Any("jobID", jobID), zap.Any("workflowType", workflowType))
	return nil
}

func (r *workflowRepository) SetCheckpoint(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobCheckpoint models.JobCheckpoint) error {
	r.logger.Debug("setting checkpoint for job", zap.Any("jobID", jobID), zap.Any("workflowType", workflowType))
	return nil
}

func (r *workflowRepository) List(ctx context.Context, workflowType models.WorkflowType, jobStatus models.JobStatus) ([]models.Job, error) {
	r.logger.Debug("listing jobs", zap.Any("workflowType", workflowType), zap.Any("jobStatus", jobStatus))
	return nil, nil
}
