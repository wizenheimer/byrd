// ./src/internal/service/workflow/utils.go
package workflow

import (
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
)

func (s *workflowService) getExecutor(wfType models.WorkflowType) (executor.WorkflowExecutor, error) {
	if exc, ok := s.executors.Load(wfType); ok {
		return exc.(executor.WorkflowExecutor), nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) executorExists(wfType models.WorkflowType) bool {
	_, ok := s.executors.Load(wfType)
	return ok
}
