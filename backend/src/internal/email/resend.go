package email

import (
	"context"

	"github.com/resend/resend-go/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type resendClient struct {
	client            *resend.Client
	logger            *logger.Logger
	notificationEmail string
}

func NewResendClient(ctx context.Context, resendKey, notificationEmail string, logger *logger.Logger) (EmailClient, error) {

	rc := resendClient{
		client:            resend.NewClient(resendKey),
		notificationEmail: notificationEmail,
		logger: logger.WithFields(map[string]any{
			"module": "resend_client",
		}),
	}

	return &rc, nil
}

func (rc *resendClient) Send(ctx context.Context, email models.Email) error {
	params := &resend.SendEmailRequest{
		From:    rc.notificationEmail,
		To:      email.To,
		Subject: email.EmailSubject,
	}

	switch email.EmailFormat {
	case models.EmailFormatHTML:
		params.Html = email.EmailContent
	case models.EmailFormatText:
		fallthrough
	default:
		params.Text = email.EmailContent
	}

	sent, err := rc.client.Emails.Send(params)
	if err != nil {
		rc.logger.Error("failed to send email", zap.Error(err))
		return err
	}

	rc.logger.Debug("email sent", zap.Any("sent", sent))
	return err
}
