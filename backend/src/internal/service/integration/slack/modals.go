package slackworkspace

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"math/rand"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// showSupportModal displays a friendly troubleshooting modal
func (svc *slackWorkspaceService) showSupportModal(client *slack.Client, triggerID string, issueMeta string, issueList []string) {
	// Ensure issueTitle is within Slack’s 24-character limit

	issuesTitles := []string{
		"Something went wrong",
		"Whoops, hit a snag",
		"Well that's awkward",
		"Yikes, didn't work",
		"Aw snap, try again?",
		"Not our finest moment",
		"That was unexpected",
		"Slight hiccup there",
		"We hit a speed bump",
		"Looks like we goofed",
		"That's not quite right",
		"Oops, not quite there",
		"Time for Plan B",
		"This wasn't in the plan",
		"That's not ideal",
	}

	// pick random title
	issueTitle := issuesTitles[rand.Intn(len(issuesTitles))]

	// Ensure issueList is not nil
	if issueList == nil {
		issueList = []string{}
	}

	// Clean up issue list items to prevent malformed Markdown
	for i, item := range issueList {
		issueList[i] = strings.TrimSpace(item)
	}

	// Append friendly troubleshooting message
	issueList = append(issueList, "No worries! Our team is just a text away.\n\nLet us know what's going on and we'll get back to you right away.")

	// Construct issue body with Markdown
	issueBody := strings.Join(issueList, "\n\n")

	// Get priority options safely
	priorityOptions := getPriorityOptions()
	if len(priorityOptions) == 0 {
		svc.logger.Warn("getPriorityOptions() returned an empty list!")
	}

	// Construct the Slack modal view
	supportModalView := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject("plain_text", issueTitle, false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false), // Required for form submission
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Issue Explanation
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", issueBody, false, false),
					nil, nil,
				),

				// Divider for clean formatting
				slack.NewDividerBlock(),

				// Text Input Field for Issue Description
				slack.NewInputBlock(
					"support_issue_input",
					slack.NewTextBlockObject("plain_text", "Tell us more", false, false),
					nil,
					slack.NewPlainTextInputBlockElement(
						slack.NewTextBlockObject("plain_text", "Type your issue here...", false, false),
						"issue_description",
					),
				),

				// Divider for clean formatting
				slack.NewDividerBlock(),

				// Priority Selection Dropdown
				slack.NewInputBlock(
					"support_priority",
					slack.NewTextBlockObject("plain_text", "Select priority", false, false),
					nil, // No hint text
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeStatic,
						slack.NewTextBlockObject("plain_text", "Choose priority...", false, false),
						"priority_selection",
						priorityOptions...,
					),
				),
			},
		},
		PrivateMetadata: issueMeta, // Store issue metadata for context
		CallbackID:      "support_submission",
	}

	// Debug log: Print full modal payload
	jsonView, _ := json.MarshalIndent(supportModalView, "", "  ")
	svc.logger.Debug("ModalViewRequest JSON", zap.String("modal", string(jsonView)))

	// Open the modal in Slack
	_, err := client.OpenView(triggerID, supportModalView)
	if err != nil {
		svc.logger.Error("Failed to open support modal", zap.Error(err))
	}
}

// getPriorityOptions returns valid Slack priority dropdown options
func getPriorityOptions() []*slack.OptionBlockObject {
	return []*slack.OptionBlockObject{
		slack.NewOptionBlockObject("high",
			slack.NewTextBlockObject(slack.PlainTextType, "High", false, false), nil),
		slack.NewOptionBlockObject("medium",
			slack.NewTextBlockObject(slack.PlainTextType, "Medium", false, false), nil),
		slack.NewOptionBlockObject("low",
			slack.NewTextBlockObject(slack.PlainTextType, "Low", false, false), nil),
	}
}

