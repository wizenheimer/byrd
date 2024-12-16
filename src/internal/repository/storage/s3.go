package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// S3Storage is a storage repository that uses S3 as the backend
type s3Storage struct {
	// client is the S3 client
	client *s3.Client
	// bucket is the S3 bucket name
	bucket string
	// logger is the logger
	logger *logger.Logger
}

// NewS3Storage creates a new S3 storage repository
// It requires the access key, secret key, bucket name, account ID, and a logger
func NewS3Storage(accessKey, secretKey, bucket, accountID, session, region string, logger *logger.Logger) (interfaces.ScreenshotRepository, error) {
	if logger == nil {
		return nil, fmt.Errorf("can't initialize r2, logger is required")
	}

	logger.Debug("creating new r2 storage", zap.Any("bucket", bucket))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, session)),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &s3Storage{
		client: client,
		bucket: bucket,
		logger: logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

// StoreScreenshot stores a screenshot in S3 storage
func (s *s3Storage) StoreScreenshotImage(ctx context.Context, data image.Image, path string, metadata models.ScreenshotMetadata) error {
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
		ContentType: aws.String("image/png"),
		Metadata:    metadata.ToMap(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

// StoreContent stores text content in S3 storage
func (s *s3Storage) StoreScreenshotContent(ctx context.Context, content string, path string, metadata models.ScreenshotMetadata) error {
	s.logger.Debug("storing content", zap.String("path", path))

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        strings.NewReader(content),
		ContentType: aws.String("text/plain"),
		Metadata:    metadata.ToMap(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload content: %w", err)
	}

	return nil
}

// GetContent retrieves text content from S3 storage
func (s *s3Storage) GetScreenshotContent(ctx context.Context, path string) (string, models.ScreenshotMetadata, error) {
	s.logger.Debug("getting content", zap.String("path", path))

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return "", models.ScreenshotMetadata{}, err
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return "", models.ScreenshotMetadata{}, err
	}

	return string(data), screenshotMetadata, nil
}

// GetScreenshot retrieves a screenshot from S3 storage
func (s *s3Storage) GetScreenshotImage(ctx context.Context, path string) (image.Image, models.ScreenshotMetadata, error) {
	s.logger.Debug("getting screenshot", zap.String("path", path))

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return nil, models.ScreenshotMetadata{}, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, models.ScreenshotMetadata{}, fmt.Errorf("failed to decode image: %w", err)
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return nil, models.ScreenshotMetadata{}, err
	}

	return img, screenshotMetadata, nil
}

// Get retrieves binary data from S3 storage
func (s *s3Storage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
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

// Delete removes a file from S3 storage
func (s *s3Storage) Delete(ctx context.Context, path string) error {
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

func (s *s3Storage) List(ctx context.Context, prefix string, maxItems int) ([]models.ScreenshotListResponse, error) {
	// hash, err := path.GenerateURLHash(url)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate URL hash: %w", err)
	// }

	// var prefix string
	// switch contentType {
	// case "image", "screenshot":
	// 	prefix = "images/"
	// case "content", "text":
	// 	prefix = "text/"
	// default:
	// 	return nil, fmt.Errorf("invalid content type: %s", contentType)
	// }

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int32(int32(maxItems)),
	}

	var results []models.ScreenshotListResponse

	output, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	for _, obj := range output.Contents {
		results = append(results, models.ScreenshotListResponse{
			Key:          *obj.Key,
			LastModified: *obj.LastModified,
		})
		if len(results) >= maxItems {
			break
		}
	}

	// Sort by LastModified in descending order to get newest first
	sort.Slice(results, func(i, j int) bool {
		return results[i].LastModified.After(results[j].LastModified)
	})

	return results, nil
}
