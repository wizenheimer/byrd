package workflow

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
)

// WorkflowService represents the service for managing workflows
type WorkflowService interface {
	// Initialize initializes the workflow service
	// This would start the service and start accepting new jobs
	// Additionally, it would recover any pre-empted jobs
	Initialize(ctx context.Context) error

	// Shutdown stops all running jobs and shuts down the workflow service
	// This would disable the service from accepting new jobs
	Shutdown(ctx context.Context) error

	// Recover recovers all pre-empted jobs
	// This would be called during the initialization of the service
	Recover(ctx context.Context) error

	// AddExecutor registers a new executor to the workflow service
	// This would be called during the initialization of the service
	// Raises an error if the executor already exists
	AddExecutor(workflowType models.WorkflowType, executor executor.WorkflowExecutor) error

	// Submits a new job to the workflow
	// This would be called by the client to submit a new job
	Submit(ctx context.Context, workflowType models.WorkflowType) (uuid.UUID, error)

	// Stops a running job in the workflow
	// This would be called by the client to stop a running job
	Stop(ctx context.Context, workflowType models.WorkflowType, jobID uuid.UUID) error

	// Gets a running job in the workflow
	// This would be called by the client to get the status of a running job
	State(ctx context.Context, workflowType models.WorkflowType, jobID uuid.UUID) (*models.Job, error)

	// List returns the list of workflows
	// This would be called by the client to get the list of running jobs
	List(ctx context.Context, workflowType models.WorkflowType, jobStatus models.JobStatus) ([]models.Job, error)
}
