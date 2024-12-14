package alert

import (
	"context"
	"fmt"
	"time"

	"github.com/wizenheimer/iris/src/internal/client"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

// SlackWorkflowClient implements WorkflowAlertClient
type slackWorkflowClient struct {
	client.AlertClient
}

// NewSlackWorkflowClient creates a new Slack workflow client
func NewSlackWorkflowClient(config models.SlackConfig, logger *logger.Logger) (client.WorkflowAlertClient, error) {
	baseClient, err := NewSlackAlertClient(config, logger)
	if err != nil {
		return nil, err
	}
	return &slackWorkflowClient{
		baseClient,
	}, nil
}

func (s *slackWorkflowClient) SendWorkflowStarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["Workflow ID"] = id.Serialize()

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Started",
			Description: fmt.Sprintf("Workflow %s has started", id.Serialize()),
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return s.Send(ctx, alert.Alert)
}

// SendWorkflowRestarted implements WorkflowAlertClient interface
func (s *slackWorkflowClient) SendWorkflowRestarted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	metadata["Workflow ID"] = id.Serialize()

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Restarted",
			Description: fmt.Sprintf("Workflow %s has been restarted", id.Serialize()),
			Timestamp:   time.Now(),
			Severity:    models.SeverityWarning,
			Metadata:    metadata,
		},
	}
	return s.Send(ctx, alert.Alert)
}

// SendWorkflowProgress implements WorkflowAlertClient interface
func (s *slackWorkflowClient) SendWorkflowProgress(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	metadata["Workflow ID"] = id.Serialize()

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Current Progress",
			Description: fmt.Sprintf("Workflow %s progress update", id.Serialize()),
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return s.Send(ctx, alert.Alert)
}

// SendWorkflowCompleted implements WorkflowAlertClient interface
func (s *slackWorkflowClient) SendWorkflowCompleted(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	metadata["Workflow ID"] = id.Serialize()

	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Completed",
			Description: fmt.Sprintf("Workflow %s has completed successfully", id.Serialize()),
			Timestamp:   time.Now(),
			Severity:    models.SeverityInfo,
			Metadata:    metadata,
		},
	}
	return s.Send(ctx, alert.Alert)
}

// SendWorkflowFailed implements WorkflowAlertClient interface
func (s *slackWorkflowClient) SendWorkflowFailed(ctx context.Context, id models.WorkflowIdentifier, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}

	metadata["Workflow ID"] = id.Serialize()
	alert := models.WorkflowAlert{
		Alert: models.Alert{
			Title:       "Workflow Failed",
			Description: fmt.Sprintf("Workflow %s has failed", id.Serialize()),
			Timestamp:   time.Now(),
			Severity:    models.SeverityCritical,
			Metadata:    metadata,
		},
	}
	return s.Send(ctx, alert.Alert)
}
