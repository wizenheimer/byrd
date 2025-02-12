// ./src/internal/service/workspace/interface.go
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

	ListWorkspacesForUser(ctx context.Context, workspaceMember *clerk.User, membershipStatus *models.MembershipStatus, limit, offset *int) ([]models.WorkspaceWithMembership, bool, error)

	CountUserWorkspaces(context.Context, uuid.UUID) (int, error)

	CountWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) (activeCount, pendingCount int, err error)

	CountWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID) (int, error)

	CountWorkspacePages(ctx context.Context, workspaceID uuid.UUID) (int, error)

	GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error)

	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceProps models.WorkspaceProps) error

	UpdateWorkspacePlan(ctx context.Context, workspaceID uuid.UUID, plan models.WorkspacePlan) error

	DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error)

	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, limit, offset *int, roleFilter *models.WorkspaceRole) ([]models.WorkspaceUser, bool, error)

	AddUsersToWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error)

	AddSlackUserToWorkspace(ctx context.Context, workspaceMember string, workspaceID uuid.UUID, emails []string) ([]models.WorkspaceUser, error)

	LeaveWorkspace(ctx context.Context, workspaceMember *clerk.User, workspaceID uuid.UUID) error

	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID, role models.WorkspaceRole) error

	GetCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Competitor, error)

	UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error)

	RemoveUserFromWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceMemberID uuid.UUID) error

	JoinWorkspace(ctx context.Context, invitedMember *clerk.User, workspaceID uuid.UUID) error

	WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	GetClerkWorkspaceUser(ctx context.Context, workspaceID uuid.UUID, clerkUser *clerk.User) (*models.PartialWorkspaceUser, error)

	WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error)

	WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error)

	// AddCompetitorToWorkspace adds a competitor to a workspace
	// It creates a single competitor for a workspace using multiple pages
	AddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (*models.Competitor, error)

	// BatchAddCompetitorToWorkspace adds multiple competitors to a workspace
	// It flattens the pages and creates a competitor for each page
	BatchAddCompetitorToWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error)

	AddPageToCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error)

	ListPagesForCompetitor(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error)

	ListHistoryForPage(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error)

	RemovePageFromWorkspace(ctx context.Context, competitorID, pageID uuid.UUID) error

	RemoveCompetitorFromWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID) error

	UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)

	GetPageForCompetitor(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error)

	// ListReports lists the reports for a competitor.
	ListReports(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error)

	// CreateReport creates a report for a competitor.
	CreateReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) (*models.Report, error)

	// DispatchReport dispatches a report for a competitor.
	DispatchReportToWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) error

	// DispatchReport dispatches a report for a competitor to an email list.
	DispatchReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID, subscriberEmails []string) error

	ListActiveWorkspaces(ctx context.Context, batchSize int, lastWorkspaceID *uuid.UUID) (<-chan []uuid.UUID, <-chan error)

	CanAddUsers(ctx context.Context, workspaceID uuid.UUID, totalIncomingUsers int) (bool, models.WorkspacePlan, error)

	CanCreateCompetitor(ctx context.Context, workspaceID uuid.UUID, totalIncomingCompetitors int, totalIncomingPages int) (bool, models.WorkspacePlan, error)

	CanCreatePage(ctx context.Context, workspaceID uuid.UUID, totalIncomingPages int) (bool, models.WorkspacePlan, error)
}
