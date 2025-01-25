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

	// ClerkID references Clerk's user ID
	// Could be null if the user has not logged in yet
	ClerkID *string `json:"clerk_id"`

	// Email is the user's email address
	Email *string `json:"email"`

	// Name is the user's name = first name + last name
	Name *string `json:"name"`

	// Status is the user's account status
	Status AccountStatus `json:"status"`

	// CreatedAt is the time the user was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProps struct {
	// Email is the email address of the user to create
	Email string `json:"email" validate:"required,email"`

	// Name is the name of the user to create
	Name string `json:"name" validate:"required"`

	// Status is the status of the user to create
	Status AccountStatus `json:"status" validate:"required,oneof=pending active inactive"`
}
