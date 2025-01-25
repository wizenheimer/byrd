// ./src/internal/service/page/interface.go
package page

import (
	"context"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

var (
	maxPageBatchSize int = 25
)

type PageService interface {
	CreatePage(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error)

	GetPage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID) (*models.Page, error)

	ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error)

	UpdatePage(ctx context.Context, competitorID uuid.UUID, pageID uuid.UUID, page models.PageProps) (*models.Page, error)

	ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error)

	ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (<-chan []uuid.UUID, <-chan error)

	RemovePage(ctx context.Context, competitorIDs []uuid.UUID, pageIDs []uuid.UUID) error

	PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error)

	RefreshPage(ctx context.Context, pageID uuid.UUID) error
}
