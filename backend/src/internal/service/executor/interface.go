// ./src/internal/interfaces/executor/workflow.go
package executor

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowObserver represents the executor for managing workflows
type WorkflowObserver interface {
	// Recover recovers any pre-empted workflow jobs
	Recover(ctx context.Context) error

	// Submit submits a new job to the workflow
	Submit(ctx context.Context) (uuid.UUID, error)

	// Status returns the status of the job submitted to the workflow
	Status(ctx context.Context, jobID uuid.UUID) (*models.JobStatus, error)

	// State returns the state of the job submitted to the workflow
	State(ctx context.Context, jobID uuid.UUID) (*models.JobState, error)

	// Get returns the job submitted to the workflow
	Get(ctx context.Context, jobID uuid.UUID) (*models.Job, error)

	// Cancel cancels the job submitted to the workflow
	Cancel(ctx context.Context, jobID uuid.UUID) error

	// List returns the list of jobs filtered by status
	List(ctx context.Context, status models.JobStatus) ([]models.Job, error)

	// Shutdown stops the workflow executor
	Shutdown(ctx context.Context) error

	// History returns the history of job runs
	History(ctx context.Context, limit, offset *int) ([]models.JobRecord, error)
}

// JobExecutor represents the executor for performing jobs
type JobExecutor interface {
	// Execute executes the task
	Execute(ctx context.Context, jobState models.JobState) (<-chan models.JobUpdate, <-chan models.JobError)

	// Terminate terminates the task
	// It handles cleanup and termination of shared resources for jobs
	Terminate(ctx context.Context) error
}
