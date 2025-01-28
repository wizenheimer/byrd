// ./src/server/startup/services.go
package startup

import (
	"context"
	"time"

	"github.com/wizenheimer/byrd/src/internal/alert"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	"github.com/wizenheimer/byrd/src/internal/event"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/schedule"
	workflow_repo "github.com/wizenheimer/byrd/src/internal/repository/workflow"
	scheduler "github.com/wizenheimer/byrd/src/internal/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
	"github.com/wizenheimer/byrd/src/internal/service/history"
	"github.com/wizenheimer/byrd/src/internal/service/notification"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	scheduler_svc "github.com/wizenheimer/byrd/src/internal/service/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type Services struct {
	History             history.PageHistoryService
	Page                page.PageService
	Competitor          competitor.CompetitorService
	User                user.UserService
	Workspace           workspace.WorkspaceService
	Workflow            workflow.WorkflowService
	Scheduler           scheduler_svc.SchedulerService
	NotificationService notification.NotificationService
	TokenManager        *utils.TokenManager
}

func SetupServices(
	cfg *config.Config,
	repos *Repositories,
	diffService diff.DiffService,
	screenshotService screenshot.ScreenshotService,
	templateLibrary template.TemplateLibrary,
	tm *transaction.TxManager,
	logger *logger.Logger,
) (*Services, error) {
	historyService := history.NewPageHistoryService(repos.History, logger)
	pageService := page.NewPageService(repos.Page, historyService, diffService, screenshotService, logger)
	competitorService := competitor.NewCompetitorService(repos.Competitor, pageService, tm, logger)

	alertClient, err := setupAlertClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	emailClient, err := setupEmailClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	eventClient, err := setupEventClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	notificationService := notification.NewNotificationService(alertClient, eventClient, emailClient, logger)

	workflowService, err := setupWorkflowService(
		cfg,
		repos.Workflow,
		notificationService,
		pageService,
		logger,
	)
	if err != nil {
		return nil, err
	}

	schedulerSvc := setupSchedulerService(
		repos.Schedule,
		workflowService,
		logger,
	)

	if err := schedulerSvc.Start(context.Background(), true); err != nil {
		return nil, err
	}

	tokenManager := utils.NewTokenManager(cfg.Services.ManagementAPIKey, cfg.Services.ManagementAPIRefreshInterval)

	userService, err := user.NewUserService(notificationService, repos.User, templateLibrary, logger)
	if err != nil {
		return nil, err
	}

	workspaceService, err := workspace.NewWorkspaceService(repos.Workspace, competitorService, userService, notificationService, templateLibrary, tm, logger)
	if err != nil {
		return nil, err
	}

	return &Services{
		History:             historyService,
		Page:                pageService,
		Competitor:          competitorService,
		User:                userService,
		Workspace:           workspaceService,
		Workflow:            workflowService,
		NotificationService: notificationService,
		Scheduler:           schedulerSvc,
		TokenManager:        tokenManager,
	}, nil
}

func setupEmailClient(cfg *config.Config, logger *logger.Logger) (email.EmailClient, error) {
	if cfg.Environment.EnvProfile == "development" {
		logger.Debug("using local email client")
		return email.NewLocalEmailClient(context.Background(), logger)
	}

	return email.NewResendClient(context.Background(), cfg.Services.ResendAPIKey, cfg.Services.ResendNotificationEmail, logger)
}

func setupAlertClient(cfg *config.Config, logger *logger.Logger) (alert.AlertClient, error) {
	clientConfig := models.DefaultSlackConfig()
	clientConfig.Token = cfg.Workflow.SlackAlertToken
	clientConfig.ChannelID = cfg.Workflow.SlackWorkflowChannelId

	if cfg.Environment.EnvProfile == "development" {
		logger.Debug("using local workflow alert client")
		return alert.NewLocalWorkflowClient(clientConfig, logger), nil
	}

	return alert.NewSlackAlertClient(clientConfig, logger)
}

func setupEventClient(cfg *config.Config, logger *logger.Logger) (event.EventClient, error) {
	if cfg.Environment.EnvProfile == "development" {
		return event.NewLocalEventClient(logger), nil
	}
	return event.NewPostHogEventClient(cfg.Services.PostHogAPIKey, logger)
}

func setupWorkflowService(
	cfg *config.Config,
	workflowRepo workflow_repo.WorkflowRepository,
	notificationService notification.NotificationService,
	pageService page.PageService,
	logger *logger.Logger,
) (workflow.WorkflowService, error) {
	runtimeConfig := models.JobExecutorConfig{
		Parallelism: cfg.Workflow.ExecutorParallelism,
		LowerBound:  time.Duration(cfg.Workflow.ExecutorLowerBound) * time.Second,
		UpperBound:  time.Duration(cfg.Workflow.ExecutorUpperBound) * time.Second,
	}

	screenshotTaskExecutor, err := executor.NewPageExecutor(pageService, runtimeConfig, logger)
	if err != nil {
		return nil, err
	}

	screenshotWorkflowExecutor, err := executor.NewWorkflowExecutor(
		models.ScreenshotWorkflowType,
		workflowRepo,
		notificationService,
		screenshotTaskExecutor,
		logger,
	)
	if err != nil {
		return nil, err
	}

	workflowService, err := workflow.NewWorkflowService(logger)
	if err != nil {
		return nil, err
	}

	if err := workflowService.Register(models.ScreenshotWorkflowType, screenshotWorkflowExecutor); err != nil {
		return nil, err
	}

	if err := workflowService.Initialize(context.Background()); err != nil {
		return nil, err
	}

	return workflowService, nil
}

func setupSchedulerService(
	scheduleRepo schedule.ScheduleRepository,
	workflowService workflow.WorkflowService,
	logger *logger.Logger,
) scheduler_svc.SchedulerService {
	return scheduler_svc.NewSchedulerService(
		scheduleRepo,
		scheduler.NewScheduler(logger),
		workflowService,
		logger,
	)
}
