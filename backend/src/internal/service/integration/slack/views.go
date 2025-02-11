package slackworkspace

import (
	"log"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

// showSuccessModal displays a celebratory confirmation modal
func (svc *slackWorkspaceService) showSuccessModal(client *slack.Client, triggerID string) {
	_, err := client.OpenView(triggerID, slack.ModalViewRequest{
		Type:  slack.VTModal,
		Title: slack.NewTextBlockObject("plain_text", "🎉 Workspace Updated!", false, false),
		Close: slack.NewTextBlockObject("plain_text", "Got it!", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Celebration Message
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn",
						"*Your Byrd workspace is now ready!*\n\n:tada: *Awesome!* Byrd is now active in this channel and will keep you updated.",
						false, false),
					nil, nil,
				),

				// Next Steps
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn",
						"*Next Steps:*\n- `/watch` - Monitor a competitor's URL.\n- `/invite` - Add teammates to Byrd.\n\nWant more? Type `/byrd help` anytime!",
						false, false),
					nil, nil,
				),
			},
		},
	})
	if err != nil {
		log.Printf("Failed to open success modal: %v", err)
	}
}

// showSupportModal displays a friendly troubleshooting modal
func (svc *slackWorkspaceService) showSupportModal(client *slack.Client, triggerID string, channelID string) {
	view := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject("plain_text", "Something went wrong", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false), // Required for form submission
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Issue Explanation
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn",
						"*We ran into an issue while updating Byrd.*\n\nNo worries! We have your back:",
						false, false),
					nil, nil,
				),

				// Divider for clean formatting
				slack.NewDividerBlock(),

				// Text Input Field for Issue Description
				slack.NewInputBlock(
					"support_issue_input",
					slack.NewTextBlockObject("plain_text", "Tell us more", false, false),
					nil,
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "Type your issue here...", false, false), "issue_description"),
				),

				// Divider for clean formatting
				slack.NewDividerBlock(),

				// Priority Selection Dropdown (Fixed)
				slack.NewInputBlock(
					"support_priority",
					slack.NewTextBlockObject("plain_text", "Select priority", false, false),
					nil, // No hint text
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeStatic,
						slack.NewTextBlockObject("plain_text", "Choose priority...", false, false),
						"priority_selection",
						getPriorityOptions()...,
					),
				),
			},
		},
		PrivateMetadata: channelID, // Store the channel ID for later use
	}

	// Open the modal
	_, err := client.OpenView(triggerID, view)
	if err != nil {
		svc.logger.Error("Failed to open support modal", zap.Error(err))
	}
}

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
