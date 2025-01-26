// ./src/internal/service/screenshot/url.go
package screenshot

import (
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type ContentType string

const (
	ContentTypeImage   ContentType = "image"
	ContentTypeContent ContentType = "content"
)

func DeterminePath(opts models.ScreenshotRequestOptions, contentType ContentType, backDate bool) (string, error) {
	switch contentType {
	case ContentTypeImage:
		// Get the current screenshot path
		if backDate {
			return getPreviousScreenshotPath(opts)
		} else {
			return getCurrentScreenshotPath(opts)
		}
	case ContentTypeContent:
		// Get the current content path
		if backDate {
			return getPreviousContentPath(opts)
		} else {
			return getCurrentContentPath(opts)
		}
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}
