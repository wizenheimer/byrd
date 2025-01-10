// ./src/internal/service/workflow/types.go
package workflow

import (
	"sync"

	repo "github.com/wizenheimer/byrd/src/internal/interfaces/repository"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workflowService struct {
	executors  sync.Map // map[WorkflowType]WorkflowExecutor
	logger     *logger.Logger
	repository repo.WorkflowRepository
}
