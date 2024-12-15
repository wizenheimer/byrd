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
) (interfaces.WorkflowService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the workflow service
	service := workflowService{
		workflowRepo:               workflowRepo,
		alertClient:                alertClient,
		logger:                     logger.WithFields(map[string]interface{}{"module": "workflow_service"}),
		screenshotWorkflowExecutor: screenshotExecutor,
		reportWorkflowExecutor:     reportExecutor,
		activeWorkflows:            make(map[string]*models.WorkflowState),
		serviceCtx:                 ctx,
		serviceCancel:              cancel,
	}

	// Recover existing workflows
	if err := service.RecoverWorkflow(ctx); err != nil {
		service.logger.Error("failed to recover workflows", zap.Error(err))
		return nil, err
	}

	return &service, nil
}

// StartWorkflow starts a new workflow
func (s *workflowService) StartWorkflow(ctx context.Context, req models.WorkflowRequest) (*models.WorkflowResponse, error) {
	// Validate the request
	if err := req.Validate(true); err != nil {
		return nil, err
	}

	s.logger.Debug("starting workflow", zap.Any("workflow_type", req.Type), zap.Any("year", req.Year), zap.Any("week_number", *req.WeekNumber), zap.Any("bucket_number", *req.BucketNumber))

	// Parse the workflow type, and get the corresponding executor
	workflowType, err := s.parseWorkflowType(req.Type)
	if err != nil {
		return nil, err
	}

	executor, err := s.getExecutorForType(workflowType)
	if err != nil {
		return nil, err
	}

	// Create workflow identifier using the parsed workflow type
	workflowID := models.WorkflowIdentifier{
		Type:         workflowType,
		Year:         *req.Year,
		WeekNumber:   *req.WeekNumber,
		BucketNumber: *req.BucketNumber,
	}

	// Generate workflow key
	prefix := workflowType.Prefix()
	workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

	s.workflowsMu.RLock()
	_, exists := s.activeWorkflows[workflowKey]
	s.workflowsMu.RUnlock()

	// Check if the workflow is already running
	if exists {
		return nil, fmt.Errorf("workflow already running: %s", workflowKey)
	}

	// Check if the workflow is already started
	status, err := s.workflowRepo.GetStatus(ctx, &workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow status: %w", err)
	}

	if status.Status == models.WorkflowStatusRunning {
		return nil, fmt.Errorf("workflow already running: %s", workflowKey)
	}

	// Start background processing
	go s.handleWorkflowBackgroundProcessing(
		workflowID,
		executor,
	)

	return &models.WorkflowResponse{
		Type:         workflowType,
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
		workflowID.Type.Prefix(),
		models.WorkflowStatusRunning,
	)

	s.logger.Debug("stopping workflow", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey))

	s.workflowsMu.RLock()
	workflowState, exists := s.activeWorkflows[workflowKey]
	s.workflowsMu.RUnlock()

	if !exists {
		return fmt.Errorf("workflow not found or not running: %s", workflowKey)
	} else {
		workflowState.Mutex.Lock()
		delete(s.activeWorkflows, workflowKey)
		s.logger.Debug("deleted workflow from local state", zap.Any("workflow_id", workflowID), zap.Any("workflow_key", workflowKey))
		workflowState.Mutex.Unlock()
	}

	workflowState.Cancel()

	// Set the status of the workflow to aborted
	if err := s.workflowRepo.StopWorkflow(ctx, &workflowID); err != nil {
		s.logger.Error("failed to stop workflow", zap.Any("workflow_id", workflowID), zap.Error(err))
	}

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
func (s *workflowService) ListWorkflows(ctx context.Context, workflowStatus models.WorkflowStatus, workflowType models.WorkflowType, scanLimit int) ([]models.WorkflowResponse, int, error) {
	s.logger.Debug("listing workflows", zap.Any("status", workflowStatus), zap.Any("scan_limit", scanLimit))

	var aggregateResp []models.WorkflowResponse
	resp, count, err := s.workflowRepo.ListWorkflows(ctx, workflowStatus, workflowType, scanLimit)
	if err != nil {
		return nil, 0, err
	}

	aggregateResp = append(aggregateResp, resp...)

	return aggregateResp, count, nil
}

// Helper to collect all running workflows
func (s *workflowService) collectRunningWorkflows(ctx context.Context) ([]models.WorkflowResponse, error) {
	var workflows []models.WorkflowResponse

	// Collect screenshot workflows
	screenshotWorkflows, count, err := s.ListWorkflows(ctx, "running", models.ScreenshotWorkflowType, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to list screenshot workflows: %w", err)
	}
	s.logger.Debug("found screenshot workflows", zap.Int("count", count))
	workflows = append(workflows, screenshotWorkflows...)

	// Collect report workflows
	reportWorkflows, count, err := s.ListWorkflows(ctx, "running", models.ReportWorkflowType, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to list report workflows: %w", err)
	}
	s.logger.Debug("found report workflows", zap.Int("count", count))
	workflows = append(workflows, reportWorkflows...)

	return workflows, nil
}

// Helper to get executor for workflow type
func (s *workflowService) getExecutorForType(wfType models.WorkflowType) (interfaces.WorkflowExecutor, error) {
	switch wfType {
	case models.ScreenshotWorkflowType:
		return s.screenshotWorkflowExecutor, nil
	case models.ReportWorkflowType:
		return s.reportWorkflowExecutor, nil
	default:
		return nil, fmt.Errorf("unknown workflow type: %s", wfType)
	}
}

