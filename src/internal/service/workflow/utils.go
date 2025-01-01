package workflow

import (
	"fmt"

	exc "github.com/wizenheimer/iris/src/internal/interfaces/executor"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

func (s *workflowService) registerExecutor(wfType core_models.WorkflowType, executor exc.WorkflowExecutor) {
	s.executors.Store(wfType, executor)
}

func (s *workflowService) getExecutor(wfType core_models.WorkflowType) (exc.WorkflowExecutor, error) {
	if executor, ok := s.executors.Load(wfType); ok {
		return executor.(exc.WorkflowExecutor), nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) validateRequest(req *api_models.WorkflowRequest) error {
	return req.Validate(false)
}
