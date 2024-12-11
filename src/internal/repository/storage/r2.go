package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type r2Storage struct {
	client *s3.Client
	bucket string
	logger *logger.Logger
}

func NewR2Storage(accessKey, secretKey, bucket, region string, logger *logger.Logger) (interfaces.StorageRepository, error) {
	logger.Debug("creating new r2 storage", zap.Any("bucket", bucket), zap.Any("region", region))
	return &r2Storage{
		client: nil,
		bucket: bucket,
		logger: logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

func (s *r2Storage) StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error {
	s.logger.Debug("storing screenshot into r2 storage", zap.Any("path", path), zap.Any("metadata", metadata))
	// Implementation
	return nil
}

func (s *r2Storage) StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error {
	s.logger.Debug("storing content into r2 storage", zap.Any("path", path), zap.Any("metadata", metadata))
	// Implementation
	return nil
}

func (s *r2Storage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	s.logger.Debug("getting file into r2 storage", zap.Any("path", path))
	// Implementation
	return nil, nil, nil
}

func (s *r2Storage) Delete(ctx context.Context, path string) error {
	s.logger.Debug("deleting file into r2 storage", zap.Any("path", path))
	// Implementation
	return nil
}
