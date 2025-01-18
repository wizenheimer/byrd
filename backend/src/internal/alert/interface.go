package alert

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// AlertClient defines alerting capabilities
type AlertClient interface {
	// Send sends an alert
	SendAlert(ctx context.Context, alert models.Alert) error

	// SendBatch sends a batch of alerts
	SendBatchAlert(ctx context.Context, alerts []models.Alert) error
}
