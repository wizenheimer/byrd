// ./src/internal/service/competitor/service.go
package competitor

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/internal/service/report"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type competitorService struct {
	competitorRepository competitor.CompetitorRepository
	reportService        report.ReportService
	pageService          page.PageService
	tm                   *transaction.TxManager
	nameFinder           *CompanyNameFinder
	logger               *logger.Logger
}

func NewCompetitorService(pageService page.PageService, reportService report.ReportService, tm *transaction.TxManager, competitorRepository competitor.CompetitorRepository, logger *logger.Logger) CompetitorService {
	return &competitorService{
		competitorRepository: competitorRepository,
		reportService:        reportService,
		pageService:          pageService,
		logger:               logger,
		tm:                   tm,
		nameFinder:           NewCompanyNameFinder(logger),
	}
}

func (cs *competitorService) CreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) (models.Competitor, error) {
	var competitor models.Competitor
	var urls []string
	for _, page := range pages {
		urls = append(urls, page.URL)
	}
	competitorName := cs.nameFinder.FindCompanyName(urls)

	// Utility function to create a competitor
	createCompetitor := func(ctx context.Context) (*models.Competitor, error) {
		// Create a new competitor using the competitor's name
		c, err := cs.competitorRepository.CreateCompetitorForWorkspace(
			ctx,
			workspaceID,
			competitorName,
		)
		if err != nil {
			return nil, err
		}

		if c == nil {
			return nil, errors.New("failed to create competitor")
		}

		// Create a page, and associate it with the created competitor
		if _, err = cs.pageService.CreatePage(
			ctx,
			c.ID,
			pages,
		); err != nil {
			return nil, err
		}

		return c, nil
	}

	// Run the transaction
	err := cs.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		c, err := createCompetitor(ctx)
		if err != nil {
			return err
		}
		competitor = *c
		return nil
	})

	if err != nil {
		return models.Competitor{}, err
	}

	return competitor, nil
}

func (cs *competitorService) BatchCreateCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, pages []models.PageProps) ([]models.Competitor, error) {
	if len(pages) == 0 {
		return nil, errors.New("non-fatal: pages unspecified for creating competitors")
	}

	var competitors []models.Competitor
	for _, page := range pages {
		// Figure out competitor's name using the url
		competitorName := cs.nameFinder.FindCompanyName([]string{
			page.URL,
		})

		var competitor *models.Competitor
		err := cs.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
			// Create a new competitor using the competitor's name
			var err error
			competitor, err = cs.competitorRepository.CreateCompetitorForWorkspace(
				ctx,
				workspaceID,
				competitorName,
			)
			if err != nil {
				return err
			}

			// Create a page, and associate it with the created competitor
			_, err = cs.pageService.CreatePage(
				ctx,
				competitor.ID,
				[]models.PageProps{
					page,
				},
			)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			cs.logger.Error("failed to create competitor for workspace", zap.Error(err), zap.Any("workspaceID", workspaceID), zap.Any("pageURL", page.URL))
			continue
		}
		if competitor == nil {
			cs.logger.Error("failed to create competitor for workspace", zap.Any("workspaceID", workspaceID), zap.Any("pageURL", page.URL))
			continue
		}

		competitors = append(competitors, *competitor)
	}

	totalPages := len(pages)
	totalCompetitors := len(competitors)

	if totalCompetitors == 0 {
		return competitors, errors.New("failed to create competitors")
	}

	if totalPages > totalCompetitors {
		return competitors, errors.New("non-fatal: failed to create some competitors")
	}

	return competitors, nil
}

func (cs *competitorService) GetCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) ([]models.Competitor, error) {
	if len(competitorIDs) == 0 {
		return nil, errors.New("non-fatal: no competitorIDs provided")
	} else if len(competitorIDs) > maxCompetitorBatchSize {
		return nil, errors.New("non-fatal: too many competitorIDs provided")
	}
	return cs.competitorRepository.BatchGetCompetitorsForWorkspace(
		ctx,
		workspaceID,
		competitorIDs,
	)
}

func (cs *competitorService) ListCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset *int) ([]models.Competitor, bool, error) {
	return cs.competitorRepository.ListCompetitorsForWorkspace(
		ctx,
		workspaceID,
		limit,
		offset,
	)
}

func (cs *competitorService) UpdateCompetitorForWorkspace(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string) (*models.Competitor, error) {
	return cs.competitorRepository.UpdateCompetitorForWorkspace(
		ctx,
		workspaceID,
		competitorID,
		competitorName,
	)
}

