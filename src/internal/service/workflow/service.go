package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type workflowService struct {
	executors  sync.Map // map[WorkflowType]WorkflowExecutor
	logger     *logger.Logger
	repository interfaces.WorkflowRepository
}

func NewWorkflowService(logger *logger.Logger, repository interfaces.WorkflowRepository, screenshotWorkflowExecutor, reportWorkflowExecutor interfaces.WorkflowExecutor) (interfaces.WorkflowService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	ws := workflowService{
		logger:     logger,
		repository: repository,
	}

	// Register executors
	ws.registerExecutor(models.ScreenshotWorkflowType, screenshotWorkflowExecutor)
	ws.registerExecutor(models.ReportWorkflowType, reportWorkflowExecutor)

	return &ws, nil
}

func (s *workflowService) Initialize(ctx context.Context) []error {
	var errors []error
	// Initialize each executor
	s.executors.Range(func(key, value interface{}) bool {
		executor := value.(interfaces.WorkflowExecutor)
		if err := executor.Initialize(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to initialize %s executor: %w", key, err))
		}
		return true
	})

	return errors
}

func (s *workflowService) StartWorkflow(ctx context.Context, req models.WorkflowRequest) (models.WorkflowResponse, error) {
	if err := s.validateRequest(&req); err != nil {
		return models.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return models.WorkflowResponse{}, err
	}

	workflowID := models.WorkflowIdentifier(req)

	if err := executor.Start(ctx, workflowID); err != nil {
		return models.WorkflowResponse{}, err
	}

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return models.WorkflowResponse{}, err
	}

	return models.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) StopWorkflow(ctx context.Context, req models.WorkflowRequest) error {
	if err := s.validateRequest(&req); err != nil {
		return err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return err
	}

	workflowID := models.WorkflowIdentifier(req)

	return executor.Stop(ctx, workflowID)
}

func (s *workflowService) GetWorkflow(ctx context.Context, req models.WorkflowRequest) (models.WorkflowResponse, error) {
	if err := s.validateRequest(&req); err != nil {
		return models.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return models.WorkflowResponse{}, err
	}

	workflowID := models.WorkflowIdentifier(req)

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return models.WorkflowResponse{}, err
	}

	return models.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) ListWorkflows(ctx context.Context, status models.WorkflowStatus, wfType models.WorkflowType) ([]models.WorkflowResponse, error) {
	return s.repository.List(ctx, status, wfType)
}
