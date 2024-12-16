package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

// WorkflowService represents the service for managing workflows
type WorkflowService interface {
	// Initialize initializes the workflow service
	Initialize(ctx context.Context) []error

	// StartWorkflow starts a new workflow
	StartWorkflow(ctx context.Context, wr models.WorkflowRequest) (models.WorkflowResponse, error)

	// StopWorkflow stops the workflow
	StopWorkflow(ctx context.Context, wr models.WorkflowRequest) error

	// GetWorkflow returns the workflow from workflow repository
	GetWorkflow(ctx context.Context, wr models.WorkflowRequest) (models.WorkflowResponse, error)

	// ListWorkflows returns the list of workflows
	ListWorkflows(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]models.WorkflowResponse, error)
}

// WorkflowRepository represents the repository for managing workflows
type WorkflowRepository interface {
	// GetState returns the stored state of the workflow from the repository
	GetState(ctx context.Context, wi models.WorkflowIdentifier) (models.WorkflowState, error)

	// SetCheckpoint stores the checkpoint of the workflow in the repository
	SetCheckpoint(ctx context.Context, wi models.WorkflowIdentifier, ws models.WorkflowStatus, wc models.WorkflowCheckpoint) error

	// SetState stores the state of the workflow in the repository
	SetState(ctx context.Context, wi models.WorkflowIdentifier, ws models.WorkflowState) error

	// List returns the list of workflows from the repository
	List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]models.WorkflowState, error)
}

// WorkflowAlertClient represents the client for sending alerts for workflows
type WorkflowAlertClient interface {
	// AlertClient represents the client for sending alerts
	AlertClient

	// SendWorkflowStarted sends an alert when a workflow is started
	SendWorkflowStarted(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowCompleted sends an alert when a workflow is completed
	SendWorkflowCompleted(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowFailed sends an alert when a workflow fails
	SendWorkflowFailed(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowCancelled sends an alert when a workflow is cancelled
	SendWorkflowCancelled(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error
}

// WorkflowExecutor represents the executor for managing workflows
type WorkflowExecutor interface {
	// Initialize initializes the workflow executor
	Initialize(ctx context.Context) error

	// Start starts the workflow
	Start(ctx context.Context, wi models.WorkflowIdentifier) error

	// Stop stops the workflow
	Stop(ctx context.Context, wi models.WorkflowIdentifier) error

	// List returns the list of workflows
	List(ctx context.Context, ws models.WorkflowStatus, wt models.WorkflowType) ([]models.WorkflowState, error)

	// GetState returns the stored state of the workflow
	Get(ctx context.Context, wi models.WorkflowIdentifier) (models.WorkflowState, error)
}

// TaskExecutor represents the executor for managing tasks
type TaskExecutor interface {
	// Execute executes the task
	Execute(ctx context.Context, t models.Task) (<-chan models.TaskUpdate, <-chan models.TaskError)
	// Terminate terminates the task
	Terminate(ctx context.Context) error
}
