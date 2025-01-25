// ./src/internal/alert/local.go
package alert

import (
	"context"
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrFailedToSendBatchAlert = fmt.Errorf("failed to send batch alert")
)

// localWorkflowClient implements WorkflowAlertClient for local development
type localWorkflowClient struct {
	logger *logger.Logger
}

// NewLocalWorkflowClient creates a new local workflow client that logs alerts
func NewLocalWorkflowClient(_ models.SlackConfig, logger *logger.Logger) AlertClient {
	return &localWorkflowClient{
		logger: logger.WithFields(map[string]interface{}{"module": "local_workflow_alert_client"}),
	}
}

// Send implements AlertClient interface
func (c *localWorkflowClient) SendAlert(ctx context.Context, alert models.Alert) error {
	c.logger.Warn("Sending alert", zap.Any("alert", alert))
	return nil
}

// SendBatch implements AlertClient interface
func (c *localWorkflowClient) SendBatchAlert(ctx context.Context, alerts []models.Alert) error {
	for _, alert := range alerts {
		if err := c.SendAlert(ctx, alert); err != nil {
			return ErrFailedToSendBatchAlert
		}
	}
	return nil
}
