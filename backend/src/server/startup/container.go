// ./src/server/startup/container.go
package startup

import (
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	slackworkspace "github.com/wizenheimer/byrd/src/internal/service/integration/slack"
	"github.com/wizenheimer/byrd/src/internal/service/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/screenshot"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

func SetupHandlerContainer(
	screenshotService screenshot.ScreenshotService,
	aiService ai.AIService,
	userService user.UserService,
	workspaceService workspace.WorkspaceService,
	workflowService workflow.WorkflowService,
	schedulerService scheduler.SchedulerService,
	slackWorkspaceService slackworkspace.SlackWorkspaceService,
	library template.TemplateLibrary,
	emailClient email.EmailClient,
	tm *transaction.TxManager,
	logger *logger.Logger,
) (*routes.HandlerContainer, error) {
	return routes.NewHandlerContainer(
		screenshotService,
		aiService,
		userService,
		workspaceService,
		workflowService,
		schedulerService,
		slackWorkspaceService,
		library,
		emailClient,
		tm,
		logger,
	)
}
