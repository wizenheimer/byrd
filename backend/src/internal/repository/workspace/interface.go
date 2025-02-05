// ./src/internal/repository/workspace/interface.go
package workspace

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkspaceRepository is the interface that provides workspace operations
// This is used to interact with the workspace repository

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID, workspacePlan models.WorkspacePlan) (*models.Workspace, error)

	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error)

	BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error)

	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, workspaceRole *models.WorkspaceRole) ([]models.PartialWorkspaceUser, bool, error)

	GetWorkspaceMemberByUserID(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error)

	GetWorkspaceUserCountsByRoleAndStatus(ctx context.Context, workspaceID uuid.UUID) (int, int, int, int, error)

	PromoteRandomUserToAdmin(ctx context.Context, workspaceID uuid.UUID) error

	ListWorkspacesForUser(ctx context.Context, userID uuid.UUID, membershipStatus *models.MembershipStatus, limit, offset *int) ([]models.Workspace, bool, error)

	BatchAddUsersToWorkspace(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID) ([]models.PartialWorkspaceUser, error)

	AddUserToWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error)

	BatchRemoveUsersFromWorkspace(ctx context.Context, workspaceID uuid.UUID, userIDs []uuid.UUID) error

	RemoveUserFromWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) error

	UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.WorkspaceRole) error

	UpdateUserMembershipStatusForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, status models.MembershipStatus) error

	UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error

	UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error

	UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error

	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error

	ListActiveWorkspaces(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActiveWorkspaceBatch, error)

	// UpdateWorkspacePlan updates the plan of a workspace
	UpdateWorkspacePlan(ctx context.Context, workspaceID uuid.UUID, plan models.WorkspacePlan) error

	// GetWorkspaceCountForUser returns the total number of workspaces for a user
	GetWorkspaceCountForUser(ctx context.Context, userID uuid.UUID) (int, error)

	// GetTotalActiveAndPendingMembers returns the total number of active and pending members in a workspace
	GetActivePendingMemberCounts(ctx context.Context, workspaceID uuid.UUID) (activeCount, pendingCount int, err error)
}