// RecoverWorkflow recovers existing workflows
func (s *workflowService) RecoverWorkflow(ctx context.Context) error {
	workflows, err := s.collectRunningWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to collect running workflows: %w", err)
	}

	s.logger.Debug("recovering workflows", zap.Int("total_count", len(workflows)))
	// Recover each workflow
	for _, workflow := range workflows {
		executor, err := s.getExecutorForType(workflow.Type)
		if err != nil {
			s.logger.Error("skipping workflow recovery - unknown type",
				zap.String("type", string(workflow.Type)),
				zap.Error(err))
			continue
		}

		workflowID := models.WorkflowIdentifier{
			Type:         workflow.Type,
			Year:         workflow.Year,
			WeekNumber:   workflow.WeekNumber,
			BucketNumber: workflow.BucketNumber,
		}

		// Create a checkpoint
		checkpoint := models.Checkpoint{
			BatchID: workflow.BatchID,
			Stage:   workflow.Stage,
		}

		// Set the active workflows
		workflowKey := workflowID.Serialize(workflow.Type.Prefix(), models.WorkflowStatusRunning)
		s.workflowsMu.Lock()
		s.activeWorkflows[workflowKey] = &models.WorkflowState{
			Status: models.WorkflowStatusRunning,
		}
		s.workflowsMu.Unlock()

		// Start recovery background processing
		go s.handleWorkflowRecoveryProcessing(
			workflowID,
			executor,
			checkpoint,
		)

		s.logger.Info("started workflow recovery",
			zap.String("type", string(workflow.Type)),
			zap.Any("id", workflowID),
			zap.Any("checkpoint", checkpoint))
	}

	return nil
}

func (s *workflowService) handleWorkflowRecoveryProcessing(
	workflowID models.WorkflowIdentifier,
	executor interfaces.WorkflowExecutor,
	checkpoint models.Checkpoint,
) {
	// Create background context
	backgroundCtx := context.Background()
	workflowCtx, cancel := context.WithCancel(backgroundCtx)

	// Initialize workflow state
	executorID := uuid.New()
	workflowState := &models.WorkflowState{
		Cancel:     cancel,
		ExecutorID: executorID,
		Status:     models.WorkflowStatusRunning,
	}

	// Generate workflow key
	prefix := workflowID.Type.Prefix()
	workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

	// Register workflow
	s.workflowsMu.Lock()
	s.activeWorkflows[workflowKey] = workflowState
	s.workflowsMu.Unlock()

	// Start the executor with recovery checkpoint
	updateChan, errorChan := executor.Recover(workflowCtx, &workflowID, &checkpoint)

	// Send alert for workflow recovery
	s.alertClient.SendWorkflowRestarted(backgroundCtx, workflowID, map[string]string{
		"workflowID": workflowKey,
		"timestamp":  fmt.Sprintf("%v", time.Now()),
	})

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

// handleWorkflowBackgroundProcessing is a background processing loop for a workflow
func (s *workflowService) handleWorkflowBackgroundProcessing(
	workflowID models.WorkflowIdentifier,
	executor interfaces.WorkflowExecutor,
) {
	// Create background context
	backgroundCtx := context.Background()
	workflowCtx, cancel := context.WithCancel(backgroundCtx)

	// Initialize workflow state
	executorID := uuid.New()
	workflowState := &models.WorkflowState{
		Cancel:     cancel,
		ExecutorID: executorID,
		Status:     models.WorkflowStatusRunning,
	}

	// Generate workflow key
	prefix := workflowID.Type.Prefix()
	workflowKey := workflowID.Serialize(prefix, models.WorkflowStatusRunning)

	// Register workflow
	s.workflowsMu.Lock()
	s.activeWorkflows[workflowKey] = workflowState
	s.workflowsMu.Unlock()

	// Start the executor with the background context
	updateChan, errorChan := executor.Start(workflowCtx, &workflowID)

	// Set the status of the workflow to running
	if err := s.workflowRepo.SetStatus(backgroundCtx, &workflowID, models.WorkflowStatusRunning, nil, nil); err != nil {
		s.logger.Error("failed to set workflow status", zap.Any("workflow_id", workflowID), zap.Error(err))
	}

	s.alertClient.SendWorkflowStarted(backgroundCtx, workflowID, map[string]string{
		"workflowID": workflowKey,
		"timestamp":  fmt.Sprintf("%v", time.Now()),
	})

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

			s.alertClient.SendWorkflowCancelled(backgroundCtx, workflowID, map[string]string{
				"workflow_key": workflowKey,
				"status":       string(models.WorkflowStatusAborted),
			})
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
				"worflow_key":  workflowKey,
				"worflow_type": string(workflowID.Type),
				"worflow_id":   fmt.Sprintf("%v", workflowID),
				"executor_id":  fmt.Sprintf("%v", executorID),
				"error":        err.Error.Error(),
				"timestamp":    fmt.Sprintf("%v", err.Timestamp),
			}
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
					"worflow_key": workflowKey,
					"status":      string(models.WorkflowStatusAborted),
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
func (s *workflowService) parseWorkflowType(typeStr string) (models.WorkflowType, error) {
	var wfType models.WorkflowType
	var err error
	switch typeStr {
	case "screenshot":
		wfType = models.ScreenshotWorkflowType
		err = nil
	case "report":
		wfType = models.ReportWorkflowType
		err = nil
	default:
		err = fmt.Errorf("unknown workflow type: %s", typeStr)
	}
	return wfType, err
}
