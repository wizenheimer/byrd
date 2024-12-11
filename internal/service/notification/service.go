package notification

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/pkg/logger"
	"go.uber.org/zap"
)

type notificationService struct {
	emailClient interfaces.EmailClient
	templates   *TemplateManager
	logger      *logger.Logger
}

func NewNotificationService(
	emailClient interfaces.EmailClient,
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
