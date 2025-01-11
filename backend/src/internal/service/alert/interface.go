package alert

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// AlertClient defines alerting capabilities
type AlertClient interface {
	// Send sends an alert
	Send(ctx context.Context, alert models.Alert) error

	// SendBatch sends a batch of alerts
	SendBatch(ctx context.Context, alerts []models.Alert) error
}

// WorkflowAlertClient represents the client for sending alerts for workflows
type WorkflowAlertClient interface {
	// AlertClient represents the client for sending alerts
	AlertClient

	// SendWorkflowStarted sends an alert when a workflow is started
	SendWorkflowStarted(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowCompleted sends an alert when a workflow is completed
	SendWorkflowCompleted(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowFailed sends an alert when a workflow fails
	SendWorkflowFailed(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error

	// SendWorkflowCancelled sends an alert when a workflow is cancelled
	SendWorkflowCancelled(ctx context.Context, id models.WorkflowIdentifier, details map[string]string) error
}
