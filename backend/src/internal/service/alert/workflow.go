// ./src/internal/service/alert/workflow.go
package alert

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

type workflowAlertClient struct {
	// client for Slack
	AlertClient
	// logger for logging
	logger *logger.Logger
}

func NewWorkflowAlertClient(embeddedClient AlertClient, logger *logger.Logger) (WorkflowAlertClient, error) {

	workflowAlertClient := workflowAlertClient{
		AlertClient: embeddedClient,
		logger:      logger.WithFields(map[string]interface{}{"module": "workflow_alert_client"}),
	}

	return &workflowAlertClient, nil
}

// SendWorkflowStarted sends an alert when a workflow starts
func (w *workflowAlertClient) SendWorkflowStarted(
	ctx context.Context,
	id models.WorkflowIdentifier,
	details map[string]string,
) error {
	return w.sendWorkflowAlert(ctx, id, EventStarted, details)
}

// SendWorkflowCompleted sends an alert when a workflow completes
func (w *workflowAlertClient) SendWorkflowCompleted(
	ctx context.Context,
	id models.WorkflowIdentifier,
	details map[string]string,
) error {
	return w.sendWorkflowAlert(ctx, id, EventCompleted, details)
}

// SendWorkflowFailed sends an alert when a workflow fails
func (w *workflowAlertClient) SendWorkflowFailed(
	ctx context.Context,
	id models.WorkflowIdentifier,
	details map[string]string,
) error {
	return w.sendWorkflowAlert(ctx, id, EventFailed, details)
}

// SendWorkflowCancelled sends an alert when a workflow is cancelled
func (w *workflowAlertClient) SendWorkflowCancelled(
	ctx context.Context,
	id models.WorkflowIdentifier,
	details map[string]string,
) error {
	return w.sendWorkflowAlert(ctx, id, EventCancelled, details)
}
