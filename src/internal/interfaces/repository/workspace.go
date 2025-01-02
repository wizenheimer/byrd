package interfaces

import (
	"context"

	"github.com/google/uuid"
	api "github.com/wizenheimer/iris/src/internal/models/api"
	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// WorkspaceRepository is the interface that provides workspace operations
// This is used to interact with the workspace repository

type WorkspaceRepository interface {
	// CreateWorkspace creates a new workspace
	// The workspace is created with the provided name and billing email
	CreateWorkspace(ctx context.Context, workspaceName, billingEmail string) (models.Workspace, error)

	// GetWorkspace gets a workspace by its ID
	// This is used to get the workspace details
	GetWorkspaces(ctx context.Context, workspaceID []uuid.UUID) ([]models.Workspace, []error)

	// WorkspaceExists checks if a workspace exists
	// This is optimized for quick lookups over the workspace table
	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	// UpdateWorkspaceBillingEmail updates the billing email of the workspace
	// This is used to update the billing email of the workspace
	UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error

	// UpdateWorkspaceName updates the name of the workspace
	// This is used to update the name of the workspace
	UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error

	// UpdateWorkspaceStatus updates the status of the workspace
	// This is used to update the status of the workspace
	UpdateWorkspaceStatus(ctx context.Context, workspaceID uuid.UUID, status models.WorkspaceStatus) error

	// UpdateWorkspace is used to update the workspace
	// This is used to update the workspace details
	UpdateWorkspace(ctx context.Context, workspaceReq api.WorkspaceUpdateRequest) error
}
