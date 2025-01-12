// ./src/internal/service/workflow/service.go
package workflow

import (
	"context"
	"errors"
	"fmt"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowService struct {
	executors  map[models.WorkflowType]executor.WorkflowExecutor
	logger     *logger.Logger
	repository workflow.WorkflowRepository
}

func NewWorkflowService(logger *logger.Logger, repository workflow.WorkflowRepository, screenshotWorkflowExecutor, reportWorkflowExecutor executor.WorkflowExecutor) (WorkflowService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	ws := workflowService{
		logger:     logger.WithFields(map[string]interface{}{"module": "workflow_service"}),
		repository: repository,
		executors:  make(map[models.WorkflowType]executor.WorkflowExecutor),
	}

	// Register executors
	ws.registerExecutor(models.ScreenshotWorkflowType, screenshotWorkflowExecutor)
	ws.registerExecutor(models.ReportWorkflowType, reportWorkflowExecutor)

	// Initialize the workflow service
	if errors := ws.Initialize(context.Background()); len(errors) > 0 {
		return nil, fmt.Errorf("failed to initialize workflow service: %v", errors)
	} else {
		ws.logger.Info("workflow service initialized and ready")
	}

	return &ws, nil
}

func (s *workflowService) Initialize(ctx context.Context) []error {
	s.logger.Debug("initializing workflow service")
	var errors []error
	// Initialize each executor
	for key, executor := range s.executors {
		if err := executor.Initialize(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to initialize %s executor: %w", key, err))
		}
	}

	return errors
}

func (s *workflowService) StartWorkflow(ctx context.Context, req api.WorkflowRequest) (api.WorkflowResponse, error) {
	s.logger.Debug("starting workflow", zap.Any("request", req))
	if err := s.validateRequest(&req); err != nil {
		return api.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return api.WorkflowResponse{}, err
	}

	workflowID := models.WorkflowIdentifier(req)

	if err := executor.Start(ctx, workflowID); err != nil {
		return api.WorkflowResponse{}, err
	}

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return api.WorkflowResponse{}, err
	}

	return api.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) StopWorkflow(ctx context.Context, req api.WorkflowRequest) error {
	s.logger.Debug("stopping workflow", zap.Any("request", req))
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

func (s *workflowService) GetWorkflow(ctx context.Context, req api.WorkflowRequest) (api.WorkflowResponse, error) {
	s.logger.Debug("getting workflow", zap.Any("request", req))
	if err := s.validateRequest(&req); err != nil {
		return api.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return api.WorkflowResponse{}, err
	}

	workflowID := models.WorkflowIdentifier(req)

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return api.WorkflowResponse{}, err
	}

	return api.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) ListWorkflows(ctx context.Context, status models.WorkflowStatus, wfType models.WorkflowType) ([]api.WorkflowResponse, error) {
	s.logger.Debug("listing workflows", zap.Any("status", status), zap.Any("type", wfType))
	return s.repository.List(ctx, status, wfType)
}
