// ./src/internal/models/core/workspace_user.go
package models

import "github.com/google/uuid"

type WorkspaceRole string

const (
	RoleAdmin WorkspaceRole = "admin"

	RoleUser WorkspaceRole = "user"
)

type MembershipStatus string

const (
	// PendingMember is a user who has been invited to join the workspace but has not yet accepted the invitation
	PendingMember MembershipStatus = "pending"

	// ActiveMember is a user who has accepted the invitation to join the workspace
	ActiveMember MembershipStatus = "active"

	// InactiveMember is a user who has been removed from the workspace.
	// This is used to implement soft deletes.
	InactiveMember MembershipStatus = "inactive"
)

// WorkspaceUser is a user and workspace composite
type WorkspaceUser struct {
	// ID is the user's unique identifier
	ID uuid.UUID `json:"user_id" validate:"required"`

	// WorkspaceID is the workspace's unique identifier
	WorkspaceID uuid.UUID `json:"workspace_id" validate:"required"`

	// Name is the user's name
	Name string `json:"name"`

	// Email is the user's email
	Email string `json:"email" validate:"required,email"`

	// Role is the role of the user in the workspace
	Role WorkspaceRole `json:"workspace_role" validate:"required,oneof=admin user viewer" default:"user"`

	// MembershipStatus is the status of the user's membership in the workspace
	MembershipStatus MembershipStatus `json:"membership_status" validate:"required,oneof=pending active inactive" default:"pending"`
}

type PartialWorkspaceUser struct {
	// ID is the user's unique identifier
	ID uuid.UUID `json:"user_id" validate:"required"`

	// Role is the role of the user in the workspace
	Role WorkspaceRole `json:"workspace_role" validate:"required,oneof=admin user viewer" default:"user"`

	// MembershipStatus is the status of the user's membership in the workspace
	MembershipStatus MembershipStatus `json:"membership_status" validate:"required,oneof=pending active inactive" default:"pending"`
}

// WorkspaceUserProps records essential properties of a workspace user
// This is used to identify the user to create in the workspace
type WorkspaceUserProps struct {
	// Email is the email address of the user to create
	Email string `json:"email" validate:"required,email"`

	// Role is the role of the user in the workspace
	// If not specified, defaults to "user"
	Role WorkspaceRole `json:"workspace_role" validate:"required,oneof=admin user viewer" default:"user"`
}
