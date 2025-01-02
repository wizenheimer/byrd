package screenshot

import (
	clf "github.com/wizenheimer/iris/src/internal/interfaces/client"
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type screenshotService struct {
	storage    repo.ScreenshotRepository
	httpClient clf.HTTPClient
	config     *models.ScreenshotServiceConfig
	logger     *logger.Logger
}

// compile time check if the interface is implemented
// TODO: reduce overhead by passing stuff by reference
var _ svc.ScreenshotService = (*screenshotService)(nil)
