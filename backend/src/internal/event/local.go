package event

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type localEventClient struct {
	logger *logger.Logger
}

func NewLocalEventClient(logger *logger.Logger) EventClient {
	return &localEventClient{
		logger: logger.WithFields(map[string]interface{}{
			"module": "local_event_client",
		}),
	}
}

func (lec *localEventClient) SendEvent(ctx context.Context, event models.Event) error {
	lec.logger.Debug("received an event", zap.Any("event", event))
	return nil
}

func (lec *localEventClient) SendBatchEvent(ctx context.Context, events []models.Event) error {
	for _, event := range events {
		_ = lec.SendEvent(ctx, event)
	}
	return nil
}
