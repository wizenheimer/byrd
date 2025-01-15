// ./src/internal/models/core/workflow.go
package models

import (
	"fmt"
	"time"
)

// WorkflowType is an enum for the type of workflow
type WorkflowType string

const (
	// ScreenshotWorkflowType is a workflow that takes a screenshot
	ScreenshotWorkflowType WorkflowType = "screenshot"
	// ReportWorkflowType is a workflow that generates a report
	ReportWorkflowType WorkflowType = "report"
)

// Parses a string into a WorkflowType
func ParseWorkflowType(s string) (WorkflowType, error) {
	switch WorkflowType(s) {
	case ScreenshotWorkflowType, ReportWorkflowType:
		return WorkflowType(s), nil
	default:
		return "", fmt.Errorf("invalid workflow type: %s", s)
	}
}

// WorkflowSchedule represents a scheduled workflow
type WorkflowSchedule struct {
	// ID is the unique identifier for the scheduled workflow
	ID ScheduleID `json:"id"`

	// WorkflowType is the type of the workflow
	WorkflowType WorkflowType `json:"workflow_type"`

	// About is the description of the workflow
	About string `json:"about"`

	// Spec is the cron specification for the workflow
	Spec string `json:"spec"`

	// LastRun is the time when the workflow was last run
	LastRun time.Time `json:"last_run"`

	// NextRun is the time when the workflow is scheduled to run next
	NextRun time.Time `json:"next_run"`

	// CreatedAt is the time when the workflow was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time when the workflow was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkflowScheduleProps struct {
	// WorkflowType is the type of the workflow
	WorkflowType WorkflowType `json:"workflow_type" default:"screenshot" validate:"required,oneof=screenshot report"`

	// About is the description of the workflow
	About string `json:"about" default:""`

	// Spec is the cron specification for the workflow
	Spec string `json:"spec" required:"true" validate:"required"`
}
