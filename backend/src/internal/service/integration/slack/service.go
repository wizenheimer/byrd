package slackworkspace

import (
	"context"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/integration/slack"
	repository "github.com/wizenheimer/byrd/src/internal/repository/integration/slack"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type slackWorkspaceService struct {
	// repo is the repository for Slack workspace data
	repo repository.SlackWorkspaceRepository

	// ws is the workspace service for managing Byrd workspaces
	ws workspace.WorkspaceService

	// logger is the logger for the Slack workspace service
	logger *logger.Logger
}

// NewSlackWorkspaceService creates a new Slack workspace service
func NewSlackWorkspaceService(repo repository.SlackWorkspaceRepository, ws workspace.WorkspaceService, logger *logger.Logger) (SlackWorkspaceService, error) {
	svc := slackWorkspaceService{
		repo:   repo,
		ws:     ws,
		logger: logger,
	}
	return &svc, nil
}

// Creates and associates an existing Byrd workspace with a Slack workspace
func (svc *slackWorkspaceService) CreateSlackWorkspace(ctx context.Context, workspaceID uuid.UUID, teamID, accessToken string) (*models.SlackWorkspace, error) {
	return svc.repo.CreateSlackWorkspace(ctx, workspaceID, teamID, accessToken)
}

// UpdateSlackWorkspace updates the access token for a Slack workspace
func (svc *slackWorkspaceService) UpdateSlackWorkspace(ctx context.Context, cmd slack.SlashCommand) (*models.SlackWorkspace, error) {
	slackWorkspace, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, cmd.TeamID)
	if err != nil {
		return nil, err
	}

	if slackWorkspace.AccessToken == nil {
		return nil, err
	}

	client := slack.New(*slackWorkspace.AccessToken)
	if client == nil {
		return nil, err
	}

	// Get Channel Info
	input := &slack.GetConversationInfoInput{ChannelID: cmd.ChannelID}
	channel, err := client.GetConversationInfo(input)
	if err != nil {
		return nil, err
	}

	var canvasID string
	if channel.Properties != nil && channel.Properties.Canvas.FileId != "" {
		canvasID = channel.Properties.Canvas.FileId
	} else {
		canvasID, err = client.CreateChannelCanvas(cmd.ChannelID, slack.DocumentContent{
			Type:     "markdown",
			Markdown: "Welcome to Byrd! This is your team's shared canvas. Use `/byrd` to interact with Byrd.",
		})
		if err != nil {
			svc.showSupportModal(client, cmd.TriggerID, cmd.ChannelID, "This is embarassing", []string{"Seems like we're having trouble associating the canvas with this channel."})
			return nil, err
		}
	}

	ws, err := svc.repo.UpdateSlackWorkspace(ctx, cmd.TeamID, cmd.ChannelID, canvasID)
	if err != nil {
		svc.showSupportModal(client, cmd.TriggerID, cmd.ChannelID, "Something went wrong", []string{"Seems like we're having trouble updating the workspace."})
		return nil, err
	}

	// Show Success Modal
	// TODO: Implement showSuccessModal
	// svc.showSuccessModal(client, cmd.TriggerID)

	// svc.showSupportModal(client, cmd.TriggerID, cmd.ChannelID, "Something went wrong", []string{"Seems like we're having trouble updating the workspace."})

	return ws, nil
}

// Handles the bookkeeping for a Slack integration that has been removed
func (svc *slackWorkspaceService) DeleteSlackWorkspace(ctx context.Context, teamID string) error {
	return svc.repo.DeleteSlackWorkspace(ctx, teamID)
}

// GetSlackWorkspaceByTeamID retrieves a Slack workspace by its team ID
func (svc *slackWorkspaceService) GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*models.SlackWorkspace, error) {
	return svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
}

// GetSlackWorkspaceByWorkspaceID retrieves a Slack workspace by its Byrd workspace ID
func (svc *slackWorkspaceService) GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.SlackWorkspace, error) {
	return svc.repo.GetSlackWorkspaceByWorkspaceID(ctx, workspaceID)
}

// BatchGetWorkspaceByWorkspaceID retrieves a list of Slack workspaces by their Byrd workspace IDs
func (svc *slackWorkspaceService) BatchGetWorkspaceByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*models.SlackWorkspace, error) {
	return svc.repo.BatchGetSlackWorkspacesByWorkspaceIDs(ctx, workspaceIDs)
}

// ------ USER MANAGEMENT ------ //

// AddUserToSlackWorkspace adds a user to a Slack workspace
func (svc *slackWorkspaceService) AddUserToSlackWorkspace(ctx context.Context, teamID string, userEmail string) error {
	return nil
}

// RemoveUserFromSlackWorkspace removes a user from a Slack workspace
func (svc *slackWorkspaceService) RemoveUserFromSlackWorkspace(ctx context.Context, teamID string, userEmail string) error {
	return nil
}

// UserExistsInSlackWorkspace checks if a user exists in a Slack workspace
func (svc *slackWorkspaceService) UserExistsInSlackWorkspace(ctx context.Context, teamID string, userEmail string) (bool, error) {
	return false, nil
}

// ------ COMPETITOR MANAGEMENT ------ //

// CreateCompetitor creates a competitor in a Slack workspace
func (svc *slackWorkspaceService) CreateCompetitor(ctx context.Context, teamID string, pageURLs []string) error {
	return nil
}

// AddPageToCompetitor adds a page to a competitor in a Slack workspace
func (svc *slackWorkspaceService) AddPageToCompetitor(ctx context.Context, teamID string, competitorID uuid.UUID, pageURLs []string) error {
	return nil
}

func (svc *slackWorkspaceService) HandleSlackInteractionPayload(ctx context.Context, payload slack.InteractionCallback) error {
	// TODO: figure out how to handle routing of different payloads
	return svc.handleSupportSubmission(ctx, payload)
}
