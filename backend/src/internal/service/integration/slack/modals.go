package slackworkspace

import (
	"encoding/json"
	"strings"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

// showSupportModal displays a friendly troubleshooting modal
func (svc *slackWorkspaceService) showSupportModal(client *slack.Client, triggerID string, channelID string, issueTitle string, issueList []string) {
	// Ensure issueTitle is within Slack’s 24-character limit
	if issueTitle == "" {
		issueTitle = "Something went wrong"
	} else if len(issueTitle) > 24 {
		issueTitle = issueTitle[:21] + "..." // Truncate and add ellipsis
	}

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
		PrivateMetadata: channelID, // Store the channel ID for later use
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
