package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
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
func (r *historyRepo) CreatePageHistory(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, err.Error) {
	pageHistoryErr := err.New()
	runner := r.tm.GetRunner(ctx)

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
			screenshot_url_2
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, page_id, week_number_1, week_number_2, year_number_1, year_number_2,
				  bucket_id_1, bucket_id_2, diff_content, screenshot_url_1, screenshot_url_2,
				  created_at
	`

	row := runner.QueryRowContext(ctx, query,
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
	)

	if err := row.Err(); err != nil {
		pageHistoryErr.Add(err, map[string]interface{}{
			"pageID": pageID,
		})
		return models.PageHistory{}, pageHistoryErr
	}

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
	)
	if err != nil {
		pageHistoryErr.Add(repo.ErrFailedToScanPageHistory, map[string]interface{}{
			"pageID": pageID,
		})
		return models.PageHistory{}, pageHistoryErr
	}

	return created, nil
}

// ListPageHistory lists page history for a page with optional pagination
func (r *historyRepo) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, err.Error) {
	runner := r.tm.GetRunner(ctx)
	pageHistoryErr := err.New()

	// Build the query with optional pagination
	query := `
		SELECT id, page_id, week_number_1, week_number_2, year_number_1, year_number_2,
			   bucket_id_1, bucket_id_2, diff_content, screenshot_url_1, screenshot_url_2,
			   created_at
		FROM page_history
		WHERE page_id = $1
		ORDER BY created_at DESC
	`
	args := []interface{}{pageID}

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
			return nil, pageHistoryErr
		}
		query += " LIMIT $2 OFFSET $3"
		args = append(args, *limit, *offset)
	}

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		pageHistoryErr.Add(err, map[string]interface{}{
			"pageID": pageID,
		})
		return nil, pageHistoryErr
	}
	defer rows.Close()

	var history []models.PageHistory
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

	if len(history) == 0 {
		pageHistoryErr.Add(repo.ErrPageHistoryNotFound, map[string]interface{}{
			"pageID": pageID,
		})
	}

	return history, pageHistoryErr
}

// RemovePageHistory removes page history for a list of pages
func (r *historyRepo) RemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) err.Error {
	pageHistoryErr := err.New()
	if len(pageIDs) == 0 {
		pageHistoryErr.Add(repo.ErrPageIDsUnspecified, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr
	}

	runner := r.tm.GetRunner(ctx)

	// Create placeholders for the IN clause
	placeholders := make([]string, len(pageIDs))
	args := make([]interface{}, len(pageIDs))
	for i, id := range pageIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		DELETE FROM page_history
		WHERE page_id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		pageHistoryErr.Add(err, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		pageHistoryErr.Add(err, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr
	}

	if rowsAffected == 0 {
		pageHistoryErr.Add(repo.ErrPageHistoryNotFound, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageHistoryErr
	}

	return nil
}