// showSuccessModal displays a confirmation modal in Slack
func (svc *slackWorkspaceService) showSuccessModal(client *slack.Client, triggerID string, channelID string, successTitle string, successMessage string, successBody []string) {
	// Ensure title is within Slack’s 24-character limit
	if successTitle == "" {
		successMessages := []string{
			"That worked nicely",
			"All set to go",
			"Looking good now",
			"Smooth sailing ahead",
			"Done and done",
			"You're all set",
			"Everything's in place",
			"Good to go",
			"That did the trick",
			"Nailed it",
			"That's a wrap",
			"All good here",
			"Ready to roll",
			"Perfect, all done",
		}

		// pick random title
		successTitle = successMessages[rand.Intn(len(successMessages))]
	} else if len(successTitle) > 24 {
		successTitle = successTitle[:21] + "..." // Truncate and add ellipsis
	}

	// Default success message if empty
	if successMessage == "" {
		successMessage = "Your request was processed successfully! 🎉"
	}

	// Construct the Slack modal view
	successModalView := slack.ModalViewRequest{
		Type:  slack.VTModal,
		Title: slack.NewTextBlockObject("plain_text", successTitle, false, false),
		Close: slack.NewTextBlockObject("plain_text", "Close", false, false), // No Submit button needed
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Success Message Section
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", successMessage, false, false),
					nil, nil,
				),
			},
		},
		PrivateMetadata: channelID, // Store channel ID for context
	}

	// Add success body if provided
	if len(successBody) > 0 {
		// Clean up success body items to prevent malformed Markdown
		for i, item := range successBody {
			successBody[i] = strings.TrimSpace(item)
		}
		successBodyString := strings.Join(successBody, "\n\n")
		successModalView.Blocks.BlockSet = append(successModalView.Blocks.BlockSet, slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", successBodyString, false, false),
			nil, nil,
		))
	}

	// Debug log: Print full modal payload

	// Open the modal in Slack
	_, err := client.OpenView(triggerID, successModalView)
	if err != nil {
		svc.logger.Error("Failed to open success modal", zap.Error(err))
	}
}

func (svc *slackWorkspaceService) showUsageLimitModal(client *slack.Client, triggerID string, workspacePlan models.WorkspacePlan, workspaceResource models.WorkspaceResource) error {
	// Get the max limit for the current plan and resource
	limitCount, err := workspacePlan.GetMaxLimit(workspaceResource)
	if err != nil {
		return fmt.Errorf("failed to get max limit: %w", err)
	}

	// Get the next plan
	nextPlan, err := workspacePlan.NextPlan()
	if err != nil {
		return fmt.Errorf("failed to get next plan: %w", err)
	}

	// Get the resource limit for the next plan
	nextPlanLimit, err := nextPlan.GetMaxLimit(workspaceResource)
	if err != nil {
		return fmt.Errorf("failed to get next plan limit: %w", err)
	}

	// Format the resource name for display
	resourceName := strings.TrimPrefix(string(workspaceResource), "WorkspaceResource")
	resourceName = strings.ToLower(resourceName)

	// Caserize the resource name
	caser := cases.Title(language.English)

	limitCategoryString := fmt.Sprintf("*You've reached the limit for %s*", resourceName)
	limitCountString := fmt.Sprintf("Your current plan allows tracking up to %v %s. Upgrade to track up to %v %s.",
		limitCount, resourceName, nextPlanLimit, resourceName)
	currentPlanString := fmt.Sprintf("*Current Plan*\n%s", caser.String(workspacePlan.ToString()))
	currentUsageString := fmt.Sprintf("*%s Used*\n%v of %v", caser.String(resourceName), limitCount, limitCount)
	nextPlanString := fmt.Sprintf("Upgrade to %s", nextPlan.ToString())

	limitModal := slack.ModalViewRequest{
		Type:  "modal",
		Title: slack.NewTextBlockObject(slack.PlainTextType, "Usage Limit Reached", false, false),
		Close: slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, limitCategoryString, false, false),
					nil,
					nil,
				),
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, limitCountString, false, false),
					nil,
					nil,
				),
				slack.NewSectionBlock(
					nil,
					[]*slack.TextBlockObject{
						slack.NewTextBlockObject(slack.MarkdownType, currentPlanString, false, false),
						slack.NewTextBlockObject(slack.MarkdownType, currentUsageString, false, false),
					},
					nil,
				),
				slack.NewDividerBlock(),
				slack.NewActionBlock(
					"upgrade_actions",
					slack.NewButtonBlockElement(
						"upgrade_plan",
						"upgrade",
						slack.NewTextBlockObject(slack.PlainTextType, nextPlanString, false, false),
					).WithStyle(slack.StylePrimary).WithURL("https://byrdhq.com/plans"),
				),
			},
		},
	}

	_, err = client.OpenView(triggerID, limitModal)
	if err != nil {
		return fmt.Errorf("failed to open limit modal: %w", err)
	}

	return nil
}

