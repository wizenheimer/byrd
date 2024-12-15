package alert

import (
	"context"
	"fmt"
	"time"

	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

// localWorkflowClient implements WorkflowAlertClient for local development
type localWorkflowClient struct {
	logger *logger.Logger
}

// NewLocalWorkflowClient creates a new local workflow client that logs alerts
func NewLocalWorkflowClient(_ models.SlackConfig, logger *logger.Logger) client.WorkflowAlertClient {
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

func (c *localWorkflowClient) SendWorkflowStarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Started",
			Description: "Workflow has started",
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}

func (c *localWorkflowClient) SendWorkflowRestarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Restarted",
			Description: "Workflow has restarted",
			Timestamp:   time.Now(),
			Severity:    models.SeverityWarning,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}

func (c *localWorkflowClient) SendWorkflowProgress(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Current Progress",
			Description: "Workflow progress update",
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}

func (c *localWorkflowClient) SendWorkflowCompleted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Completed",
			Description: "Workflow marked as complete",
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}

func (c *localWorkflowClient) SendWorkflowCancelled(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Cancelled",
			Description: "Workflow has been cancelled by host",
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}

func (c *localWorkflowClient) SendWorkflowFailed(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Failed",
			Description: "Workflow has failed",
			Timestamp:   time.Now(),
			Severity:    models.SeverityCritical,
			Metadata:    metadata,
		},
	}
	return c.Send(ctx, alert.Alert)
}
