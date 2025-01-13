package workflow

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowRepository represents the repository for managing workflows
// This is used to interact with the workflow repository

type WorkflowRepository interface {
	// InitializeWorkflow initializes a new workflow in the repository
	Initialize(ctx context.Context, jobProps models.Job, workflowType models.WorkflowType) error

	// GetState returns the stored state of the workflow from the repository
	GetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType) (models.JobState, error)

	// SetState stores the state of the workflow in the repository
	SetState(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobState models.JobState) error

	// SetCheckpoint stores the checkpoint of the workflow in the repository
	SetCheckpoint(ctx context.Context, jobID uuid.UUID, workflowType models.WorkflowType, jobCheckpoint models.JobCheckpoint) error

	// List returns the list of workflows from the repository
	List(ctx context.Context, workflowType models.WorkflowType, jobStatus models.JobStatus) ([]models.Job, error)
}
