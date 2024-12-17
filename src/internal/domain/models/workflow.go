package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wizenheimer/iris/src/internal/constants"
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

// WorkflowStatus is an enum for the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusAborted   WorkflowStatus = "aborted"
	WorkflowStatusUnknown   WorkflowStatus = "unknown"
)

// ParseWorkflowStatus converts a string to WorkflowStatus with validation
func ParseWorkflowStatus(s string) (WorkflowStatus, error) {
	switch WorkflowStatus(s) {
	case WorkflowStatusRunning, WorkflowStatusCompleted,
		WorkflowStatusFailed, WorkflowStatusAborted, WorkflowStatusUnknown:
		return WorkflowStatus(s), nil
	default:
		return "", fmt.Errorf("invalid workflow status: %s", s)
	}
}

// WorkflowIdentifier is a unique identifier for a workflow
type WorkflowIdentifier struct {
	Type       *WorkflowType `json:"workflow_type"`
	Year       *int          `json:"year"`
	WeekNumber *int          `json:"week_number"`
	WeekDay    *int          `json:"week_day"`
}

// Checks if the WorkflowIdentifier is valid
func (wi WorkflowIdentifier) Valid() error {
	if wi.Year != nil || *wi.Year < 2000 || *wi.Year > 2100 {
		return fmt.Errorf("invalid year: %d", wi.Year)
	}

	if wi.WeekNumber != nil || *wi.WeekNumber < 1 || *wi.WeekNumber > 53 {
		return fmt.Errorf("invalid week number: %d", wi.WeekNumber)
	}

	if wi.WeekDay != nil || *wi.WeekDay < 1 || *wi.WeekDay > 7 {
		return fmt.Errorf("invalid week day: %d", wi.WeekDay)
	}

	return nil
}

// Returns a string representation of the WorkflowIdentifier
func (wi WorkflowIdentifier) Serialize() (string, error) {
	if err := wi.Valid(); err != nil {
		return "", fmt.Errorf("invalid workflow identifier")
	}
	bucketID, err := convertWeekDayToBucketID(wi.WeekDay)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d-%d-%s", *wi.Type, *wi.Year, *wi.WeekNumber, bucketID), nil
}

// Deserializes a string into a WorkflowIdentifier
func Deserialize(serialized string) (*WorkflowIdentifier, error) {
	parts := strings.Split(serialized, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid workflow split count: %s", serialized)
	}

	var err error

	prefix, yearString, weekNumberString, bucketIDString := parts[0], parts[1], parts[2], parts[3]

	workflowType := WorkflowType(prefix)

	workflowYear, err := strconv.Atoi(yearString)
	if err != nil {
		return nil, fmt.Errorf("invalid year: %s", yearString)
	} else if workflowYear < 2000 || workflowYear > 2100 {
		return nil, fmt.Errorf("invalid year: %d", workflowYear)
	}

	workflowWeek, err := strconv.Atoi(weekNumberString)
	if err != nil {
		return nil, fmt.Errorf("invalid week number: %s", weekNumberString)
	} else if workflowWeek < 1 || workflowWeek > 53 {
		return nil, fmt.Errorf("invalid week number: %d", workflowWeek)
	}

	workflowWeekDay, err := convertBucketIDtoWeekDay(bucketIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid bucket id: %s", bucketIDString)
	}

	wi := &WorkflowIdentifier{
		Type:       &workflowType,
		Year:       &workflowYear,
		WeekNumber: &workflowWeek,
		WeekDay:    &workflowWeekDay,
	}

	return wi, nil
}

// ----------------- Helper Functions -----------------

// Converts a week day to a bucket id
func convertWeekDayToBucketID(weekDay *int) (string, error) {
	if weekDay == nil {
		return "", fmt.Errorf("week day is nil")
	}

	if *weekDay < 0 || *weekDay > 6 {
		return "", fmt.Errorf("week day is invalid")
	}

	if *weekDay <= 3 {
		return constants.FirstRunID, nil
	}

	return constants.LastRunID, nil
}

// Converts a bucket id to a week day
func convertBucketIDtoWeekDay(bucketID string) (int, error) {
	switch bucketID {
	case constants.FirstRunID:
		return 1, nil
	case constants.LastRunID:
		return 7, nil
	default:
		return 0, fmt.Errorf("invalid bucket id: %s", bucketID)
	}
}
