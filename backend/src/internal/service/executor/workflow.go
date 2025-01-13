// ./src/internal/service/executor/workflow.go
package executor

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/alert"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type workflowExecutor struct {
	// workflowType represents the type of workflow
	workflowType models.WorkflowType

	// repository represents the repository for managing workflows
	repository workflow.WorkflowRepository

	// alertClient represents the alert client for the workflow
	alertClient alert.WorkflowAlertClient

	// eventClient represents the event client for the workflow
	// eventClient    event.WorkflowEventClient

	// jobExecutor represents the job executor for the workflow
	// this would be used to execute the jobs in the workflow in a background
	jobExecutor JobExecutor

	// logger represents the logger for the workflow
	logger *logger.Logger

	// activeJobs represents the active jobs in the workflow
	activeJobs sync.Map // map[uuid.UUID]*jobContext
}

func NewWorkflowExecutor(
	workflowType models.WorkflowType,
	repository workflow.WorkflowRepository,
	alertClient alert.WorkflowAlertClient,
	// eventClient    event.WorkflowEventClient,
	jobExecutor JobExecutor,
	logger *logger.Logger,
) (WorkflowExecutor, error) {

	workflowExecutor := &workflowExecutor{
		workflowType: workflowType,
		repository:   repository,
		alertClient:  alertClient,
		jobExecutor:  jobExecutor,
		logger: logger.WithFields(
			map[string]interface{}{
				"module": "workflow_executor",
			}),
	}

	return workflowExecutor, nil
}

func (e *workflowExecutor) Recover(ctx context.Context) error {
	e.logger.Debug("recovering workflows")

	// List the workflows from the repository
	jobs, err := e.repository.List(ctx, e.workflowType, models.JobStatusRunning)
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

func (e *workflowExecutor) Submit(ctx context.Context) (uuid.UUID, error) {
	e.logger.Debug("submitting workflow")

	// Create a new job
	job := models.NewJob()

	// Create a new job context
	jobContext, executionContext := models.NewJobContextForJob(job)

	// Persist the job state in the repository
	if err := e.repository.SetState(ctx, job.JobID, e.workflowType, jobContext.JobState); err != nil {
		e.logger.Error("failed to persist job state", zap.Error(err))
		return uuid.Nil, err
	}

	// Store the job context in the active jobs
	e.activeJobs.Store(job.JobID, jobContext)

	// Start the job execution
	go e.executeJob(executionContext, jobContext)

	return job.JobID, nil
}

func (e *workflowExecutor) executeJob(executionContext context.Context, jobContext *models.JobContext) {
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

func (e *workflowExecutor) handleJobCancellation(jobContext *models.JobContext) {
	e.logger.Debug("cancelling job", zap.Any("job_id", jobContext.JobID))
	jobContext.HandleCancellation()
	e.activeJobs.Delete(jobContext.JobID)
}

func (e *workflowExecutor) handleJobCompletion(ctx context.Context, jobContext *models.JobContext) {
	e.logger.Debug("completing job", zap.Any("job_id", jobContext.JobID))

	if err := e.repository.SetState(ctx, jobContext.JobID, e.workflowType, jobContext.JobState); err != nil {
		e.logger.Error("failed to handle job completion", zap.Error(err))
		return
	}
	jobContext.HandleCompletion()
	e.activeJobs.Delete(jobContext.JobID)
	// TODO: inject alert client
}

func (e *workflowExecutor) handleJobError(jobContext *models.JobContext, jobError *models.JobError) {
	e.logger.Debug("handling job error", zap.Error(jobError.Error))
	jobContext.IncrementFailed(1)
	// TODO: inject event client
}

func (e *workflowExecutor) handleJobUpdate(ctx context.Context, jobContext *models.JobContext, jobUpdate models.JobUpdate) {
	e.logger.Debug("handling job update", zap.Any("update", jobUpdate))

	if err := e.repository.SetState(ctx, jobContext.JobID, e.workflowType, jobContext.JobState); err != nil {
		e.logger.Error("failed to persist job state", zap.Error(err))
		return
	}

	jobContext.HandleUpdate(&jobUpdate)
}

func (e *workflowExecutor) Status(ctx context.Context, jobID uuid.UUID) (*models.JobStatus, error) {
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

func (e *workflowExecutor) State(ctx context.Context, jobID uuid.UUID) (*models.JobState, error) {
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

func (e *workflowExecutor) Get(ctx context.Context, jobID uuid.UUID) (*models.Job, error) {
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

func (e *workflowExecutor) Cancel(ctx context.Context, jobID uuid.UUID) error {
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

func (e *workflowExecutor) List(ctx context.Context, status models.JobStatus) ([]models.Job, error) {
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

func (e *workflowExecutor) Shutdown(ctx context.Context) error {
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
