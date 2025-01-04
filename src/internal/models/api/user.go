package models

import models "github.com/wizenheimer/iris/src/internal/models/core"

// InviteUserToWorkspaceRequest is the request to invite a user to a workspace
type InviteUserToWorkspaceRequest struct {
	// Email is the email address of the user to invite
	Email string `json:"email" validate:"required,email"`

	// Role is the role of the user in the workspace
	Role models.UserWorkspaceRole `json:"role" validate:"required,oneof=admin user viewer" default:"user"`
}

// UpdateUserRequest is the request to update a user
// This request can be used to update the user's name, email, role, and status
// Users can update their own name and email
type UpdateUserRequest struct {
	Name   string               `json:"name"`
	Email  string               `json:"email" validate:"email"`
	Status models.AccountStatus `json:"status" validate:"required,oneof=pending active inactive" default:"active"`
}

// CreateWorkspaceUserRequest is the request to create a user in a workspace
// This request can be used to create a user in a workspace with a specific role
// Admins can create users in the workspace with any role
type CreateWorkspaceUserRequest struct {
	// Email is the email address of the user to create
	Email string `json:"email" validate:"required,email"`

	// Role is the role of the user in the workspace
	// If not specified, defaults to "user"
	Role models.UserWorkspaceRole `json:"role" validate:"required,oneof=admin user viewer" default:"user"`
}

// CreateWorkspaceUserResponse is the response to creating a user in a workspace
type CreateWorkspaceUserResponse struct {
	// User is the user that was created
	User *models.WorkspaceUser `json:"user"`

	// Error is the error that occurred while creating the user
	Error error `json:"error"`
}
