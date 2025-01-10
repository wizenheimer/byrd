// ./src/internal/service/workflow/utils.go
package workflow

import (
	"fmt"

	exc "github.com/wizenheimer/byrd/src/internal/interfaces/executor"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

func (s *workflowService) registerExecutor(wfType models.WorkflowType, executor exc.WorkflowExecutor) {
	s.executors.Store(wfType, executor)
}

func (s *workflowService) getExecutor(wfType models.WorkflowType) (exc.WorkflowExecutor, error) {
	if executor, ok := s.executors.Load(wfType); ok {
		return executor.(exc.WorkflowExecutor), nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) validateRequest(req *api.WorkflowRequest) error {
	return req.Validate(false)
}
