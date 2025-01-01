package notification

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/config"
	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type resendEmailClient struct {
	authKey string
	client  clf.HTTPClient
	logger  *logger.Logger
}

func NewResendEmailClient(config *config.Config, client clf.HTTPClient, logger *logger.Logger) (clf.EmailClient, error) {
	logger.Debug("creating new resend email client")

	return &resendEmailClient{
		authKey: config.Services.ResendAPIKey,
		client:  client,
		logger:  logger.WithFields(map[string]interface{}{"module": "resend_email_client"}),
	}, nil
}

func (c *resendEmailClient) Send(ctx context.Context, params core_models.EmailParams) error {
	c.logger.Debug("sending email", zap.Any("from", params.From), zap.Any("to", params.To), zap.Any("subject", params.Subject))
	// Implementation
	return nil
}
