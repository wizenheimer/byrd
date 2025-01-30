// ./src/internal/service/executor/observer.go
package executor

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/notification"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// WorkflowObserver represents the observer for the workflow
// This is used to oversee lifecycle of the workflow
type workflowObserver struct {
	// workflowType represents the type of workflow
	workflowType models.WorkflowType

	// repository represents the repository for checkpoint related operations
	repository workflow.WorkflowRepository

	// alertClient represents the alert client for the workflow
	// alertClient alert.AlertClient

	// eventClient represents the event client for the workflow
	// eventClient event.EventClient

	// alertChannel represents the alert channel for the workflow
	alertChannel chan models.Alert

	// eventChannel represents the event channel for the workflow
	eventChannel chan models.Event

	// jobExecutor represents the job executor for the workflow
	// this would be used to execute the jobs in the workflow in a background
	jobExecutor JobExecutor

	// logger represents the logger for the workflow
	logger *logger.Logger

	// activeJobs represents the active jobs in the workflow
	activeJobs sync.Map // map[uuid.UUID]*jobContext
}

func NewWorkflowObserver(
	workflowType models.WorkflowType,
	repository workflow.WorkflowRepository,
	notificationService notification.NotificationService,
	// alertClient alert.AlertClient,
	// eventClient event.EventClient,
	jobExecutor JobExecutor,
	logger *logger.Logger,
) (WorkflowObserver, error) {

	alertChannel, err := notificationService.GetAlertChannel(context.TODO(), 1, 25)
	if err != nil {
		return nil, err
	}

	eventChannel, err := notificationService.GetEventChannel(context.TODO(), 1, 25)
	if err != nil {
		return nil, err
	}

	workflowObserver := &workflowObserver{
		workflowType: workflowType,
		repository:   repository,
		alertChannel: alertChannel,
		eventChannel: eventChannel,
		jobExecutor:  jobExecutor,
		logger: logger.WithFields(
			map[string]interface{}{
				"module": "workflow_executor",
			}),
	}

	return workflowObserver, nil
}

func (e *workflowObserver) Recover(ctx context.Context) error {
	e.logger.Debug("recovering workflows")

	// List the workflows from the repository
	jobs, err := e.repository.ListActiveJobs(ctx, e.workflowType)
	if err != nil {
		e.logger.Error("failed to list workflows", zap.Error(err))
		return err
	}

	// Recover the workflows
	for _, job := range jobs {
		jobContext, executionContext := models.NewJobContextForJob(&job)
		e.activeJobs.Store(job.JobID, jobContext)
		go e.executeJob(executionContext, jobContext)
	}

	return nil
}

func (e *workflowObserver) Submit(ctx context.Context) (uuid.UUID, error) {
	e.logger.Debug("submitting workflow")

	// Create a new job
	job := models.NewJob()

	// Create a new job context
	jobContext, executionContext := models.NewJobContextForJob(job)

	// Start the job in the repository
	if err := e.repository.StartJob(ctx, job.JobID, e.workflowType); err != nil {
		e.logger.Error("failed to persist job status", zap.Error(err))
		return uuid.Nil, err
	}

	// Store the job context in the active jobs
	e.activeJobs.Store(job.JobID, jobContext)

	// Start the job execution
	go e.executeJob(executionContext, jobContext)

	return job.JobID, nil
}

func (e *workflowObserver) executeJob(executionContext context.Context, jobContext *models.JobContext) {
	e.logger.Debug("executing job", zap.Any("job_id", jobContext.JobID))

	// Execute the job
	jobUpdateCh, jobErrorCh := e.jobExecutor.Execute(executionContext, jobContext.JobState)

	// Wait for the job to complete
	for {
		select {
		case <-executionContext.Done():
			e.handleJobCancellation(jobContext)
			return
		case jobUpdate, ok := <-jobUpdateCh:
			if !ok {
				e.handleJobCompletion(executionContext, jobContext)
				return
			}
			e.handleJobUpdate(executionContext, jobContext, jobUpdate)
		case jobError, ok := <-jobErrorCh:
			if !ok {
				e.handleJobCompletion(executionContext, jobContext)
				return
			}
			e.handleJobError(jobContext, &jobError)
		}
	}
}

func (e *workflowObserver) handleJobCancellation(jobContext *models.JobContext) {
	e.logger.Debug("cancelling job", zap.Any("job_id", jobContext.JobID))
	jobContext.HandleCancellation()

	// Refresh remote state
	if err := e.repository.CancelJob(context.Background(), jobContext.JobID, &jobContext.JobState, e.workflowType); err != nil {
		e.logger.Error("failed to persist job cancellation", zap.Error(err))
		return
	}

	// Refresh local state
	e.activeJobs.Delete(jobContext.JobID)
}

