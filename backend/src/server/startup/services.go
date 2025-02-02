// ./src/server/startup/services.go
package startup

import (
	"context"
	"time"

	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/repository/schedule"
	workflow_repo "github.com/wizenheimer/byrd/src/internal/repository/workflow"
	scheduler "github.com/wizenheimer/byrd/src/internal/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/internal/service/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/diff"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
	"github.com/wizenheimer/byrd/src/internal/service/history"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/internal/service/report"
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
	History      history.PageHistoryService
	Page         page.PageService
	Competitor   competitor.CompetitorService
	User         user.UserService
	Workspace    workspace.WorkspaceService
	Workflow     workflow.WorkflowService
	Scheduler    scheduler_svc.SchedulerService
	TokenManager *utils.TokenManager
}

func SetupServices(
	cfg *config.Config,
	repos *Repositories,
	aiService ai.AIService,
	diffService diff.DiffService,
	screenshotService screenshot.ScreenshotService,
	templateLibrary template.TemplateLibrary,
	emailClient email.EmailClient,
	tm *transaction.TxManager,
	logger *logger.Logger,
	errorRecorder *recorder.ErrorRecorder,
) (*Services, error) {
	historyService := history.NewPageHistoryService(repos.History, logger)
	pageService := page.NewPageService(repos.Page, historyService, diffService, screenshotService, logger)

	reportService, err := report.NewReportService(aiService, emailClient, templateLibrary, repos.Report, logger, errorRecorder)
	if err != nil {
		return nil, err
	}

	competitorService := competitor.NewCompetitorService(pageService, reportService, tm, repos.Competitor, logger)

	tokenManager := utils.NewTokenManager(cfg.Services.ManagementAPIKey, cfg.Services.ManagementAPIRefreshInterval)

	userService, err := user.NewUserService(repos.User, templateLibrary, logger, errorRecorder)
	if err != nil {
		return nil, err
	}

	workspaceService, err := workspace.NewWorkspaceService(repos.Workspace, competitorService, userService, templateLibrary, tm, logger, errorRecorder)
	if err != nil {
		return nil, err
	}

	workflowService, err := setupWorkflowService(
		cfg,
		repos.Workflow,
		pageService,
		workspaceService,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}

	schedulerSvc, err := setupSchedulerService(
		repos.Schedule,
		workflowService,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}

	if err := schedulerSvc.Start(context.Background(), true); err != nil {
		return nil, err
	}

	return &Services{
		History:      historyService,
		Page:         pageService,
		Competitor:   competitorService,
		User:         userService,
		Workspace:    workspaceService,
		Workflow:     workflowService,
		Scheduler:    schedulerSvc,
		TokenManager: tokenManager,
	}, nil
}

func setupEmailClient(cfg *config.Config, logger *logger.Logger) (email.EmailClient, error) {
	if cfg.Environment.EnvProfile == "development" {
		return email.NewLocalEmailClient(context.Background(), logger)
	}

	return email.NewResendClient(context.Background(), cfg.Services.ResendAPIKey, cfg.Services.ResendNotificationEmail, logger)
}

func setupWorkflowService(
	cfg *config.Config,
	workflowRepo workflow_repo.WorkflowRepository,
	pageService page.PageService,
	workspaceService workspace.WorkspaceService,
	logger *logger.Logger,
	errorRecorder *recorder.ErrorRecorder,
) (workflow.WorkflowService, error) {
	screenshotTaskRuntimeConfig := models.JobExecutorConfig{
		Parallelism: cfg.Workflow.ScreenshotExecutorParallelism,
		LowerBound:  time.Duration(cfg.Workflow.ScreenshotExecutorLowerBound) * time.Second,
		UpperBound:  time.Duration(cfg.Workflow.ScreenshotExecutorUpperBound) * time.Second,
	}

	screenshotTaskExecutor, err := executor.NewPageExecutor(pageService, screenshotTaskRuntimeConfig, logger)
	if err != nil {
		return nil, err
	}

	screenshotWorkflowObserver, err := executor.NewWorkflowObserver(
		models.ScreenshotWorkflowType,
		workflowRepo,
		screenshotTaskExecutor,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}

	reportTaskRuntimeConfig := models.JobExecutorConfig{
		Parallelism: cfg.Workflow.ReportExecutorParallelism,
		LowerBound:  time.Duration(cfg.Workflow.ReportExecutorLowerBound) * time.Second,
		UpperBound:  time.Duration(cfg.Workflow.ReportExecutorUpperBound) * time.Second,
	}

	reportTaskExecutor, err := executor.NewReportExecutor(workspaceService, logger, reportTaskRuntimeConfig)
	if err != nil {
		return nil, err
	}

	reportWorkflowObserver, err := executor.NewWorkflowObserver(
		models.ReportWorkflowType,
		workflowRepo,
		reportTaskExecutor,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}

	dispatchTaskExecutor, err := executor.NewDispatchExecutor(workspaceService, logger, reportTaskRuntimeConfig)
	if err != nil {
		return nil, err
	}

	dispatchWorkflowObserver, err := executor.NewWorkflowObserver(
		models.DispatchWorkflowType,
		workflowRepo,
		dispatchTaskExecutor,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}

	workflowService, err := workflow.NewWorkflowService(logger, errorRecorder)
	if err != nil {
		return nil, err
	}

	if err := workflowService.Register(models.ScreenshotWorkflowType, screenshotWorkflowObserver); err != nil {
		return nil, err
	}

	if err := workflowService.Register(models.ReportWorkflowType, reportWorkflowObserver); err != nil {
		return nil, err
	}

	if err := workflowService.Register(models.DispatchWorkflowType, dispatchWorkflowObserver); err != nil {
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
	errorRecorder *recorder.ErrorRecorder,
) (scheduler_svc.SchedulerService, error) {
	s, err := scheduler_svc.NewSchedulerService(
		scheduleRepo,
		scheduler.NewScheduler(logger),
		workflowService,
		logger,
		errorRecorder,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}
