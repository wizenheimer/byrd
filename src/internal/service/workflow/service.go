package workflow

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowService struct {
	workflowRepo interfaces.WorkflowRepository
	logger       *logger.Logger
}

// NewWorkflowService creates a new WorkflowService
func NewWorkflowService(workflowRepo interfaces.WorkflowRepository, logger *logger.Logger) interfaces.WorkflowService {
	return &workflowService{
		workflowRepo: workflowRepo,
		logger:       logger.WithFields(map[string]interface{}{"module": "workflow_service"}),
	}
}

// StartWorkflow starts a new workflow
func (s *workflowService) StartWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
	workflowID := models.WorkflowIdentifier{
		Type:         req.Type,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
	}
	if err := s.workflowRepo.SetStatus(ctx, &workflowID, models.WorkflowStatusRunning, nil); err != nil {
		return nil, err
	}
	return &models.WorkflowResponse{
		Type:         req.Type,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
		BatchID:      nil,
		Status:       models.WorkflowStatusRunning,
	}, nil
}

// GetWorkflowStatus retrieves the status of a workflow
func (s *workflowService) GetWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
	workflowID := models.WorkflowIdentifier{
		Type:         req.Type,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
	}

	status, err := s.workflowRepo.GetStatus(ctx, &workflowID)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// ListWorkflows lists the workflows
func (s *workflowService) ListWorkflows(ctx context.Context, scanLimit int) ([]models.WorkflowResponse, int, error) {
	resp, count, err := s.workflowRepo.ListWorkflows(ctx, scanLimit)
	if err != nil {
		return nil, 0, err
	}
	s.logger.Debug("ListWorkflows", zap.Any("count", count), zap.Any("workflows", resp))
	return resp, count, nil
}