func (e *workflowObserver) handleJobCompletion(ctx context.Context, jobContext *models.JobContext) {
	e.logger.Debug("completing job", zap.Any("job_id", jobContext.JobID))

	// Handle job completion
	jobContext.HandleCompletion()

	// Refresh remote state
	if err := e.repository.CompleteJob(ctx, jobContext.JobID, &jobContext.JobState, e.workflowType); err != nil {
		e.logger.Error("failed to persist job completion", zap.Error(err))
		return
	}

	// Refresh local state
	e.activeJobs.Delete(jobContext.JobID)
	// TODO: inject alert client
}

func (e *workflowObserver) handleJobError(jobContext *models.JobContext, jobError *models.JobError) {
	e.logger.Debug("handling job error", zap.Error(jobError.Error))
	jobContext.IncrementFailed(1)

	// Sync the job error with the job context
	jobErrorEvent := models.NewJobErrorEvent(jobContext, jobError)

	// Send the event
	go func() {
		e.eventChannel <- jobErrorEvent
	}()

}

func (e *workflowObserver) handleJobUpdate(ctx context.Context, jobContext *models.JobContext, jobUpdate models.JobUpdate) {
	e.logger.Debug("handling job update", zap.Any("update", jobUpdate))

	if err := e.repository.SetState(ctx, jobContext.JobID, e.workflowType, jobContext.JobState); err != nil {
		e.logger.Error("failed to persist job state", zap.Error(err))
		return
	}

	jobContext.HandleUpdate(&jobUpdate)
}

func (e *workflowObserver) Status(ctx context.Context, jobID uuid.UUID) (*models.JobStatus, error) {
	e.logger.Debug("getting workflow status", zap.Any("job_id", jobID))

	// Get the job from the active jobs
	jobContext, ok := e.activeJobs.Load(jobID)
	if !ok {
		return nil, errors.New("job not found")
	}

	// Get the job status
	status := jobContext.(*models.JobContext).GetStatus()
	return &status, nil
}

func (e *workflowObserver) State(ctx context.Context, jobID uuid.UUID) (*models.JobState, error) {
	e.logger.Debug("getting workflow state", zap.Any("job_id", jobID))

	// Get the job from the active jobs
	jobContext, ok := e.activeJobs.Load(jobID)
	if !ok {
		return nil, errors.New("job not found")
	}

	// Get the job state
	status, checkpoint := jobContext.(*models.JobContext).GetState()
	return &models.JobState{
		Status:     status,
		Checkpoint: checkpoint,
	}, nil
}

func (e *workflowObserver) Get(ctx context.Context, jobID uuid.UUID) (*models.Job, error) {
	e.logger.Debug("getting workflow", zap.Any("job_id", jobID))

	// Get the job from the active jobs
	jobContext, ok := e.activeJobs.Load(jobID)
	if !ok {
		return nil, errors.New("job not found")
	}

	// Get the job
	job := jobContext.(*models.JobContext).Job
	return &job, nil
}

func (e *workflowObserver) Cancel(ctx context.Context, jobID uuid.UUID) error {
	e.logger.Debug("cancelling workflow", zap.Any("job_id", jobID))

	// Get the job from the active jobs
	jobContext, ok := e.activeJobs.Load(jobID)
	if !ok {
		return errors.New("job not found")
	}

	// Cancel the job
	jobContext.(*models.JobContext).HandleCancellation()
	e.activeJobs.Delete(jobID)
	return nil
}

func (e *workflowObserver) List(ctx context.Context, status models.JobStatus) ([]models.Job, error) {
	e.logger.Debug("listing workflows", zap.Any("status", status))

	// List the jobs from the active jobs
	var jobs []models.Job
	e.activeJobs.Range(func(key, value interface{}) bool {
		job := value.(*models.JobContext).Job
		if job.Status == status {
			jobs = append(jobs, job)
		}
		return true
	})

	return jobs, nil
}

func (e *workflowObserver) Shutdown(ctx context.Context) error {
	e.logger.Debug("shutting down workflow")

	// Iterate over the active jobs and cancel them
	e.activeJobs.Range(func(key, value interface{}) bool {
		jobContext := value.(*models.JobContext)
		jobContext.HandleCancellation()
		return true
	})

	// Shutdown the job executor
	if err := e.jobExecutor.Terminate(ctx); err != nil {
		e.logger.Error("failed to shutdown job executor", zap.Error(err))
		return err
	}

	return nil
}

func (e *workflowObserver) History(ctx context.Context, limit, offset *int) ([]models.JobRecord, error) {
	e.logger.Debug("getting workflow history", zap.Any("limit", limit), zap.Any("offset", offset))

	// Get the history from the repository
	jobRecords, err := e.repository.ListRecords(ctx, &e.workflowType, limit, offset)
	if err != nil {
		e.logger.Error("failed to get history", zap.Error(err))
		return nil, err
	}

	return jobRecords, nil
}
