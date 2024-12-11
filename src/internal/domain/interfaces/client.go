package interfaces

import (
	"context"
	"net/http"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type EmailClient interface {
	Send(ctx context.Context, params models.EmailParams) error
}
