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
	SaveDiff(ctx context.Context, diff *models.DiffAnalysis, metadata models.DiffMetadata) error
	GetDiffHistory(ctx context.Context, params models.DiffHistoryParams) ([]models.DiffReport, error)
	GetLatestDiff(ctx context.Context, url string) (*models.DiffReport, error)
}

type StorageRepository interface {
	StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error
	StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error
	Get(ctx context.Context, path string) ([]byte, map[string]string, error)
	Delete(ctx context.Context, path string) error
}

type SubscriptionRepository interface {
	Subscribe(ctx context.Context, competitorID int, email string) error
	Unsubscribe(ctx context.Context, competitorID int, email string) error
	GetSubscribersByCompetitor(ctx context.Context, competitorID int) ([]string, error)
}
