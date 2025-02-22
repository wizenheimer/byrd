package workspace

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

func (ws *workspaceService) CreateWorkspace(ctx context.Context, workspaceCreatorEmail string, pages []models.PageProps, userEmails []string) (*models.Workspace, error) {
	// Step 0: Create workspace along with the owner
	var workspace *models.Workspace
	var workspaceCreator *models.User
	var err error

	// Step 1: Get or create the workspace owner
	workspaceCreator, err = ws.userService.GetOrCreateUser(ctx, workspaceCreatorEmail)
	if err != nil {
		return nil, err
	}

	// Step 2: Check if the user can create a workspace
	canCreate, err := ws.CanCreateWorkspace(ctx, workspaceCreator.ID)
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, errors.New("user cannot create workspace")
	}

	// Step 3: Create the workspace
	workspaceName := utils.GenerateWorkspaceName()

	err = ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		workspace, err = ws.workspaceRepo.CreateWorkspace(ctx, workspaceName, *workspaceCreator.Email, workspaceCreator.ID, models.WorkspaceTrial)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("failed to create workspace")
	}
	// Step 4: Check if the user creation request is valid
	// If not, truncate the list to the maximum allowed limit
	creationLimit, err := workspace.GetMaxCompetitors()
	if err != nil {
		return nil, err
	}
	if len(pages) > creationLimit {
		pages = pages[:creationLimit]
	}

	// Step 5: Add competitors to the workspace
	if competitors, err := ws.competitorService.BatchCreateCompetitorForWorkspace(ctx, workspace.ID, pages); err != nil {
		ws.logger.Error("failed to create competitors for workspace", zap.Error(err), zap.Any("workspaceID", workspace.ID), zap.Any("pages", len(pages)), zap.Any("competitors", len(competitors)))
	} else if len(competitors) != len(pages) {
		ws.logger.Error("failed to create all competitors for workspace", zap.Any("workspaceID", workspace.ID), zap.Any("pages", len(pages)), zap.Any("competitors", len(competitors)))
	}

	// Step 6: Check if the user creation request is valid
	// If not, truncate the list to the maximum allowed limit
	userLimit, err := workspace.GetMaxUsers()
	if err != nil {
		return nil, err
	}
	if len(userEmails) > userLimit {
		userEmails = userEmails[:userLimit]
	}

	// Step 7: Add users to the workspace
	// If the user is the owner, add them as an admin
	workspaceUsers, err := ws.AddUsersToWorkspace(ctx, workspaceCreatorEmail, workspace.ID, userEmails)
	if err != nil {
		ws.logger.Error("failed to add users to workspace", zap.Error(err), zap.Any("workspaceID", workspace.ID), zap.Any("userEmails", len(userEmails)))
	} else if len(workspaceUsers) != len(userEmails) {
		ws.logger.Error("failed to add all users to workspace", zap.Any("workspaceID", workspace.ID), zap.Any("userEmails", len(userEmails)), zap.Any("workspaceUsers", len(workspaceUsers)))
	}

	// Step 8: Return the workspace
	return workspace, nil
}

func (ws *workspaceService) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (*models.Workspace, error) {
	workspace, err := ws.workspaceRepo.GetWorkspaceByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (ws *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, workspaceProps models.WorkspaceProps) error {
	// Check if the workspace name or billing email is being updated
	updatedNameRequiresUpdate := workspaceProps.Name != ""
	billingEmailRequiresUpdate := workspaceProps.BillingEmail != ""

	// Update the workspace
	if updatedNameRequiresUpdate && billingEmailRequiresUpdate {
		// Update both the workspace name and billing email
		return ws.workspaceRepo.UpdateWorkspaceDetails(ctx, workspaceID, workspaceProps.Name, workspaceProps.BillingEmail)
	} else if updatedNameRequiresUpdate {
		// Update the workspace name
		return ws.workspaceRepo.UpdateWorkspaceName(ctx, workspaceID, workspaceProps.Name)
	} else if billingEmailRequiresUpdate {
		// Update the billing email
		return ws.workspaceRepo.UpdateWorkspaceBillingEmail(ctx, workspaceID, workspaceProps.BillingEmail)
	}

	return nil
}

func (ws *workspaceService) UpdateWorkspacePlan(ctx context.Context, workspaceID uuid.UUID, plan models.WorkspacePlan) error {
	return ws.workspaceRepo.UpdateWorkspacePlan(ctx, workspaceID, plan)
}

func (ws *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) (models.WorkspaceStatus, error) {
	err := ws.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		// Step 1: Remove all competitors from the workspace
		err := ws.competitorService.RemoveCompetitorForWorkspace(ctx, workspaceID, nil)
		if err != nil {
			return err
		}

		// Step 2: Remove all users from the workspace
		err = ws.workspaceRepo.DeleteWorkspace(ctx, workspaceID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.WorkspaceActive, err
	}

	// Step 3: Delete the workspace
	return models.WorkspaceInactive, nil
}

func (ws *workspaceService) ListActiveWorkspaces(ctx context.Context, batchSize int, lastWorkspaceID *uuid.UUID) (<-chan []uuid.UUID, <-chan error) {
	workspaceChan := make(chan []uuid.UUID)
	errorChan := make(chan error)

	go func() {
		defer close(workspaceChan)
		defer close(errorChan)

		currentLastID := lastWorkspaceID
		hasMore := true

		for hasMore {
			activeWorkspaces, err := ws.workspaceRepo.ListActiveWorkspaces(ctx, batchSize, currentLastID)
			if err != nil {
				errorChan <- err
				return
			}

			if len(activeWorkspaces.WorkspaceIDs) == 0 {
				// No more workspaces to process
				return
			}

			workspaceChan <- activeWorkspaces.WorkspaceIDs

			// Update cursor for next iteration
			currentLastID = activeWorkspaces.LastSeen
			hasMore = activeWorkspaces.HasMore
		}
	}()

	return workspaceChan, errorChan
}
