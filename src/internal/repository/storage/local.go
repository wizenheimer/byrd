package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// localStorage is a storage repository that uses the local filesystem as the backend
type localStorage struct {
	// directory is the directory where the files are stored
	directory string
	// logger is the logger
	logger *logger.Logger
}

// NewLocalStorage creates a new local storage repository
// It requires the directory where the files will be stored and a logger
func NewLocalStorage(directory string, logger *logger.Logger) (interfaces.StorageRepository, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	// Append "data" to the directory
	parentDirector := "data"
	directory = filepath.Join(parentDirector, directory)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	logger.Debug("creating new local storage", zap.String("directory", directory))
	return &localStorage{
		directory: directory,
		logger:    logger.WithFields(map[string]interface{}{"module": "storage"}),
	}, nil
}

// ensurePath ensures the directory structure exists for a given path
func (s *localStorage) ensurePath(path string) error {
	dir := filepath.Dir(filepath.Join(s.directory, path))
	return os.MkdirAll(dir, 0755)
}

// getMetadataPath returns the path for the metadata file
func (s *localStorage) getMetadataPath(path string) string {
	return filepath.Join(s.directory, path+".metadata.json")
}

// saveMetadata saves metadata to a separate file
func (s *localStorage) saveMetadata(path string, metadata models.ScreenshotMetadata) error {
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

// loadMetadata loads metadata from the metadata file
func (s *localStorage) loadMetadata(path string) (models.ScreenshotMetadata, error) {
	metadataPath := s.getMetadataPath(path)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.ScreenshotMetadata{}, nil
		}
		return models.ScreenshotMetadata{}, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata models.ScreenshotMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return models.ScreenshotMetadata{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return metadata, nil
}

// StoreScreenshot stores a screenshot in the local storage
func (s *localStorage) StoreScreenshotImage(ctx context.Context, data image.Image, path string, metadata models.ScreenshotMetadata) error {
	s.logger.Debug("storing screenshot", zap.String("path", path))

	if err := s.ensurePath(path); err != nil {
		return err
	}

	fullPath := filepath.Join(s.directory, path)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := encodeImage(data, file); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return s.saveMetadata(path, metadata)
}

// StoreContent stores a text content in the local storage
func (s *localStorage) StoreScreenshotContent(ctx context.Context, content string, path string, metadata models.ScreenshotMetadata) error {
	s.logger.Debug("storing content", zap.String("path", path))

	if err := s.ensurePath(path); err != nil {
		return err
	}

	fullPath := filepath.Join(s.directory, path)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return s.saveMetadata(path, metadata)
}

// GetContent retrieves a text content from the local storage
func (s *localStorage) GetScreenshotContent(ctx context.Context, path string) (string, models.ScreenshotMetadata, error) {
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

// GetScreenshot retrieves a screenshot from the local storage
func (s *localStorage) GetScreenshotImage(ctx context.Context, path string) (image.Image, models.ScreenshotMetadata, error) {
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

// Get retrieves a binary from the local storage
func (s *localStorage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	s.logger.Debug("getting binary", zap.String("path", path))

	fullPath := filepath.Join(s.directory, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	metadata, err := s.loadMetadata(path)
	if err != nil {
		return nil, nil, err
	}

	return data, metadata.ToMap(), nil
}

// Delete deletes a file from the local storage
func (s *localStorage) Delete(ctx context.Context, path string) error {
	s.logger.Debug("deleting file", zap.String("path", path))

	fullPath := filepath.Join(s.directory, path)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Try to delete metadata file if it exists
	metadataPath := s.getMetadataPath(path)
	_ = os.Remove(metadataPath) // Ignore error as metadata file might not exist

	return nil
}
