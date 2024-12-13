package interfaces

import (
	"context"
	"image"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type CompetitorRepository interface {
	Create(ctx context.Context, competitor *models.Competitor) error
	Update(ctx context.Context, competitor *models.Competitor) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.Competitor, error)
	List(ctx context.Context, limit, offset int) ([]models.Competitor, int, error)
	FindByURLHash(ctx context.Context, hash string) ([]models.Competitor, error)
	AddURL(ctx context.Context, competitorID int, url string) error
	RemoveURL(ctx context.Context, competitorID int, url string) error
}

type DiffRepository interface {
	// SaveDiff saves the diff analysis for the given URL
	SaveDiff(ctx context.Context, url string, diff *models.URLDiffAnalysis) error

	// GetDiff retrieves the diff analysis for the given URL, week day, and week number
	GetDiff(ctx context.Context, url, weekDay, weekNumber string) (*models.URLDiffAnalysis, error)
}

type StorageRepository interface {
	// StoreScreenshotImage stores screenshot image in the storage
	StoreScreenshotImage(ctx context.Context, data image.Image, path string, metadata models.ScreenshotMetadata) error

	// StoreScreenshotContent stores screenshot content in the storage
	StoreScreenshotContent(ctx context.Context, content string, path string, metadata models.ScreenshotMetadata) error

	// GetContent retrieves a text content from the storage
	// Serialize the content to a string and return it
	GetScreenshotContent(ctx context.Context, path string) (string, models.ScreenshotMetadata, error)

	// GetScreenshot retrieves a screenshot from the storage
	// Deserialize the content to an image and return it
	GetScreenshotImage(ctx context.Context, path string) (image.Image, models.ScreenshotMetadata, error)

	// Get retrieves a binary from the storage
	// Return the binary content and the metadata
	Get(ctx context.Context, path string) ([]byte, map[string]string, error)

	// Delete deletes a file from the storage
	// Return an error if the file does not exist or cannot be deleted
	Delete(ctx context.Context, path string) error
}

type SubscriptionRepository interface {
	Subscribe(ctx context.Context, competitorID int, email string) error
	Unsubscribe(ctx context.Context, competitorID int, email string) error
	GetSubscribersByCompetitor(ctx context.Context, competitorID int) ([]string, error)
}
