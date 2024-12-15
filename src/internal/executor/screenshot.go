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

	executorID := uuid.New()
	e.activeJobs.Store(executorID.String(), ctx)

	checkpoint := models.Checkpoint{
		BatchID: nil,
		Stage:   nil,
	}
	go func() {
		defer e.activeJobs.Delete(executorID.String())
		defer close(updateChan)
		defer close(errorChan)
		e.execute(updateChan, errorChan, &checkpoint)
	}()

	return updateChan, errorChan
}

// Recover recovers the workflow
func (e *screenshotExecutor) Recover(ctx context.Context, workflowID *models.WorkflowIdentifier, checkpoint *models.Checkpoint) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError) {
	errorChan := make(chan models.WorkflowError)
	updateChan := make(chan models.WorkflowUpdate)

	executorID := uuid.New()
	e.activeJobs.Store(executorID.String(), ctx)

	go func() {
		defer e.activeJobs.Delete(executorID.String())
		defer close(updateChan)
		defer close(errorChan)
		e.execute(updateChan, errorChan, checkpoint)
	}()
	return updateChan, errorChan
}

// Stop stops the workflow
func (e *screenshotExecutor) Stop(ctx context.Context, workflowID *models.WorkflowIdentifier, executorID *uuid.UUID) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError) {
	errorChan := make(chan models.WorkflowError)
	updateChan := make(chan models.WorkflowUpdate)

	go func() {
		if executorCtx, ok := e.activeJobs.Load(executorID.String()); ok {
			if ctx, ok := executorCtx.(context.Context); ok {
				select {
				case <-ctx.Done():
					return
				default:
					// Perform any cleanup needed
					e.activeJobs.Delete(executorID.String())
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
	return activeJobs
}

// Execute executes the workflow
// This is the main logic of the workflow
func (e *screenshotExecutor) execute(errorChan chan models.WorkflowUpdate, updateChan chan models.WorkflowError, checkpoint *models.Checkpoint) {
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
	e.execute(errorChan, updateChan, nextCheckpoint)
}
