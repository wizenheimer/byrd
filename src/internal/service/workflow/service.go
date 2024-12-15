package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// workflowService implements WorkflowService
// It is responsible for starting and managing workflows
type workflowService struct {
	// Logger for logging
	logger *logger.Logger

	// Repository for checkpointing workflow status
	workflowRepo interfaces.WorkflowRepository

	// Alert client for sending alerts
	alertClient client.WorkflowAlertClient

	// Executor for different workflow types
	screenshotWorkflowExecutor interfaces.WorkflowExecutor
	reportWorkflowExecutor     interfaces.WorkflowExecutor

	// Service State
	activeWorkflows map[string]*models.WorkflowState
	workflowsMu     sync.RWMutex
	serviceCtx      context.Context
	serviceCancel   context.CancelFunc
}

// NewWorkflowService creates a new WorkflowService
func NewWorkflowService(workflowRepo interfaces.WorkflowRepository,
	alertClient client.WorkflowAlertClient,
	screenshotExecutor interfaces.WorkflowExecutor,
	reportExecutor interfaces.WorkflowExecutor,
	logger *logger.Logger,
) interfaces.WorkflowService {
	ctx, cancel := context.WithCancel(context.Background())

	return &workflowService{
		workflowRepo:               workflowRepo,
		alertClient:                alertClient,
		logger:                     logger.WithFields(map[string]interface{}{"module": "workflow_service"}),
		screenshotWorkflowExecutor: screenshotExecutor,
		reportWorkflowExecutor:     reportExecutor,
		activeWorkflows:            make(map[string]*models.WorkflowState),
		serviceCtx:                 ctx,
		serviceCancel:              cancel,
	}
}

// StartWorkflow starts a new workflow
func (s *workflowService) StartWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
	// Validate the request
	if err := req.Validate(true); err != nil {
		return nil, err
	}

	s.logger.Debug("starting workflow", zap.Any("workflow_type", req.Type), zap.Any("year", req.Year), zap.Any("week_number", *req.WeekNumber), zap.Any("bucket_number", *req.BucketNumber))

	// Parse the workflow type
	wfType, executor, err := s.parseWorkflowType(req.Type)
	if err != nil {
		return nil, err
	}

	// Create workflow identifier
	workflowID := models.WorkflowIdentifier{
		Type:         wfType,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
	}

	s.logger.Debug("starting executor", zap.Any("workflow_id", workflowID), zap.Any("workflow_type", wfType))

	// Start the executor
	updateChan, errorChan := executor.Start(ctx, &workflowID)

	// Generate workflow key
	prefix := models.GetWorkflowPrefixFromWorkflowType(wfType)
	workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

	// Start background processing
	go s.handleWorkflowBackgroundProcessing(
		workflowID,
		workflowKey,
		uuid.New(),
		updateChan,
		errorChan,
	)

	return &models.WorkflowResponse{
		Type:         wfType,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
		BatchID:      nil,
		Status:       models.WorkflowStatusRunning,
	}, nil
}

// StopWorkflow stops a running workflow
func (s *workflowService) StopWorkflow(ctx context.Context, workflowID models.WorkflowIdentifier) error {

	workflowKey := workflowID.Serialize(
		models.GetWorkflowPrefixFromWorkflowType(workflowID.Type),
		models.WorkflowStatusRunning,
	)

	s.logger.Debug("stopping workflow", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey))

	s.workflowsMu.RLock()
	workflowState, exists := s.activeWorkflows[workflowKey]
	s.workflowsMu.RUnlock()

	if !exists {
		return fmt.Errorf("workflow not found or not running: %s", workflowKey)
	}

	workflowState.Cancel()
	return nil
}

// Shutdown stops all running workflows
func (s *workflowService) Shutdown(ctx context.Context) error {
	s.logger.Debug("shutting down workflow service", zap.Any("workflows", s.activeWorkflows))

	s.serviceCancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("shutdown interupted", zap.Any("workflows", s.activeWorkflows))

			return ctx.Err()
		case <-ticker.C:
			s.logger.Debug("waiting for workflows to finish", zap.Any("workflows", s.activeWorkflows))

			s.workflowsMu.RLock()
			count := len(s.activeWorkflows)
			s.workflowsMu.RUnlock()

			if count == 0 {
				return nil
			}
		}
	}
}

