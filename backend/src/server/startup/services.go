// ./src/server/startup/services.go
package startup

import (
	"context"
	"fmt"
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
	"go.uber.org/zap"
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
	tm *transaction.TxManager,
	logger *logger.Logger,
) (*Services, error) {
	templateLibrary := setupLibrary(logger)

	historyService := history.NewPageHistoryService(repos.History, logger)
	pageService := page.NewPageService(repos.Page, historyService, diffService, screenshotService, logger)
	competitorService := competitor.NewCompetitorService(repos.Competitor, pageService, tm, logger)
	userService := user.NewUserService(repos.User, logger)
	workspaceService := workspace.NewWorkspaceService(repos.Workspace, competitorService, userService, templateLibrary, tm, logger)

	alertClient, err := setupAlertClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	workflowService, err := setupWorkflowService(
		cfg,
		repos.Workflow,
		pageService,
		alertClient,
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

	emailClient, err := email.NewResendClient(context.Background(), cfg.Services.ResendAPIKey, cfg.Services.ResendNotificationEmail, logger)
	if err != nil {
		return nil, err
	}

	eventClient, err := setupEventClient(logger)
	if err != nil {
		return nil, err
	}

	notificationService := notification.NewNotificationService(alertClient, eventClient, emailClient, logger)

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

func setupLibrary(logger *logger.Logger) template.TemplateLibrary {
	library := template.NewTemplateLibrary(logger)
	if err := registerDefaultTemplates(library); err != nil {
		logger.Error("failed to register default templates", zap.Error(err))
	}

	return library
}

// registerDefaultTemplates pre-registers all the default email templates
func registerDefaultTemplates(lib template.TemplateLibrary) error {
	templates := map[string]template.Template{
		"welcome": &template.CommonTemplate{
			PreviewText: "Your Competition's Worst Nightmare Just Got Real",
			Title:       "Your Competition's Worst Nightmare Just Got Real",
			Subtitle:    "The wait is over - We are officially yours!",
			Body: []string{
				"Hi there,",
				"Wish you could see every move your competitors make? Well, now you can. With Byrd, every change, no matter how small, gets flagged. From pricing shifts to product updates, you'll be the first to know.",
			},
			BulletTitle: "WHAT'S INCLUDED",
			Bullets: []string{
				"Page Monitoring - Stay ahead of every product move",
				"Inbox Monitoring - Monitor the direct line to their customers",
				"Social Monitoring - Keep a pulse on their community",
				"Review Monitoring - They churn. You learn.",
			},
			CTA: &template.CallToAction{
				ButtonText: "Get Started Now",
				ButtonURL:  "https://byrdhq.com/dashboard",
				FooterText: "They're hoping this invite expires. Disappoint them.",
			},
			Footer: template.Footer{
				ContactMessage: "Need help? We've got your back:",
				ContactEmail:   "hey@byrd.com",
			},
			GeneratedAt: time.Now(),
		},
		"halfway": &template.CommonTemplate{
			PreviewText: "Don't get back to playing catch-up. Keep calling the shots",
			Title:       "Don't get back to playing catch-up",
			Subtitle:    "Keep calling the shots.",
			Body: []string{
				"Hi there,",
				"Your trial hits the halfway mark, today. Like what you're seeing? Let's make it permanent with a paid plan.",
				"Subscribe to keep your front-row seat to every competitor move. The view's better when you can see everything coming.",
			},
			CTA: &template.CallToAction{
				ButtonText: "Subscribe Now",
				ButtonURL:  "https://byrdhq.com/upgrade",
				FooterText: "Because busy work is for your competitors",
			},
			Footer: template.Footer{
				ContactMessage: "Want to discuss a quote that works for your team?",
				ContactEmail:   "hey@byrd.com",
			},
			GeneratedAt: time.Now(),
		},
		"expiration_warning": &template.CommonTemplate{
			PreviewText: "Lock In Your Advantage. Keep dealing from the top.",
			Title:       "Lock In Your Advantage",
			Subtitle:    "Stack the deck in your favor",
			Body: []string{
				"Hi there,",
				"The countdown is on. In 3 days, your competitors get their blindspots back. Unless you choose differently.",
				"If you're still evaluating plans, let's talk. We've helped dozens of teams pick the right setup. Yearly plan saves you 20% - many teams find it's a no-brainer for the budget.",
			},
			CTA: &template.CallToAction{
				ButtonText: "Subscribe Now",
				ButtonURL:  "https://byrdhq.com/upgrade",
				FooterText: "Because knowing first means moving first.",
			},
			Footer: template.Footer{
				ContactMessage: "Want to discuss a quote that works for your team?",
				ContactEmail:   "hey@byrd.com",
			},
			GeneratedAt: time.Now(),
		},
		"trial_conversion": &template.CommonTemplate{
			PreviewText: "That's not marketing. That's just math",
			Title:       "Save One Deal, and We're Effectively Free",
			Subtitle:    "That's not marketing. That's just math",
			Body: []string{
				"Hi there,",
				"Your trial has been running strong for 14 days now, and your competitors have been keeping us busy.",
				"Ready to make all of this Official? Subscribe to a paid plan to lock in everything you're using now, plus some features tailored just for you.",
			},
			CTA: &template.CallToAction{
				ButtonText: "Subscribe Now",
				ButtonURL:  "https://byrdhq.com/upgrade",
				FooterText: "Lock in early-adopter pricing",
			},
			Footer: template.Footer{
				ContactMessage: "Want to discuss a quote that works for your team?",
				ContactEmail:   "hey@byrd.com",
			},
			GeneratedAt: time.Now(),
		},
		"post_trial_feedback": &template.CommonTemplate{
			PreviewText: "Your Spot's Still Warm. Let's Figure This Out Together.",
			Title:       "Your Spot's Still Warm",
			Subtitle:    "Let's Figure This Out Together",
			Body: []string{
				"Hi there,",
				"Not going to lie - our team's going to miss having you around. Wanted to drop by and see if we can do right by you.",
				"Mind if I ask what held you back from an upgrade? Hit reply with any feedback (good or bad), and we'll extend your access for a month.",
			},
			CTA: &template.CallToAction{
				ButtonText: "Claim an Extension",
				ButtonURL:  "https://byrdhq.com/upgrade",
				FooterText: "Your first month on us. No strings attached.",
			},
			Footer: template.Footer{
				ContactMessage: "Want to discuss a quote that works for your team?",
				ContactEmail:   "hey@byrd.com",
			},
			GeneratedAt: time.Now(),
		},
		"renewal_success": &template.CommonTemplate{
			PreviewText: "You made our nights and weekends worth it (and yes, your access is confirmed)",
			Title:       "More Than Just a Renewal",
			Subtitle:    "Team just did a happy dance (your renewal triggered it)",
			Body: []string{
				"Hi there,",
				"We know this is supposed to be a standard payment confirmation email, but instead we wanted to drop by and say a heartfelt thanks - not just for the wire, but for all the feedback that's helped make Byrd better.",
				"Teams like yours are why we get excited about competitive intelligence. To make this better, we've snuck a few extra user seats to your account (on the house). Here's to more winning moves, together with Byrd.",
			},
			Footer:      template.Footer{},
			GeneratedAt: time.Now(),
		},
		"waitlist": &template.CommonTemplate{
			PreviewText: "Good things come to those who... actually, let's speed this up",
			Title:       "You're Almost There",
			Body: []string{
				"Hi there,",
				"We're thrilled to have you on board! We're gradually rolling out access for new teams, and will reach out to you with your onboarding details as soon as your spot opens up.",
				"While we know waiting isn't ideal, we've prepared some excellent resources to help you make the most of this time. From swipe files to sales battlecards, you'll have access to the exact tools our top users rely on - and they're yours to keep!",
			},
			ClosingText: "Can't Wait? Don't Wait\nIncase waiting doesn't work for you (we understand!), reach out to us. We're founders too and we occasionally fast-track access for teams who are ready to dive straight in.",
			Footer:      template.Footer{},
			GeneratedAt: time.Now(),
		},
		"weekly_roundup": &template.SectionedTemplate{
			Competitor:  "Competitor X",
			FromDate:    time.Now().AddDate(0, 0, -7),
			ToDate:      time.Now(),
			Summary:     "This week has been particularly active with major updates across branding and pricing. Here's what you need to know.",
			GeneratedAt: time.Now(),
			Sections: map[string]template.Section{
				"branding": {
					Title:   "BRANDING",
					Summary: "Major brand refresh and positioning updates",
					Bullets: []template.BulletPoint{
						{
							Text:    "Updated logo and visual identity",
							LinkURL: "https://example.com/brand-update",
						},
						{
							Text:    "New brand guidelines released",
							LinkURL: "https://example.com/guidelines",
						},
					},
				},
				"pricing": {
					Title:   "PRICING",
					Summary: "New pricing structure implemented",
					Bullets: []template.BulletPoint{
						{
							Text:    "Introduced new enterprise tier",
							LinkURL: "https://example.com/enterprise",
						},
					},
				},
			},
		},
	}

	// Register each template
	for name, tmpl := range templates {
		if err := lib.RegisterTemplate(name, tmpl); err != nil {
			return fmt.Errorf("failed to register template %s: %w", name, err)
		}
	}

	return nil
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

func setupEventClient(logger *logger.Logger) (event.EventClient, error) {
	eventClient := event.NewLocalEventClient(logger)
	// TODO: add environment specific event client
	return eventClient, nil
}

func setupWorkflowService(
	cfg *config.Config,
	workflowRepo workflow_repo.WorkflowRepository,
	pageService page.PageService,
	alertClient alert.AlertClient,
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
		alertClient,
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
