package interfaces

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

type CompetitorService interface {
	Create(ctx context.Context, input models.CompetitorInput) (*models.Competitor, error)
	Update(ctx context.Context, id int, input models.CompetitorInput) (*models.Competitor, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Competitor, error)
	List(ctx context.Context, limit, offset int) ([]models.Competitor, int, error)
	FindByURLHash(ctx context.Context, hash string) ([]models.Competitor, error)
	AddURL(ctx context.Context, id int, url string) error
	RemoveURL(ctx context.Context, id int, url string) error
}
