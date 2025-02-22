package slackworkspace

import (
	"context"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	core_models "github.com/wizenheimer/byrd/src/internal/models/core"
	models "github.com/wizenheimer/byrd/src/internal/models/integration/slack"
)

type SlackWorkspaceService interface {
	// ------ WORKSPACE MANAGEMENT ------ //

	// Creates and associates an existing Byrd workspace with a Slack workspace
	CreateWorkspace(ctx context.Context, pages []core_models.PageProps, channelID, channelWebhookURL, userID, teamID, accessToken string) (*models.SlackWorkspace, error)

	// Handles SlackWorkspace deletions
	DeleteSlackWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	// ------ COMPETITOR MANAGEMENT ------ //

	// CreateCompetitor creates a competitor in a Slack workspace
	CreateCompetitorForWorkspace(ctx context.Context, cmd slack.SlashCommand) error

	// ------ USER MANAGEMENT ------ //

	// AddUserToSlackWorkspace adds a user to a Slack workspace
	AddUserToSlackWorkspace(ctx context.Context, cmd slack.SlashCommand) error

	// ----- RESOURCE MANAGEMENT ----- //

	// GetSlackWorkspaceByTeamID retrieves a Slack workspace by its team ID
	GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*models.SlackWorkspace, error)

	// GetSlackWorkspaceByWorkspaceID retrieves a Slack workspace by its Byrd workspace ID
	GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.SlackWorkspace, error)

	// BatchGetWorkspaceByWorkspaceIDs retrieves a list of Slack workspaces by their Byrd workspace IDs
	BatchGetWorkspaceByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*models.SlackWorkspace, error)

	// IntegrationExists checks if a Slack integration exists for a Byrd workspace
	IntegrationExistsForWorkspace(ctx context.Context, workspaceID uuid.UUID) (bool, error)

	// UserExistsInSlackWorkspace checks if a user exists in a Slack workspace
	UserExistsInSlackWorkspace(ctx context.Context, teamID string, userEmail string) (bool, error)

	// ----- Slack Interaction Payload Router ---- //

	// HandleSlackInteractionPayload handles slack interaction payloads
	HandleSlackInteractionPayload(ctx context.Context, payload slack.InteractionCallback) error

	// ----- Slack Report Management ----- //
	DispatchReportToWorkspaceMembers(ctx context.Context, workspaceID, competitorID uuid.UUID) error
}
