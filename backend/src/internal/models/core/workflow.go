// ./src/internal/models/core/workflow.go
package models

import (
	"fmt"
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
