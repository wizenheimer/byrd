package interfaces

import (
	"context"

	models "github.com/wizenheimer/iris/src/internal/models/core"
)

type EmailClient interface {
	Send(ctx context.Context, params models.EmailParams) error
}
