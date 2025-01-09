package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/transaction"
	"github.com/wizenheimer/byrd/src/pkg/err"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type pageRepo struct {
	tm     *transaction.TxManager
	logger *logger.Logger
}

func NewPageRepository(tm *transaction.TxManager, logger *logger.Logger) repo.PageRepository {
	return &pageRepo{
		tm:     tm,
		logger: logger.WithFields(map[string]interface{}{"module": "page_repository"}),
	}
}

// AddPagesToCompetitor adds pages to a competitor
func (r *pageRepo) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, err.Error) {
	pageErr := err.New()
	if len(pages) == 0 {
		pageErr.Add(repo.ErrPagesUnspecified, map[string]any{
			"competitor_id": competitorID,
		})
		return nil, pageErr
	}

	runner := r.tm.GetRunner(ctx)

	// Build batch insert query
	valueStrings := make([]string, 0, len(pages))
	valueArgs := make([]interface{}, 0, len(pages)*5)

	for i, page := range pages {
		// Set defaults for CaptureProfile if needed
		if err := defaults.Set(&page.CaptureProfile); err != nil {
			pageErr.Add(err, map[string]any{
				"pageURL": page.URL,
				"error":   "failed to set capture profile defaults",
			})
			return nil, pageErr
		}

		// Initialize empty diff profile if nil
		if page.DiffProfile == nil {
			page.DiffProfile = make(map[string]interface{})
		}

		// Add the placeholder for the prepared statement
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))

		// Marshal the capture profile
		captureProfileJSON, err := json.Marshal(page.CaptureProfile)
		if err != nil {
			pageErr.Add(repo.ErrFailedToMarshallCaptureProfile, map[string]any{
				"pageURL": page.URL,
				"error":   err.Error(),
			})
			return nil, pageErr
		}

		// Marshal the diff profile
		diffProfileJSON, err := json.Marshal(page.DiffProfile)
		if err != nil {
			pageErr.Add(repo.ErrFailedToMarshallDiffProfile, map[string]any{
				"pageURL": page.URL,
				"error":   err.Error(),
			})
			return nil, pageErr
		}

		valueArgs = append(valueArgs,
			competitorID,
			page.URL,
			captureProfileJSON,
			diffProfileJSON,
			models.PageStatusActive) // Default status for new pages
	}

	query := fmt.Sprintf(`
		INSERT INTO pages (competitor_id, url, capture_profile, diff_profile, status)
		VALUES %s
		RETURNING id, competitor_id, url, capture_profile::jsonb, diff_profile::jsonb, last_checked_at, status, created_at, updated_at
	`, strings.Join(valueStrings, ","))

	rows, err := runner.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitor_id": competitorID,
			"error":         err.Error(),
		})
		return nil, pageErr
	}
	defer rows.Close()

	createdPages, err := scanPages(rows)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitor_id": competitorID,
			"error":         err.Error(),
		})
		return nil, pageErr
	}

	return createdPages, nil
}

// RemovePagesFromCompetitor removes pages from a competitor
func (r *pageRepo) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) err.Error {
	pageErr := err.New()
	runner := r.tm.GetRunner(ctx)

	var query string
	var args []interface{}

	if pageIDs == nil {
		// Soft delete all active pages for the competitor
		query = `
			UPDATE pages
			SET status = $2, updated_at = CURRENT_TIMESTAMP
			WHERE competitor_id = $1 AND status = $3
		`
		args = []interface{}{competitorID, models.PageStatusInactive, models.PageStatusActive}
	} else {
		// Soft delete specific active pages
		placeholders := make([]string, len(pageIDs))
		args = make([]interface{}, len(pageIDs)+3) // +3 for competitorID, new status, and current status
		args[0] = competitorID
		args[1] = models.PageStatusInactive
		args[2] = models.PageStatusActive

		for i, id := range pageIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+4) // Start from $4 since we have 3 fixed params
			args[i+3] = id
		}

		query = fmt.Sprintf(`
			UPDATE pages
			SET status = $2, updated_at = CURRENT_TIMESTAMP
			WHERE competitor_id = $1
			AND status = $3
			AND id IN (%s)
		`, strings.Join(placeholders, ","))
	}

	result, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"query": query,
		})
		return pageErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		pageErr.Add(err, map[string]any{
			"query": query,
		})
	}

	if rowsAffected == 0 {
		pageErr.Add(repo.ErrNoCompetitorPages, map[string]any{
			"pageIDs": pageIDs,
		})
		return pageErr
	}

	return nil
}

