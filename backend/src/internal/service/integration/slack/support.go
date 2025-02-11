package slackworkspace

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
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

func (svc *slackWorkspaceService) handlePageSubmission(ctx context.Context, payload slack.InteractionCallback) error {
	// Extract selected competitor
	selectedCompetitor := payload.View.State.Values["competitor_selection"]["select_competitor"].SelectedOption.Value

	// Extract selected diff profiles (multi-select)
	selectedDiffProfiles := []string{}
	if multiSelectBlock, ok := payload.View.State.Values["diff_profile_selection"]["select_diff_profiles"]; ok {
		for _, option := range multiSelectBlock.SelectedOptions {
			selectedDiffProfiles = append(selectedDiffProfiles, option.Value)
		}
	}

	// Decode private metadata
	base64String := payload.View.PrivateMetadata
	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	var competitorData CompetitorData
	err = json.Unmarshal(decodedBytes, &competitorData)
	if err != nil {
		return err
	}

	// Fetch workspace details
	ws, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, payload.User.TeamID)
	if err != nil {
		return err
	}

	if ws.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	// Use user-selected diff profiles instead of default ones
	diffProfiles := selectedDiffProfiles
	if len(diffProfiles) == 0 { // Fallback to default if nothing selected
		diffProfiles = models.GetDefaultDiffProfile()
	}

	// Process URLs into PageProps
	var pages []models.PageProps
	for _, u := range competitorData.URLs {
		pageProp, err := models.NewPageProps(u, diffProfiles)
		if err != nil {
			svc.logger.Error("Failed to create page props", zap.Error(err))
			continue
		}
		pages = append(pages, pageProp)
	}

	// Handle competitor logic
	var competitorUUID uuid.UUID
	if selectedCompetitor == uuid.Nil.String() {
		_, err = svc.ws.AddCompetitorToWorkspace(ctx, ws.WorkspaceID, pages)
	} else {
		competitorUUID, err = uuid.Parse(selectedCompetitor)
		if err != nil {
			return err
		}
		_, err = svc.ws.AddPageToCompetitor(ctx, ws.WorkspaceID, competitorUUID, pages)
	}

	if err != nil {
		return err
	}

	// Notify user in Slack
	client := slack.New(*ws.AccessToken)
	channelID := competitorData.ChannelID

	_, err = client.PostEphemeral(
		channelID,       // Channel where the interaction happened
		payload.User.ID, // User who triggered the action
		slack.MsgOptionText("URL is now getting tracked", false),
	)

	if err != nil {
		svc.logger.Error("Failed to post ephemeral message", zap.Error(err))
	}

	return nil
}
