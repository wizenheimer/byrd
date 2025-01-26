// ./src/internal/repository/page/repository.go
package page

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type pageRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageRepository(tm *transaction.TxManager, logger *logger.Logger) PageRepository {
	return &pageRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "page_repository"}),
	}
}

func (r *pageRepo) getQuerier(ctx context.Context) interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
} {
	return r.tm.GetQuerier(ctx)
}

func (r *pageRepo) AddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, page models.PageProps) (*models.Page, error) {
	result := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
      INSERT INTO pages (competitor_id, url, title, capture_profile, diff_profile, status)
      VALUES ($1, $2, $3, $4, $5, $6)
      RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		competitorID, page.URL, page.Title, page.CaptureProfile, page.DiffProfile, models.PageStatusActive,
	).Scan(
		&result.ID,
		&result.CompetitorID,
		&result.URL,
		&result.Title,
		&result.CaptureProfile,
		&result.DiffProfile,
		&result.LastCheckedAt,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add page: %w", err)
	}

	return result, nil
}

func (r *pageRepo) BatchAddPageToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, error) {
	if len(pages) == 0 {
		return []models.Page{}, nil
	}

	// Create values for bulk insert
	valueStrings := make([]string, 0, len(pages))
	valueArgs := make([]interface{}, 0, len(pages)*6)
	for i, page := range pages {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs,
			competitorID,
			page.URL,
			page.Title,
			page.CaptureProfile,
			page.DiffProfile,
			models.PageStatusActive,
		)
	}

	query := fmt.Sprintf(`
        INSERT INTO pages (competitor_id, url, title, capture_profile, diff_profile, status)
        VALUES %s
        RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		strings.Join(valueStrings, ","))

	rows, err := r.getQuerier(ctx).Query(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to batch add pages: %w", err)
	}
	defer rows.Close()

	var results []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(
			&page.ID,
			&page.CompetitorID,
			&page.URL,
			&page.Title,
			&page.CaptureProfile,
			&page.DiffProfile,
			&page.LastCheckedAt,
			&page.Status,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page: %w", err)
		}
		results = append(results, page)
	}

	return results, rows.Err()
}

func (r *pageRepo) GetCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) (*models.Page, error) {
	page := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
        SELECT id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
        FROM pages
        WHERE competitor_id = $1 AND id = $2 AND status != $3`,
		competitorID, pageID, models.PageStatusInactive,
	).Scan(
		&page.ID,
		&page.CompetitorID,
		&page.URL,
		&page.Title,
		&page.CaptureProfile,
		&page.DiffProfile,
		&page.LastCheckedAt,
		&page.Status,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	r.logger.Debug("page", zap.Any("page", page))
	return page, nil
}

func (r *pageRepo) BatchGetCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	if len(pageIDs) == 0 {
		return []models.Page{}, false, nil
	}

	query := `
		SELECT id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
    FROM pages
		WHERE competitor_id = $1 AND id = ANY($2) AND status != $3
		ORDER BY created_at DESC`

	args := []interface{}{competitorID, pageIDs, models.PageStatusInactive}

	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *limit)
	}

	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *offset)
	}

	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to batch get pages: %w", err)
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(
			&page.ID,
			&page.CompetitorID,
			&page.URL,
			&page.Title,
			&page.CaptureProfile,
			&page.DiffProfile,
			&page.LastCheckedAt,
			&page.Status,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan page: %w", err)
		}
		pages = append(pages, page)
	}

	hasMore := limit != nil && len(pages) == *limit

	return pages, hasMore, rows.Err()
}

func (r *pageRepo) GetCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, bool, error) {
	query := `
		SELECT id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
    FROM pages
		WHERE competitor_id = $1 AND status != $2
		ORDER BY created_at DESC`

	args := []interface{}{competitorID, models.PageStatusInactive}

	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *limit)
	}

	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *offset)
	}

	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get competitor pages: %w", err)
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(
			&page.ID,
			&page.CompetitorID,
			&page.URL,
			&page.Title,
			&page.CaptureProfile,
			&page.DiffProfile,
			&page.LastCheckedAt,
			&page.Status,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan page: %w", err)
		}
		pages = append(pages, page)
	}

	hasMore := limit != nil && len(pages) == *limit

	return pages, hasMore, rows.Err()
}

