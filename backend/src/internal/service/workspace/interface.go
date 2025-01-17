package workspace

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// Note: WorkspaceService is user-facing and handler-owned service
// It's an entry point for all workspace-related operations

// It embeds other services and repositories
// Repository:
// - WorkspaceRepository
// Service:
// - UserService
// - CompetitorService (embeds PageService)
// - BillingService

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, workspaceOwner *clerk.User, pages []models.PageProps, userEmails []string) (*models.Workspace, error)

	ListUserWorkspaces(ctx context.Context, workspaceMember *clerk.User) ([]models.Workspace, error)

	GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error)

	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceProps models.WorkspaceProps) error

	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error)

	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, roleFilter *models.WorkspaceRole) ([]models.WorkspaceUser, bool, error)

	AddUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error)

	LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error

	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.WorkspaceRole) error

	RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error

	JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error

	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	ClerkUserIsWorkspaceAdmin(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

	ClerkUserIsWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

  ClerkUserIsActiveWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

  ClerkUserIsPendingWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (bool, error)

	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)

	WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error)

	AddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (*models.Competitor, error)

	AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error)

	ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error)

	ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error)

	RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error

	RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error

	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)
}
