// ./src/internal/interfaces/repository/workspace.go
package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// WorkspaceRepository is the interface that provides workspace operations
// This is used to interact with the workspace repository

type WorkspaceRepository interface {
	// CreateWorkspace creates a new workspace
	// The workspace is created with the provided name and billing email
	CreateWorkspace(ctx context.Context, workspaceName, billingEmail string) (models.Workspace, errs.Error)

	// GetWorkspace gets a workspace by its ID
	// This is used to get the workspace details
	GetWorkspaces(ctx context.Context, workspaceID []uuid.UUID) ([]models.Workspace, errs.Error)

	// WorkspaceExists checks if a workspace exists
	// This is optimized for quick lookups over the workspace table
	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, errs.Error)

	// UpdateWorkspaceBillingEmail updates the billing email of the workspace
	// This is used to update the billing email of the workspace
	UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) errs.Error

	// UpdateWorkspaceName updates the name of the workspace
	// This is used to update the name of the workspace
	UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) errs.Error

	// RemoveWorkspaces removes the workspaces by their IDs
	// This is used to remove the workspaces
	RemoveWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) errs.Error

	// UpdateWorkspaceStatus updates the status of the workspace
	// This is used to update the status of the workspace
	UpdateWorkspaceStatus(ctx context.Context, workspaceID uuid.UUID, status models.WorkspaceStatus) errs.Error

	// UpdateWorkspace is used to update the workspace
	// This is used to update the workspace details
	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceReq models.WorkspaceProps) errs.Error
}

var (
    // ErrFailedToCreateWorkspaceInWorkspaceRepository is returned when the workspace creation fails
	ErrFailedToCreateWorkspaceInWorkspaceRepository = errors.New("failed to create workspace in workspace repository")

    // ErrFailedToGetWorkspaceFromWorkspaceRepository is returned when the workspace get fails
    ErrFailedToGetWorkspaceFromWorkspaceRepository = errors.New("failed to get workspace from workspace repository")

    // ErrFailedToCheckIfWorkspaceExistsInWorkspaceRepository is returned when the workspace exists check fails
    ErrFailedToCheckIfWorkspaceExistsInWorkspaceRepository = errors.New("failed to check if workspace exists in workspace repository")

    // ErrFailedToUpdateWorkspaceBillingEmailInWorkspaceRepository is returned when the workspace billing email update fails
    ErrFailedToUpdateWorkspaceBillingEmailInWorkspaceRepository = errors.New("failed to update workspace billing email in workspace repository")

    // ErrFailedToUpdateWorkspaceNameInWorkspaceRepository is returned when the workspace name update fails
    ErrFailedToUpdateWorkspaceNameInWorkspaceRepository = errors.New("failed to update workspace name in workspace repository")

    // ErrFailedToUpdateWorkspaceStatusInWorkspaceRepository is returned when the workspace status update fails
    ErrFailedToUpdateWorkspaceStatusInWorkspaceRepository = errors.New("failed to update workspace status in workspace repository")

    // ErrFailedToUpdateWorkspaceInWorkspaceRepository is returned when the workspace update fails
    ErrFailedToUpdateWorkspaceInWorkspaceRepository = errors.New("failed to update workspace in workspace repository")

    // ErrFailedToRemoveWorkspacesInWorkspaceRepository is returned when the workspace removal fails
    ErrFailedToRemoveWorkspacesInWorkspaceRepository = errors.New("failed to remove workspaces in workspace repository")
)
