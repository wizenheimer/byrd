package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceStatus string

const (
	WorkspaceStatusActive   WorkspaceStatus = "active"
	WorkspaceStatusInactive WorkspaceStatus = "inactive"
)

type Workspace struct {
	// ID is the unique identifier of the workspace
	ID uuid.UUID `json:"id"`
	// Name is the name of the workspace, defaults to the name of the user who created the workspace
	Name string `json:"name"`
	// Slug is the unique identifier of the workspace
	Slug string `json:"slug"`
	// BillingEmail is the email address to which billing information is sent, defaults to the email of the user who created the workspace
	BillingEmail string `json:"billing_email"`
	// Status is the status of the workspace
	Status WorkspaceStatus `json:"status"`
	// CreatedAt is the timestamp when the workspace was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the workspace was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkspaceProps records essential properties of a workspace
type WorkspaceProps struct {
	// Name is the name of the workspace
	Name string `json:"name"`

	// BillingEmail is the email address to which billing information is sent
	BillingEmail string `json:"billing_email"`
}
