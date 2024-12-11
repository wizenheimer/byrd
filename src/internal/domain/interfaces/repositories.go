package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/models"
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
	// StoreScreenshot stores a screenshot in the storage
	StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error

	// GetScreenshot retrieves a screenshot from the storage
	StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error

	// Get retrieves the content of a file from the storage
	Get(ctx context.Context, path string) ([]byte, map[string]string, error)

	// Delete deletes a file from the storage
	Delete(ctx context.Context, path string) error
}

type SubscriptionRepository interface {
	Subscribe(ctx context.Context, competitorID int, email string) error
	Unsubscribe(ctx context.Context, competitorID int, email string) error
	GetSubscribersByCompetitor(ctx context.Context, competitorID int) ([]string, error)
}
