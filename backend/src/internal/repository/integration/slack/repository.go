package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wizenheimer/byrd/src/internal/models/integration/slack"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type swr struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewSlackWorkspaceRepository(
	tm *transaction.TxManager,
	logger *logger.Logger,
) (SlackWorkspaceRepository, error) {
	return &swr{
		tm:     tm,
		logger: logger,
	}, nil
}

func (repo *swr) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return repo.tm.GetQuerier(ctx)
}

// scanSlackWorkspace scans a row into a SlackWorkspace
func scanSlackWorkspace(row pgx.Row) (*slack.SlackWorkspace, error) {
	workspace := slack.SlackWorkspace{}
	var channelID, canvasID sql.NullString

	err := row.Scan(
		&workspace.WorkspaceID,
		&workspace.TeamID,
		&channelID,
		&canvasID,
		&workspace.AccessToken,
		&workspace.Status,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("workspace not found") // Return proper error
		}
		return nil, fmt.Errorf("error scanning slack workspace: %w", err)
	}

	// Handle nullable fields
	if channelID.Valid {
		workspace.ChannelID = &channelID.String
	}
	if canvasID.Valid {
		workspace.CanvasID = &canvasID.String
	}

	// Values from DB are always encoded
	workspace.IsDecoded = false

	// Automatically decode for the caller
	if err := workspace.Decode(); err != nil {
		return nil, fmt.Errorf("error decoding workspace: %w", err)
	}

	return &workspace, nil
}

// CreateSlackWorkspace creates a new slack workspace in pending status
func (repo *swr) CreateSlackWorkspace(ctx context.Context, workspaceID uuid.UUID, teamID, accessToken string) (*slack.SlackWorkspace, error) {
	// Create workspace with decoded values
	workspace := &slack.SlackWorkspace{
		WorkspaceID: workspaceID,
		TeamID:      teamID,
		AccessToken: &accessToken,
		Status:      slack.SlackWorkspaceStatusPending,
		IsDecoded:   true,
	}

	// Encode before saving to DB
	if err := workspace.Encode(); err != nil {
		return nil, fmt.Errorf("failed to encode workspace: %w", err)
	}

	query := `
        INSERT INTO slack_workspaces (workspace_id, team_id, access_token, status)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (workspace_id) DO UPDATE
        SET team_id = EXCLUDED.team_id,
            access_token = EXCLUDED.access_token,
            status = EXCLUDED.status
        RETURNING workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at`

	row := repo.getQuerier(ctx).QueryRow(ctx, query,
		workspace.WorkspaceID,
		workspace.TeamID,
		*workspace.AccessToken,
		workspace.Status,
	)

	return scanSlackWorkspace(row)
}

// UpdateSlackWorkspace sets channel and canvas, activating the workspace
func (repo *swr) UpdateSlackWorkspace(ctx context.Context, teamID, channelID, canvasID string) (*slack.SlackWorkspace, error) {
	query := `
        UPDATE slack_workspaces
        SET channel_id = $2,
            canvas_id = $3,
            status = $4
        WHERE team_id = $1 AND status != $5
        RETURNING workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at`

	row := repo.getQuerier(ctx).QueryRow(ctx, query,
		teamID,
		channelID,
		canvasID,
		slack.SlackWorkspaceStatusActive,
		slack.SlackWorkspaceStatusInactive,
	)

	return scanSlackWorkspace(row)
}

// GetSlackWorkspaceByTeamID gets a slack workspace by team ID
func (repo *swr) GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*slack.SlackWorkspace, error) {
	query := `
        SELECT workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at
        FROM slack_workspaces
        WHERE team_id = $1 AND status != $2`

	row := repo.getQuerier(ctx).QueryRow(ctx, query, teamID, slack.SlackWorkspaceStatusInactive)
	return scanSlackWorkspace(row)
}

// GetSlackWorkspaceByWorkspaceID gets a slack workspace by workspace ID
func (repo *swr) GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*slack.SlackWorkspace, error) {
	query := `
        SELECT workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at
        FROM slack_workspaces
        WHERE workspace_id = $1 AND status != $2`

	row := repo.getQuerier(ctx).QueryRow(ctx, query, workspaceID, slack.SlackWorkspaceStatusInactive)
	return scanSlackWorkspace(row)
}

// BatchGetSlackWorkspacesByWorkspaceIDs gets slack workspaces by workspace IDs
func (repo *swr) BatchGetSlackWorkspacesByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*slack.SlackWorkspace, error) {
	query := `
        SELECT workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at
        FROM slack_workspaces
        WHERE workspace_id = ANY($1) AND status != $2`

	rows, err := repo.getQuerier(ctx).Query(ctx, query, workspaceIDs, slack.SlackWorkspaceStatusInactive)
	if err != nil {
		return nil, fmt.Errorf("error querying slack workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []*slack.SlackWorkspace
	for rows.Next() {
		workspace, err := scanSlackWorkspace(rows)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

// WorkspaceExists checks if a non-inactive workspace exists
func (repo *swr) WorkspaceExists(ctx context.Context, teamID string) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM slack_workspaces
            WHERE team_id = $1 AND status != $2
        )`

	var exists bool
	err := repo.getQuerier(ctx).QueryRow(ctx, query, teamID, slack.SlackWorkspaceStatusInactive).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking workspace existence: %w", err)
	}

	return exists, nil
}

// DeleteSlackWorkspace marks a workspace as inactive
func (repo *swr) DeleteSlackWorkspace(ctx context.Context, teamID string) error {
	query := `
        UPDATE slack_workspaces
        SET status = $2
        WHERE team_id = $1 AND status != $2`

	result, err := repo.getQuerier(ctx).Exec(ctx, query,
		teamID,
		slack.SlackWorkspaceStatusInactive,
	)
	if err != nil {
		return fmt.Errorf("error deleting workspace: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no active workspace found with team_id: %s", teamID)
	}

	return nil
}

// UpdateSlackWorkspaceAccessToken updates the access token of an active workspace
func (repo *swr) UpdateSlackWorkspaceAccessToken(ctx context.Context, teamID, accessToken string) (*slack.SlackWorkspace, error) {
	// Create temporary workspace to encode the new token
	tempWorkspace := &slack.SlackWorkspace{
		AccessToken: &accessToken,
		IsDecoded:   true,
	}

	// Encode for storage
	if err := tempWorkspace.Encode(); err != nil {
		return nil, fmt.Errorf("failed to encode access token: %w", err)
	}

	query := `
        UPDATE slack_workspaces
        SET access_token = $2
        WHERE team_id = $1 AND status = $3
        RETURNING workspace_id, team_id, channel_id, canvas_id, access_token, status, created_at, updated_at`

	row := repo.getQuerier(ctx).QueryRow(ctx, query,
		teamID,
		*tempWorkspace.AccessToken,
		slack.SlackWorkspaceStatusActive,
	)

	return scanSlackWorkspace(row)
}
