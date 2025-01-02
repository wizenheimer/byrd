package models

import models "github.com/wizenheimer/iris/src/internal/models/core"

// WorkspaceCreationRequest is the request to create a workspace
type WorkspaceCreationRequest struct {
	// Competitor is the competitors to create the workspace for
	CompetitorCreationRequest CreateCompetitorRequest `json:"competitors"`

	// Users is the users to create the workspace for
	// This excludes the user who is creating the workspace
	WorkspaceUserCreationRequest []CreateWorkspaceUserRequest `json:"users"`
}

// WorkspaceCreationResponse is the response to create a workspace
type WorkspaceCreationResponse struct {
	// Workspace is the workspace that was created
	Workspace models.Workspace `json:"workspace"`

	// Users is the list of users that are part of the workspace
	Users []models.WorkspaceUser `json:"users"`
}

// WorkspaceUpdateRequest is the request to update a workspace
type WorkspaceUpdateRequest struct {
	// Name is the name of the workspace
	Name string `json:"name"`

	// BillingEmail is the email address to which billing information is sent
	BillingEmail string `json:"billing_email"`
}

// WorkspaceMembersListingParams is the parameters for listing workspace members
type WorkspaceMembersListingParams struct {
	// PaginationParams is the pagination parameters
	PaginationParams

	// Filtering parameters for listing users
	IncludeAdmins  bool `query:"include_admins"`
	IncludeMembers bool `query:"include_members"`
}
