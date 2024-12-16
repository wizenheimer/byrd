package alert

import (
	"context"
	"fmt"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// localWorkflowClient implements WorkflowAlertClient for local development
type localWorkflowClient struct {
	logger *logger.Logger
}

// NewLocalWorkflowClient creates a new local workflow client that logs alerts
func NewLocalWorkflowClient(_ models.SlackConfig, logger *logger.Logger) interfaces.AlertClient {
	return &localWorkflowClient{
		logger: logger.WithFields(map[string]interface{}{"module": "local_workflow_alert_client"}),
	}
}

// Send implements AlertClient interface
func (c *localWorkflowClient) Send(ctx context.Context, alert models.Alert) error {
	c.logger.Warn("Sending alert", zap.Any("alert", alert))
	return nil
}

// SendBatch implements AlertClient interface
func (c *localWorkflowClient) SendBatch(ctx context.Context, alerts []models.Alert) error {
	for _, alert := range alerts {
		if err := c.Send(ctx, alert); err != nil {
			return fmt.Errorf("failed to send alert in batch: %w", err)
		}
	}
	return nil
}
