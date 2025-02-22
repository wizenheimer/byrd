package slackworkspace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"go.uber.org/zap"
)

func formatSlackReportMarkdown(report models.Report) string {
	// Your existing formatting code remains the same
	year, week := report.Time.ISOWeek()

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

	// Format the report into Markdown
	reportMarkdown := formatSlackReportMarkdown(*report)

	// Create the webhook message payload
	payload := map[string]interface{}{
		"text":   reportMarkdown,
		"mrkdwn": true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		svc.logger.Error("Failed to marshal webhook payload", zap.Error(err))
		return err
	}

	// Send the webhook request
	resp, err := http.Post(
		slackWorkspace.ChannelWebhookURL,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		svc.logger.Error("failed to send webhook request", zap.Error(err))
		return err
	}
	if resp == nil {
		return errors.New("received nil response from slack")
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
		svc.logger.Error("Webhook request failed",
			zap.Error(err),
			zap.Int("statusCode", resp.StatusCode),
		)
		return err
	}

	svc.logger.Info("Successfully sent report via webhook",
		zap.Any("reportMarkdown", reportMarkdown),
		zap.String("webhookURL", slackWorkspace.ChannelWebhookURL),
	)
	return nil
}

type competitorDTO struct {
	ChannelID string   `json:"channel_id"`
	URLs      []string `json:"urls"`
}