func (r *pageRepo) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (*models.Page, error) {
	result := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
      UPDATE pages
      SET url = $1, title = $2, capture_profile = $3, diff_profile = $4
      WHERE competitor_id = $5 AND id = $6 AND status != $7
      RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		page.URL, page.Title, page.CaptureProfile, page.DiffProfile, competitorID, pageID, models.PageStatusInactive,
	).Scan(
		&result.ID,
		&result.CompetitorID,
		&result.URL,
		&result.Title,
		&result.CaptureProfile,
		&result.DiffProfile,
		&result.LastCheckedAt,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return result, nil
}

func (r *pageRepo) UpdateCompetitorPageURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error) {
	result := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
      UPDATE pages
      SET url = $1
      WHERE competitor_id = $2 AND id = $3 AND status != $4
      RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		url, competitorID, pageID, models.PageStatusInactive,
	).Scan(
		&result.ID,
		&result.CompetitorID,
		&result.URL,
		&result.Title,
		&result.CaptureProfile,
		&result.DiffProfile,
		&result.LastCheckedAt,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to update page URL: %w", err)
	}

	return result, nil
}

func (r *pageRepo) UpdateCompetitorCaptureProfile(ctx context.Context, competitorID, pageID uuid.UUID, captureProfile *models.CaptureProfile, url string) (*models.Page, error) {
	result := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
      UPDATE pages
      SET capture_profile = $1, url = $2
      WHERE competitor_id = $3 AND id = $4 AND status != $5
      RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		captureProfile, url, competitorID, pageID, models.PageStatusInactive,
	).Scan(
		&result.ID,
		&result.CompetitorID,
		&result.URL,
		&result.Title,
		&result.CaptureProfile,
		&result.DiffProfile,
		&result.LastCheckedAt,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("page not found")
		}
		return nil, fmt.Errorf("failed to update capture profile: %w", err)
	}

	return result, nil
}

func (r *pageRepo) UpdateCompetitorDiffProfile(ctx context.Context, competitorID, pageID uuid.UUID, diffProfile []string) (*models.Page, error) {
	result := &models.Page{}

	err := r.getQuerier(ctx).QueryRow(ctx, `
      UPDATE pages
      SET diff_profile = $1
      WHERE competitor_id = $2 AND id = $3 AND status != $4
      RETURNING id, competitor_id, url, title, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at`,
		diffProfile, competitorID, pageID, models.PageStatusInactive,
	).Scan(
		&result.ID,
		&result.CompetitorID,
		&result.URL,
		&result.Title,
		&result.CaptureProfile,
		&result.DiffProfile,
		&result.LastCheckedAt,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to update diff profile: %w", err)
	}

	return result, nil
}

func (r *pageRepo) UpdateCompetitorURL(ctx context.Context, competitorID, pageID uuid.UUID, url string) (*models.Page, error) {
	if competitorID == uuid.Nil || pageID == uuid.Nil || url == "" {
		return nil, errors.New("invalid competitor ID, page ID, or URL")
	}

	page := &models.Page{}
	err := r.getQuerier(ctx).QueryRow(ctx, `
      UPDATE pages
      SET
          url = $1,
          capture_profile = jsonb_set(
              COALESCE(capture_profile::jsonb, '{}'::jsonb),
              '{url}',
              $1::text::jsonb
          ),
          updated_at = CURRENT_TIMESTAMP
      WHERE
          competitor_id = $2
          AND id = $3
          AND status != $4
      RETURNING
          id,
          competitor_id,
          url,
          title,
          capture_profile,
          diff_profile,
          last_checked_at,
          status,
          created_at,
          updated_at`,
		url,
		competitorID,
		pageID,
		models.PageStatusInactive,
	).Scan(
		&page.ID,
		&page.CompetitorID,
		&page.URL,
		&page.Title,
		&page.CaptureProfile,
		&page.DiffProfile,
		&page.LastCheckedAt,
		&page.Status,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to update competitor URL: %w", err)
	}

	return page, nil
}

