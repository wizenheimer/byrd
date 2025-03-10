// ./src/server/startup/initializer.go
package startup

import (
	"context"

	"github.com/wizenheimer/byrd/src/internal/api/middleware"
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"github.com/wizenheimer/byrd/src/server/startup/services"
)

func Initialize(
	ctx context.Context,
	cfg *config.Config,
	logger *logger.Logger,
	errorRecorder *recorder.ErrorRecorder,
) (*routes.HandlerContainer, *middleware.ResourceMiddleware, *middleware.AccessMiddleware, error) {
	// Initialize utilities
	utils.InitializeValidator()

	// Initialize database
	sqlDb, err := SetupDB(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	// Initialize transaction manager
	tm := transaction.NewTxManager(sqlDb, logger)

	// Set up screenshot client
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

	// Setup email client
	emailClient, err := setupEmailClient(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up repositories
	repos, err := SetupRepositories(ctx, cfg, tm, redisClient, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up template library
	templateLibrary, err := template.NewTemplateLibrary(logger)
	if err != nil {
		return nil, nil, nil, err
	}

	// Set up all services
	services, err := SetupServices(cfg, repos, aiService, diffService, screenshotService, templateLibrary, emailClient, tm, logger, errorRecorder)
	if err != nil {
		return nil, nil, nil, err
	}

	resourceMiddleware := middleware.NewResourceMiddleware(services.Workspace, logger)
	accessMiddleware := middleware.NewAccessMiddleware(services.Workspace, services.User, services.TokenManager, logger)

	// Initialize handlers
	handlers, err := SetupHandlerContainer(
		screenshotService,
		aiService,
		services.User,
		services.Workspace,
		services.Workflow,
		services.Scheduler,
		services.SlackWorkspace,
		templateLibrary,
		emailClient,
		tm,
		logger,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return handlers, resourceMiddleware, accessMiddleware, nil
}
