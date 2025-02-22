// ./src/internal/models/api/workspace.go
package models

import (
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// WorkspaceCreationRequest is the request to create a workspace
type WorkspaceCreationRequest struct {
	// Competitors is the list of competitors to track
	// This is a list of competitor URLs
	Competitors []string `json:"competitors" validate:"required,dive,url"`

	// Profiles is the list of profiles to track for the workspace
	// This is a list of profile strings
	Profiles []string `json:"profiles" validate:"dive,oneof=branding customers integration product pricing partnerships messaging" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`

	// Features is the list of features to track for the workspace
	// This is a list of feature strings
	Features []string `json:"features" validate:"required"`

	// Team is the team to create the workspace for
	// This is a list of user emails
	Team []string `json:"team" validate:"required,dive,email"` // TODO: remove this
}

// WorkspaceUpdateRequest is the request to update a workspace
// type WorkspaceUpdateRequest = models.WorkspaceProps
type WorkspaceUpdateRequest struct {
	BillingEmail *string `json:"billing_email,omitempty" validate:"omitempty,email"`
	Name         *string `json:"name,omitempty"`
}

// ToProps converts the request to workspace properties.
// This is used to update a workspace.
// It ensures that the billing email is normalized.
// If the name is not provided, it is not set.
func (r *WorkspaceUpdateRequest) ToProps() models.WorkspaceProps {
	props := models.WorkspaceProps{}
	if r.BillingEmail != nil {
		email := *r.BillingEmail
		props.BillingEmail = utils.NormalizeEmail(email)
	}
	if r.Name != nil {
		props.Name = *r.Name
	}
	return props
}
