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
