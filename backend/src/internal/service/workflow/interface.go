package workflow

import (
	"context"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowService represents the service for managing workflows
type WorkflowService interface {
	// Initialize initializes the workflow service
	Initialize(ctx context.Context) []error

	// StartWorkflow starts a new workflow
	StartWorkflow(ctx context.Context, wr api.WorkflowRequest) (api.WorkflowResponse, error)

	// StopWorkflow stops the workflow
	StopWorkflow(ctx context.Context, wr api.WorkflowRequest) error

	// GetWorkflow returns the workflow from workflow repository
	GetWorkflow(ctx context.Context, wr api.WorkflowRequest) (api.WorkflowResponse, error)

	// ListWorkflows returns the list of workflows
	ListWorkflows(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]api.WorkflowResponse, error)
}
