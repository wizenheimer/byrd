package slackworkspace

import (
	"context"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/integration/slack"
)

type SlackWorkspaceService interface {
	// ------ WORKSPACE MANAGEMENT ------ //

	// Creates and associates an existing Byrd workspace with a Slack workspace
	CreateSlackWorkspace(ctx context.Context, workspaceID uuid.UUID, teamID, accessToken string) (*models.SlackWorkspace, error)

	// UpdateSlackWorkspace updates the access token for a Slack workspace
	UpdateSlackWorkspace(ctx context.Context, cmd slack.SlashCommand) (*models.SlackWorkspace, error)

	// Handles the bookkeeping for a Slack integration that has been removed
	DeleteSlackWorkspace(ctx context.Context, teamID string) error

	// GetSlackWorkspaceByTeamID retrieves a Slack workspace by its team ID
	GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*models.SlackWorkspace, error)

	// GetSlackWorkspaceByWorkspaceID retrieves a Slack workspace by its Byrd workspace ID
	GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.SlackWorkspace, error)

	// BatchGetWorkspaceByWorkspaceIDs retrieves a list of Slack workspaces by their Byrd workspace IDs
	BatchGetWorkspaceByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*models.SlackWorkspace, error)

	// ------ USER MANAGEMENT ------ //

	// AddUserToSlackWorkspace adds a user to a Slack workspace
	AddUserToSlackWorkspace(ctx context.Context, teamID string, userEmail string) error

	// RemoveUserFromSlackWorkspace removes a user from a Slack workspace
	RemoveUserFromSlackWorkspace(ctx context.Context, teamID string, userEmail string) error

	// UserExistsInSlackWorkspace checks if a user exists in a Slack workspace
	UserExistsInSlackWorkspace(ctx context.Context, teamID string, userEmail string) (bool, error)

	// ------ COMPETITOR MANAGEMENT ------ //

	// CreateCompetitor creates a competitor in a Slack workspace
	CreateCompetitor(ctx context.Context, teamID string, pageURLs []string) error

	// AddPageToCompetitor adds a page to a competitor in a Slack workspace
	AddPageToCompetitor(ctx context.Context, teamID string, competitorID uuid.UUID, pageURLs []string) error

	// ----- Slack Interaction Payload Router ---- //

	// HandleSlackInteractionPayload handles slack interaction payloads
	HandleSlackInteractionPayload(ctx context.Context, payload slack.InteractionCallback) error
}
