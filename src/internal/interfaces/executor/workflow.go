package executor

import (
	"context"

	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// WorkflowExecutor represents the executor for managing workflows
type WorkflowExecutor interface {
	// Initialize initializes the workflow executor
	Initialize(ctx context.Context) error

	// Start starts the workflow
	Start(ctx context.Context, wi models.WorkflowIdentifier) error

	// Stop stops the workflow
	Stop(ctx context.Context, wi models.WorkflowIdentifier) error

	// Restart restarts the workflow
	Restart(ctx context.Context, workflowID models.WorkflowIdentifier, errChan chan error)

	// List returns the list of workflows
	List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]api_models.WorkflowState, error)

	// GetState returns the stored state of the workflow
	Get(ctx context.Context, wi models.WorkflowIdentifier) (api_models.WorkflowState, error)
}
