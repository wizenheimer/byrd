package slackworkspace

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (svc *slackWorkspaceService) handleSupportSubmission(ctx context.Context, payload slack.InteractionCallback) error {
	// Ensure it's a modal submission
	if payload.Type != slack.InteractionTypeViewSubmission {
		return errors.New("invalid interaction type")
	}

	// Extract issue description
	issueDescription := payload.View.State.Values["support_issue_input"]["issue_description"].Value

	// Extract priority selection
	prioritySelection := payload.View.State.Values["support_priority"]["priority_selection"].SelectedOption.Value

	// Format response message
	confirmationMessage := fmt.Sprintf(
		"📝 *Support Request Submitted!*\n\n*Issue:* %s\n*Priority:* %s\n\nOur team will review it soon",
		issueDescription, strings.ToUpper(prioritySelection),
	)

	teamID := payload.User.TeamID

	slackWorkspace, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
	if err != nil {
		return err
	}

	if slackWorkspace.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	// Send a confirmation message to the user
	client := slack.New(*slackWorkspace.AccessToken) // Replace with actual token

	channelID := payload.View.PrivateMetadata
	_, err = client.PostEphemeral(
		channelID,
		payload.User.ID,
		slack.MsgOptionText(confirmationMessage, false),
	)
	if err != nil {
		svc.logger.Error("Failed to send confirmation message", zap.Error(err))
	}

	return nil
}
