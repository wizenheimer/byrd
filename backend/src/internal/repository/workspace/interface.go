package workspace

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkspaceRepository is the interface that provides workspace operations
// This is used to interact with the workspace repository

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID) (*models.Workspace, error)

	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error)

	BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error)

	GetWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) ([]models.PartialWorkspaceUser, error)

	GetWorkspaceMemberByUserID(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error)

	GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.WorkspaceRole]int, error)

	GetWorkspacesForUserID(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error)

	BatchAddUsersToWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) ([]models.PartialWorkspaceUser, error)

	AddUserToWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (*models.PartialWorkspaceUser, error)

	BatchRemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error

	RemoveUserFromWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error

	UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.WorkspaceRole) error

	UpdateUserMembershipStatusForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, status models.MembershipStatus) error

	UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error

	UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error

	UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error

	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error
}
