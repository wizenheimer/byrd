package email

import (
	"context"
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type localEmailClient struct {
	logger *logger.Logger
}

func NewLocalEmailClient(ctx context.Context, logger *logger.Logger) (EmailClient, error) {
	lc := localEmailClient{
		logger: logger.WithFields(map[string]any{
			"module": "local_email_client",
		}),
	}

	return &lc, nil
}

func (lc *localEmailClient) Send(ctx context.Context, email models.Email) error {
	// Mock email sending latency
	time.Sleep(1 * time.Second)
	return nil
}
