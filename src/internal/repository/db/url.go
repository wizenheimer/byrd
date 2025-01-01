package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type urlRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewURLRepository(db *sql.DB, logger *logger.Logger) interfaces.URLRepository {
	return &urlRepository{
		db:     db,
		logger: logger.WithFields(map[string]interface{}{"module": "url_repository"}),
	}
}

// AddURL: adds a new URL if it does not exist
func (r *urlRepository) AddURL(ctx context.Context, url string) (*core_models.URL, error) {
	r.logger.Debug("adding new URL", zap.String("url", url))

	const query = `
        INSERT INTO urls (url)
        VALUES ($1)
        RETURNING id, url, created_at`

	var result core_models.URL
	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&result.ID,
		&result.URL,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUrls: lists all URLs in batches
func (r *urlRepository) ListURLs(ctx context.Context, batchSize int, lastSeenID *uuid.UUID) (*core_models.URLBatch, error) {
	const query = `
        SELECT id, url, created_at
        FROM urls
        WHERE $1::uuid IS NULL OR id < $1
        ORDER BY id DESC
        LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, lastSeenID, batchSize+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []core_models.URL
	for rows.Next() {
		var url core_models.URL
		err := rows.Scan(
			&url.ID,
			&url.URL,
			&url.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := &core_models.URLBatch{
		URLs:    urls,
		HasMore: false,
	}

	if len(urls) > batchSize {
		result.HasMore = true
		result.URLs = urls[:len(urls)-1]
		result.LastSeen = result.URLs[len(result.URLs)-1].ID
	}

	return result, nil
}

// DeleteURL: deletes a URL
func (r *urlRepository) DeleteURL(ctx context.Context, url string) error {
	r.logger.Debug("deleting URL", zap.String("url", url))

	const query = `
		DELETE FROM urls
		WHERE url = $1`

	commandTag, err := r.db.ExecContext(ctx, query, url)
	if err != nil {
		return err
	}

	rowAffected, err := commandTag.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *urlRepository) URLExists(ctx context.Context, url string) (*core_models.URL, bool, error) {
	r.logger.Debug("checking if URL exists", zap.String("url", url))

	const query = `
        SELECT id, url, created_at
        FROM urls
        WHERE url = $1`

	var result core_models.URL
	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&result.ID,
		&result.URL,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &result, true, nil
}
