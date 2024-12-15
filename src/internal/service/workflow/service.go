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

// StartWorkflow starts a new workflow
func (s *workflowService) StartWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
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

func (s *workflowService) handleWorkflowBackgroundProcessing(
	workflowID models.WorkflowIdentifier,
	workflowKey string,
	executorID uuid.UUID,
	updateChan <-chan models.WorkflowUpdate,
	errorChan <-chan models.WorkflowError,
) {
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
		cancel()
		s.workflowsMu.Lock()
		delete(s.activeWorkflows, workflowKey)
		s.workflowsMu.Unlock()
	}()

	for {
		select {
		case <-s.serviceCtx.Done():
			// Service-wide shutdown
			s.workflowRepo.SetStatus(backgroundCtx, &workflowID, models.WorkflowStatusAborted, nil, nil)
			return

		case err, ok := <-errorChan:
			if !ok {
				return
			}

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

			workflowState.Mutex.Lock()
			workflowState.Status = update.Status
			workflowState.Mutex.Unlock()

			s.workflowRepo.SetStatus(backgroundCtx, &workflowID, update.Status, update.Checkpoint.BatchID, update.Checkpoint.Stage)

			if update.Status == models.WorkflowStatusCompleted {
				alertParam := map[string]string{
					"workflowID": workflowKey,
					"status":     string(update.Status),
				}
				s.alertClient.SendWorkflowCompleted(backgroundCtx, workflowID, alertParam)
			}

		case <-workflowCtx.Done():
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

func (s *workflowService) StopWorkflow(ctx context.Context, workflowID models.WorkflowIdentifier) error {
	workflowKey := workflowID.Serialize(
		models.GetWorkflowPrefixFromWorkflowType(workflowID.Type),
		models.WorkflowStatusRunning,
	)

	s.workflowsMu.RLock()
	workflowState, exists := s.activeWorkflows[workflowKey]
	s.workflowsMu.RUnlock()

	if !exists {
		return fmt.Errorf("workflow not found or not running: %s", workflowKey)
	}

	workflowState.Cancel()
	return nil
}

func (s *workflowService) Shutdown(ctx context.Context) error {
	s.serviceCancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
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

	status, err := s.workflowRepo.GetStatus(ctx, &workflowID)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// ListWorkflows lists the workflows
func (s *workflowService) ListWorkflows(ctx context.Context, status string, scanLimit int) ([]models.WorkflowResponse, int, error) {
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

	var aggregateResp []models.WorkflowResponse
	var aggregateCount int
	for _, workflowStatus := range workflowStatuses {
		resp, count, err := s.workflowRepo.ListWorkflows(ctx, workflowStatus, scanLimit)
		if err != nil {
			return nil, 0, err
		}

		aggregateResp = append(aggregateResp, resp...)
		aggregateCount += count

		s.logger.Debug("ListWorkflows", zap.Any("count", count), zap.Any("workflows", resp))
	}

	return aggregateResp, aggregateCount, nil
}

// RecoverWorkflow recovers existing workflows
func (s *workflowService) RecoverWorkflow(ctx context.Context) error {
	// List all running workflows
	workflows, _, err := s.ListWorkflows(ctx, "running", 100)
	if err != nil {
		return err
	}

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

		// Start the executor with the checkpoint
		updateChan, errorChan := executor.Recover(ctx, &workflowID, &checkpoint)

		// Generate workflow key
		prefix := models.GetWorkflowPrefixFromWorkflowType(workflow.Type)
		workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

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
