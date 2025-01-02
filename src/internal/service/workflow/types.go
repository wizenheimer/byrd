package workflow

import (
	"sync"

	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type workflowService struct {
	executors  sync.Map // map[WorkflowType]WorkflowExecutor
	logger     *logger.Logger
	repository repo.WorkflowRepository
}
