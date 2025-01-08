package diff

import (
	"context"

	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/err"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils"
)

type diffService struct {
	aiService svc.AIService
	processor *utils.MarkdownProcessor
	logger    *logger.Logger
}

// NewDiffService creates a new diff service
func NewDiffService(aiService svc.AIService, logger *logger.Logger) (svc.DiffService, error) {
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

func (d *diffService) Compare(ctx context.Context, content1, content2 *models.ScreenshotHTMLContentResponse, profileStr string, persist bool) (*models.DynamicChanges, err.Error) {
	return nil, nil
}
