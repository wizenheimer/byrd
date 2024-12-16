package executor

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type screenshotExecutor struct {
	// urlService is an interface that is defined in the domain layer
	urlService interfaces.URLService
	// screenshotService is an interface that is defined in the domain layer
	screenshotService interfaces.ScreenshotService
	// diffService is an interface that is defined in the domain layer
	diffService interfaces.DiffService // Note: currently not used
	// logger is a structured logger for logging
	logger     *logger.Logger
	activeJobs sync.Map
}

// NewScreenshotWorkflowExecutor creates a new ScreenshotWorkflowService
func NewScreenshotWorkflowExecutor(screenshotService interfaces.ScreenshotService, diffService interfaces.DiffService, urlService interfaces.URLService, logger *logger.Logger) interfaces.WorkflowExecutor {
	return &screenshotExecutor{
		urlService:        urlService,
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
		e.execute(ctx, updateChan, errorChan, &checkpoint, workflowID)
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
		e.execute(ctx, updateChan, errorChan, checkpoint, workflowID)
	}()
	return updateChan, errorChan
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
func (e *screenshotExecutor) execute(ctx context.Context, updateChan chan models.WorkflowUpdate, errorChan chan models.WorkflowError, checkpoint *models.Checkpoint, workflowID *models.WorkflowIdentifier) {
	const batchWaitTime = 60 * time.Second

	updateChan <- models.WorkflowUpdate{
		ID:         workflowID,
		Checkpoint: checkpoint,
		Timestamp:  time.Now(),
		Status:     models.WorkflowStatusRunning,
	}

	// Check if the context is done
	select {
	case <-ctx.Done():
		e.logger.Debug("Context done", zap.Any("checkpoint", checkpoint))
		return
	default:
		e.logger.Debug("Context not done, executing workflow", zap.Any("checkpoint", checkpoint))
	}

	urlChan, errChan := e.urlService.ListURLs(ctx, 40, checkpoint.BatchID)

	// Create a ticker for batch processing
	batchTicker := time.NewTicker(batchWaitTime)
	defer batchTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Debug("Context done", zap.Any("checkpoint", checkpoint))
			updateChan <- models.WorkflowUpdate{
				ID:         workflowID,
				Checkpoint: checkpoint,
				Timestamp:  time.Now(),
				Status:     models.WorkflowStatusAborted,
			}
			return

		case err := <-errChan:
			e.logger.Error("Error processing URL Batch", zap.Error(err))
			// Send error to error channel but continue processing
			errorChan <- models.WorkflowError{
				ID:        workflowID,
				Error:     err,
				Timestamp: time.Now(),
			}

		case urlBatch, ok := <-urlChan:
			if !ok {
				// Channel closed, all URLs marked for processing
				e.logger.Debug("URL channel closed, workflow completed")
				updateChan <- models.WorkflowUpdate{
					ID:         workflowID,
					Checkpoint: checkpoint,
					Timestamp:  time.Now(),
					Status:     models.WorkflowStatusCompleted,
				}
				return
			}

			// Process the batch
			var wg sync.WaitGroup
			for _, url := range urlBatch.URLs {
				wg.Add(1)
				go e.processURLs(ctx, url, workflowID, errorChan, &wg)
			}

			// Wait for all URLs in this batch to be processed
			wg.Wait()

			// Update checkpoint and send status update
			checkpoint.BatchID = urlBatch.LastSeen
			checkpoint.Stage = nil // Reset the stage

			updateChan <- models.WorkflowUpdate{
				ID:         workflowID,
				Checkpoint: checkpoint,
				Timestamp:  time.Now(),
				Status:     models.WorkflowStatusRunning,
			}

			// Wait for the batch interval before processing next batch
			<-batchTicker.C
		}
	}
}

// processURLs processes the URLs
func (e *screenshotExecutor) processURLs(ctx context.Context, url models.URL, workflowID *models.WorkflowIdentifier, errorChan chan models.WorkflowError, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a sub-context with timeout for the screenshot process
	screenshotCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Get current screenshot image and content
	screenshotOptions := models.ScreenshotRequestOptions{
		URL: url.URL,
	}
	_, currentScreenshotImage, _, err := e.screenshotService.CaptureScreenshot(screenshotCtx, screenshotOptions)
	if err != nil {
		e.logger.Error("Error capturing screenshot",
			zap.Error(err),
			zap.String("url", url.URL))
		errorChan <- models.WorkflowError{
			ID:        workflowID,
			Error:     err,
			Timestamp: time.Now(),
		}
		return
	}

	// Get previous screenshot image and content
	screenshotImageResponse, err := e.screenshotService.GetPreviousScreenshotImage(screenshotCtx, url.URL)
	if err != nil {
		e.logger.Error("Error getting previous screenshot",
			zap.Error(err),
			zap.String("url", url.URL))
		errorChan <- models.WorkflowError{
			ID:        workflowID,
			Error:     err,
			Timestamp: time.Now(),
		}
		return
	}
	previousScreenshotImage := screenshotImageResponse.Image

	// Perform diff analysis
	if _, err := e.diffService.CreateCurrentDiffFromScreenshotImages(screenshotCtx, url.URL, currentScreenshotImage, previousScreenshotImage); err != nil {
		e.logger.Error("Error creating diff",
			zap.Error(err),
			zap.String("url", url.URL))
		errorChan <- models.WorkflowError{
			ID:        workflowID,
			Error:     err,
			Timestamp: time.Now(),
		}
		return
	}

}
