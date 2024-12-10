package storage

import (
	"context"

	"github.com/wizenheimer/iris/internal/domain/interfaces"
)

type localStorage struct {
	directory string
}

func NewlocalStorage(directory string) (interfaces.StorageRepository, error) {
	// Implementation
	return &localStorage{
		directory: directory,
	}, nil
}

func (s *localStorage) StoreScreenshot(ctx context.Context, data []byte, path string, metadata map[string]string) error {
	// Implementation
	return nil
}

func (s *localStorage) StoreContent(ctx context.Context, content string, path string, metadata map[string]string) error {
	// Implementation
	return nil
}

func (s *localStorage) Get(ctx context.Context, path string) ([]byte, map[string]string, error) {
	// Implementation
	return nil, nil, nil
}

func (s *localStorage) Delete(ctx context.Context, path string) error {
	// Implementation
	return nil
}
