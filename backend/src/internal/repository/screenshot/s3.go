package screenshot

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// s3ScreenshotRepo is a storage repository that uses S3 as the backend
type s3ScreenshotRepo struct {
	// client is the S3 client
	client *s3.Client
	// bucket is the S3 bucket name
	bucket string
	// logger is the logger
	logger *logger.Logger
}

// NewS3ScreenshotRepo creates a new S3 storage repository
func NewS3ScreenshotRepo(baseEndpoint, accessKey, secretKey, bucket, region string, logger *logger.Logger) (ScreenshotRepository, error) {
	if logger == nil {
		return nil, fmt.Errorf("can't initialize S3, logger is required")
	}

	logger.Debug("creating new S3 storage", zap.Any("bucket", bucket))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(baseEndpoint)
	})

	return &s3ScreenshotRepo{
		client: client,
		bucket: bucket,
		logger: logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

// StoreScreenshotImage stores a screenshot in local storage
func (s *s3ScreenshotRepo) StoreScreenshotImage(ctx context.Context, data *models.ScreenshotImage) error {
	return nil
}

// GetScreenshotImage retrieves a screenshot from local storage
func (s *s3ScreenshotRepo) StoreScreenshotContent(ctx context.Context, data *models.ScreenshotContent) error {
	return nil
}

// RetrieveScreenshotImage retrieves screenshot image from the storage
func (s *s3ScreenshotRepo) RetrieveScreenshotImage(ctx context.Context, path string) (*models.ScreenshotImage, error) {
	return nil, nil
}

// RetrieveScreenshotContent retrieves screenshot content from the storage
func (s *s3ScreenshotRepo) RetrieveScreenshotContent(ctx context.Context, path string) (*models.ScreenshotContent, error) {
	return nil, nil
}