// ListCompetitorPages gets the pages for a competitor with optional pagination
func (r *pageRepo) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, err.Error) {
	pageErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
		FROM pages
		WHERE competitor_id = $1 AND status = 'active'
		ORDER BY created_at DESC
	`
	args := []interface{}{competitorID}
	paramCount := 1

	// Only add LIMIT if limit is provided and valid
	if limit != nil {
		if *limit < 0 {
			pageErr.Add(repo.ErrInvalidLimit, nil)
			return nil, pageErr
		}
		paramCount++
		query += fmt.Sprintf(" LIMIT $%d", paramCount)
		args = append(args, *limit)
	}

	// Only add OFFSET if offset is provided and valid
	if offset != nil {
		if *offset < 0 {
			pageErr.Add(repo.ErrInvalidOffset, nil)
			return nil, pageErr
		}
		paramCount++
		query += fmt.Sprintf(" OFFSET $%d", paramCount)
		args = append(args, *offset)
	}

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitorID": competitorID,
		})
		return nil, pageErr
	}
	defer rows.Close()

	pages, err := scanPages(rows)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitorID": competitorID,
		})
		return nil, pageErr
	}

	if len(pages) == 0 {
		pageErr.Add(repo.ErrNoCompetitorPages, map[string]any{
			"competitorID": competitorID,
		})
		return nil, pageErr
	}

	return pages, nil
}

// ListActivePages lists all active pages in batches using cursor-based pagination
func (r *pageRepo) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, err.Error) {
	pageErr := err.New()
	if batchSize <= 0 {
		pageErr.Add(repo.ErrInvalidBatchSize, nil)
		return models.ActivePageBatch{}, pageErr
	}

	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
		FROM pages
		WHERE status = 'active'
	`
	args := make([]interface{}, 0)

	if lastPageID != nil {
		query += " AND created_at > (SELECT created_at FROM pages WHERE id = $1)"
		args = append(args, lastPageID)
	}

	query += fmt.Sprintf(`
		ORDER BY created_at ASC
		LIMIT %d
	`, batchSize+1) // Get one extra to determine if there are more pages

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"batchSize":  batchSize,
			"lastPageID": lastPageID,
		})
		return models.ActivePageBatch{}, pageErr
	}
	defer rows.Close()

	var pages []models.Page
	pages, err = scanPages(rows)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"lastPageID": lastPageID,
			"batchSize":  batchSize,
		})
		return models.ActivePageBatch{}, pageErr
	}

	result := models.ActivePageBatch{
		HasMore: len(pages) > batchSize,
		Pages:   pages,
	}

	if result.HasMore {
		result.Pages = pages[:batchSize]
		result.LastSeen = &result.Pages[len(result.Pages)-1].ID
	}

	return result, nil
}

// GetCompetitorPage gets a page for a competitor
func (r *pageRepo) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (models.Page, err.Error) {
	pageErr := err.New()
	runner := r.tm.GetRunner(ctx)

	query := `
    SELECT id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
    FROM pages
    WHERE id = $1 AND competitor_id = $2 AND status = $3
`
	row := runner.QueryRowContext(ctx, query, pageID, competitorID, models.PageStatusActive)
	page, err := scanPage(row)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitorID": competitorID,
			"pageID":       pageID,
		})
		return models.Page{}, pageErr
	}

	if page == nil {
		pageErr.Add(repo.ErrPageNotFound, nil)
		return models.Page{}, pageErr
	}

	return *page, nil
}