// GetWorkflowStatus retrieves the status of a workflow
func (s *workflowService) GetWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
	// Validate the request
	if err := req.Validate(true); err != nil {
		return nil, err
	}

	s.logger.Debug("getting workflow status", zap.Any("workflow_type", req.Type), zap.Any("year", req.Year), zap.Any("week_number", *req.WeekNumber), zap.Any("bucket_number", *req.BucketNumber))

	var wfType models.WorkflowType
	switch req.Type {
	case "screenshot":
		wfType = models.ScreenshotWorkflowType
	case "report":
		wfType = models.ReportWorkflowType
	default:
		return nil, fmt.Errorf("unknown workflow type: %s", req.Type)
	}

	workflowID := models.WorkflowIdentifier{
		Type:         wfType,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
	}

	s.logger.Debug("getting workflow status", zap.Any("workflow_id", workflowID), zap.Any("workflow_type", wfType))

	status, err := s.workflowRepo.GetStatus(ctx, &workflowID)
	if err != nil {
		s.logger.Error("failed to get workflow status", zap.Any("workflow_id", workflowID), zap.Error(err))
		return nil, err
	}

	s.logger.Debug("workflow status", zap.Any("workflow_id", workflowID), zap.Any("status", status))

	return status, nil
}

// ListWorkflows lists the workflows
func (s *workflowService) ListWorkflows(ctx context.Context, status string, scanLimit int) ([]models.WorkflowResponse, int, error) {
	s.logger.Debug("listing workflows", zap.Any("status", status), zap.Any("scan_limit", scanLimit))

	// Convert status string to WorkflowStatus
	var workflowStatuses []models.WorkflowStatus
	switch status {
	case "running":
		workflowStatuses = append(workflowStatuses, models.WorkflowStatusRunning)
	case "completed":
		workflowStatuses = append(workflowStatuses, models.WorkflowStatusCompleted)
	case "failed":
		workflowStatuses = append(workflowStatuses, models.WorkflowStatusFailed)
	case "aborted":
		workflowStatuses = append(workflowStatuses, models.WorkflowStatusAborted)
	case "expired":
		workflowStatuses = append(workflowStatuses, models.WorkflowStatusExpired)
	case "all":
		workflowStatuses = []models.WorkflowStatus{
			models.WorkflowStatusRunning,
			models.WorkflowStatusCompleted,
			models.WorkflowStatusFailed,
			models.WorkflowStatusAborted,
			models.WorkflowStatusExpired,
		}
	default:
		return nil, 0, fmt.Errorf("invalid workflow status: %s", status)
	}

	s.logger.Debug("workflow statuses", zap.Any("statuses", workflowStatuses))

	var aggregateResp []models.WorkflowResponse
	var aggregateCount int
	for index, workflowStatus := range workflowStatuses {
		resp, count, err := s.workflowRepo.ListWorkflows(ctx, workflowStatus, scanLimit)
		if err != nil {
			return nil, 0, err
		}

		aggregateResp = append(aggregateResp, resp...)
		aggregateCount += count

		s.logger.Debug("aggregating listing run", zap.Any("count", count), zap.Any("workflows", resp), zap.Any("index", index), zap.Any("aggregate_count", aggregateCount))
	}

	return aggregateResp, aggregateCount, nil
}

// RecoverWorkflow recovers existing workflows
func (s *workflowService) RecoverWorkflow(ctx context.Context) error {
	// List all running workflows
	workflows, count, err := s.ListWorkflows(ctx, "running", 100)
	if err != nil {
		return err
	}

	s.logger.Debug("recovering workflows", zap.Any("workflows", workflows), zap.Any("count", count))

	// Iterate through workflows
	for _, workflow := range workflows {
		// Find the executor
		var executor interfaces.WorkflowExecutor
		switch workflow.Type {
		case models.ScreenshotWorkflowType:
			executor = s.screenshotWorkflowExecutor
		case models.ReportWorkflowType:
			executor = s.reportWorkflowExecutor
		default:
			s.logger.Error("unknown workflow type", zap.Any("workflow_type", workflow.Type))
			continue
		}

		// Prepare checkpoint
		checkpoint := models.Checkpoint{
			BatchID: workflow.BatchID,
			Stage:   workflow.Stage,
		}

		// Prepare identifier
		workflowID := models.WorkflowIdentifier{
			Type:         workflow.Type,
			Year:         workflow.Year,
			WeekNumber:   workflow.WeekNumber,
			BucketNumber: workflow.BucketNumber,
		}

		s.logger.Debug("recovering workflow", zap.Any("workflow_id", workflowID), zap.Any("workflow_type", workflow.Type), zap.Any("checkpoint", checkpoint))

		// Start the executor with the checkpoint
		updateChan, errorChan := executor.Recover(ctx, &workflowID, &checkpoint)

		// Generate workflow key
		prefix := models.GetWorkflowPrefixFromWorkflowType(workflow.Type)
		workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

        s.logger.Debug("recovering workflow", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey))
		// Start background processing
		go s.handleWorkflowBackgroundProcessing(
			workflowID,
			workflowKey,
			uuid.New(),
			updateChan,
			errorChan,
		)
	}
	return nil
}

