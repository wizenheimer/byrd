package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// InviteUserToWorkspaceRequest is the request to invite a user to a workspace
type InviteUserToWorkspaceRequest = models.WorkspaceUserProps

// UpdateUserRequest is the request to update a user
// This request can be used to update the user's name, email, role, and status
// Users can update their own name and email
type UpdateUserRequest = models.UserProps

// CreateWorkspaceUserRequest is the request to create a user in a workspace
// This request can be used to create a user in a workspace with a specific role
// Admins can create users in the workspace with any role
type CreateWorkspaceUserRequest = models.WorkspaceUserProps

// CreateWorkspaceUserResponse is the response to creating a user in a workspace
type CreateWorkspaceUserResponse = models.WorkspaceUser

// UpdateWorkspaceUserRoleRequest is the request to update a user's role in a workspace
type UpdateWorkspaceUserRoleRequest struct {
	Role models.UserWorkspaceRole `json:"role" validate:"required,oneof=admin user viewer" default:"user"`
}
