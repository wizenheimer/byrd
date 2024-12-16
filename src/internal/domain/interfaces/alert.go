package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

// AlertClient defines alerting capabilities
type AlertClient interface {
	// Send sends an alert
	Send(ctx context.Context, alert models.Alert) error

	// SendBatch sends a batch of alerts
	SendBatch(ctx context.Context, alerts []models.Alert) error
}
