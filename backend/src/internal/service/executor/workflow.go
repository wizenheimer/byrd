// ./src/internal/service/executor/workflow.go
package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/alert"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowExecutor struct {
	workflowType models.WorkflowType
	config       models.ExecutorConfig
	repository   workflow.WorkflowRepository
	alertClient  alert.WorkflowAlertClient
	taskExecutor TaskExecutor
	logger       *logger.Logger

	activeWorkflows sync.Map // map[string]*workflowContext
}

type workflowContext struct {
	cancel context.CancelFunc
	task   models.Task
	state  api.WorkflowState
	mutex  sync.RWMutex
}

func NewWorkflowExecutor(
	wfType models.WorkflowType,
	repository workflow.WorkflowRepository,
	alertClient alert.WorkflowAlertClient,
	taskExecutor TaskExecutor,
	logger *logger.Logger,
) (WorkflowExecutor, error) {
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
		logger:       logger.WithFields(map[string]interface{}{"module": "workflow_executor"}),
	}

	return workflowExecutor, nil
}

func (e *workflowExecutor) Initialize(ctx context.Context) error {
	e.logger.Debug("initializing workflow executor")
	// Iterate over all active workflows
	workflowList, err := e.repository.List(ctx, models.WorkflowStatusRunning, e.workflowType)
	if err != nil {
		return fmt.Errorf("failed to list active workflows: %w", err)
	}

	errChan := make(chan error, len(workflowList))
	defer close(errChan)

	go func() {
		for err := range errChan {
			e.logger.Error("failed to restart workflow", zap.Error(err))
		}
	}()

	for _, workflow := range workflowList {
		// Start workflow execution
		go e.Restart(ctx, workflow.WorkflowID, errChan)
	}

	return nil
}

func (e *workflowExecutor) Restart(ctx context.Context, workflowID models.WorkflowIdentifier, errChan chan error) {
	e.logger.Debug("restarting workflow", zap.Any("workflow_id", workflowID))
	// Create task ID
	taskID := uuid.New().String()

	// Get the workflow state
	state, err := e.repository.GetState(ctx, workflowID)
	if err != nil {
		errChan <- fmt.Errorf("failed to get workflow state: %w", err)
		return
	}

	task := models.Task{
		TaskID:     taskID,
		WorkflowID: workflowID,
		Checkpoint: state.Checkpoint,
	}

	// Create workflow context
	workflowCtx, cancel := context.WithCancel(ctx)

	wfCtx := &workflowContext{
		cancel: cancel,
		task:   task,
		state: api.WorkflowState{
			Status: models.WorkflowStatusRunning,
		},
	}

	// Store in active workflows
	e.activeWorkflows.Store(taskID, wfCtx)

	// Send start alert
	if err := e.alertClient.SendWorkflowStarted(ctx, workflowID, map[string]string{
		"task_id": taskID,
	}); err != nil {
		e.logger.Error("failed to send start alert", zap.Error(err))
	}

	// Start execution
	go e.executeWorkflow(workflowCtx, wfCtx)
}

func (e *workflowExecutor) Start(ctx context.Context, workflowID models.WorkflowIdentifier) error {
	e.logger.Debug("starting workflow", zap.Any("workflow_id", workflowID))
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
		state: api.WorkflowState{
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
	e.logger.Debug("stopping workflow", zap.Any("workflow_id", workflowID))
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

func (e *workflowExecutor) List(ctx context.Context, status models.WorkflowStatus, wfType models.WorkflowType) ([]api.WorkflowState, error) {
	e.logger.Debug("listing workflows", zap.Any("status", status), zap.Any("type", wfType))
	workflowList, err := e.repository.List(ctx, status, wfType)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	workflowStateList := make([]api.WorkflowState, 0, len(workflowList))
	for _, workflow := range workflowList {
		workflowStateList = append(workflowStateList, workflow.WorkflowState)
	}
	return workflowStateList, nil
}

func (e *workflowExecutor) Get(ctx context.Context, workflowID models.WorkflowIdentifier) (api.WorkflowState, error) {
	e.logger.Debug("getting workflow", zap.Any("workflow_id", workflowID))
	return e.repository.GetState(ctx, workflowID)
}

func (e *workflowExecutor) executeWorkflow(ctx context.Context, wfCtx *workflowContext) {
	e.logger.Debug("executing workflow", zap.Any("workflow_id", wfCtx.task.WorkflowID))
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
	e.logger.Debug("handling task update", zap.Any("workflow_id", wfCtx.task.WorkflowID), zap.Any("update", update))
	wfCtx.mutex.Lock()
	wfCtx.state.Checkpoint = update.NewCheckpoint
	wfCtx.mutex.Unlock()

	if err := e.repository.SetCheckpoint(ctx, wfCtx.task.WorkflowID, wfCtx.state.Status, update.NewCheckpoint); err != nil {
		e.logger.Error("failed to update checkpoint", zap.Error(err))
	}
}

func (e *workflowExecutor) handleTaskError(ctx context.Context, wfCtx *workflowContext, taskErr models.TaskError) {
	e.logger.Debug("handling task error", zap.Any("workflow_id", wfCtx.task.WorkflowID), zap.Any("error", taskErr))
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusFailed
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	if err := e.alertClient.SendWorkflowFailed(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
		"error":   taskErr.Error.Error(),
	}); err != nil {
		e.logger.Error("failed to send error alert", zap.Error(err))
	}

}

func (e *workflowExecutor) handleWorkflowCompletion(ctx context.Context, wfCtx *workflowContext) {
	e.logger.Debug("handling workflow completion", zap.Any("workflow_id", wfCtx.task.WorkflowID))
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusCompleted
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	if err := e.alertClient.SendWorkflowCompleted(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
	}); err != nil {
		e.logger.Error("failed to send completion alert", zap.Error(err))
	}
}

func (e *workflowExecutor) handleWorkflowCancellation(ctx context.Context, wfCtx *workflowContext) {
	e.logger.Debug("handling workflow cancellation", zap.Any("workflow_id", wfCtx.task.WorkflowID))
	wfCtx.mutex.Lock()
	wfCtx.state.Status = models.WorkflowStatusAborted
	wfCtx.mutex.Unlock()

	if err := e.repository.SetState(ctx, wfCtx.task.WorkflowID, wfCtx.state); err != nil {
		e.logger.Error("failed to update state", zap.Error(err))
	}

	if err := e.alertClient.SendWorkflowCancelled(ctx, wfCtx.task.WorkflowID, map[string]string{
		"task_id": wfCtx.task.TaskID,
	}); err != nil {
		e.logger.Error("failed to send cancel alert", zap.Error(err))
	}
}
