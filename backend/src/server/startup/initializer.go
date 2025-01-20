package startup

import (
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"github.com/wizenheimer/byrd/src/server/startup/services"
)

func Initialize(cfg *config.Config, tm *transaction.TxManager, logger *logger.Logger) (*routes.HandlerContainer, *middleware.ResourceMiddleware, *middleware.AccessMiddleware, error) {
	// Initialize utilities
	utils.InitializeValidator()

	// Set up HTTP client
	screenshotClient, err := SetupScreenshotClient(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up services
	screenshotService, err := services.SetupScreenshotService(cfg, screenshotClient, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	aiService, err := services.SetupAIService(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	diffService, err := diff.NewDiffService(aiService, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up Redis
	redisClient, err := SetupRedis(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up repositories
	repos, err := SetupRepositories(tm, redisClient, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up all services
	services, err := SetupServices(cfg, repos, diffService, screenshotService, tm, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	resourceMiddleware := middleware.NewResourceMiddleware(services.Workspace, logger)
	accessMiddleware := middleware.NewAccessMiddleware(services.Workspace, services.User, services.TokenManager, logger)

	// Initialize handlers
	handlers := SetupHandlerContainer(
		screenshotService,
		aiService,
		services.User,
		services.Workspace,
		services.Workflow,
		services.Scheduler,
		tm,
		logger,
	)

	return handlers, resourceMiddleware, accessMiddleware, nil
}
