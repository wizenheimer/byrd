// ./src/internal/interfaces/service/page.go
package page

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type PageService interface {
	CreatePage(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error)

	ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error)

	UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)

	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, error)

	ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []models.Page, <-chan error)

	RemovePage(ctx context.Context, competitorIDs []uuid.UUID, pageIDs []uuid.UUID) error

	PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error)
}
