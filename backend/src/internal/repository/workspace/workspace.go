package workspace

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewWorkspaceRepository(tm *transaction.TxManager, logger *logger.Logger) WorkspaceRepository {
	return &workspaceRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "workspace_repository"}),
	}
}

// CreateWorkspace creates a new workspace.
// The workspace is created with the provided name and billing email.
func (wr *workspaceRepo) CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID) (models.Workspace, error) {
	return models.Workspace{}, nil
}

// ---- Read operations for workspace ----

// GetWorkspaceByWorkspaceID gets a workspace by its ID.
// This is used to get the workspace details.
func (wr *workspaceRepo) GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (models.Workspace, error) {
	return models.Workspace{}, nil
}

// BatchGetWorkspacesByIDs gets workspaces by their IDs.
// This is used to get the workspace details.
func (wr *workspaceRepo) BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error) {
	return []models.Workspace{}, nil
}

// GetWorkspaceUserIDs gets the users of the workspace.
// This is used to get the users of the workspace.
func (wr *workspaceRepo) GetWorkspaceUserIDs(ctx context.Context, workspaceID uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

// GetWorkspaceUserCountByRole gets the count of users by role in the workspace.
// This is used to get the count of users by role in the workspace.
func (wr *workspaceRepo) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.UserWorkspaceRole]int, error) {
	return map[models.UserWorkspaceRole]int{}, nil
}

// GetWorkspaceForUser gets the workspace for a user.
// This is used to get the workspace for a user.
func (wr *workspaceRepo) GetWorkspaceForUser(ctx context.Context, userID, workspaceID uuid.UUID) (models.Workspace, error) {
	return models.Workspace{}, nil
}

// ---- Update operations for workspace ----

// User Specific Update Operations

// BatchAddUsersToWorkspace adds a batch of users to workspace
func (wr *workspaceRepo) BatchAddUsersToWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) ([]models.WorkspaceUser, error) {
	return []models.WorkspaceUser{}, nil
}

// AddUserToWorkspace adds a user to a workspace
func (wr *workspaceRepo) AddUserToWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (models.WorkspaceUser, error) {
	return models.WorkspaceUser{}, nil
}

// BatchRemoveUsersFromWorkspace removes a batch of users from a workspace
func (wr *workspaceRepo) BatchRemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error {
	return nil
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wr *workspaceRepo) RemoveUserFromWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error {
	return nil
}

// Membership Specific Update Operations

// UpdateUserRoleForWorkspace updates the role of a user in the workspace
func (wr *workspaceRepo) UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.UserWorkspaceRole) error {
	return nil
}

// Workspace Specific Update Operations

// UpdateWorkspaceBillingEmail updates the billing email of the workspace.
// This is used to update the billing email of the workspace.
func (wr *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error {
	return nil
}

// UpdateWorkspaceName updates the name of the workspace.
// This is used to update the name of the workspace.
func (wr *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error {
	return nil
}

// UpdateWorkspaceDetails updates the details of the workspace.
// This is used to update the details of the workspace.
func (wr *workspaceRepo) UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error {
	return nil
}

// ---- Delete operations for workspace ----

// DeleteWorkspace deletes a workspace by its ID.
// This is used to delete the workspace.
func (wr *workspaceRepo) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	return nil
}

// BatchDeleteWorkspaces deletes workspaces by their IDs.
// This is used to delete the workspaces.
func (wr *workspaceRepo) BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error {
	return nil
}
