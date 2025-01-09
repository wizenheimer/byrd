package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/slack-go/slack"
	clf "github.com/wizenheimer/byrd/src/internal/interfaces/client"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
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
	// For deduplication
	mu     sync.RWMutex
	recent map[string]time.Time
}

func NewSlackAlertClient(config models.SlackConfig, logger *logger.Logger) (clf.AlertClient, error) {
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
		recent:  make(map[string]time.Time),
	}

	go client.cleanupLoop()
	return client, nil
}

func (s *slackAlertClient) Send(ctx context.Context, alert models.Alert) error {
	// Rate limiting
	if err := s.limiter.Wait(ctx); err != nil {
		return ErrEncounteredSlackRateLimit
	}

	// Check for duplicates
	key := fmt.Sprintf("%s-%s-%s", alert.Title, alert.Description, alert.Severity)
	if s.isDuplicate(key) {
		s.logger.Debug("skipping duplicate alert", zap.Any("alert", alert))
		return nil
	}

	// Send with retries
	err := s.sendWithRetries(ctx, alert)
	if err != nil {
		return err
	}

	// Record successful message
	s.recordMessage(key)
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

func (s *slackAlertClient) isDuplicate(key string) bool {
	s.mu.RLock()
	lastSent, exists := s.recent[key]
	s.mu.RUnlock()

	if !exists {
		return false
	}
	return time.Since(lastSent) < 10*time.Minute
}

func (s *slackAlertClient) recordMessage(key string) {
	s.mu.Lock()
	s.recent[key] = time.Now()
	s.mu.Unlock()
}

func (s *slackAlertClient) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, timestamp := range s.recent {
			if now.Sub(timestamp) > 1*time.Minute {
				delete(s.recent, key)
			}
		}
		s.mu.Unlock()
	}
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
		Footer:     "Byrd Alerts",
		FooterIcon: "https://platform.slack-edge.com/img/default_application_icon.png",
	}
}

func (s *slackAlertClient) SendBatch(ctx context.Context, alerts []models.Alert) error {
	for _, alert := range alerts {
		if err := s.Send(ctx, alert); err != nil {
			return ErrFailedToSendBatchAlert
		}
	}
	return nil
}
