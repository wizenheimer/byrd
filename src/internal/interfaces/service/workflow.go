package interfaces

import (
	"context"

	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// WorkflowService represents the service for managing workflows
type WorkflowService interface {
	// Initialize initializes the workflow service
	Initialize(ctx context.Context) []error

	// StartWorkflow starts a new workflow
	StartWorkflow(ctx context.Context, wr api_models.WorkflowRequest) (api_models.WorkflowResponse, error)

	// StopWorkflow stops the workflow
	StopWorkflow(ctx context.Context, wr api_models.WorkflowRequest) error

	// GetWorkflow returns the workflow from workflow repository
	GetWorkflow(ctx context.Context, wr api_models.WorkflowRequest) (api_models.WorkflowResponse, error)

	// ListWorkflows returns the list of workflows
	ListWorkflows(ctx context.Context, ws core_models.WorkflowStatus, wt core_models.WorkflowType) ([]api_models.WorkflowResponse, error)
}
