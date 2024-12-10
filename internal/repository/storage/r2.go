package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
)

type r2Storage struct {
	client *s3.Client
	bucket string
}

func NewR2Storage(accessKey, secretKey, bucket, region string) (interfaces.StorageRepository, error) {
	return &r2Storage{
		client: nil,
		bucket: bucket,
	}, nil
}

func (s *r2Storage) StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error {
	// Implementation
	return nil
}

func (s *r2Storage) StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error {
	// Implementation
	return nil
}

func (s *r2Storage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	// Implementation
	return nil, nil, nil
}

func (s *r2Storage) Delete(ctx context.Context, path string) error {
	// Implementation
	return nil
}
