package workspace

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

func (ws *workspaceService) ListReports(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	return ws.competitorService.ListReports(ctx, workspaceID, competitorID, limit, offset)
}

func (ws *workspaceService) CreateReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) (*models.Report, error) {
	return ws.competitorService.CreateReport(ctx, workspaceID, competitorID)
}

func (ws *workspaceService) DispatchReportToWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) error {
	members, hasMore, err := ws.workspaceRepo.ListWorkspaceMembers(ctx, workspaceID, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return errors.New("no members found")
	}
	if hasMore {
		ws.logger.Warn("more members found", zap.Any("workspaceID", workspaceID))
	}

	userIDs := make([]uuid.UUID, 0)
	for _, member := range members {
		userIDs = append(userIDs, member.ID)
	}

	users, err := ws.userService.ListUsersByUserIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	subscriberEmails := make([]string, 0)
	for _, user := range users {
		if user.Email != nil {
			subscriberEmails = append(subscriberEmails, *user.Email)
		}
	}

	return ws.competitorService.DispatchReport(ctx, workspaceID, competitorID, subscriberEmails)
}

func (ws *workspaceService) DispatchReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID, subscriberEmails []string) error {
	// Clean up the email list for duplicates and nil values
	subscriberEmails = utils.CleanEmailList(subscriberEmails, nil)

	// Dispatch the report
	return ws.competitorService.DispatchReport(ctx, workspaceID, competitorID, subscriberEmails)
}
