package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// r2Storage is a storage repository that uses R2 as the backend
type r2Storage struct {
	// client is the S3 client
	client *s3.Client
	// bucket is the S3 bucket name
	bucket string
	// logger is the logger
	logger *logger.Logger
}

// NewR2Storage creates a new R2 storage repository
func NewR2Storage(accessKey, secretKey, bucket, accountID string, logger *logger.Logger) (interfaces.StorageRepository, error) {
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

	return &r2Storage{
		client: client,
		bucket: bucket,
		logger: logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

// StoreScreenshot stores a screenshot in R2 storage
func (s *r2Storage) StoreScreenshot(ctx context.Context, data image.Image, path string, metadata map[string]string) error {
	s.logger.Debug("storing screenshot",
		zap.String("path", path))

	buf := new(bytes.Buffer)
	if err := encodeImage(data, buf); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        buf,
		ContentType: aws.String(metadata["Content-Type"]),
		Metadata:    metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

// StoreContent stores text content in R2 storage
func (s *r2Storage) StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error {
	s.logger.Debug("storing content", zap.String("path", path))

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        strings.NewReader(content),
		ContentType: aws.String(metadata["Content-Type"]),
		Metadata:    metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload content: %w", err)
	}

	return nil
}

// GetContent retrieves text content from R2 storage
func (s *r2Storage) GetContent(ctx context.Context, path string) (string, map[string]string, error) {
	s.logger.Debug("getting content", zap.String("path", path))

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return "", nil, err
	}
	return string(data), metadata, nil
}

// GetScreenshot retrieves a screenshot from R2 storage
func (s *r2Storage) GetScreenshot(ctx context.Context, path string) (image.Image, map[string]string, error) {
	s.logger.Debug("getting screenshot", zap.String("path", path))

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, metadata, nil
}

// Get retrieves binary data from R2 storage
func (s *r2Storage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	s.logger.Debug("getting binary", zap.String("path", path))

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

// Delete removes a file from R2 storage
func (s *r2Storage) Delete(ctx context.Context, path string) error {
	s.logger.Debug("deleting file", zap.String("path", path))

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// encodeImage encodes an image.Image to the specified writer
func encodeImage(img image.Image, w io.Writer) error {
	switch v := img.(type) {
	case *image.NRGBA, *image.RGBA:
		return png.Encode(w, v)
	case *image.YCbCr:
		return jpeg.Encode(w, v, &jpeg.Options{Quality: 90})
	default:
		// Default to PNG for unknown types
		return png.Encode(w, img)
	}
}
