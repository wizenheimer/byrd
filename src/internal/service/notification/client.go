package notification

import (
	"context"
	"net/http"

	"github.com/wizenheimer/iris/src/internal/config"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type resendEmailClient struct {
	authKey string
	client  *http.Client
	logger  *logger.Logger
}

func NewResendEmailClient(config *config.Config, client *http.Client, logger *logger.Logger) (interfaces.EmailClient, error) {
	logger.Debug("creating new resend email client")

	return &resendEmailClient{
		authKey: config.Services.ResendAPIKey,
		client:  client,
		logger:  logger.WithFields(map[string]interface{}{"module": "resend_email_client"}),
	}, nil
}

func (c *resendEmailClient) Send(ctx context.Context, params models.EmailParams) error {
	c.logger.Debug("sending email", zap.Any("from", params.From), zap.Any("to", params.To), zap.Any("subject", params.Subject))
	// Implementation
	return nil
}
