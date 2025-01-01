package notification

import (
	"context"

	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type notificationService struct {
	emailClient clf.EmailClient
	templates   *TemplateManager
	logger      *logger.Logger
}

func NewNotificationService(
	emailClient clf.EmailClient,
	templates *TemplateManager,
	logger *logger.Logger,
) (svc.NotificationService, error) {
	logger.Debug("creating new notification service")

	return &notificationService{
		emailClient: emailClient,
		templates:   templates,
		logger:      logger.WithFields(map[string]interface{}{"module": "notification_service"}),
	}, nil
}

func (s *notificationService) SendNotification(ctx context.Context, req core_models.NotificationRequest) (*core_models.NotificationResults, error) {
	s.logger.Debug("sending notification", zap.Any("emails", req.Emails))
	// Implementation
	emailNotificationResults := core_models.EmailNotificationResults{
		Successful: []string{},
		Failed:     []string{},
	}
	return &core_models.NotificationResults{
		EmailNotificationResults: emailNotificationResults,
	}, nil
}
