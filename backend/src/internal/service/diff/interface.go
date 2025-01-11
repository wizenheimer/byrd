// ./src/internal/interfaces/service/diff.go
package diff

import (
	"context"
	"errors"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// DiffService is the interface that provides diff operations
type DiffService interface {
	// Compare: compares two HTML contents and returns the differences using the given profile
	Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, error)
}

var (
	ErrFailedToListPageHistoryForPage = errors.New("failed to list page history for page")

	ErrFailedToClearPageHistoryForPage = errors.New("failed to clear page history for page")
)
