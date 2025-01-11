package workspace

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkspaceRepository is the interface that provides workspace operations
// This is used to interact with the workspace repository

type WorkspaceRepository interface {
	// ------ CRUD operations for workspace ------

	// ---- Create operations for workspace ----

	// CreateWorkspace creates a new workspace.
	// The workspace is created with the provided name and billing email.
	CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID) (models.Workspace, error)

	// ---- Read operations for workspace ----

	// GetWorkspaceByWorkspaceID gets a workspace by its ID.
	// This is used to get the workspace details.
	GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (models.Workspace, error)

	// BatchGetWorkspacesByIDs gets workspaces by their IDs.
	// This is used to get the workspace details.
	BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error)

	// GetWorkspaceUserIDs gets the users of the workspace.
	// This is used to get the users of the workspace.
	GetWorkspaceUserIDs(ctx context.Context, workspaceID uuid.UUID) ([]uuid.UUID, error)

	// GetWorkspaceUserCountByRole gets the count of users by role in the workspace.
	// This is used to get the count of users by role in the workspace.
	GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.UserWorkspaceRole]int, error)

	// GetWorkspaceForUser gets the workspace for a user.
	// This is used to get the workspace for a user.
	GetWorkspaceForUser(ctx context.Context, userID, workspaceID uuid.UUID) (models.Workspace, error)

	// ---- Update operations for workspace ----

	// User Specific Update Operations

	// BatchAddUsersToWorkspace adds a batch of users to workspace
	BatchAddUsersToWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) ([]models.WorkspaceUser, error)

	// AddUserToWorkspace adds a user to a workspace
	AddUserToWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error)

	// BatchRemoveUsersFromWorkspace removes a batch of users from a workspace
	BatchRemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error

	// RemoveUserFromWorkspace removes a user from a workspace
	RemoveUserFromWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error

	// Membership Specific Update Operations

	// UpdateUserRoleForWorkspace updates the role of a user in the workspace
	UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.UserWorkspaceRole) error

	// Workspace Specific Update Operations

	// UpdateWorkspaceBillingEmail updates the billing email of the workspace.
	// This is used to update the billing email of the workspace.
	UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error

	// UpdateWorkspaceName updates the name of the workspace.
	// This is used to update the name of the workspace.
	UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error

	// UpdateWorkspaceDetails updates the details of the workspace.
	// This is used to update the details of the workspace.
	UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error

	// ---- Delete operations for workspace ----

	// DeleteWorkspace deletes a workspace by its ID.
	// This is used to delete the workspace.
	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	// BatchDeleteWorkspaces deletes workspaces by their IDs.
	// This is used to delete the workspaces.
	BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error
}
