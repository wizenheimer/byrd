// ./src/internal/service/workflow/utils.go
package workflow

import (
	"fmt"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
)

func (s *workflowService) registerExecutor(wfType models.WorkflowType, exc executor.WorkflowExecutor) {
	s.executors.Store(wfType, exc)
}

func (s *workflowService) getExecutor(wfType models.WorkflowType) (executor.WorkflowExecutor, error) {
	if exc, ok := s.executors.Load(wfType); ok {
		return exc.(executor.WorkflowExecutor), nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) validateRequest(req *api.WorkflowRequest) error {
	return req.Validate(false)
}
