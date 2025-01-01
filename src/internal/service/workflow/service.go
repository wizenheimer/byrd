package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"

	exc "github.com/wizenheimer/iris/src/internal/interfaces/executor"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type workflowService struct {
	executors  sync.Map // map[WorkflowType]WorkflowExecutor
	logger     *logger.Logger
	repository repo.WorkflowRepository
}

func NewWorkflowService(logger *logger.Logger, repository repo.WorkflowRepository, screenshotWorkflowExecutor, reportWorkflowExecutor exc.WorkflowExecutor) (svc.WorkflowService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	ws := workflowService{
		logger:     logger,
		repository: repository,
	}

	// Register executors
	ws.registerExecutor(core_models.ScreenshotWorkflowType, screenshotWorkflowExecutor)
	ws.registerExecutor(core_models.ReportWorkflowType, reportWorkflowExecutor)

	// Initialize the workflow service
	if errors := ws.Initialize(context.Background()); len(errors) > 0 {
		return nil, fmt.Errorf("failed to initialize workflow service: %v", errors)
	} else {
		ws.logger.Info("workflow service initialized and ready")
	}

	return &ws, nil
}

func (s *workflowService) Initialize(ctx context.Context) []error {
	var errors []error
	// Initialize each executor
	s.executors.Range(func(key, value interface{}) bool {
		executor := value.(exc.WorkflowExecutor)
		if err := executor.Initialize(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to initialize %s executor: %w", key, err))
		}
		return true
	})

	return errors
}

func (s *workflowService) StartWorkflow(ctx context.Context, req api_models.WorkflowRequest) (api_models.WorkflowResponse, error) {
	if err := s.validateRequest(&req); err != nil {
		return api_models.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return api_models.WorkflowResponse{}, err
	}

	workflowID := core_models.WorkflowIdentifier(req)

	if err := executor.Start(ctx, workflowID); err != nil {
		return api_models.WorkflowResponse{}, err
	}

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return api_models.WorkflowResponse{}, err
	}

	return api_models.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) StopWorkflow(ctx context.Context, req api_models.WorkflowRequest) error {
	if err := s.validateRequest(&req); err != nil {
		return err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return err
	}

	workflowID := core_models.WorkflowIdentifier(req)

	return executor.Stop(ctx, workflowID)
}

func (s *workflowService) GetWorkflow(ctx context.Context, req api_models.WorkflowRequest) (api_models.WorkflowResponse, error) {
	if err := s.validateRequest(&req); err != nil {
		return api_models.WorkflowResponse{}, err
	}

	executor, err := s.getExecutor(*req.Type)
	if err != nil {
		return api_models.WorkflowResponse{}, err
	}

	workflowID := core_models.WorkflowIdentifier(req)

	state, err := executor.Get(ctx, workflowID)
	if err != nil {
		return api_models.WorkflowResponse{}, err
	}

	return api_models.WorkflowResponse{
		WorkflowID:    workflowID,
		WorkflowState: state,
	}, nil
}

func (s *workflowService) ListWorkflows(ctx context.Context, status core_models.WorkflowStatus, wfType core_models.WorkflowType) ([]api_models.WorkflowResponse, error) {
	return s.repository.List(ctx, status, wfType)
}
