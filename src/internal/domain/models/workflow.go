package models

import "fmt"

// WorkflowIdentifier uniquely identifies a workflow
type WorkflowIdentifier struct {
	// Type is the type of the workflow
	Type string
	// Year is the year of the workflow
	Year int
	// WeekNumber is the week number of the workflow
	WeekNumber int
	// BucketNumber is the bucket number of the workflow
	BucketNumber int
}

// formatWorkflowID creates a consistent string representation
func (id *WorkflowIdentifier) Serialize() string {
	return fmt.Sprintf("%s-%d-%d-%d", id.Type, id.Year, id.WeekNumber, id.BucketNumber)
}

// parseWorkflowID parses the string representation into a WorkflowIdentifier
func ParseWorkflowID(serialized string) (*WorkflowIdentifier, error) {
	var id WorkflowIdentifier
	_, err := fmt.Sscanf(serialized, "%s-%d-%d-%d", &id.Type, &id.Year, &id.WeekNumber, &id.BucketNumber)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
