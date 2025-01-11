// ./src/internal/interfaces/service/diff.go
package interfaces

import (
	"context"
	"errors"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

// DiffService is the interface that provides diff operations
type DiffService interface {
	// Compare: compares two HTML contents and returns the differences using the given profile
	Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, errs.Error)
}

var (
	ErrFailedToListPageHistoryForPage = errors.New("failed to list page history for page")

	ErrFailedToClearPageHistoryForPage = errors.New("failed to clear page history for page")
)
