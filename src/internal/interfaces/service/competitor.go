package interfaces

import (
	"context"

	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

type CompetitorService interface {
	Create(ctx context.Context, input core_models.CompetitorInput) (*core_models.Competitor, error)
	Update(ctx context.Context, id int, input core_models.CompetitorInput) (*core_models.Competitor, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*core_models.Competitor, error)
	List(ctx context.Context, limit, offset int) ([]core_models.Competitor, int, error)
	FindByURLHash(ctx context.Context, hash string) ([]core_models.Competitor, error)
	AddURL(ctx context.Context, id int, url string) error
	RemoveURL(ctx context.Context, id int, url string) error
}
