package alert

import (
	"context"
	"fmt"
	"strings"
	"time"

	models "github.com/wizenheimer/iris/src/internal/models/core"
	"go.uber.org/zap"
)

// EventType represents the type of workflow event
type EventType string

const (
	EventStarted   EventType = "started"
	EventCompleted EventType = "completed"
	EventFailed    EventType = "failed"
	EventCancelled EventType = "cancelled"
)

// AlertTemplate defines the structure for alert templates
type AlertTemplate struct {
	TitleTemplate       string
	DescriptionTemplate string
	Severity            models.Severity
}

// AlertTemplates stores workflow-specific alert templates
var AlertTemplates = map[string]map[string]AlertTemplate{
	"screenshot": {
		"started": {
			TitleTemplate:       "Screenshot Workflow Started - Week %d/%d",
			DescriptionTemplate: "Screenshot capture workflow initiated for week %d of %d. Bucket: %d",
			Severity:            models.SeverityInfo,
		},
		"completed": {
			TitleTemplate:       "Screenshot Workflow Completed - Week %d/%d",
			DescriptionTemplate: "Successfully completed screenshot captures for week %d of %d. Bucket: %d. Total processed: %s",
			Severity:            models.SeverityInfo,
		},
		"failed": {
			TitleTemplate:       "Screenshot Workflow Failed - Week %d/%d",
			DescriptionTemplate: "Screenshot workflow failed for week %d of %d. Bucket: %d. Error: %s",
			Severity:            models.SeverityError,
		},
		"cancelled": {
			TitleTemplate:       "Screenshot Workflow Cancelled - Week %d/%d",
			DescriptionTemplate: "Screenshot workflow cancelled for week %d of %d. Bucket: %d. Reason: %s",
			Severity:            models.SeverityWarning,
		},
	},
	"report": {
		"started": {
			TitleTemplate:       "Report Generation Started - Week %d/%d",
			DescriptionTemplate: "Report generation initiated for week %d of %d. Bucket: %d",
			Severity:            models.SeverityInfo,
		},
		"completed": {
			TitleTemplate:       "Report Generation Completed - Week %d/%d",
			DescriptionTemplate: "Successfully generated reports for week %d of %d. Bucket: %d. Total processed: %s",
			Severity:            models.SeverityInfo,
		},
		"failed": {
			TitleTemplate:       "Report Generation Failed - Week %d/%d",
			DescriptionTemplate: "Report generation failed for week %d of %d. Bucket: %d. Error: %s",
			Severity:            models.SeverityError,
		},
		"cancelled": {
			TitleTemplate:       "Report Generation Cancelled - Week %d/%d",
			DescriptionTemplate: "Report generation cancelled for week %d of %d. Bucket: %d. Reason: %s",
			Severity:            models.SeverityWarning,
		},
	},
}

// workflowAlertClient implementation
func (w *workflowAlertClient) prepareAlert(id models.WorkflowIdentifier, details map[string]string) (models.Alert, error) {
	// Determine alert type from details
	alertType := w.getAlertTypeFromDetails(details)
	if alertType == "" {
		return models.Alert{}, fmt.Errorf("alert type not found in details")
	}

	// Get workflow type as string
	workflowType := string(*id.Type)

	// Get template for this workflow type and alert type
	template, err := w.getTemplate(workflowType, alertType)
	if err != nil {
		return models.Alert{}, err
	}

	// Prepare format arguments for templates
	args := w.prepareTemplateArgs(id, details)

	// Format title and description
	title := fmt.Sprintf(template.TitleTemplate, args...)
	description := fmt.Sprintf(template.DescriptionTemplate, args...)

	// Prepare enriched metadata
	metadata := w.enrichMetadata(id, details)

	return models.Alert{
		Title:       title,
		Description: description,
		Timestamp:   time.Now(),
		Severity:    template.Severity,
		Metadata:    metadata,
	}, nil
}

func (w *workflowAlertClient) getTemplate(workflowType, alertType string) (AlertTemplate, error) {
	templates, exists := AlertTemplates[workflowType]
	if !exists {
		return AlertTemplate{}, fmt.Errorf("no templates found for workflow type: %s", workflowType)
	}

	template, exists := templates[alertType]
	if !exists {
		return AlertTemplate{}, fmt.Errorf("no template found for alert type: %s", alertType)
	}

	return template, nil
}

func (w *workflowAlertClient) getAlertTypeFromDetails(details map[string]string) string {
	// Look for standard alert indicators in the details
	if strings.Contains(details["status"], "completed") {
		return "completed"
	}
	if strings.Contains(details["status"], "failed") {
		return "failed"
	}
	if strings.Contains(details["status"], "cancelled") ||
		strings.Contains(details["status"], "aborted") {
		return "cancelled"
	}
	if _, hasError := details["error"]; hasError {
		return "failed"
	}
	return "started"
}

func (w *workflowAlertClient) prepareTemplateArgs(id models.WorkflowIdentifier, details map[string]string) []interface{} {
	// Common arguments for all templates
	args := []interface{}{
		*id.WeekNumber,
		*id.Year,
		*id.WeekNumber,
		*id.Year,
		*id.WeekDay,
	}

	// Add specific arguments based on alert type
	alertType := w.getAlertTypeFromDetails(details)
	switch alertType {
	case "completed":
		args = append(args, details["processed_count"])
	case "failed":
		args = append(args, details["error"])
	case "cancelled":
		reason := details["reason"]
		if reason == "" {
			reason = "Manual cancellation"
		}
		args = append(args, reason)
	}

	return args
}

func (w *workflowAlertClient) enrichMetadata(id models.WorkflowIdentifier, details map[string]string) map[string]string {
	metadata := make(map[string]string)

	// Add workflow identifier information
	metadata["workflow_type"] = string(*id.Type)
	metadata["year"] = fmt.Sprintf("%d", *id.Year)
	metadata["week_number"] = fmt.Sprintf("%d", *id.WeekNumber)
	metadata["week_day"] = fmt.Sprintf("%d", *id.WeekDay)

	// Add all provided details
	for k, v := range details {
		metadata[k] = v
	}

	// Add additional context
	metadata["alert_generated_at"] = time.Now().Format(time.RFC3339)
	if taskID, ok := details["task_id"]; ok {
		metadata["task_id"] = taskID
	}

	return metadata
}

// sendWorkflowAlert is a helper function to reduce code duplication
func (w *workflowAlertClient) sendWorkflowAlert(
	ctx context.Context,
	id models.WorkflowIdentifier,
	eventType EventType,
	details map[string]string,
) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Initialize details map if nil
	if details == nil {
		details = make(map[string]string)
	}

	// Add event type to details
	details["event_type"] = string(eventType)
	details["workflow_type"] = string(*id.Type)

	// Prepare and send alert
	alert, err := w.prepareAlert(id, details)
	if err != nil {
		return fmt.Errorf("failed to prepare alert: %w", err)
	}

	if err := w.Send(ctx, alert); err != nil {
		return fmt.Errorf("failed to send alert: %w", err)
	}

	w.logger.Debug("sent workflow alert",
		zap.Any("event_type", eventType),
		zap.Any("workflow_id", id),
		zap.Any("details", details))

	return nil
}