// handleWorkflowBackgroundProcessing is a background processing loop for a workflow
func (s *workflowService) handleWorkflowBackgroundProcessing(
	workflowID models.WorkflowIdentifier,
	workflowKey string,
	executorID uuid.UUID,
	updateChan <-chan models.WorkflowUpdate,
	errorChan <-chan models.WorkflowError,
) {
	s.logger.Debug("handling processing in the background", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID))

	// Create background context
	backgroundCtx := context.Background()
	workflowCtx, cancel := context.WithCancel(backgroundCtx)

	// Initialize workflow state
	workflowState := &models.WorkflowState{
		Cancel:     cancel,
		ExecutorID: executorID,
		Status:     models.WorkflowStatusRunning,
	}

	// Register workflow
	s.workflowsMu.Lock()
	s.activeWorkflows[workflowKey] = workflowState
	s.workflowsMu.Unlock()

	// Cleanup on exit
	defer func() {
		s.logger.Debug("cleaning up workflow", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID))
		cancel()
		s.workflowsMu.Lock()
		delete(s.activeWorkflows, workflowKey)
		s.workflowsMu.Unlock()
	}()

	for {
		select {
		case <-s.serviceCtx.Done():
			// Service-wide shutdown
			s.logger.Debug("workflow service-wide shutdown triggered", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID))

			s.workflowRepo.SetStatus(backgroundCtx, &workflowID, models.WorkflowStatusAborted, nil, nil)
			return

		case err, ok := <-errorChan:
			if !ok {
				return
			}

			s.logger.Error("workflow error", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID), zap.Error(err.Error))

			workflowState.Mutex.Lock()
			workflowState.Status = models.WorkflowStatusFailed
			workflowState.Mutex.Unlock()

			alertParam := map[string]string{
				"workflowID": workflowKey,
				"error":      err.Error.Error(),
				"timestamp":  fmt.Sprintf("%v", err.Timestamp),
			}
			s.workflowRepo.SetStatus(backgroundCtx, &workflowID, models.WorkflowStatusFailed, nil, nil)
			s.alertClient.SendWorkflowFailed(backgroundCtx, workflowID, alertParam)

		case update, ok := <-updateChan:
			if !ok {
				return
			}

			s.logger.Debug("workflow update", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID), zap.Any("status", update.Status), zap.Any("batch_id", update.Checkpoint.BatchID), zap.Any("stage", update.Checkpoint.Stage))

			workflowState.Mutex.Lock()
			workflowState.Status = update.Status
			workflowState.Mutex.Unlock()

			s.workflowRepo.SetStatus(backgroundCtx, &workflowID, update.Status, update.Checkpoint.BatchID, update.Checkpoint.Stage)

			if update.Status == models.WorkflowStatusCompleted {
				alertParam := map[string]string{
					"worflow_key": workflowKey,
					"status":      string(update.Status),
				}
				s.alertClient.SendWorkflowCompleted(backgroundCtx, workflowID, alertParam)
			}

		case <-workflowCtx.Done():
			s.logger.Debug("workflow cancelled", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey), zap.Any("executor_id", executorID))

			if workflowCtx.Err() == context.Canceled {
				s.workflowRepo.SetStatus(backgroundCtx, &workflowID, models.WorkflowStatusAborted, nil, nil)
				alertParam := map[string]string{
					"workflowID": workflowKey,
					"status":     string(models.WorkflowStatusAborted),
				}
				s.alertClient.SendWorkflowCancelled(backgroundCtx, workflowID, alertParam)
			}
			return
		}
	}
}

// parseWorkflowType parses the workflow type string and returns the corresponding WorkflowType and WorkflowExecutor
// If the workflow type is unknown, an error is returned
// This function is used internally by the workflow service
func (s *workflowService) parseWorkflowType(typeStr string) (models.WorkflowType, interfaces.WorkflowExecutor, error) {
	var wfType models.WorkflowType
	var executor interfaces.WorkflowExecutor
	var err error
	switch typeStr {
	case "screenshot":
		wfType = models.ScreenshotWorkflowType
		executor = s.screenshotWorkflowExecutor
		err = nil
	case "report":
		wfType = models.ReportWorkflowType
		executor = s.reportWorkflowExecutor
		err = nil
	default:
		err = fmt.Errorf("unknown workflow type: %s", typeStr)
	}
	return wfType, executor, err
}
