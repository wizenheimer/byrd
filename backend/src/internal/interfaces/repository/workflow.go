package interfaces

import (
	"context"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowRepository represents the repository for managing workflows
// This is used to interact with the workflow repository

type WorkflowRepository interface {
	// InitializeWorkflow initializes a new workflow in the repository
	InitializeWorkflow(ctx context.Context, wi models.WorkflowIdentifier) error

	// GetState returns the stored state of the workflow from the repository
	GetState(ctx context.Context, wi models.WorkflowIdentifier) (api.WorkflowState, error)

	// SetCheckpoint stores the checkpoint of the workflow in the repository
	SetCheckpoint(ctx context.Context, wi models.WorkflowIdentifier, ws models.WorkflowStatus, wc models.WorkflowCheckpoint) error

	// SetState stores the state of the workflow in the repository
	SetState(ctx context.Context, wi models.WorkflowIdentifier, ws api.WorkflowState) error

	// List returns the list of workflows from the repository
	List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]api.WorkflowResponse, error)
}
