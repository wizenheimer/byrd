// ./src/internal/interfaces/executor/workflow.go
package executor

import (
	"context"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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
	List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]api.WorkflowState, error)

	// GetState returns the stored state of the workflow
	Get(ctx context.Context, wi models.WorkflowIdentifier) (api.WorkflowState, error)
}
