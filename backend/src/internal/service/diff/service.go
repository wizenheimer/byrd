// ./src/internal/service/diff/service.go
package diff

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type diffService struct {
	aiService ai.AIService
	processor *utils.MarkdownProcessor
	logger    *logger.Logger
}

// NewDiffService creates a new diff service
func NewDiffService(aiService ai.AIService, logger *logger.Logger) (DiffService, error) {
	processor, err := utils.NewMarkdownProcessor()
	if err != nil {
		return nil, err
	}
	return &diffService{
		aiService: aiService,
		processor: processor,
		logger:    logger,
	}, nil
}

func (d *diffService) Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, error) {
	return nil, nil
}