// showUserInviteModal displays a modal to invite a user to the workspace
func (svc *slackWorkspaceService) showUserInviteModal(client *slack.Client, cmd slack.SlashCommand) error {
	creatorID := cmd.UserID
	modalRequest := slack.ModalViewRequest{
		Type:       "modal",
		Title:      slack.NewTextBlockObject(slack.PlainTextType, "Invite Users", false, false),
		Submit:     slack.NewTextBlockObject(slack.PlainTextType, "Send Invites", false, false),
		Close:      slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		CallbackID: "invite_users",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"user_selection",
					slack.NewTextBlockObject(slack.PlainTextType, "Users", false, false),
					nil,
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeUser,
						slack.NewTextBlockObject(slack.PlainTextType, "Select users", false, false),
						"user_select",
					),
				),
			},
		},
		PrivateMetadata: creatorID,
	}

	_, err := client.OpenView(cmd.TriggerID, modalRequest)
	return err
}

func (svc *slackWorkspaceService) showPageAddModal(client *slack.Client, cmd slack.SlashCommand, workspaceID uuid.UUID, urls []string) error {
	// --- Dropdown (Single Select) ---
	competitorSelect := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, "Select a competitor", false, false),
		"select_competitor",
		svc.getCompetitorOptions(context.Background(), workspaceID)...,
	)

	competitorBlock := slack.NewInputBlock(
		"competitor_selection",
		slack.NewTextBlockObject(slack.PlainTextType, "Assign to Competitor", false, false),
		nil, // No hint
		competitorSelect,
	)

	// --- Multi-Select (DiffProfile) ---
	caser := cases.Title(language.English)
	diffProfileOptions := models.GetDefaultDiffProfile()
	var multiSelectOptions []*slack.OptionBlockObject
	for _, profile := range diffProfileOptions {
		multiSelectOptions = append(multiSelectOptions, slack.NewOptionBlockObject(
			profile,
			slack.NewTextBlockObject(slack.PlainTextType, caser.String(profile), false, false),
			nil,
		))
	}

	diffProfileMultiSelect := slack.NewOptionsMultiSelectBlockElement(
		slack.MultiOptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, "Product, Pricing, Partnerships etc.", false, false),
		"select_diff_profiles",
		multiSelectOptions...,
	)

	diffProfileBlock := slack.NewInputBlock(
		"diff_profile_selection",
		slack.NewTextBlockObject(slack.PlainTextType, "Select Competitor Profiles", false, false),
		nil, // No hint
		diffProfileMultiSelect,
	)

	competitorData := competitorDTO{
		ChannelID: cmd.ChannelID,
		URLs:      urls,
	}

	jsonBytes, err := json.Marshal(competitorData)
	if err != nil {
		return err
	}
	base64String := base64.StdEncoding.EncodeToString(jsonBytes)

	modal := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject(slack.PlainTextType, "Assign Competitor", false, false),
		Submit: slack.NewTextBlockObject(slack.PlainTextType, "Save", false, false),
		Close:  slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				competitorBlock,
				diffProfileBlock,
			}, // Ensure correct structure
		},
		CallbackID:      "save_competitor",
		PrivateMetadata: base64String,
	}

	_, err = client.OpenView(cmd.TriggerID, modal)

	return err
}
