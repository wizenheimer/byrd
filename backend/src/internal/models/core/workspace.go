// ./src/internal/models/core/workspace.go
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// WorkspaceStatus is the status of a workspace
type WorkspaceStatus string

const (
	// Active is the status of a workspace that is active
	WorkspaceActive WorkspaceStatus = "active"

	// Inactive is the status of a workspace that is inactive
	// This is used to implement soft deletes
	WorkspaceInactive WorkspaceStatus = "inactive"
)

// WorkspacePlan is the plan of a workspace
type WorkspacePlan string

const (
	// Trial is the plan of a workspace that is on a trial
	// Has a limited number of users, pages and competitors
	WorkspaceTrial WorkspacePlan = "trial"

	// Starter is the plan of a workspace that is starter
	// Has a limited number of users, pages and competitors
	WorkspaceStarter WorkspacePlan = "starter"

	// Scaler is the plan of a workspace that is scaler
	// Has a higher limit of users, pages and competitors
	WorkspaceScaler WorkspacePlan = "scaler"

	// Enterprise is the plan of a workspace that is enterprise
	// Has the highest limit of users, pages and competitors
	WorkspaceEnterprise WorkspacePlan = "enterprise"
)

// ToString returns the string representation of a workspace status
func (w WorkspacePlan) ToString() string {
	return string(w)
}

// NewWorkspacePlan returns a new workspace plan
func NewWorkspacePlan(plan string) (WorkspacePlan, error) {
	plan = strings.TrimSpace(strings.ToLower(plan))
	switch plan {
	case "trial":
		return WorkspaceTrial, nil
	case "starter":
		return WorkspaceStarter, nil
	case "scaler":
		return WorkspaceScaler, nil
	case "enterprise":
		return WorkspaceEnterprise, nil
	default:
		return "", fmt.Errorf("invalid workspace plan: %s", plan)
	}
}

// GetMaxCompetitors returns the maximum number of competitors allowed for a workspace plan
func (w Workspace) GetMaxCompetitors() (int, error) {
	if w.WorkspaceStatus != WorkspaceActive {
		return 0, fmt.Errorf("workspace is not active")
	}

	plan := w.WorkspacePlan
	switch plan {
	case WorkspaceTrial: // same as starter but capped
		return 3, nil
	case WorkspaceStarter: // smaller teams (seed stage startups) - upto 5 competitors
		return 3, nil
	case WorkspaceScaler: // medium teams (early stage startups) - upto 10 competitors
		return 10, nil
	case WorkspaceEnterprise: // larger teams (growth stage startups) - upto 20 competitors
		return 20, nil
	default:
		return 0, fmt.Errorf("invalid workspace plan: %s", plan)
	}
}

// GetMaxUsers returns the maximum number of users allowed for a workspace plan
func (w Workspace) GetMaxUsers() (int, error) {
	if w.WorkspaceStatus != WorkspaceActive {
		return 0, fmt.Errorf("workspace is not active")
	}

	switch w.WorkspacePlan {
	case WorkspaceTrial: // same as starter but capped
		return 10, nil
	case WorkspaceStarter: // smaller teams (seed stage startups) - upto 10 users
		return 10, nil
	case WorkspaceScaler: // medium teams (early stage startups) - upto 25 users
		return 25, nil
	case WorkspaceEnterprise: // larger teams (growth stage startups) - upto 50 users
		return 50, nil
	default:
		return 0, fmt.Errorf("invalid workspace plan: %s", w.WorkspacePlan)
	}
}

// GetMaxPages returns the maximum number of pages allowed for a workspace plan
func (w Workspace) GetMaxPages() (int, error) {
	if w.WorkspaceStatus != WorkspaceActive {
		return 0, fmt.Errorf("workspace is not active")
	}

	switch w.WorkspacePlan {
	case WorkspaceTrial: // same as starter but capped
		return 15, nil
	case WorkspaceStarter: // smaller teams (seed stage startups) - upto 15 pages
		return 15, nil
	case WorkspaceScaler: // medium teams (early stage startups) - upto 50 pages
		return 50, nil
	case WorkspaceEnterprise: // larger teams (growth stage startups) - upto 100 pages
		return 100, nil
	default:
		return 0, fmt.Errorf("invalid workspace plan: %s", w.WorkspacePlan)
	}
}

type Workspace struct {
	// ID is the unique identifier of the workspace
	ID uuid.UUID `json:"id"`

	// Name is the name of the workspace, defaults to the name of the user who created the workspace
	Name string `json:"name"`

	// Slug is the unique identifier of the workspace
	Slug string `json:"slug"`

	// BillingEmail is the email address to which billing information is sent, defaults to the email of the user who created the workspace
	BillingEmail string `json:"billing_email"`

	// Status is the status of the workspace
	WorkspaceStatus WorkspaceStatus `json:"workspace_status" validate:"required,oneof=active inactive" default:"pending" omitempty:"true"`

	// Plan is the plan of the workspace
	WorkspacePlan WorkspacePlan `json:"workspace_plan" validate:"required,oneof=trial starter scaler enterprise" default:"trial" omitempty:"true"`

	// CreatedAt is the timestamp when the workspace was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the timestamp when the workspace was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkspaceProps records essential properties of a workspace
type WorkspaceProps struct {
	// Name is the name of the workspace
	Name string `json:"name,omitempty"`

	// BillingEmail is the email address to which billing information is sent
	BillingEmail string `json:"billing_email,omitempty" validate:"omitempty,email"`
}

// ActiveWorkspaceBatch is a batch of active workspaces
type ActiveWorkspaceBatch struct {
	WorkspaceIDs []uuid.UUID `json:"workspaces"`

	// HasMore is true if there are more pages
	HasMore bool `json:"has_more"`

	// LastSeen is the last seen page
	LastSeen *uuid.UUID `json:"last_seen,omitempty"`
}
