package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowExecutor struct {
	workflowType models.WorkflowType
	config       models.ExecutorConfig
	repository   interfaces.WorkflowRepository
	alertClient  interfaces.WorkflowAlertClient
	taskExecutor interfaces.TaskExecutor
	logger       *logger.Logger

	activeWorkflows sync.Map // map[string]*workflowContext
}

type workflowContext struct {
	cancel context.CancelFunc
	task   models.Task
	state  models.WorkflowState
	mutex  sync.RWMutex
}

func NewWorkflowExecutor(
	wfType models.WorkflowType,
	repository interfaces.WorkflowRepository,
	alertClient interfaces.WorkflowAlertClient,
	taskExecutor interfaces.TaskExecutor,
	logger *logger.Logger,
) (interfaces.WorkflowExecutor, error) {
	config, err := models.GetExecutorConfig(wfType)
	if err != nil {
		return nil, err
	}

	workflowExecutor := &workflowExecutor{
		workflowType: wfType,
		config:       config,
		repository:   repository,
		alertClient:  alertClient,
		taskExecutor: taskExecutor,
		logger:       logger,
	}

	return workflowExecutor, nil
}

func (e *workflowExecutor) Initialize(ctx context.Context) error {
	// Iterate over all active workflows
	// For each active workflow, create a new context and start execution
	return nil
}

func (e *workflowExecutor) Start(ctx context.Context, workflowID models.WorkflowIdentifier) error {
	// Create task ID
	taskID := uuid.New().String()

	// Create workflow context
	workflowCtx, cancel := context.WithCancel(ctx)

	task := models.Task{
		TaskID:     taskID,
		WorkflowID: workflowID,
	}

	wfCtx := &workflowContext{
		cancel: cancel,
		task:   task,
		state: models.WorkflowState{
			Status: models.WorkflowStatusRunning,
		},
	}

	// Store in active workflows
	e.activeWorkflows.Store(taskID, wfCtx)

	// Initialize state in repository
	if err := e.repository.SetState(ctx, workflowID, wfCtx.state); err != nil {
		return fmt.Errorf("failed to set initial state: %w", err)
	}

	// Send start alert
	if err := e.alertClient.SendWorkflowStarted(ctx, workflowID, map[string]string{
		"task_id": taskID,
	}); err != nil {
		e.logger.Error("failed to send start alert", zap.Error(err))
	}

	// Start execution
	go e.executeWorkflow(workflowCtx, wfCtx)

	return nil
}

func (e *workflowExecutor) Stop(ctx context.Context, workflowID models.WorkflowIdentifier) error {
	var foundCtx *workflowContext

	// Find the workflow context
	e.activeWorkflows.Range(func(key, value interface{}) bool {
		wfCtx := value.(*workflowContext)
		if wfCtx.task.WorkflowID == workflowID {
			foundCtx = wfCtx
			return false
		}
		return true
	})

	if foundCtx == nil {
		return fmt.Errorf("workflow not found: %v", workflowID)
	}

	// Cancel the context
	foundCtx.cancel()

	// Update state
	foundCtx.mutex.Lock()
	foundCtx.state.Status = models.WorkflowStatusAborted
	foundCtx.mutex.Unlock()

	// Update repository
	if err := e.repository.SetState(ctx, workflowID, foundCtx.state); err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}

	// Send cancel alert
	if err := e.alertClient.SendWorkflowCancelled(ctx, workflowID, map[string]string{
		"task_id": foundCtx.task.TaskID,
	}); err != nil {
		e.logger.Error("failed to send cancel alert", zap.Error(err))
	}

	return nil
}

func (e *workflowExecutor) List(ctx context.Context, status models.WorkflowStatus, wfType models.WorkflowType) ([]models.WorkflowState, error) {
	return e.repository.List(ctx, status, wfType)
}

func (e *workflowExecutor) Get(ctx context.Context, workflowID models.WorkflowIdentifier) (models.WorkflowState, error) {
	return e.repository.GetState(ctx, workflowID)
}

func (e *workflowExecutor) executeWorkflow(ctx context.Context, wfCtx *workflowContext) {
	updates, errors := e.taskExecutor.Execute(ctx, wfCtx.task)

	for {
		select {
		case <-ctx.Done():
			e.handleWorkflowCancellation(ctx, wfCtx)
			return

		case err, ok := <-errors:
			if !ok {
				continue
			}
			e.handleTaskError(ctx, wfCtx, err)

		case update, ok := <-updates:
			if !ok {
				e.handleWorkflowCompletion(ctx, wfCtx)
				return
			}
			e.handleTaskUpdate(ctx, wfCtx, update)
		}
	}
}

func (e *workflowExecutor) handleTaskUpdate(ctx context.Context, wfCtx *workflowContext, update models.TaskUpdate) {
	wfCtx.mutex.Lock()
	wfCtx.state.Checkpoint = update.NewCheckpoint
	wfCtx.mutex.Unlock()

	if err := e.repository.SetCheckpoint(ctx, wfCtx.task.WorkflowID, wfCtx.state.Status, update.NewCheckpoint); err != nil {
		e.logger.Error("failed to update checkpoint", zap.Error(err))
	}
}

func (e *workflowExecutor) handleTaskError(ctx context.Context, wfCtx *workflowContext, taskErr models.TaskError) {
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusFailed
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	e.alertClient.SendWorkflowFailed(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
		"error":   taskErr.Error.Error(),
	})
}

func (e *workflowExecutor) handleWorkflowCompletion(ctx context.Context, wfCtx *workflowContext) {
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusCompleted
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	e.alertClient.SendWorkflowCompleted(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
	})
}

func (e *workflowExecutor) handleWorkflowCancellation(ctx context.Context, wfCtx *workflowContext) {
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusAborted
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	e.alertClient.SendWorkflowCancelled(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
	})
}
