// ./src/server/startup/container.go
package startup

import (
	"github.com/wizenheimer/byrd/src/internal/api/routes"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
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
	tm *transaction.TxManager,
	logger *logger.Logger,
) *routes.HandlerContainer {
	return routes.NewHandlerContainer(
		screenshotService,
		aiService,
		userService,
		workspaceService,
		workflowService,
		schedulerService,
		tm,
		logger,
	)
}
