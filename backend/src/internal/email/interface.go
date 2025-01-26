// ./src/internal/email/interface.go
// ./src/internal/interfaces/client/email.go
package email

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type EmailClient interface {
	Send(ctx context.Context, email models.Email) error
}
