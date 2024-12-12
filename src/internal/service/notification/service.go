package notification

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type notificationService struct {
	emailClient client.EmailClient
	templates   *TemplateManager
	logger      *logger.Logger
}

func NewNotificationService(
	emailClient client.EmailClient,
	templates *TemplateManager,
	logger *logger.Logger,
) (interfaces.NotificationService, error) {
	logger.Debug("creating new notification service")

	return &notificationService{
		emailClient: emailClient,
		templates:   templates,
		logger:      logger.WithFields(map[string]interface{}{"module": "notification_service"}),
	}, nil
}

func (s *notificationService) SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResults, error) {
	s.logger.Debug("sending notification", zap.Any("emails", req.Emails))
	// Implementation
	emailNotificationResults := models.EmailNotificationResults{
		Successful: []string{},
		Failed:     []string{},
	}
	return &models.NotificationResults{
		EmailNotificationResults: emailNotificationResults,
	}, nil
}
