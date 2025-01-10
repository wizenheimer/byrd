// ./src/internal/models/core/workspace_user.go
package models

type UserWorkspaceRole string

const (
	UserRoleAdmin  UserWorkspaceRole = "admin"
	UserRoleUser   UserWorkspaceRole = "user"
	UserRoleViewer UserWorkspaceRole = "viewer"
)

type UserWorkspaceStatus string

const (
	// UserWorkspaceStatusPending is the status for a user that has been invited but has not yet accepted the invitation
	UserWorkspaceStatusPending UserWorkspaceStatus = "pending"

	// UserWorkspaceStatusActive is the status for a user that has accepted the invitation and is active in the workspace
	UserWorkspaceStatusActive UserWorkspaceStatus = "active"

	// UserWorkspaceStatusInactive is the status for a user that has been removed from the workspace
	UserWorkspaceStatusInactive UserWorkspaceStatus = "inactive"
)

// WorkspaceUser is a user in a workspace
type WorkspaceUser struct {
	// Embeds the User struct
	User

	// WorkspaceID is the ID of the workspace
	WorkspaceID string `json:"workspace_id"`

	// Role is the user's role in the workspace
	Role UserWorkspaceRole `json:"role"`

	// Status is the user's status in the workspace
	WorkspaceUserStatus UserWorkspaceStatus `json:"workspace_status"`
}

// WorkspaceUserProps is the properties for creating a workspace user
type WorkspaceUserProps struct {
	// Email is the email address of the user to create
	Email string `json:"email" validate:"required,email"`

	// Role is the role of the user in the workspace
	// If not specified, defaults to "user"
	Role UserWorkspaceRole `json:"role" validate:"required,oneof=admin user viewer" default:"user"`

	// Status is the status of the user in the workspace
	// If not specified, defaults to "pending"
	Status UserWorkspaceStatus `json:"status" validate:"required,oneof=pending active inactive" default:"pending"`
}
