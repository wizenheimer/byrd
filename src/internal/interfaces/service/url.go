package interfaces

import (
	"context"

	"github.com/google/uuid"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// URLService is the interface that provides URL management operations
type URLService interface {
	// AddURL: adds a new URL if it does not exist
	AddURL(ctx context.Context, url string) (*core_models.URL, error)

	// ListURLs: lists all URLs
	ListURLs(ctx context.Context, batchSize int, lastSeenID *uuid.UUID) (<-chan core_models.URLBatch, <-chan error)

	// DeleteURL: deletes a URL
	DeleteURL(ctx context.Context, url string) error

	// URLExists: checks if a URL exists
	URLExists(ctx context.Context, url string) (bool, error)
}