// UpdateCompetitorPage updates a page for a competitor
func (r *pageRepo) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (models.Page, err.Error) {
	pageErr := err.New()
	runner := r.tm.GetRunner(ctx)

	// Verify the page belongs to the competitor
	var exists bool
	err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pages WHERE id = $1 AND competitor_id = $2)", pageID, competitorID).Scan(&exists)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitorID": competitorID,
			"pageID":       pageID,
		})
		return models.Page{}, pageErr
	}
	if !exists {
		pageErr.Add(repo.ErrPageNotFound, map[string]any{
			"competitorID": competitorID,
			"pageID":       pageID,
		})
		return models.Page{}, pageErr
	}

	// Set defaults for CaptureProfile if needed
	if err := defaults.Set(&page.CaptureProfile); err != nil {
		pageErr.Add(err, map[string]any{
			"pageID": pageID,
			"error":  "failed to set capture profile defaults",
		})
		return models.Page{}, pageErr
	}

	// Marshal the capture profile
	captureProfileJSON, err := json.Marshal(page.CaptureProfile)
	if err != nil {
		pageErr.Add(repo.ErrFailedToMarshallCaptureProfile, map[string]any{
			"pageID": pageID,
			"error":  err.Error(),
		})
		return models.Page{}, pageErr
	}

	// Initialize empty diff profile if nil
	if page.DiffProfile == nil {
		page.DiffProfile = make(map[string]interface{})
	}

	// Marshal the diff profile
	diffProfileJSON, err := json.Marshal(page.DiffProfile)
	if err != nil {
		pageErr.Add(repo.ErrFailedToMarshallDiffProfile, map[string]any{
			"pageID": pageID,
			"error":  err.Error(),
		})
		return models.Page{}, pageErr
	}

	query := `
		UPDATE pages
		SET url = $1, capture_profile = $2, diff_profile = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND competitor_id = $5
		RETURNING id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
	`

	row := runner.QueryRowContext(ctx, query,
		page.URL,
		captureProfileJSON,
		diffProfileJSON,
		pageID,
		competitorID,
	)

	updatedPage, err := scanPage(row)
	if err != nil {
		pageErr.Add(err, map[string]any{
			"competitorID": competitorID,
			"pageID":       pageID,
		})
		return models.Page{}, pageErr
	}

	return *updatedPage, nil
}

// scanPage scans a single row into a Page object
func scanPage(row dbScanner) (*models.Page, error) {
	var page models.Page
	var captureProfileBytes, diffProfileBytes []byte

	err := row.Scan(
		&page.ID,
		&page.CompetitorID,
		&page.URL,
		&captureProfileBytes,
		&diffProfileBytes,
		&page.LastCheckedAt,
		&page.Status,
		&page.CreatedAt,
		&page.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repo.ErrPageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repo.ErrFailedToScanPages, err)
	}

	// Initialize empty CaptureProfile with defaults
	if err := defaults.Set(&page.CaptureProfile); err != nil {
		return nil, fmt.Errorf("failed to set capture profile defaults: %w", err)
	}

	// Initialize empty DiffProfile if needed
	if page.DiffProfile == nil {
		page.DiffProfile = make(map[string]interface{})
	}

	// Unmarshal CaptureProfile with proper error handling
	if len(captureProfileBytes) > 0 {
		if err := json.Unmarshal(captureProfileBytes, &page.CaptureProfile); err != nil {
			return nil, fmt.Errorf("%w: %v", repo.ErrFailedToUnmarshalCaptureProfile, err)
		}
		// Set defaults after unmarshaling to ensure all fields are properly initialized
		if err := defaults.Set(&page.CaptureProfile); err != nil {
			return nil, fmt.Errorf("failed to set capture profile defaults after unmarshal: %w", err)
		}
	}

	// Unmarshal DiffProfile with proper error handling
	if len(diffProfileBytes) > 0 {
		if err := json.Unmarshal(diffProfileBytes, &page.DiffProfile); err != nil {
			return nil, fmt.Errorf("%w: %v", repo.ErrFailedToUnmarshalDiffProfile, err)
		}
	} else {
		// Ensure DiffProfile is initialized even when no bytes are present
		page.DiffProfile = make(map[string]interface{})
	}

	return &page, nil
}

// scanPages scans multiple rows into a slice of Page objects
func scanPages(rows *sql.Rows) ([]models.Page, error) {
	var pages []models.Page
	for rows.Next() {
		page, err := scanPage(rows)
		if err != nil {
			return nil, err
		}
		pages = append(pages, *page)
	}
	if err := rows.Err(); err != nil {
		return nil, repo.ErrFailedToIterateOverPagesForScan
	}

	return pages, nil
}

// dbScanner is an interface that abstracts the Scan method
// This allows us to use the same scanning logic for both sql.Row and sql.Rows
type dbScanner interface {
	Scan(dest ...interface{}) error
}
