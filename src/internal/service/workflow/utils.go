package workflow

import (
	"fmt"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
)

func (s *workflowService) registerExecutor(wfType models.WorkflowType, executor interfaces.WorkflowExecutor) {
	s.executors.Store(wfType, executor)
}

func (s *workflowService) getExecutor(wfType models.WorkflowType) (interfaces.WorkflowExecutor, error) {
	if executor, ok := s.executors.Load(wfType); ok {
		return executor.(interfaces.WorkflowExecutor), nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) validateRequest(req *models.WorkflowRequest) error {
	return req.Validate(false)
}