func (r *pageRepo) DeleteCompetitorPageByID(ctx context.Context, competitorID, pageID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE pages
		SET status = $1
		WHERE competitor_id = $2 AND id = $3 AND status != $1`,
		models.PageStatusInactive, competitorID, pageID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("page not found")
	}

	return nil
}

func (r *pageRepo) BatchDeleteCompetitorPagesByIDs(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) error {
	if len(pageIDs) == 0 {
		return nil
	}

	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE pages
		SET status = $1
		WHERE competitor_id = $2 AND id = ANY($3) AND status != $1`,
		models.PageStatusInactive, competitorID, pageIDs,
	)

	if err != nil {
		return fmt.Errorf("failed to batch delete pages: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("pages not found")
	}

	return nil
}

func (r *pageRepo) DeleteAllCompetitorPages(ctx context.Context, competitorID uuid.UUID) error {
	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE pages
		SET status = $1
		WHERE competitor_id = $2 AND status != $1`,
		models.PageStatusInactive, competitorID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete all competitor pages: %w", err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn("pages not found", zap.Any("competitorID", competitorID))
	}

	return nil
}

func (r *pageRepo) BatchDeleteAllCompetitorPages(ctx context.Context, competitorIDs []uuid.UUID) error {
	if len(competitorIDs) == 0 {
		return nil
	}

	result, err := r.getQuerier(ctx).Exec(ctx, `
		UPDATE pages
		SET status = $1
		WHERE competitor_id = ANY($2) AND status != $1`,
		models.PageStatusInactive, competitorIDs,
	)

	if err != nil {
		return fmt.Errorf("failed to batch delete all competitor pages: %w", err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn("pages not found", zap.Any("competitorIDs", competitorIDs))
	}

	return nil
}

func (r *pageRepo) GetActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error) {
	if batchSize <= 0 {
		return models.ActivePageBatch{}, errors.New("invalid batch size")
	}

	// Build base query
	query := `
		SELECT id
		FROM pages
		WHERE status = $1`
	args := []interface{}{models.PageStatusActive}

	// Add cursor-based pagination using lastPageID
	if lastPageID != nil {
		query += ` AND id > $2` // Using > to get next batch after the last seen ID
		args = append(args, *lastPageID)
	}

	// Ensure deterministic ordering
	query += ` ORDER BY id ASC`

	// Add limit
	query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
	args = append(args, batchSize+1) // Request one extra to determine if there are more pages

	// Execute query
	rows, err := r.getQuerier(ctx).Query(ctx, query, args...)
	if err != nil {
		return models.ActivePageBatch{}, fmt.Errorf("failed to query active pages: %w", err)
	}
	defer rows.Close()

	// Collect results
	var pageIDs []uuid.UUID
	for rows.Next() {
		var pageID uuid.UUID
		if err := rows.Scan(&pageID); err != nil {
			return models.ActivePageBatch{}, fmt.Errorf("failed to scan page ID: %w", err)
		}
		pageIDs = append(pageIDs, pageID)
	}

	if err = rows.Err(); err != nil {
		return models.ActivePageBatch{}, fmt.Errorf("error iterating pages: %w", err)
	}

	// Determine if there are more pages
	hasMore := len(pageIDs) > batchSize
	if hasMore {
		pageIDs = pageIDs[:batchSize] // Remove the extra item we requested
	}

	// Set the last seen ID
	var lastSeen *uuid.UUID
	if len(pageIDs) > 0 {
		lastSeen = &pageIDs[len(pageIDs)-1]
	}

	return models.ActivePageBatch{
		PageIDs:  pageIDs,
		HasMore:  hasMore,
		LastSeen: lastSeen,
	}, nil
}

func (r *pageRepo) GetPageByPageID(ctx context.Context, pageID uuid.UUID) (*models.Page, error) {
	if pageID == uuid.Nil {
		return nil, errors.New("invalid page ID")
	}

	page := &models.Page{}
	err := r.getQuerier(ctx).QueryRow(ctx, `
    SELECT
        id,
        competitor_id,
        url,
        title,
        capture_profile,
        diff_profile,
        last_checked_at,
        status,
        created_at,
        updated_at
    FROM pages
    WHERE id = $1 AND status != $2`,
		pageID,
		models.PageStatusInactive,
	).Scan(
		&page.ID,
		&page.CompetitorID,
		&page.URL,
		&page.Title,
		&page.CaptureProfile,
		&page.DiffProfile,
		&page.LastCheckedAt,
		&page.Status,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to get page by ID: %w", err)
	}

	return page, nil
}
