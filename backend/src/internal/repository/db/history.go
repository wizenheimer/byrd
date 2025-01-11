// ./src/internal/repository/db/history.go
package db

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type historyRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageHistoryRepository(tm *transaction.TxManager, logger *logger.Logger) repo.PageHistoryRepository {
	return &historyRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "history_repository"}),
	}
}

// CreatePageHistory creates a new page history entry
func (r *historyRepo) CreatePageHistory(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, errs.Error) {
	pageHistoryErr := errs.New()
	querier := r.tm.GetQuerier(ctx)

	if err := utils.SetDefaultsAndValidate(&pageHistory); err != nil {
		pageHistoryErr.Add(repo.ErrInvalidPageHistory, map[string]interface{}{
			"pageHistory": pageHistory,
		})
		return models.PageHistory{}, pageHistoryErr.Propagate(repo.ErrFailedToCreatePageHistoryInPageHistoryRepository)
	}

	query := `
        INSERT INTO page_history (
            page_id,
            week_number_1,
            week_number_2,
            year_number_1,
            year_number_2,
            bucket_id_1,
            bucket_id_2,
            diff_content,
            screenshot_url_1,
            screenshot_url_2,
            status
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, page_id, week_number_1, week_number_2, year_number_1, year_number_2,
                  bucket_id_1, bucket_id_2, diff_content, screenshot_url_1, screenshot_url_2,
                  created_at, status
    `

	row := querier.QueryRow(ctx, query,
		pageID,
		pageHistory.WeekNumber1,
		pageHistory.WeekNumber2,
		pageHistory.YearNumber1,
		pageHistory.YearNumber2,
		pageHistory.BucketID1,
		pageHistory.BucketID2,
		pageHistory.DiffContent,
		pageHistory.ScreenshotURL1,
		pageHistory.ScreenshotURL2,
		models.HistoryStatusActive,
	)

	var created models.PageHistory
	err := row.Scan(
		&created.ID,
		&created.PageID,
		&created.WeekNumber1,
		&created.WeekNumber2,
		&created.YearNumber1,
		&created.YearNumber2,
		&created.BucketID1,
		&created.BucketID2,
		&created.DiffContent,
		&created.ScreenshotURL1,
		&created.ScreenshotURL2,
		&created.CreatedAt,
		&created.Status,
	)
	if err != nil {
		pageHistoryErr.Add(repo.ErrFailedToScanPageHistory, map[string]interface{}{
			"pageID": pageID,
			"error":  err,
		})
		return models.PageHistory{}, pageHistoryErr.Propagate(repo.ErrFailedToCreatePageHistoryInPageHistoryRepository)
	}

	return created, nil
}

// ListPageHistory lists page history for a page with optional pagination
func (r *historyRepo) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, errs.Error) {
	querier := r.tm.GetQuerier(ctx)
	pageHistoryErr := errs.New()

	// Build the query with optional pagination
	query := `
        SELECT id, page_id, week_number_1, week_number_2, year_number_1, year_number_2, bucket_id_1, bucket_id_2, diff_content, screenshot_url_1, screenshot_url_2, created_at, status
        FROM page_history
        WHERE page_id = $1
        AND status = $2
        ORDER BY created_at DESC
    `
	args := []interface{}{pageID, models.HistoryStatusActive}

	// Add pagination if provided
	if limit != nil && offset != nil {
		if *limit < 0 {
			pageHistoryErr.Add(repo.ErrInvalidLimit, map[string]interface{}{
				"limit": *limit,
			})
		}
		if *offset < 0 {
			pageHistoryErr.Add(repo.ErrInvalidOffset, map[string]interface{}{
				"offset": *offset,
			})
		}
		if pageHistoryErr.HasErrors() {
			return nil, pageHistoryErr.Propagate(repo.ErrFailedToListPageHistoryFromPageHistoryRepository)
		}
		query += " LIMIT $3 OFFSET $4"
		args = append(args, *limit, *offset)
	}

	rows, err := querier.Query(ctx, query, args...)
	if err != nil {
		pageHistoryErr.Add(err, map[string]interface{}{
			"pageID": pageID,
		})
		return nil, pageHistoryErr.Propagate(repo.ErrFailedToListPageHistoryFromPageHistoryRepository)
	}
	defer rows.Close()

	history := make([]models.PageHistory, 0)
	for rows.Next() {
		var h models.PageHistory
		err := rows.Scan(
			&h.ID,
			&h.PageID,
			&h.WeekNumber1,
			&h.WeekNumber2,
			&h.YearNumber1,
			&h.YearNumber2,
			&h.BucketID1,
			&h.BucketID2,
			&h.DiffContent,
			&h.ScreenshotURL1,
			&h.ScreenshotURL2,
			&h.CreatedAt,
			&h.Status,
		)
		if err != nil {
			pageHistoryErr.Add(err, map[string]interface{}{
				"pageID": pageID,
			})
			continue
		}
		history = append(history, h)
	}

	if err = rows.Err(); err != nil {
		pageHistoryErr.Add(err, map[string]interface{}{
			"pageID": pageID,
		})
	}

	return history, pageHistoryErr.Propagate(repo.ErrFailedToListPageHistoryFromPageHistoryRepository)
}

// RemovePageHistory removes page history for a list of pages
func (r *historyRepo) RemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) errs.Error {
	pageHistoryErr := errs.New()
	if len(pageIDs) == 0 {
		pageHistoryErr.Add(repo.ErrPageIDsUnspecified, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr.Propagate(repo.ErrFailedToRemovePageHistoryFromPageHistoryRepository)
	}

	querier := r.tm.GetQuerier(ctx)

	query := `
        UPDATE page_history
        SET status = $1
        WHERE page_id = ANY($2)
        AND status = $3
    `

	_, err := querier.Exec(ctx, query,
		models.HistoryStatusInactive,
		pageIDs, // pgx handles array parameters natively
		models.HistoryStatusActive,
	)
	if err != nil {
		pageHistoryErr.Add(err, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr.Propagate(repo.ErrFailedToRemovePageHistoryFromPageHistoryRepository)
	}

	return nil
}
