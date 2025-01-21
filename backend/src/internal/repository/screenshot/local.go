package screenshot

import (
	"context"
	"fmt"
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
	return nil
}

// GetScreenshotImage retrieves a screenshot from local storage
func (s *localScreenshotRepo) StoreScreenshotContent(ctx context.Context, data *models.ScreenshotContent) error {
	return nil
}

// RetrieveScreenshotImage retrieves screenshot image from the storage
func (s *localScreenshotRepo) RetrieveScreenshotImage(ctx context.Context, path string) (*models.ScreenshotImage, error) {
	return nil, nil
}

// RetrieveScreenshotContent retrieves screenshot content from the storage
func (s *localScreenshotRepo) RetrieveScreenshotContent(ctx context.Context, path string) (*models.ScreenshotContent, error) {
	return nil, nil
}
