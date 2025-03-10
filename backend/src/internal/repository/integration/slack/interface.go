package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/models/integration/slack"
)

type SlackWorkspaceRepository interface {
	// CreateSlackWorkspace creates a new slack workspace
	// This is the immediate outcome of linking a workspace
	// The status of the workspace is pending
	CreateSlackWorkspace(ctx context.Context, workspaceID uuid.UUID, channelID, channelWebhookURL, teamID, accessToken string) (*slack.SlackWorkspace, error)

	// GetSlackWorkspaceByTeamID gets a slack workspace by team ID
	GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*slack.SlackWorkspace, error)

	// GetSlackWorkspaceByWorkspaceID gets a slack workspace by workspace ID
	GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*slack.SlackWorkspace, error)

	// BatchGetSlackWorkspacesByWorkspaceIDs gets slack workspaces by workspace IDs
	BatchGetSlackWorkspacesByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*slack.SlackWorkspace, error)

	// WorkspaceExists checks if a workspace exists
	WorkspaceExists(ctx context.Context, teamID string) (bool, error)

	// Delete deletes a slack workspace
	DeleteSlackWorkspace(ctx context.Context, teamID string) error

	// UpdateSlackWorkspaceAccessToken updates the access token of a slack workspace
	UpdateSlackWorkspaceAccessToken(ctx context.Context, teamID, accessToken string) (*slack.SlackWorkspace, error)
}
