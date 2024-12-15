package client

import (
	"context"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

// AlertClient defines alerting capabilities
type AlertClient interface {
	// Send sends an alert
	Send(ctx context.Context, alert models.Alert) error

	// SendBatch sends a batch of alerts
	SendBatch(ctx context.Context, alerts []models.Alert) error
}

// WorkflowAlertClient defines workflow-specific alerting capabilities
type WorkflowAlertClient interface {
	// WorkflowAlertClient is an interface for sending alerts about workflow progress
	AlertClient

	// SendWorkflowStarted sends a notification that a workflow has started
	SendWorkflowStarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error

	// SendWorkflowFailed sends a notification that a workflow has failed
	SendWorkflowFailed(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error

	// SendWorkflowRestarted sends a notification that a workflow has restarted
	SendWorkflowRestarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error

	// SendWorkflowProgress sends a notification that a workflow has made significant progress
	SendWorkflowProgress(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error

	// SendWorkflowCompleted sends a notification that a workflow has completed
	SendWorkflowCompleted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error

	// SendWorkflowCompleted sends a notification that a workflow has cancelled
	SendWorkflowCancelled(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error
}
