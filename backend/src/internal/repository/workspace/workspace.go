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

func (r *workspaceRepo) CreateWorkspace(ctx context.Context, workspaceName, billingEmail string, workspaceCreatorUserID uuid.UUID) (*models.Workspace, error) {
	return nil, nil
}

func (r *workspaceRepo) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	return false, nil
}

func (r *workspaceRepo) GetWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	return nil, nil
}

func (r *workspaceRepo) BatchGetWorkspacesByIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Workspace, error) {
	return nil, nil
}

func (r *workspaceRepo) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, workspaceRole *models.WorkspaceRole) ([]models.PartialWorkspaceUser, bool, error) {
	return nil, false, nil
}

func (r *workspaceRepo) GetWorkspaceMemberByUserID(ctx context.Context, workspaceID, userID uuid.UUID) (*models.PartialWorkspaceUser, error) {
	return nil, nil
}

func (r *workspaceRepo) GetWorkspaceUserCountByRole(ctx context.Context, workspaceID uuid.UUID) (map[models.WorkspaceRole]int, error) {
	return nil, nil
}

func (r *workspaceRepo) GetWorkspacesForUserID(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error) {
	return nil, nil
}

func (r *workspaceRepo) BatchAddUsersToWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) ([]models.PartialWorkspaceUser, error) {
	return nil, nil
}

func (r *workspaceRepo) AddUserToWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (*models.PartialWorkspaceUser, error) {
	return nil, nil
}

func (r *workspaceRepo) BatchRemoveUsersFromWorkspace(ctx context.Context, userIDs []uuid.UUID, workspaceID uuid.UUID) error {
	return nil
}

func (r *workspaceRepo) RemoveUserFromWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error {
	return nil
}

func (r *workspaceRepo) UpdateUserRoleForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role models.WorkspaceRole) error {
	return nil
}

func (r *workspaceRepo) UpdateUserMembershipStatusForWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, membershipStatus models.MembershipStatus) error {
	return nil
}

func (r *workspaceRepo) UpdateWorkspaceBillingEmail(ctx context.Context, workspaceID uuid.UUID, billingEmail string) error {
	return nil
}

func (r *workspaceRepo) UpdateWorkspaceName(ctx context.Context, workspaceID uuid.UUID, workspaceName string) error {
	return nil
}

func (r *workspaceRepo) UpdateWorkspaceDetails(ctx context.Context, workspaceID uuid.UUID, workspaceName, billingEmail string) error {
	return nil
}

func (r *workspaceRepo) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	return nil
}

func (r *workspaceRepo) BatchDeleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error {
	return nil
}
