// ./src/internal/repository/workflow/interface.go
package workflow

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowRepository represents the repository for managing workflows
// This is used to interact with the workflow repository
type WorkflowRepository interface {
	CheckpointRepository
	StateRepository
}

// CheckpointRepository is the interface that provides checkpoint operations
// This is used by workflow observers to store and retrieve the state of a running workflow
type CheckpointRepository interface {
	// GetState returns the stored state of the workflow from the repository
	GetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) (models.JobState, error)

	// SetState stores the state of the workflow in the repository
	SetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobState models.JobState) error

	// Active returns the list of active jobs in the workflow
	ListActiveJobs(ctx context.Context, workflowType models.WorkflowType) ([]models.Job, error)
}

// StateRepository is the interface that provides state operations
// This is used by workflow scheduler and workflow service to store and retrieve the state of a workflow
type StateRepository interface {
	// StartJob initializes a new workflow in the repository
	StartJob(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) error

	// CompleteJob completes a workflow in the repository
	CompleteJob(ctx context.Context, jobID uuid.UUID, jobContext *models.JobState, workflowType models.WorkflowType) error

	// CancelJob cancels a workflow in the repository
	CancelJob(ctx context.Context, jobID uuid.UUID, jobContext *models.JobState, workflowType models.WorkflowType) error

	// ListRecords returns the list of jobs records in the repository
	ListRecords(ctx context.Context, workflowType *models.WorkflowType, limit, offset *int) ([]models.JobRecord, error)
}
