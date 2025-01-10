package interfaces

import (
	"context"
	"errors"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
)

type ScreenshotService interface {
	// Refresh refreshes the screenshot repository with the latest screenshot and html content
	// url is the URL to take a screenshot of
	// opts are the options for the screenshot
	Refresh(ctx context.Context, url string, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, errs.Error)

	// Retrieve retrieves previous screenshot and html content from the storage
	// url is the URL to retrieve the screenshot from
	Retrieve(ctx context.Context, url string) (*models.ScreenshotImageResponse, *models.ScreenshotHTMLContentResponse, errs.Error)

	// GetCurrentImage retrieves the current screenshot from the storage if present
	// Or it will take a new screenshot and store it as an image
	GetCurrentImage(ctx context.Context, save bool, opts models.ScreenshotRequestOptions) (*models.ScreenshotImageResponse, errs.Error)

	// GetCurrentHTML retrieves the current html content from the storage if present
	// Or it will take a new screenshot and store it as html
	GetCurrentHTMLContent(ctx context.Context, save bool, opts models.ScreenshotHTMLRequestOptions) (*models.ScreenshotHTMLContentResponse, errs.Error)

	// GetPreviousImage retrieves the previous screenshot from the storage
	GetPreviousImage(ctx context.Context, url string) (*models.ScreenshotImageResponse, errs.Error)

	// GetPreviousHTML retrieves the previous content of a screenshot from the storage
	GetPreviousHTMLContent(ctx context.Context, url string) (*models.ScreenshotHTMLContentResponse, errs.Error)

	// GetImage retrieves a screenshot from the storage
	GetImage(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotImageResponse, errs.Error)

	// GetHTMLContent retrieves the content of a screenshot from the storage
	GetHTMLContent(ctx context.Context, url string, year int, weekNumber int, weekDay int) (*models.ScreenshotHTMLContentResponse, errs.Error)

	// ListScreenshots lists the latest content (images or text) for a given URL
	ListScreenshots(ctx context.Context, url string, contentType string, maxItems int) ([]models.ScreenshotListResponse, errs.Error)
}

var (
	ErrStorageRepositoryRequired = errors.New("storage repository is required")
	ErrHTTPClientRequired        = errors.New("HTTP client is required")
	ErrScreenshotKeyRequired     = errors.New("screenshot key is required")
	ErrScreenshotOriginRequired  = errors.New("screenshot origin is required")

	ErrFailedToConvertCurrentDayToInt   = errors.New("failed to convert current day to int")
	ErrFailedToGetCurrentScreenshotPath = errors.New("failed to get current screenshot path")
	ErrFailedToGetCurrentContentPath    = errors.New("failed to get current content path")
	ErrFailedToStoreScreenshotImage     = errors.New("failed to store screenshot image")

	ErrFailedToGetPreviousScreenshotPath  = errors.New("failed to get previous screenshot path")
	ErrFailedToGetPreviousScreenshotImage = errors.New("failed to get previous screenshot image")

	ErrFailedToGetPreviousContentPath = errors.New("failed to get previous content path")
	ErrFailedToGetPreviousContent     = errors.New("failed to get previous content")

	ErrFailedToGetScreenshotPath  = errors.New("failed to get screenshot path")
	ErrFailedToGetScreenshotImage = errors.New("failed to get screenshot image")

	ErrFailedToGetContentPath = errors.New("failed to get content path")
	ErrFailedToGetContent     = errors.New("failed to get content")

	ErrFailedToGetListingPrefixFromContentType = errors.New("failed to get listing prefix from content type")
)

var (
	ErrNoContentURLFoundInHeaders = errors.New("no content URL found in headers")
	ErrNon200StatusCode           = errors.New("non-200 status code")
	ErrReadingHTMLBody            = errors.New("error reading HTML body")
	ErrUnexpectedContentType      = errors.New("received unexpected content type")
	ErrCannotDecodeImage          = errors.New("cannot decode image")
	ErrFailedToGetHTMLContent     = errors.New("failed to get HTML content")
)
