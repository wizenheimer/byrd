package url

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/path"
	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
	"go.uber.org/zap"
)

// urlService is the service that provides URL management operations
type urlService struct {
	// urlRepository is an interface that provides URL management operations
	urlRepository repo.URLRepository

	// logger is a structured logger for logging
	logger *logger.Logger
}

// NewURLService creates a new URLService
// It returns an error if the logger or urlRepository is nil
func NewURLService(urlRepository repo.URLRepository, logger *logger.Logger) (svc.URLService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if urlRepository == nil {
		return nil, errors.New("urlRepository is required")
	}

	return &urlService{
		urlRepository: urlRepository,
		logger:        logger,
	}, nil
}

// AddURL adds a new URL if it does not exist
// If the URL already exists, it returns the existing URL without an error
// It returns an error if the URL is invalid or if there is an error adding the URL
func (s *urlService) AddURL(ctx context.Context, rawURL string) (*core_models.URL, error) {
	url, err := s.preprocessURL(rawURL)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("adding new URL", zap.Any("url", *url), zap.Any("rawURL", rawURL))

	if existingURL, exists, err := s.urlRepository.URLExists(ctx, *url); err == nil && exists {
		return existingURL, nil
	}

	return s.urlRepository.AddURL(ctx, *url)
}

// DeleteURL deletes a URL
// It returns an error if the URL is invalid or if there is an error deleting the URL
// It does not return an error if the URL does not exist
func (s *urlService) DeleteURL(ctx context.Context, rawURL string) error {
	url, err := s.preprocessURL(rawURL)
	if err != nil {
		return err
	}

	return s.urlRepository.DeleteURL(ctx, *url)
}

// URLExists checks if a URL exists
// It returns an error if the URL is invalid or if there is an error checking if the URL exists
func (s *urlService) URLExists(ctx context.Context, rawURL string) (bool, error) {
	url, err := s.preprocessURL(rawURL)
	if err != nil {
		return false, err
	}

	_, exists, err := s.urlRepository.URLExists(ctx, *url)
	return exists, err
}

// ListURLs lists all URLs in batches
// It returns a channel that emits URLBatch objects
// It returns a channel that emits errors
func (s *urlService) ListURLs(ctx context.Context, batchSize int, lastSeenID *uuid.UUID) (<-chan core_models.URLBatch, <-chan error) {
	result := make(chan core_models.URLBatch)
	errc := make(chan error, 1) // buffer to prevent goroutine leak

	go func() {
		defer close(result)
		defer close(errc)

		for {
			batch, err := s.urlRepository.ListURLs(ctx, batchSize, lastSeenID)
			if err != nil {
				errc <- err
				return
			}

			if len(batch.URLs) == 0 {
				return
			}

			result <- *batch
			lastSeenID = batch.URLs[len(batch.URLs)-1].ID
		}
	}()

	return result, errc
}

// preprocessURL preprocesses a URL
// It returns a pointer to the preprocessed URL
func (s *urlService) preprocessURL(rawURL string) (*string, error) {
	url, err := path.PreProcessURL(rawURL)
	if err != nil {
		return nil, err
	}

	return ptr.To(url), nil
}
