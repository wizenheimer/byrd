package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// AddUsersToWorkspaceRequest is the request to invite a user to a workspace
type AddUsersToWorkspaceRequest struct {
	Emails []string `json:"emails" validate:"required"`
}

// UpdateUserRequest is the request to update a user
// This request can be used to update the user's name, email, role, and status
// Users can update their own name and email
type UpdateUserRequest = models.UserProps

// CreateWorkspaceUserRequest is the request to create a user in a workspace
// This request can be used to create a user in a workspace with a specific role
// Admins can create users in the workspace with any role
type CreateWorkspaceUserRequest = models.WorkspaceUserProps

// UpdateWorkspaceUserRoleRequest is the request to update a user's role in a workspace
type UpdateWorkspaceUserRoleRequest struct {
	Role models.WorkspaceRole `json:"role" validate:"required,oneof=admin user" default:"user"`
}
