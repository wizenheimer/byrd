package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// WorkspaceCreationRequest is the request to create a workspace
type WorkspaceCreationRequest struct {
	// Competitors is the list of competitors to track
	// This is a list of competitor URLs
	Competitors []string `json:"competitors" validate:"required"`

	// Profiles is the list of profiles to track for the workspace
	// This is a list of profile strings
	Profiles []string `json:"profiles" validate:"required"`

	// Features is the list of features to track for the workspace
	// This is a list of feature strings
	Features []string `json:"features" validate:"required"`

	// Team is the team to create the workspace for
	// This is a list of user emails
	Team []string `json:"team" validate:"required"`
}

// WorkspaceUpdateRequest is the request to update a workspace
type WorkspaceUpdateRequest = models.WorkspaceProps
