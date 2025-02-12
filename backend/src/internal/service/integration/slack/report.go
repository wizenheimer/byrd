package slackworkspace

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type CategoryChange struct {
	Category string
	Summary  string
	Changes  []string
}

type SlackReport struct {
	CompetitorName string
	Changes        []CategoryChange
	Time           time.Time
}

func formatSlackReportMarkdown(report SlackReport) string {
	// Format the week number (assuming report.Time is in UTC)
	_, week := report.Time.ISOWeek()
	weekText := fmt.Sprintf("Week %d", week)

	// Format the report into Markdown
	reportMarkdown := fmt.Sprintf(
		"*Competitor: %s*\n*%s*\n---\n",
		report.CompetitorName, weekText,
	)

	for _, categoryChange := range report.Changes {
		reportMarkdown += fmt.Sprintf(
			"*%s*\n_%s_\n",
			categoryChange.Category, categoryChange.Summary,
		)
		for _, change := range categoryChange.Changes {
			reportMarkdown += fmt.Sprintf("- %s\n", change)
		}
		reportMarkdown += "\n"
	}

	return reportMarkdown
}

func (svc *slackWorkspaceService) AppendReportToCanvas(client *slack.Client, canvasID string, report SlackReport) error {
	// Format the report into Markdown
	reportMarkdown := formatSlackReportMarkdown(report)

	// Insert at the start of the Canvas
	err := client.EditCanvas(slack.EditCanvasParams{
		CanvasID: canvasID,
		Changes: []slack.CanvasChange{
			{
				Operation: "insert_at_start",
				DocumentContent: slack.DocumentContent{
					Type:     "markdown",
					Markdown: reportMarkdown,
				},
			},
		},
	})

	if err != nil {
		svc.logger.Error("Failed to update Slack canvas with report", zap.Error(err))
		return err
	}

	svc.logger.Info("Successfully updated Slack canvas with report")
	return nil
}
