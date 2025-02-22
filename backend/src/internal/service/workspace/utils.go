package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"go.uber.org/zap"
)

func (ws *workspaceService) WorkspaceExists(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	return ws.workspaceRepo.WorkspaceExists(ctx, workspaceID)
}

func (ws *workspaceService) WorkspaceCompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	exists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ws *workspaceService) WorkspaceCompetitorPageExists(ctx context.Context, workspaceID, competitorID, pageID uuid.UUID) (bool, error) {
	competitorExists, err := ws.competitorService.CompetitorExists(ctx, workspaceID, competitorID)
	if err != nil {
		return false, err
	}

	if !competitorExists {
		return false, nil
	}

	pageExists, err := ws.competitorService.PageExists(ctx, competitorID, pageID)
	if err != nil {
		return false, err
	}

	if !pageExists {
		return false, nil
	}

	return true, nil
}

func (ws *workspaceService) sendEmail(email models.Email) {
	// Create a context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Important to avoid context leak

	if err := ws.emailClient.Send(ctx, email); err != nil {
		ws.errorRecord.RecordError(ctx, err, zap.Any("subscriberEmails", email.To), zap.Any("emailSubject", email.EmailSubject))
	}
}

// CanCreateWorkspace checks if the user can create a workspace
// based on the user's current workspace count and the maximum workspace limit
func (ws *workspaceService) CanCreateWorkspace(ctx context.Context, userID uuid.UUID) (bool, error) {
	totalIncomingWorkspaces := 1

	currentCount, err := ws.CountUserWorkspaces(ctx, userID)
	if err != nil {
		return false, err
	}

	maxCount, err := models.GetWorkspaceCreationLimit()
	if err != nil {
		return false, err
	}

	// Check if the user can create a workspace
	return currentCount+totalIncomingWorkspaces <= maxCount, nil
}

// CanCreateCompetitor checks if the user can create a competitor
// based on the user's current competitor count and the maximum competitor limit
// it also checks if the user can create a page based on the user's current page count and the maximum page limit
func (ws *workspaceService) CanCreateCompetitor(ctx context.Context, workspaceID uuid.UUID, totalIncomingCompetitors int, totalIncomingPages int) (bool, models.WorkspacePlan, error) {
	canCreatePage, workspacePlan, err := ws.CanCreatePage(ctx, workspaceID, totalIncomingPages)
	if err != nil {
		return false, workspacePlan, err
	}
	if !canCreatePage {
		return false, workspacePlan, errors.New("user cannot create page")
	}

	// Get the workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, models.WorkspaceStarter, err
	}

	// Get the max competitors for the workspace
	competitorLimit, err := workspace.GetMaxCompetitors()
	if err != nil {
		return false, workspace.WorkspacePlan, err
	}

	// Get the current competitor count
	currentCount, err := ws.CountWorkspaceCompetitors(ctx, workspaceID)
	if err != nil {
		return false, workspace.WorkspacePlan, err
	}

	// Check if the user can create a competitor
	return currentCount+totalIncomingCompetitors <= competitorLimit, workspace.WorkspacePlan, nil
}

// CanCreatePage checks if the user can create a page
// based on the user's current page count and the maximum page limit
func (ws *workspaceService) CanCreatePage(ctx context.Context, workspaceID uuid.UUID, totalIncomingPages int) (bool, models.WorkspacePlan, error) {
	// Get the workspace
	workspace, err := ws.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, models.WorkspaceStarter, err
	}

	// Get the max pages for the competitor
	limit, err := workspace.GetMaxPages()
	if err != nil {
		return false, workspace.WorkspacePlan, err
	}

	// Get the current page count
	pageCount, err := ws.CountWorkspacePages(ctx, workspaceID)
	if err != nil {
		return false, workspace.WorkspacePlan, err
	}

	// Check if the user can create a page
	return pageCount+totalIncomingPages <= limit, workspace.WorkspacePlan, nil
}

// CountUserWorkspaces counts the number of workspaces for a user
func (ws *workspaceService) CountUserWorkspaces(ctx context.Context, userID uuid.UUID) (totalCount int, err error) {
	totalCount, err = ws.workspaceRepo.GetWorkspaceCountForUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// CountWorkspaceMembers counts the number of active and pending members for a workspace
func (ws *workspaceService) CountWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) (activeCount, pendingCount int, err error) {
	activeCount, pendingCount, err = ws.workspaceRepo.GetActivePendingMemberCounts(ctx, workspaceID)
	if err != nil {
		return 0, 0, err
	}

	return activeCount, pendingCount, nil
}

// CountWorkspaceCompetitors counts the number of competitors for a workspace
func (ws *workspaceService) CountWorkspaceCompetitors(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	count, err := ws.competitorService.CountCompetitorsForWorkspace(ctx, workspaceID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountWorkspacePages counts the number of pages for a workspace
func (ws *workspaceService) CountWorkspacePages(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	competitors, _, err := ws.ListCompetitorsForWorkspace(ctx, workspaceID, nil, nil)
	if err != nil {
		return 0, err
	}

	competitorsIDs := make([]uuid.UUID, 0)
	for _, competitor := range competitors {
		competitorsIDs = append(competitorsIDs, competitor.ID)
	}

	count, err := ws.competitorService.CountPagesForCompetitors(ctx, competitorsIDs)
	if err != nil {
		return 0, err
	}

	return count, nil
}
