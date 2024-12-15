package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type WorkflowExecutor interface {
	// Start starts a workflow
	Start(ctx context.Context, workflowID *models.WorkflowIdentifier) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError)
	// Stop stops a workflow
	Stop(ctx context.Context, workflowID *models.WorkflowIdentifier, uuid *uuid.UUID) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError)
	// Recover recovers a workflow from a checkpoint
	Recover(ctx context.Context, workflowID *models.WorkflowIdentifier, checkpoint *models.Checkpoint) (<-chan models.WorkflowUpdate, <-chan models.WorkflowError)
	// List lists the workflows
	List() map[string]context.Context
}
