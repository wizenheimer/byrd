// ./src/internal/service/workspace/service.go
package workspace

import (
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/repository/workspace"
	"github.com/wizenheimer/byrd/src/internal/service/competitor"
	"github.com/wizenheimer/byrd/src/internal/service/user"
	"github.com/wizenheimer/byrd/src/internal/transaction"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workspaceService struct {
	workspaceRepo     workspace.WorkspaceRepository
	competitorService competitor.CompetitorService
	library           template.TemplateLibrary
	emailClient       email.EmailClient
	userService       user.UserService
	logger            *logger.Logger
	errorRecord       *recorder.ErrorRecorder
	tm                *transaction.TxManager
}

func NewWorkspaceService(
	workspaceRepo workspace.WorkspaceRepository,
	competitorService competitor.CompetitorService,
	userService user.UserService,
	library template.TemplateLibrary,
	tm *transaction.TxManager,
	emailClient email.EmailClient,
	logger *logger.Logger,
	errorRecord *recorder.ErrorRecorder,
) (WorkspaceService, error) {

	ws := workspaceService{
		workspaceRepo:     workspaceRepo,
		competitorService: competitorService,
		userService:       userService,
		library:           library,
		logger: logger.WithFields(map[string]any{
			"module": "workspace_service",
		}),
		emailClient: emailClient,
		errorRecord: errorRecord,
		tm:          tm,
	}

	return &ws, nil
}
