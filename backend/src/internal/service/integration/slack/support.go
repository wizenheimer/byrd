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

	svc.logger.Info("Handling support submission", zap.Any(
		"payload", payload,
	))

	// Extract issue description
	issueDescription := payload.View.State.Values["support_issue_input"]["issue_description"].Value

	// Extract priority selection
	prioritySelection := payload.View.State.Values["support_priority"]["priority_selection"].SelectedOption.Value

	// Format response message
	confirmationMessage := fmt.Sprintf(
		"📝 *Request Submitted!*\nThanks for letting us know - we're right here with you. \nWe'll take good care of this!\n\n*Issue:* %s\n*Priority:* %s",
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

	// Open a DM with the user
	_, _, err = client.PostMessage(
		payload.User.ID,
		slack.MsgOptionText(confirmationMessage, false),
	)
	if err != nil {
		return err
	}

	return nil
}
