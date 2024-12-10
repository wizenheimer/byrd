package notification

import (
	"context"
	"net/http"

	"github.com/wizenheimer/iris/internal/config"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type resendEmailClient struct {
	authKey string
	client  *http.Client
}

func NewResendEmailClient(config *config.Config, client *http.Client) (interfaces.EmailClient, error) {
	return &resendEmailClient{
		authKey: config.Services.ResendAPIKey,
		client:  client,
	}, nil
}

func (c *resendEmailClient) Send(ctx context.Context, params models.EmailParams) error {
	// Implementation
	return nil
}
