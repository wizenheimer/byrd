// ./src/server/startup/repositories.go
package startup

import (
	"github.com/redis/go-redis/v9"
	"github.com/wizenheimer/byrd/src/internal/repository/competitor"
	"github.com/wizenheimer/byrd/src/internal/repository/history"
	"github.com/wizenheimer/byrd/src/internal/repository/page"
	"github.com/wizenheimer/byrd/src/internal/repository/schedule"
	"github.com/wizenheimer/byrd/src/internal/repository/user"
	"github.com/wizenheimer/byrd/src/internal/repository/workflow"
	"github.com/wizenheimer/byrd/src/internal/repository/workspace"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type Repositories struct {
	Competitor competitor.CompetitorRepository
	Workspace  workspace.WorkspaceRepository
	User       user.UserRepository
	Page       page.PageRepository
	History    history.PageHistoryRepository
	Schedule   schedule.ScheduleRepository
	Workflow   workflow.WorkflowRepository
}

func SetupRepositories(tm *transaction.TxManager, redisClient *redis.Client, logger *logger.Logger) (*Repositories, error) {
	workflowRepo, err := workflow.NewWorkflowRepository(redisClient, tm, logger)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Competitor: competitor.NewCompetitorRepository(tm, logger),
		Workspace:  workspace.NewWorkspaceRepository(tm, logger),
		User:       user.NewUserRepository(tm, logger),
		Page:       page.NewPageRepository(tm, logger),
		History:    history.NewPageHistoryRepository(tm, logger),
		Schedule:   schedule.NewScheduleRepo(tm, logger),
		Workflow:   workflowRepo,
	}, nil
}
