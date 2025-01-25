// ./src/internal/alert/slack.go
// ./src/internal/service/alert/slack.go
package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/slack-go/slack"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"golang.org/x/time/rate"
)

var (
	ErrSlackTokenMissing         = fmt.Errorf("slack token missing")
	ErrSlackChannelIDMissing     = fmt.Errorf("slack channel ID missing")
	ErrEncounteredSlackRateLimit = fmt.Errorf("encountered slack rate limit")
	ErrExhaustedSlackRetries     = fmt.Errorf("exhausted slack retries")
)

// SlackAlertClient implements AlertClient interface for Slack
type slackAlertClient struct {
	// client for Slack
	client *slack.Client
	// config for Slack
	config models.SlackConfig
	// logger for logging
	logger *logger.Logger
	// Rate limiter for Slack API
	limiter *rate.Limiter
}

func NewSlackAlertClient(config models.SlackConfig, logger *logger.Logger) (AlertClient, error) {
	if config.Token == "" {
		return nil, ErrSlackTokenMissing
	}

	if config.ChannelID == "" {
		return nil, ErrSlackChannelIDMissing
	}

	// Set default retry count if not specified
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	client := &slackAlertClient{
		client:  slack.New(config.Token),
		config:  config,
		logger:  logger.WithFields(map[string]interface{}{"module": "slack_alert_client"}),
		limiter: rate.NewLimiter(1, 30), // 1 message per second, burst of 30
	}

	return client, nil
}

func (s *slackAlertClient) SendAlert(ctx context.Context, alert models.Alert) error {
	// Rate limiting
	if err := s.limiter.Wait(ctx); err != nil {
		return ErrEncounteredSlackRateLimit
	}

	// Send with retries
	err := s.sendWithRetries(ctx, alert)
	if err != nil {
		return err
	}

	return nil
}

func (s *slackAlertClient) sendWithRetries(ctx context.Context, alert models.Alert) error {
	attachment := s.createAttachment(alert)

	for attempt := 0; attempt <= s.config.RetryCount; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, _, err := s.client.PostMessageContext(
				ctx,
				s.config.ChannelID,
				slack.MsgOptionAttachments(attachment),
				slack.MsgOptionUsername(s.config.DefaultUser),
			)
			if err == nil {
				return nil
			}

			if attempt == s.config.RetryCount {
				return ErrExhaustedSlackRetries
			}

			// Calculate exponential backoff with jitter
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			jitter := time.Duration(100 * time.Millisecond)
			backoff = backoff + time.Duration(time.Now().UnixNano()%int64(jitter))

			// Wait for backoff duration or context cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				continue
			}
		}
	}
	return nil
}

func (s *slackAlertClient) createAttachment(alert models.Alert) slack.Attachment {
	color := s.config.ColorMapping[alert.Severity]

	fields := make([]slack.AttachmentField, 0)
	for key, value := range alert.Metadata {
		fields = append(fields, slack.AttachmentField{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	return slack.Attachment{
		Color:      color,
		Title:      alert.Title,
		Text:       alert.Description,
		Fields:     fields,
		Ts:         json.Number(fmt.Sprintf("%d", alert.Timestamp.Unix())),
		Footer:     "Byrd",
		FooterIcon: "https://raw.githubusercontent.com/egonelbre/gophers/master/.thumb/vector/computer/music.png",
	}
}

func (s *slackAlertClient) SendBatchAlert(ctx context.Context, alerts []models.Alert) error {
	for _, alert := range alerts {
		if err := s.SendAlert(ctx, alert); err != nil {
			return ErrFailedToSendBatchAlert
		}
	}
	return nil
}
