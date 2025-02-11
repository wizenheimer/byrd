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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (svc *slackWorkspaceService) handleSupportSubmission(ctx context.Context, payload slack.InteractionCallback) error {
	// Ensure it's a modal submission
	if payload.Type != slack.InteractionTypeViewSubmission {
		return errors.New("invalid interaction type")
	}

	svc.logger.Info("Handling support submission", zap.Any("payload", payload))

	// Extract issue description
	issueDescription := payload.View.State.Values["support_issue_input"]["issue_description"].Value

	// Extract priority selection
	prioritySelection := strings.ToUpper(payload.View.State.Values["support_priority"]["priority_selection"].SelectedOption.Value)

	// Get Slack workspace
	teamID := payload.User.TeamID
	slackWorkspace, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
	if err != nil {
		return err
	}
	if slackWorkspace.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	client := slack.New(*slackWorkspace.AccessToken)

	// Open a DM with the user
	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{payload.User.ID},
	})
	if err != nil {
		svc.logger.Error("failed to open DM with user", zap.Error(err))
		return err
	}

	// Construct the styled Slack message blocks
	msgBlocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject(slack.PlainTextType, "📝 Request Submitted!", false, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "Thanks for letting us know – we're right here with you. \nWe'll take good care of this! 🚀", false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "*📌 Issue:*", false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("> %s", issueDescription), false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "*⚡ Priority:*", false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("`%s`", prioritySelection), false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock(
			"tracking_info",
			slack.NewTextBlockObject(slack.MarkdownType, "💡 _We'll update you as soon as we have progress._", false, false),
		),
	}

	// Send the message as a DM to the user
	_, _, err = client.PostMessage(
		channel.ID, // DM Channel ID
		slack.MsgOptionBlocks(msgBlocks...),
	)
	if err != nil {
		svc.logger.Error("failed to send DM support confirmation", zap.Error(err))
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
			svc.logger.Error("failed to create page props", zap.Error(err))
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

	// Open a direct message with the user
	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{payload.User.ID},
	})
	if err != nil {
		svc.logger.Error("failed to open DM with user", zap.Error(err))
		return err
	}

	// Format URLs as individual sections
	var urlBlocks []slack.Block
	for _, url := range competitorData.URLs {
		urlBlocks = append(urlBlocks, slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("🔗 *%s*", url), false, false),
			nil, nil,
		))
	}

	// Format profile selections
	profileText := "None selected"
	caser := cases.Title(language.English)
	if len(diffProfiles) > 0 {
		var profileList string
		for _, profile := range diffProfiles {
			profileList += fmt.Sprintf("• %s\n", caser.String(profile))
		}
		profileText = profileList
	}

	// Construct the Slack message blocks
	msgBlocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject(slack.PlainTextType, "🚀 Tracking Started!", false, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "We've started tracking the following URLs for you:", false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "*🔍 Tracked URLs:*", false, false),
			nil, nil,
		),
	}
	msgBlocks = append(msgBlocks, urlBlocks...) // Add each URL as a separate section

	msgBlocks = append(msgBlocks,
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "*📂 Profiles Selected:*", false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, profileText, false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock(
			"tracking_info",
			slack.NewTextBlockObject(slack.MarkdownType, "⚡ _We'll notify you when changes happen._", false, false),
		),
	)

	// Send the message as a DM to the user
	_, _, err = client.PostMessage(
		channel.ID, // DM Channel ID
		slack.MsgOptionBlocks(msgBlocks...),
	)
	if err != nil {
		svc.logger.Error("failed to send DM tracking confirmation", zap.Error(err))
	}

	return nil
}
