package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

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
