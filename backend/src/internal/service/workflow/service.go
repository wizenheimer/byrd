// ./src/internal/service/workflow/service.go
package workflow

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowService struct {
	executors   sync.Map //map[models.WorkflowType]executor.WorkflowObserver
	logger      *logger.Logger
	errorRecord *recorder.ErrorRecorder
	live        atomic.Bool
}

func NewWorkflowService(logger *logger.Logger, errorRecord *recorder.ErrorRecorder) (WorkflowService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	ws := workflowService{
		logger: logger.WithFields(map[string]any{
			"module": "workflow_service",
		}),
		errorRecord: errorRecord,
	}

	return &ws, nil
}

// Initialize initializes the workflow service
// This would start the service and start accepting new jobs
// Additionally, it would recover any pre-empted jobs
func (ws *workflowService) Initialize(ctx context.Context) error {
	// Warn if there are no executors
	count := 0
	ws.executors.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	if count == 0 {
		ws.logger.Warn("no executors found, make sure to add executors before dispatching jobs")
	}

	// Recover any pre-empted workflows
	err := ws.Recover(ctx)
	if err != nil {
		return err
	}

	// Start accepting new jobs
	ws.live.Store(true)
	return nil
}

// Shutdown stops all running jobs and shuts down the workflow service
// This would disable the service from accepting new jobs
func (ws *workflowService) Shutdown(ctx context.Context) error {
	// Stop accepting new jobs
	ws.live.Store(false)

	// Stop all running jobs
	var errs []error
	ws.executors.Range(func(key, value interface{}) bool {
		exc, ok := value.(executor.WorkflowObserver)
		if !ok {
			ws.logger.Error("failed to cast executor to WorkflowObserver during shutodown operation", zap.Any("workflowType", key))
			return false
		}
		if err := exc.Shutdown(ctx); err != nil {
			ws.logger.Error("failed to shutdown executor", zap.Error(err))
			errs = append(errs, err)
		}
		return true
	})

	if len(errs) > 0 {
		return errors.New("failed to shutdown all executors")
	}
	return nil
}

// Recover recovers all pre-empted jobs
// This would be called during the initialization of the service
func (ws *workflowService) Recover(ctx context.Context) error {
	// Stop all running jobs
	var errs []error
	count := 0
	ws.executors.Range(func(key, value interface{}) bool {
		exc, ok := value.(executor.WorkflowObserver)
		if !ok {
			ws.logger.Error("failed to cast executor to WorkflowObserver")
			return false
		}
		if err := exc.Recover(ctx); err != nil {
			ws.logger.Error("failed to recover executor", zap.Error(err))
			errs = append(errs, err)
		} else {
			count++
		}
		return true
	})

	if len(errs) > 0 {
		return errors.New("failed to recover all executors")
	}

	return nil
}

// Register registers a new executor to the workflow service
// This would be called during the initialization of the service
// Raises an error if the executor already exists
func (ws *workflowService) Register(workflowType models.WorkflowType, executor executor.WorkflowObserver) error {
	if _, ok := ws.executors.LoadOrStore(workflowType, executor); ok {
		return errors.New("executor already exists")
	}
	return nil
}

// Submits a new job to the workflow
// This would be called by the client to submit a new job
func (ws *workflowService) Submit(ctx context.Context, workflowType models.WorkflowType) (uuid.UUID, error) {
	if !ws.live.Load() {
		return uuid.Nil, errors.New("service is not live")
	}

	exc, ok := ws.executors.Load(workflowType)
	if !ok {
		return uuid.Nil, errors.New("executor not found")
	}

	jobID, err := exc.(executor.WorkflowObserver).Submit(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return jobID, nil
}

// Stops a running job in the workflow
// This would be called by the client to stop a running job
func (ws *workflowService) Stop(ctx context.Context, workflowType models.WorkflowType, jobID uuid.UUID) error {
	if !ws.live.Load() {
		return errors.New("service is not live")
	}

	exc, ok := ws.executors.Load(workflowType)
	if !ok {
		return errors.New("executor not found")
	}

	executor, ok := exc.(executor.WorkflowObserver)
	if !ok {
		return errors.New("failed to cast executor to WorkflowObserver")
	}

	return executor.Cancel(ctx, jobID)
}

// Gets a running job in the workflow
// This would be called by the client to get the status of a running job
func (ws *workflowService) State(ctx context.Context, workflowType models.WorkflowType, jobID uuid.UUID) (*models.Job, error) {
	exc, ok := ws.executors.Load(workflowType)
	if !ok {
		return nil, errors.New("executor not found")
	}

	job, err := exc.(executor.WorkflowObserver).Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// List returns the list of workflows
// This would be called by the client to get the list of running jobs
func (ws *workflowService) List(ctx context.Context, workflowType models.WorkflowType, jobStatus models.JobStatus) ([]models.Job, error) {
	exc, ok := ws.executors.Load(workflowType)
	if !ok {
		return nil, errors.New("executor not found")
	}

	jobs, err := exc.(executor.WorkflowObserver).List(ctx, jobStatus)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (ws *workflowService) History(ctx context.Context, limit, offset *int, workflowType *models.WorkflowType) ([]models.JobRecord, error) {
	observers := make([]executor.WorkflowObserver, 0)
	if workflowType == nil {
		// Get all executors history
		ws.executors.Range(func(_, value interface{}) bool {
			obs, ok := value.(executor.WorkflowObserver)
			if !ok {
				ws.logger.Error("failed to cast executor to WorkflowObserver during history lookup")
				return false
			}
			observers = append(observers, obs)
			return true
		})
	} else {
		// Get specific executor history
		exc, ok := ws.executors.Load(*workflowType)
		if !ok {
			return nil, errors.New("executor not found")
		}

		obs, ok := exc.(executor.WorkflowObserver)
		if !ok {
			return nil, errors.New("failed to cast executor to WorkflowObserver")
		}
		observers = append(observers, obs)
	}

	var observerHistory []models.JobRecord
	for _, observer := range observers {
		history, err := observer.History(ctx, limit, offset)
		if err != nil {
			return nil, err
		}
		observerHistory = append(observerHistory, history...)
	}

	if observerHistory == nil {
		observerHistory = make([]models.JobRecord, 0)
	}

	return observerHistory, nil
}
