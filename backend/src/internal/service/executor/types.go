package executor

import (
	"context"
	"sync"

	clf "github.com/wizenheimer/byrd/src/internal/interfaces/client"
	exc "github.com/wizenheimer/byrd/src/internal/interfaces/executor"
	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workflowExecutor struct {
	workflowType models.WorkflowType
	config       models.ExecutorConfig
	repository   repo.WorkflowRepository
	alertClient  clf.WorkflowAlertClient
	taskExecutor exc.TaskExecutor
	logger       *logger.Logger

	activeWorkflows sync.Map // map[string]*workflowContext
}

type workflowContext struct {
	cancel context.CancelFunc
	task   models.Task
	state  api.WorkflowState
	mutex  sync.RWMutex
}
