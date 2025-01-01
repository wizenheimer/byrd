package interfaces

import (
	"context"

	"github.com/google/uuid"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// URLRepository is the interface that provides URL management operations
type URLRepository interface {
	// AddURL: adds a new URL if it does not exist
	AddURL(ctx context.Context, url string) (*core_models.URL, error)

	// ListURLs: lists all URLs in batches
	ListURLs(ctx context.Context, batchSize int, lastSeenID *uuid.UUID) (*core_models.URLBatch, error)

	// DeleteURL: deletes a URL
	DeleteURL(ctx context.Context, url string) error

	// URLExists: checks if a URL exists
	URLExists(ctx context.Context, url string) (*core_models.URL, bool, error)
}
