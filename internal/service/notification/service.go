package notification

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type notificationService struct {
	emailClient interfaces.EmailClient
	templates   *TemplateManager
}

func NewNotificationService(
	emailClient interfaces.EmailClient,
	templates *TemplateManager,
) (interfaces.NotificationService, error) {
	return &notificationService{
		emailClient: emailClient,
		templates:   templates,
	}, nil
}

func (s *notificationService) SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResults, error) {
	// Implementation
	emailNotificationResults := models.EmailNotificationResults{
		Successful: []string{},
		Failed:     []string{},
	}
	return &models.NotificationResults{
		EmailNotificationResults: emailNotificationResults,
	}, nil
}
