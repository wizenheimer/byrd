package screenshot

import (
	clf "github.com/wizenheimer/byrd/src/internal/interfaces/client"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
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
