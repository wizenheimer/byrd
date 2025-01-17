package screenshot

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
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
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

// StoreScreenshot stores a screenshot in S3 storage
func (s *s3ScreenshotRepo) StoreScreenshotImage(ctx context.Context, data models.ScreenshotImageResponse, path string) error {
	s.logger.Debug("storing screenshot",
		zap.String("path", path))

	buf := new(bytes.Buffer)
	if err := encodeImage(data.Image, buf); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	// Convert to bytes.Reader for seekable reading
	reader := bytes.NewReader(buf.Bytes())

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String("image/png"),
		Metadata:    data.Metadata.ToMap(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

// StoreContent stores text content in S3 storage
func (s *s3ScreenshotRepo) StoreScreenshotHTMLContent(ctx context.Context, data models.ScreenshotHTMLContentResponse, path string) error {
	s.logger.Debug("storing content", zap.String("path", path))

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        strings.NewReader(data.HTMLContent),
		ContentType: aws.String("text/plain"),
		Metadata:    data.Metadata.ToMap(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload content: %w", err)
	}

	return nil
}

// GetContent retrieves text content from S3 storage
func (s *s3ScreenshotRepo) GetScreenshotHTMLContent(ctx context.Context, path string) (models.ScreenshotHTMLContentResponse, []error) {
	errs := []error{}

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return models.ScreenshotHTMLContentResponse{}, append(errs, err)
	}

	screenshotMetadata, errs := models.ScreenshotMetadataFromMap(metadata)
	if errs != nil {
		return models.ScreenshotHTMLContentResponse{}, errs
	}

	resp := models.ScreenshotHTMLContentResponse{
		Status:      "success",
		HTMLContent: string(data),
		Metadata:    &screenshotMetadata,
		StoragePath: path,
	}

	return resp, nil

}

// GetScreenshot retrieves a screenshot from S3 storage
func (s *s3ScreenshotRepo) GetScreenshotImage(ctx context.Context, path string) (models.ScreenshotImageResponse, []error) {

	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return models.ScreenshotImageResponse{}, []error{err}
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return models.ScreenshotImageResponse{}, []error{fmt.Errorf("failed to decode image: %w", err)}
	}

	screenshotMetadata, errs := models.ScreenshotMetadataFromMap(metadata)
	if errs != nil {
		return models.ScreenshotImageResponse{}, errs
	}

	imgWidth, imgHeight, err := getImageDimensions(img)
	if err != nil {
		return models.ScreenshotImageResponse{}, []error{err}
	}

	resp := models.ScreenshotImageResponse{
		Status:      "success",
		Image:       img,
		Metadata:    &screenshotMetadata,
		ImageHeight: utils.ToPtr(imgHeight),
		ImageWidth:  utils.ToPtr(imgWidth),
		StoragePath: path,
	}

	return resp, nil
}

// Get retrieves binary data from S3 storage
func (s *s3ScreenshotRepo) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
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
func (s *s3ScreenshotRepo) Delete(ctx context.Context, path string) error {
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

// Listlists the latest content (images or text) for a given URL
func (s *s3ScreenshotRepo) List(ctx context.Context, prefix string, maxItems int) ([]models.ScreenshotListResponse, error) {
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
