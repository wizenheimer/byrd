// ./src/internal/event/interface.go
package event

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type EventClient interface {
	// Send sends an event
	SendEvent(ctx context.Context, event models.Event) error

	// SendBatch sends a batch of events
	SendBatchEvent(ctx context.Context, events []models.Event) error
}
