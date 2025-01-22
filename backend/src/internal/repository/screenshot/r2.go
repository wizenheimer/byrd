// ./src/internal/repository/storage/r2.go
package screenshot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// r2ScreenshotRepo is a storage repository that uses R2 as the backend
type r2ScreenshotRepo struct {
	// client is the S3 client
	client *s3.Client
	// bucket is the S3 bucket name
	bucket string
	// logger is the logger
	logger *logger.Logger
}

// NewR2ScreenshotRepo creates a new R2 storage repository
func NewR2ScreenshotRepo(accessKey, secretKey, bucket, accountID string, logger *logger.Logger) (ScreenshotRepository, error) {
	if logger == nil {
		return nil, fmt.Errorf("can't initialize r2, logger is required")
	}

	logger.Debug("creating new r2 storage", zap.Any("bucket", bucket))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &r2ScreenshotRepo{
		client: client,
		bucket: bucket,
		logger: logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

// StoreScreenshotImage stores a screenshot in local storage
func (r *r2ScreenshotRepo) StoreScreenshotImage(ctx context.Context, data *models.ScreenshotImage) error {
	if data == nil {
		return fmt.Errorf("screenshot image is required")
	}

	if data.StoragePath == "" {
		return fmt.Errorf("screenshot image storage path is required")
	}

	buf := new(bytes.Buffer)
	if err := encodeImage(data.Image, buf); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	metadata, err := data.Metadata.ToMap()
	if err != nil {
		return err
	}

	// Convert to bytes.Reader for seekable reading
	reader := bytes.NewReader(buf.Bytes())

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(data.StoragePath),
		Body:        reader,
		ContentType: aws.String("image/png"),
		Metadata:    metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

// GetScreenshotImage retrieves a screenshot from local storage
func (r *r2ScreenshotRepo) StoreScreenshotContent(ctx context.Context, data *models.ScreenshotContent) error {
	if data == nil {
		return fmt.Errorf("screenshot content is required")
	}

	if data.StoragePath == "" {
		return fmt.Errorf("screenshot content storage path is required")
	}

	metadata, err := data.Metadata.ToMap()
	if err != nil {
		return err
	}

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(data.StoragePath),
		Body:        strings.NewReader(data.Content),
		ContentType: aws.String("text/plain"),
		Metadata:    metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload content: %w", err)
	}

	return nil
}

// RetrieveScreenshotImage retrieves screenshot image from the storage
func (r *r2ScreenshotRepo) RetrieveScreenshotImage(ctx context.Context, path string) (*models.ScreenshotImage, error) {
	data, metadata, err := r.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w, %v", err, data[:min(len(data), 8)])
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return nil, err
	}

	resp := models.ScreenshotImage{
		StoragePath: path,
		Image:       img,
		Metadata:    screenshotMetadata,
	}

	return &resp, nil
}

// RetrieveScreenshotContent retrieves screenshot content from the storage
func (r *r2ScreenshotRepo) RetrieveScreenshotContent(ctx context.Context, path string) (*models.ScreenshotContent, error) {
	if path == "" {
		return nil, fmt.Errorf("screenshot image path is required")
	}

	data, metadata, err := r.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return nil, err
	}

	resp := models.ScreenshotContent{
		StoragePath: path,
		Content:     string(data),
		Metadata:    screenshotMetadata,
	}

	return &resp, nil
}

func (s *r2ScreenshotRepo) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer output.Body.Close()

	content, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return content, output.Metadata, nil
}
