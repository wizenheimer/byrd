package history

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type historyRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageHistoryRepository(tm *transaction.TxManager, logger *logger.Logger) PageHistoryRepository {
	return &historyRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "history_repository"}),
	}
}

func (r *historyRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func (r *historyRepo) CreateHistoryForPage(ctx context.Context, pageID uuid.UUID, diffContent any, prev, curr string) error {
	// Validate diffContent as well
	if pageID == uuid.Nil || diffContent == nil {
		return fmt.Errorf("page ID and diff content are required")
	}

	// Convert diffContent to JSONB
	diffContentJSON, err := json.Marshal(diffContent)
	if err != nil {
		return fmt.Errorf("failed to marshal diff content: %w", err)
	}

	query := `
        INSERT INTO page_history (
            page_id,
            diff_content,
            status,
            prev,
            curr
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	var id uuid.UUID
	err = r.getQuerier(ctx).QueryRow(ctx, query,
		pageID,
		diffContentJSON,
		models.HistoryStatusActive,
		prev,
		curr,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create page history: %w", err)
	}

	return nil
}

func (r *historyRepo) BatchGetPageHistory(ctx context.Context, pageID uuid.UUID, limit, offset *int) ([]models.PageHistory, bool, error) {
	if pageID == uuid.Nil {
		return nil, false, fmt.Errorf("page ID is required")
	}

	// Build query with pagination
	query := `
        SELECT
            id,
            page_id,
            diff_content,
            created_at,
            status,
            prev,
            curr
        FROM page_history
        WHERE page_id = $1
        AND status = $2
        ORDER BY created_at DESC`

	args := []interface{}{pageID, models.HistoryStatusActive}

	// Add limit if provided
	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *limit+1) // Fetch one extra to determine if there are more results
	}

	// Add offset if provided
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *offset)
	}

	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to query page history: %w", err)
	}
	defer rows.Close()

	var histories []models.PageHistory
	for rows.Next() {
		var history models.PageHistory
		var diffContentJSON []byte

		err := rows.Scan(
			&history.ID,
			&history.PageID,
			&diffContentJSON,
			&history.CreatedAt,
			&history.Status,
			&history.Prev,
			&history.Curr,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan page history: %w", err)
		}

		// Unmarshal JSON content
		err = json.Unmarshal(diffContentJSON, &history.DiffContent)
		if err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal diff content: %w", err)
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, false, fmt.Errorf("error iterating page history: %w", err)
	}

	hasMore := false
	if limit != nil && len(histories) > *limit {
		hasMore = true
		histories = histories[:*limit] // Remove the extra result
	}

	return histories, hasMore, nil
}

func (r *historyRepo) BatchRemovePageHistory(ctx context.Context, pageIDs []uuid.UUID) error {
	if len(pageIDs) == 0 {
		return nil
	}

	query := `
		UPDATE page_history
		SET status = $1
		WHERE page_id = ANY($2)
		AND status = $3`

	result, err := r.getQuerier(ctx).Exec(ctx, query,
		models.HistoryStatusInactive,
		pageIDs,
		models.HistoryStatusActive,
	)
	if err != nil {
		return fmt.Errorf("failed to batch remove page history: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no page history found")
	}

	return nil
}
