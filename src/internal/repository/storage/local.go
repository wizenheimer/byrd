package storage

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/pkg/logger"
	"go.uber.org/zap"
)

type localStorage struct {
	directory string
	logger    *logger.Logger
}

func NewlocalStorage(directory string, logger *logger.Logger) (interfaces.StorageRepository, error) {
	logger.Debug("creating new local storage", zap.String("directory", directory))
	// Implementation
	return &localStorage{
		directory: directory,
		logger:    logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

func (s *localStorage) StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error {
	s.logger.Debug("storing screenshot into local storage", zap.String("path", path), zap.Any("metadata", metadata))
	// Implementation
	return nil
}

func (s *localStorage) StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error {
	s.logger.Debug("storing content into local storage", zap.String("path", path), zap.Any("metadata", metadata))
	// Implementation
	return nil
}

func (s *localStorage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	s.logger.Debug("getting file from local storage", zap.String("path", path))
	// Implementation
	return nil, nil, nil
}

func (s *localStorage) Delete(ctx context.Context, path string) error {
	s.logger.Debug("deleting file from local storage", zap.String("path", path))
	// Implementation
	return nil
}
