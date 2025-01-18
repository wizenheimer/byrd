// ./src/internal/interfaces/client/email.go
package interfaces

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type EmailClient interface {
	// Send(ctx context.Context, params models.EmailParams) error

	// Send sends an alert via email
	SendAlert(ctx context.Context, alert models.Alert) error

	// SendBatch sends a batch of alerts via email
	SendBatchAlert(ctx context.Context, alerts []models.Alert) error
}
