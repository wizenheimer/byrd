package client

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type EmailClient interface {
	Send(ctx context.Context, params models.EmailParams) error
}