func (cs *competitorService) RemoveCompetitorForWorkspace(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error {
	if len(competitorIDs) > maxCompetitorBatchSize {
		return errors.New("non-fatal: too many competitorIDs provided")
	}

	if competitorIDs == nil {
		// In case competitorIDs are unspecified, get the list of competitors
		competitors, _, err := cs.competitorRepository.ListCompetitorsForWorkspace(
			ctx,
			workspaceID,
			nil,
			nil,
		)
		if err != nil {
			return err
		}
		if len(competitors) == 0 {
			// No competitors to remove, return from here
			return nil
		}
		// Populate the competitorIDs from the competitors retrieved
		for _, competitor := range competitors {
			competitorIDs = append(competitorIDs, competitor.ID)
		}
	}

	// Utility function to remove competitors
	removeCompetitor := func(ctx context.Context, workspaceID uuid.UUID, competitorIDs []uuid.UUID) error {
		if competitorIDs == nil {
			// Remove all the competitors
			err := cs.competitorRepository.RemoveAllCompetitorsForWorkspace(
				ctx,
				workspaceID,
			)
			if err != nil {
				return err
			}
		} else if len(competitorIDs) == 1 {
			err := cs.competitorRepository.RemoveCompetitorForWorkspace(
				ctx,
				workspaceID,
				competitorIDs[0],
			)
			if err != nil {
				return err
			}
		} else {
			err := cs.competitorRepository.BatchRemoveCompetitorForWorkspace(
				ctx,
				workspaceID,
				competitorIDs,
			)
			if err != nil {
				return err
			}
			return nil
		}

		return nil
	}

	// Run the transaction
	return cs.tm.RunInTx(context.Background(), nil, func(ctx context.Context) error {
		// removeCompetitor is a helper for handling competitor removal
		err := removeCompetitor(ctx, workspaceID, competitorIDs)
		if err != nil {
			return err
		}
		// Check if competitorIDs is empty
		if len(competitorIDs) == 0 {
			// This means pages for competitors are already removed, including their pages
			return nil
		}
		return cs.pageService.RemovePage(ctx, competitorIDs, nil)
	})
}

func (cs *competitorService) CompetitorExists(ctx context.Context, workspaceID, competitorID uuid.UUID) (bool, error) {
	return cs.competitorRepository.WorkspaceCompetitorExists(
		ctx,
		workspaceID,
		competitorID,
	)
}

func (cs *competitorService) PageExists(ctx context.Context, competitorID, pageID uuid.UUID) (bool, error) {
	return cs.pageService.PageExists(
		ctx,
		competitorID,
		pageID,
	)
}

func (cs *competitorService) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	return cs.pageService.CreatePage(
		ctx,
		competitorID,
		pages,
	)
}

func (cs *competitorService) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	return cs.pageService.GetPage(
		ctx,
		competitorID,
		pageID,
	)
}

func (cs *competitorService) UpdatePage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	return cs.pageService.UpdatePage(
		ctx,
		competitorID,
		pageID,
		page,
	)
}

func (cs *competitorService) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageID []uuid.UUID) error {
	return cs.pageService.RemovePage(
		ctx,
		[]uuid.UUID{competitorID},
		pageID,
	)
}

func (cs *competitorService) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	return cs.pageService.ListCompetitorPages(
		ctx,
		competitorID,
		limit,
		offset,
	)
}

func (cs *competitorService) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	return cs.pageService.ListPageHistory(
		ctx,
		pageID,
		limit,
		offset,
	)
}

// ListReports returns a list of reports for a competitor.
// The limit and offset parameters are used for pagination.
func (cs *competitorService) ListReports(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	if limit == nil {
		return nil, false, errors.New("limit is required")
	}
	if offset == nil {
		return nil, false, errors.New("offset is required")
	}

	reports, hasMore, err := cs.reportService.List(
		ctx,
		workspaceID,
		competitorID,
		limit,
		offset,
	)
	if err != nil {
		return nil, false, err
	}

	return reports, hasMore, nil
}

// CreateReport creates a new report for a competitor.
// The report is generated based on the latest page history for the competitor.
func (cs *competitorService) CreateReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) (*models.Report, error) {

	// Get the competitor
	competitor, err := cs.GetCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil {
		return nil, err
	}
	if len(competitor) == 0 {
		return nil, errors.New("competitor not found")
	}

	// List all the active pages for the competitor
	pages, _, err := cs.pageService.ListCompetitorPages(ctx, competitorID, nil, nil)
	if err != nil {
		return nil, err
	}

	var pageIDs []uuid.UUID
	for _, page := range pages {
		pageIDs = append(pageIDs, page.ID)
	}

	// Get the latest page history for the pages
	pageHistories, err := cs.pageService.GetLatestPageHistory(ctx, pageIDs)
	if err != nil {
		return nil, err
	}

	// Create a new report
	report, err := cs.reportService.Create(ctx, workspaceID, competitorID, pageHistories)
	if err != nil {
		return nil, err
	}

	// Return the created report
	return report, nil
}

// DispatchReport sends the report to the subscribers.
func (cs *competitorService) DispatchReport(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID, subscriberEmails []string) error {
	// Get the competitor
	competitor, err := cs.GetCompetitorForWorkspace(ctx, workspaceID, []uuid.UUID{competitorID})
	if err != nil {
		return err
	}
	if len(competitor) == 0 {
		return errors.New("competitor not found")
	}

	// Send the report to the subscribers
	if err := cs.reportService.Dispatch(ctx, workspaceID, competitorID, competitor[0].Name, subscriberEmails); err != nil {
		return err
	}

	return nil
}

func (cs *competitorService) CountPagesForCompetitors(ctx context.Context, competitorIDs []uuid.UUID) (int, error) {
	return cs.pageService.CountActivePagesForCompetitors(ctx, competitorIDs)
}

func (cs *competitorService) CountCompetitorsForWorkspace(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	return cs.competitorRepository.GetActiveCompetitorCount(ctx, workspaceID)
}
