package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/internal/repository/transaction"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

var (
	ErrPageHistoryNotFound = errors.New("page history not found")
	ErrNoPageHistory       = errors.New("no page history found")
	ErrInvalidPageIDs      = errors.New("invalid page IDs")
	ErrInvalidPagination   = errors.New("invalid pagination parameters")
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
func (r *historyRepo) CreatePageHistory(ctx context.Context, pageID uuid.UUID, pageHistory models.PageHistory) (models.PageHistory, error) {
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
		return models.PageHistory{}, fmt.Errorf("failed to create page history: %w", err)
	}

	return created, nil
}

// ListPageHistory lists page history for a page with optional pagination
func (r *historyRepo) ListPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, error) {
	runner := r.tm.GetRunner(ctx)

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
		if *limit < 0 || *offset < 0 {
			return nil, ErrInvalidPagination
		}
		query += " LIMIT $2 OFFSET $3"
		args = append(args, *limit, *offset)
	}

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list page history: %w", err)
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
			return nil, fmt.Errorf("failed to scan page history: %w", err)
		}
		history = append(history, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(history) == 0 {
		return []models.PageHistory{}, nil
	}

	return history, nil
}

// RemovePageHistory removes page history for a list of pages
func (r *historyRepo) RemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) []error {
	if len(pageIDs) == 0 {
		return []error{ErrInvalidPageIDs}
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
		return []error{fmt.Errorf("failed to remove page history: %w", err)}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return []error{fmt.Errorf("failed to get affected rows: %w", err)}
	}

	if rowsAffected == 0 {
		return []error{ErrNoPageHistory}
	}

	return nil
}

// PageHistoryExists checks if a page history exists with specific criteria
func (r *historyRepo) PageHistoryExists(ctx context.Context, pageID string, weekNumber1, weekNumber2, yearNumber1, yearNumber2 int, bucketID1, bucketID2 string) (bool, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM page_history
			WHERE page_id = $1
			AND week_number_1 = $2
			AND week_number_2 = $3
			AND year_number_1 = $4
			AND year_number_2 = $5
			AND bucket_id_1 = $6
			AND bucket_id_2 = $7
		)
	`

	var exists bool
	err := runner.QueryRowContext(ctx, query,
		pageID, weekNumber1, weekNumber2, yearNumber1, yearNumber2, bucketID1, bucketID2,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check page history existence: %w", err)
	}

	return exists, nil
}
