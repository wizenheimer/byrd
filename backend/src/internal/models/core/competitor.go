// ./src/internal/models/core/competitor.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type CompetitorStatus string

const (
	CompetitorStatusActive   CompetitorStatus = "active"
	CompetitorStatusInactive CompetitorStatus = "inactive"
)

// Competitor is a competitor in the competitor table
type Competitor struct {
	// ID is the competitor's unique identifier
	ID uuid.UUID `json:"id"`
	// WorkspaceID is the workspace's unique identifier
	WorkspaceID uuid.UUID `json:"workspace_id"`
	// Name is the competitor's name, this is automatically generated from the Page's URL
	Name string `json:"name"`
	// Status is the competitor's status
	Status CompetitorStatus `json:"status"`
	// CreatedAt is the time the competitor was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the competitor was last updated
	UpdatedAt time.Time `json:"updated_at"`
}
