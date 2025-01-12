// ./src/internal/service/workflow/utils.go
package workflow

import (
	"fmt"

	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/executor"
	"go.uber.org/zap"
)

func (s *workflowService) registerExecutor(wfType models.WorkflowType, executor executor.WorkflowExecutor) {
	s.logger.Debug("registering executor", zap.String("type", string(wfType)))
	s.executors[wfType] = executor
}

func (s *workflowService) getExecutor(wfType models.WorkflowType) (executor.WorkflowExecutor, error) {
	s.logger.Debug("getting executor", zap.String("type", string(wfType)))
	if executor, ok := s.executors[wfType]; ok {
		return executor, nil
	}
	return nil, fmt.Errorf("no executor found for type: %s", wfType)
}

func (s *workflowService) validateRequest(req *api.WorkflowRequest) error {
	return req.Validate(false)
}
