package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

// NotificationService defines notification capabilities
type NotificationService interface {
	// SendNotification sends a notification
	SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResults, error)
}
