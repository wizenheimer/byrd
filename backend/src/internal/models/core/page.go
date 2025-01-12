// ./src/internal/models/core/page.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// PageStatus is the status of a page
type PageStatus string

const (
	// When PageStatusActive, the page is active and will be checked for changes
	PageStatusActive PageStatus = "active"

	// When PageStatusInactive, the page is inactive and will not be checked for changes
	PageStatusInactive PageStatus = "inactive"
)

// Page is a page in the page table
type Page struct {
	// ID is the page's unique identifier
	ID uuid.UUID `json:"id" validate:"required"`

	// CompetitorID is the competitor's unique identifier
	CompetitorID uuid.UUID `json:"competitor_id" validate:"required"`

	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	CaptureProfile ScreenshotRequestOptions `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	DiffProfile []string `json:"diff_profile"`

	// LastCheckedAt is the time the page was last checked
	// this is updated after every check
	LastCheckedAt *time.Time `json:"last_checked_at"`

	// Status is the page's status
	Status PageStatus `json:"status"`

	// CreatedAt is the time the page was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time the page was last updated
	// This excludes changes made via last_checked_at
	UpdatedAt time.Time `json:"updated_at"`
}

// ActivePageBatch is a batch of pages
type ActivePageBatch struct {
	// Pages is the list of pages
	Pages []Page `json:"pages"`

	// HasMore is true if there are more pages
	HasMore bool `json:"has_more"`

	// LastSeen is the last seen page
	LastSeen *uuid.UUID `json:"last_seen,omitempty"`
}

// PageProps is struct for essential page properties
type PageProps struct {
	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	// This is optional and defaults to an default capture profile
	CaptureProfile *ScreenshotRequestOptions `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to an default diff profile
	DiffProfile []string `json:"diff_profile"`
}
