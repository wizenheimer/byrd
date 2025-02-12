package slackworkspace

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"go.uber.org/zap"
)

func formatSlackReportMarkdown(report models.Report) string {
	// Format the week number (assuming report.Time is in UTC)
	year, week := report.Time.ISOWeek()

	// Format the report into Markdown
	reportMarkdown := fmt.Sprintf(
		"## %s\n   Week: %d, %d\n",
		report.CompetitorName, week, year,
	)

	for _, categoryChange := range report.Changes {
		reportMarkdown += fmt.Sprintf(
			"### %s\n_%s_\n\n",
			categoryChange.Category, categoryChange.Summary,
		)

		for _, change := range categoryChange.Changes {
			reportMarkdown += fmt.Sprintf("- %s\n", change)
		}
		reportMarkdown += "\n"
	}

	reportMarkdown += "\n---\n"

	return reportMarkdown
}

func (svc *slackWorkspaceService) refreshReport(ctx context.Context, report *models.Report) error {
	// Get the Slack workspace
	slackWorkspace, err := svc.repo.GetSlackWorkspaceByWorkspaceID(ctx, report.WorkspaceID)
	if err != nil {
		svc.logger.Error("Failed to get Slack workspace", zap.Error(err))
		return err
	}
	if slackWorkspace.AccessToken == nil {
		return fmt.Errorf("no access token found for Slack workspace")
	}

	// Create a Slack client
	client := slack.New(*slackWorkspace.AccessToken)

	// If the slack workspace does not have a canvas, return an error
	if slackWorkspace.CanvasID == nil {
		return fmt.Errorf("no canvas found for Slack workspace")
	}

	// Format the report into Markdown
	reportMarkdown := formatSlackReportMarkdown(*report)

	// Insert at the start of the Canvas
	err = client.EditCanvas(slack.EditCanvasParams{
		CanvasID: *slackWorkspace.CanvasID,
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

type competitorDTO struct {
	ChannelID string   `json:"channel_id"`
	URLs      []string `json:"urls"`
}
