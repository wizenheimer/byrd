package screenshot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// localScreenshotRepo is a storage repository that uses the local filesystem as the backend
type localScreenshotRepo struct {
	// directory is the directory where the files are stored
	directory string
	// logger is the logger
	logger *logger.Logger
}

// NewLocalScreenshotRepo creates a new local storage repository
// It requires the directory where the files will be stored and a logger
func NewLocalScreenshotRepo(directory string, logger *logger.Logger) (ScreenshotRepository, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	// Append "data" to the directory
	parentDirector := "data"
	directory = filepath.Join(parentDirector, directory)

	s := localScreenshotRepo{
		directory: directory,
		logger:    logger.WithFields(map[string]interface{}{"module": "storage"}),
	}

	return &s, s.ensurePath(directory)
}

// ensurePath ensures the directory structure exists for a given path
func (s *localScreenshotRepo) ensurePath(path string) error {
	dir := filepath.Dir(filepath.Join(s.directory, path))
	return os.MkdirAll(dir, 0755)
}

// StoreScreenshotImage stores a screenshot in local storage
func (s *localScreenshotRepo) StoreScreenshotImage(ctx context.Context, data *models.ScreenshotImage) error {
	if data == nil {
		return fmt.Errorf("screenshot image is required")
	}

	if data.StoragePath == "" {
		return fmt.Errorf("screenshot image storage path is required")
	}

	if err := s.ensurePath(data.StoragePath); err != nil {
		return err
	}

	fullPath := filepath.Join(s.directory, data.StoragePath)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := encodeImage(data.Image, file); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return s.saveMetadata(data.StoragePath, data.Metadata)
}

// GetScreenshotImage retrieves a screenshot from local storage
func (s *localScreenshotRepo) StoreScreenshotContent(ctx context.Context, data *models.ScreenshotContent) error {
	if data == nil {
		return fmt.Errorf("screenshot content is required")
	}

	if data.StoragePath == "" {
		return fmt.Errorf("screenshot content storage path is required")
	}

	if err := s.ensurePath(data.StoragePath); err != nil {
		return err
	}

	fullPath := filepath.Join(s.directory, data.StoragePath)
	if err := os.WriteFile(fullPath, []byte(data.Content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return s.saveMetadata(data.StoragePath, data.Metadata)
}

// RetrieveScreenshotImage retrieves screenshot image from the storage
func (s *localScreenshotRepo) RetrieveScreenshotImage(ctx context.Context, path string) (*models.ScreenshotImage, error) {
	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return nil, err
	}

	screenshotResp := models.ScreenshotImage{
		Image:       img,
		Metadata:    screenshotMetadata,
		StoragePath: path,
	}

	return &screenshotResp, nil
}

// RetrieveScreenshotContent retrieves screenshot content from the storage
func (s *localScreenshotRepo) RetrieveScreenshotContent(ctx context.Context, path string) (*models.ScreenshotContent, error) {
	data, metadata, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	screenshotMetadata, err := models.ScreenshotMetadataFromMap(metadata)
	if err != nil {
		return nil, err
	}

	screenshotResp := models.ScreenshotContent{
		StoragePath: path,
		Content:     string(data),
		Metadata:    screenshotMetadata,
	}

	return &screenshotResp, nil
}

// Get retrieves a binary from the local storage
func (s *localScreenshotRepo) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	fullPath := filepath.Join(s.directory, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	metadata, err := s.loadMetadata(path)
	if err != nil {
		return nil, nil, err
	}

	screenshotMetadata, err := metadata.ToMap()
	if err != nil {
		return nil, nil, err
	}

	return data, screenshotMetadata, nil
}

// loadMetadata loads metadata from the metadata file
func (s *localScreenshotRepo) loadMetadata(path string) (*models.ScreenshotMetadata, error) {
	metadataPath := s.getMetadataPath(path)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata models.ScreenshotMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

// getMetadataPath returns the path for the metadata file
func (s *localScreenshotRepo) getMetadataPath(path string) string {
	return filepath.Join(s.directory, path+".metadata.json")
}

// saveMetadata saves metadata to a separate file
func (s *localScreenshotRepo) saveMetadata(path string, metadata *models.ScreenshotMetadata) error {
	if metadata == nil {
		metadata = &models.ScreenshotMetadata{}
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := s.getMetadataPath(path)
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}
