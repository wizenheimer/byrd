package interfaces

import (
	"context"

	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// NotificationService defines notification capabilities
type NotificationService interface {
	// SendNotification sends a notification
	SendNotification(ctx context.Context, req core_models.NotificationRequest) (*core_models.NotificationResults, error)
}
