package event

import (
	"context"
	"time"

	"github.com/posthog/posthog-go"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type posthogEventClient struct {
	client *posthog.Client
	logger *logger.Logger
}

func NewPostHogEventClient(apiKey string, logger *logger.Logger) (EventClient, error) {
	client, err := posthog.NewWithConfig(apiKey, posthog.Config{
		BatchSize: 20,
		Interval:  30 * time.Second,
		Endpoint:  "https://us.i.posthog.com",
		Logger:    logger.AsPosthogLogger(),
	})

	if err != nil {
		return nil, err
	}

	return &posthogEventClient{
		client: &client,
		logger: logger.WithFields(map[string]interface{}{
			"module": "posthog_event_client",
		}),
	}, nil
}

func (c *posthogEventClient) SendEvent(ctx context.Context, event models.Event) error {
	properties := posthog.NewProperties()
	for key, value := range event.GetProperties() {
		properties.Set(key, value)
	}

	err := (*c.client).Enqueue(posthog.Capture{
		DistinctId: event.GetDistinctID(),
		Event:      string(event.GetEventType()),
		Properties: properties,
	})

	if err != nil {
		c.logger.Error("failed to send event",
			zap.Error(err),
			zap.String("event_type", string(event.GetEventType())),
			zap.String("distinct_id", event.GetDistinctID()))
		return err
	}

	return nil
}

func (c *posthogEventClient) SendBatchEvent(ctx context.Context, events []models.Event) error {
	for _, event := range events {
		if err := c.SendEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (c *posthogEventClient) Close() error {
	return (*c.client).Close()
}
