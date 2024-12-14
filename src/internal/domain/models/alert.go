package models

import "time"

// Alert represents the base alert structure
type Alert struct {
	Title       string
	Description string
	Timestamp   time.Time
	Severity    Severity
	Metadata    map[string]string
}

// Severity represents alert severity levels
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
)

// WorkflowAlert extends base Alert for workflow-specific information
type WorkflowAlert struct {
	// Alert is the base alert
	Alert
}

type WorkflowAlertType string

const (
	// When a workflow is started
	WorkflowStartedAlert WorkflowAlertType = "STARTED"
	// When a workflow fails on startup
	WorkflowFailedAlert WorkflowAlertType = "FAILED"
	// When a workflow is restarted
	WorkflowRestartedAlert WorkflowAlertType = "RESTARTED"
	// When a workflow is in progress
	WorkflowProgressAlert WorkflowAlertType = "PROGRESS"
	// When a workflow is completed
	WorkflowCompletedAlert WorkflowAlertType = "COMPLETED"
)
