// ./src/internal/models/core/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	// AccountStatusPending is the status for an account that has been created but has not yet been activated
	// An account is said to be activated when the user logs in for the first time
	AccountStatusPending AccountStatus = "pending"

	// AccountStatusActive is the status for an account that has been activated
	AccountStatusActive AccountStatus = "active"

	// AccountStatusInactive is the status for an account that has been deactivated or deleted or declined or blocked
	AccountStatusInactive AccountStatus = "inactive"
)

// User is a user in the user table
type User struct {
	// ID is the user's unique identifier
	ID uuid.UUID `json:"id"`

	// Email is the user's email address
	Email *string `json:"email"`

	// Status is the user's account status
	Status AccountStatus `json:"status"`

	// CreatedAt is the time the user was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// GetWorkspaceCreationLimit returns the number of workspaces a user can create
// This is used to limit abuse
func GetWorkspaceCreationLimit() (int, error) {
	return 3, nil
}

type UserProps struct {
	// Email is the email address of the user to create
	Email string `json:"email" validate:"required,email"`

	// Status is the status of the user to create
	Status AccountStatus `json:"status" validate:"required,oneof=pending active inactive"`
}
