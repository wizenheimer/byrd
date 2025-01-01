package interfaces

import (
	"context"

	models "github.com/wizenheimer/iris/src/internal/models/core"
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
