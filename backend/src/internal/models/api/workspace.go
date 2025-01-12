// ./src/internal/models/api/workspace.go
package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// WorkspaceCreationRequest is the request to create a workspace
type WorkspaceCreationRequest struct {
	// Pages is the pages to create the workspace for
	Pages []models.PageProps `json:"competitors" validate:"required,dive"`

	// Users is the users to create the workspace for
	// This excludes the user who is creating the workspace
	Users []models.UserProps `json:"users" validate:"required,dive"`
}

// WorkspaceUpdateRequest is the request to update a workspace
type WorkspaceUpdateRequest = models.WorkspaceProps
