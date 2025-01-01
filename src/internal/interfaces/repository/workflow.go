package interfaces

import (
	"context"

	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// WorkflowRepository represents the repository for managing workflows
type WorkflowRepository interface {
	// GetState returns the stored state of the workflow from the repository
	GetState(ctx context.Context, wi core_models.WorkflowIdentifier) (api_models.WorkflowState, error)

	// SetCheckpoint stores the checkpoint of the workflow in the repository
	SetCheckpoint(ctx context.Context, wi core_models.WorkflowIdentifier, ws core_models.WorkflowStatus, wc core_models.WorkflowCheckpoint) error

	// SetState stores the state of the workflow in the repository
	SetState(ctx context.Context, wi core_models.WorkflowIdentifier, ws api_models.WorkflowState) error

	// List returns the list of workflows from the repository
	List(ctx context.Context, ws core_models.WorkflowStatus, wt core_models.WorkflowType) ([]api_models.WorkflowResponse, error)
}
