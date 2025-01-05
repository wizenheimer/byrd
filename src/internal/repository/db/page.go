package db

import (
	"context"
	"database/sql"
	"encoding/json"
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
	// ---- validation errors -----
	ErrPagesUnspecified                = errors.New("pages unspecified for competitor")
	ErrFailedToMarshallCaptureProfile  = errors.New("failed to marshal capture profile for page")
	ErrFailedToUnmarshalCaptureProfile = errors.New("failed to unmarshal capture profile for page")
	ErrFailedToMarshallDiffProfile     = errors.New("failed to marshal diff profile for page")
	ErrFailedToUnmarshalDiffProfile    = errors.New("failed to unmarshal diff profile for page")
	ErrInvalidBatchSize                = errors.New("invalid batch size")

	// ---- non fatal errors ----
	ErrFailedToConfirmPageRemoval = errors.New("failed to confirm page removal")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrNoCompetitorPages = errors.New("no pages found for the competitor")
	ErrPageNotFound      = errors.New("page not found")
    ErrFailedToScanPages              = errors.New("failed to scan pages from pages table")
    ErrFailedToIterateOverPagesForScan = errors.New("failed to scan pages from pages table")
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
func (r *pageRepo) AddPagesToCompetitor(ctx context.Context, competitorID uuid.UUID, pages []models.PageProps) ([]models.Page, []error) {
	if len(pages) == 0 {
		return nil, []error{ErrPagesUnspecified}
	}

	runner := r.tm.GetRunner(ctx)

	// Build batch insert query
	valueStrings := make([]string, 0, len(pages))
	valueArgs := make([]interface{}, 0, len(pages)*4)
	for i, page := range pages {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))

		captureProfileJSON, err := json.Marshal(page.CaptureProfile)
		if err != nil {
			return nil, []error{err, ErrFailedToMarshallCaptureProfile}
		}

		diffProfileJSON, err := json.Marshal(page.DiffProfile)
		if err != nil {
			return nil, []error{err, ErrFailedToMarshallDiffProfile}
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
		return nil, []error{err}
	}
	defer rows.Close()

	createdPages, err := scanPages(rows)
	if err != nil {
		return nil, []error{err}
	}

	return createdPages, nil
}

// RemovePagesFromCompetitor removes pages from a competitor
func (r *pageRepo) RemovePagesFromCompetitor(ctx context.Context, competitorID uuid.UUID, pageIDs []uuid.UUID) []error {
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
		return []error{
			err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return []error{
			ErrFailedToConfirmPageRemoval,
		}
	}

	if rowsAffected == 0 {
		return []error{
			ErrNoCompetitorPages,
		}
	}

	return nil
}

// ListCompetitorPages gets the pages for a competitor with optional pagination
func (r *pageRepo) ListCompetitorPages(ctx context.Context, competitorID uuid.UUID, limit, offset *int) ([]models.Page, []error) {
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
			return nil, []error{ErrInvalidLimit}
		}
		paramCount++
		query += fmt.Sprintf(" LIMIT $%d", paramCount)
		args = append(args, *limit)
	}

	// Only add OFFSET if offset is provided and valid
	if offset != nil {
		if *offset < 0 {
			return nil, []error{ErrInvalidOffset}
		}
		paramCount++
		query += fmt.Sprintf(" OFFSET $%d", paramCount)
		args = append(args, *offset)
	}

	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	pages, err := scanPages(rows)
	if err != nil {
		return nil, []error{err}
	}

	if len(pages) == 0 {
		return nil, []error{ErrNoCompetitorPages}
	}

	return pages, nil
}

// ListActivePages lists all active pages in batches using cursor-based pagination
func (r *pageRepo) ListActivePages(ctx context.Context, batchSize int, lastPageID *uuid.UUID) (models.ActivePageBatch, error) {
	if batchSize <= 0 {
		return models.ActivePageBatch{}, ErrInvalidBatchSize
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
		return models.ActivePageBatch{}, err
	}
	defer rows.Close()

	var pages []models.Page
	pages, err = scanPages(rows)
	if err != nil {
		return models.ActivePageBatch{}, err
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
func (r *pageRepo) GetCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID) (models.Page, error) {
	runner := r.tm.GetRunner(ctx)

	query := `
		SELECT id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
		FROM pages
		WHERE id = $1 AND competitor_id = $2
	`

	row := runner.QueryRowContext(ctx, query, pageID, competitorID)
	page, err := scanPage(row)
	if err != nil {
		return models.Page{}, err
	}

	if page == nil {
		return models.Page{}, ErrPageNotFound
	}

	return *page, nil
}

// UpdateCompetitorPage updates a page for a competitor
// UpdateCompetitorPage updates a page for a competitor
func (r *pageRepo) UpdateCompetitorPage(ctx context.Context, competitorID, pageID uuid.UUID, page models.PageProps) (models.Page, error) {
	runner := r.tm.GetRunner(ctx)

	// Verify the page belongs to the competitor
	var exists bool
	err := runner.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pages WHERE id = $1 AND competitor_id = $2)", pageID, competitorID).Scan(&exists)
	if err != nil {
		return models.Page{}, err
	}
	if !exists {
		return models.Page{}, ErrPageNotFound
	}

	// Marshal the maps to JSON
	captureProfileJSON, err := json.Marshal(page.CaptureProfile)
	if err != nil {
		return models.Page{}, ErrFailedToMarshallCaptureProfile
	}

	diffProfileJSON, err := json.Marshal(page.DiffProfile)
	if err != nil {
		return models.Page{}, ErrFailedToMarshallDiffProfile
	}

	query := `
		UPDATE pages
		SET url = $1, capture_profile = $2, diff_profile = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND competitor_id = $5
		RETURNING id, competitor_id, url, capture_profile, diff_profile, last_checked_at, status, created_at, updated_at
	`

	row := runner.QueryRowContext(ctx, query,
		page.URL,
		captureProfileJSON, // Using the marshaled JSON instead of the raw map
		diffProfileJSON,    // Using the marshaled JSON instead of the raw map
		pageID,
		competitorID,
	)

	updatedPage, err := scanPage(row)
	if err != nil {
		return models.Page{}, err
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
		return nil, ErrPageNotFound
	}
	if err != nil {
		return nil, ErrFailedToScanPages
	}

	// Initialize empty maps if needed
	if page.CaptureProfile == nil {
		page.CaptureProfile = make(map[string]interface{})
	}
	if page.DiffProfile == nil {
		page.DiffProfile = make(map[string]interface{})
	}

	// Unmarshal JSONB data into maps
	if len(captureProfileBytes) > 0 {
		if err := json.Unmarshal(captureProfileBytes, &page.CaptureProfile); err != nil {
			return nil, ErrFailedToUnmarshalCaptureProfile
		}
	}
	if len(diffProfileBytes) > 0 {
		if err := json.Unmarshal(diffProfileBytes, &page.DiffProfile); err != nil {
			return nil, ErrFailedToUnmarshalDiffProfile
		}
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
		return nil, ErrFailedToIterateOverPagesForScan
	}

	return pages, nil
}

// dbScanner is an interface that abstracts the Scan method
// This allows us to use the same scanning logic for both sql.Row and sql.Rows
type dbScanner interface {
	Scan(dest ...interface{}) error
}
