package executor

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
	"go.uber.org/zap"
)

type screenshotExecutor struct {
	// screenshotService is an interface that is defined in the domain layer
	screenshotService interfaces.ScreenshotService
	// diffService is an interface that is defined in the domain layer
	diffService interfaces.DiffService // Note: currently not used
	// logger is a structured logger for logging
	logger     *logger.Logger
	activeJobs sync.Map
}

// NewScreenshotWorkflowExecutor creates a new ScreenshotWorkflowService
func NewScreenshotWorkflowExecutor(screenshotService interfaces.ScreenshotService, diffService interfaces.DiffService, logger *logger.Logger) interfaces.WorkflowExecutor {
	return &screenshotExecutor{
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger.WithFields(map[string]interface{}{"module": "screenshot_workflow_service"}),
	}
}

// ExecuteWorkflow executes a screenshot workflow
func (e *screenshotExecutor) Start(ctx context.Context, workflowID *models.WorkflowIdentifier) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError) {
	errorChan := make(chan models.WorkflowError)
	updateChan := make(chan models.WorkflowUpdate)

	jobID := uuid.New()
	e.activeJobs.Store(jobID.String(), ctx)

	checkpoint := models.Checkpoint{
		BatchID: nil,
		Stage:   nil,
	}

	e.logger.Debug("Starting workflow", zap.Any("workflow_id", workflowID), zap.Any("job_id", jobID), zap.Any("checkpoint", checkpoint))

	go func() {
		defer e.activeJobs.Delete(jobID.String())
		defer close(updateChan)
		defer close(errorChan)
		e.logger.Debug("Executing workflow", zap.Any("workflow_id", workflowID), zap.Any("job_id", jobID), zap.Any("checkpoint", checkpoint))
		e.execute(ctx, updateChan, errorChan, &checkpoint)
		e.sendCompletionAlert(workflowID, updateChan)
	}()

	return updateChan, errorChan
}

// Recover recovers the workflow
func (e *screenshotExecutor) Recover(ctx context.Context, workflowID *models.WorkflowIdentifier, checkpoint *models.Checkpoint) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError) {
	errorChan := make(chan models.WorkflowError)
	updateChan := make(chan models.WorkflowUpdate)

	executorID := uuid.New()
	e.activeJobs.Store(executorID.String(), ctx)

	e.logger.Debug("Recovering workflow", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID), zap.Any("checkpoint", checkpoint))

	go func() {
		defer e.activeJobs.Delete(executorID.String())
		defer close(updateChan)
		defer close(errorChan)
		e.logger.Debug("Recovering workflow", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID), zap.Any("checkpoint", checkpoint))
		e.execute(ctx, updateChan, errorChan, checkpoint)
		e.sendCompletionAlert(workflowID, updateChan)
	}()
	return updateChan, errorChan
}

func (e *screenshotExecutor) sendCompletionAlert(workflowID *models.WorkflowIdentifier, updateChan chan models.WorkflowUpdate) {
	e.logger.Debug("Sending completion alert", zap.Any("workflow_id", workflowID))
	// Implement sending completion alert here
	update := models.WorkflowUpdate{
		ID:         workflowID,
		Checkpoint: nil,
		Timestamp:  time.Now(),
		Status:     models.WorkflowStatusCompleted,
	}

	e.logger.Debug("Sending completion alert", zap.Any("workflow_id", workflowID))
	// Send the update to the update channel
	updateChan <- update
	e.logger.Debug("Sent completion alert", zap.Any("workflow_id", workflowID))
}

// Stop stops the workflow
func (e *screenshotExecutor) Stop(ctx context.Context, workflowID *models.WorkflowIdentifier, executorID *uuid.UUID) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError) {
	e.logger.Debug("Stopping workflow", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID))

	errorChan := make(chan models.WorkflowError)
	updateChan := make(chan models.WorkflowUpdate)

	go func() {
		defer close(updateChan)
		defer close(errorChan)

		e.logger.Debug("Stopping workflow", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID))

		if executorCtx, ok := e.activeJobs.Load(executorID.String()); ok {
			if ctx, ok := executorCtx.(context.Context); ok {
				select {
				case <-ctx.Done():
					// Workflow already stopped
					e.activeJobs.Delete(executorID.String())
					e.logger.Debug("Workflow already stopped", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID))
					return
				default:
					// Perform any cleanup needed
					e.activeJobs.Delete(executorID.String())
					e.logger.Debug("Workflow stopped", zap.Any("workflow_id", workflowID), zap.Any("executor_id", executorID))
					return
				}
			}
		}
	}()

	return updateChan, errorChan
}

// List returns a list of active jobs
func (e *screenshotExecutor) List() map[string]context.Context {
	activeJobs := make(map[string]context.Context)
	e.activeJobs.Range(func(key, value interface{}) bool {
		if ctx, ok := value.(context.Context); ok {
			activeJobs[key.(string)] = ctx
		}
		return true
	})
	e.logger.Debug("Listing active jobs", zap.Any("active_jobs", activeJobs))
	return activeJobs
}

// Execute executes the workflow
// This is the main logic of the workflow
func (e *screenshotExecutor) execute(ctx context.Context, errorChan chan models.WorkflowUpdate, updateChan chan models.WorkflowError, checkpoint *models.Checkpoint) {

	// Check if the context is done
	select {
	case <-ctx.Done():
		e.logger.Debug("Context done", zap.Any("checkpoint", checkpoint))
		return
	default:
		e.logger.Debug("Context not done, executing workflow", zap.Any("checkpoint", checkpoint))
	}

	e.logger.Debug("Executing workflow", zap.Any("checkpoint", checkpoint))
	// Implement the workflow logic here
	step := 0
	if checkpoint.Stage != nil {
		step = *checkpoint.Stage
	}

	nextCheckpoint := &models.Checkpoint{
		BatchID: nil,
		Stage:   ptr.To(step + 1),
	}

	switch step {
	case 0:
		// Step 0: Take a screenshot
		e.logger.Info("Taking a screenshot")
		// Iterate through the batches
		time.Sleep(10 * time.Second)

	case 1:
		// Step 1: Process the screenshot
		e.logger.Info("Processing the screenshot")
		// Iterate through the batches
		time.Sleep(10 * time.Second)

	case 2:
		// Step 2: Compare the screenshot
		e.logger.Info("Comparing the screenshot")
		// Iterate through the batches
		time.Sleep(10 * time.Second)

	case 3:
		// Step 3: Send the report
		e.logger.Info("Sending the report")
		// Iterate through the batches
		time.Sleep(10 * time.Second)

	default:
		// Done
		e.logger.Info("Workflow completed")
		return
	}
	// Once done, move to the next step
	e.logger.Debug("Moving to the next step", zap.Any("checkpoint", nextCheckpoint))
	e.execute(ctx, errorChan, updateChan, nextCheckpoint)
}
