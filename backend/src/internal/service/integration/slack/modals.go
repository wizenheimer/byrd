package slackworkspace

import (
	"encoding/json"
	"strings"

	"math/rand"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
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
