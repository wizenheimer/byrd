package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
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

	// Title is the page's title
	Title string `json:"title"`

	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	CaptureProfile ScreenshotRequestOptions `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	DiffProfile []string `json:"diff_profile" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`

	// LastCheckedAt is the time the page was last checked
	// this is updated after every check
	LastCheckedAt sql.NullTime `json:"last_checked_at,omitempty"`

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
	// PageIDs is the list of pages
	PageIDs []uuid.UUID `json:"pages"`

	// HasMore is true if there are more pages
	HasMore bool `json:"has_more"`

	// LastSeen is the last seen page
	LastSeen *uuid.UUID `json:"last_seen,omitempty"`
}

// PageProps is struct for essential page properties
type PageProps struct {
	// Title is the page's title
	Title string `json:"title"`

	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	// This is optional and defaults to an default capture profile
	CaptureProfile *ScreenshotRequestOptions `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to an default diff profile
	DiffProfile []string `json:"diff_profile" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`
}

// pageJSON is an internal type for JSON marshaling/unmarshaling
type pageJSON struct {
	ID             string                   `json:"id"`
	CompetitorID   string                   `json:"competitor_id"`
	URL            string                   `json:"url"`
	Title          string                   `json:"title"`
	CaptureProfile ScreenshotRequestOptions `json:"capture_profile"`
	DiffProfile    []string                 `json:"diff_profile"`
	LastCheckedAt  *time.Time               `json:"last_checked_at,omitempty"`
	Status         PageStatus               `json:"status"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// MarshalJSON implements custom JSON marshaling for Page
func (p Page) MarshalJSON() ([]byte, error) {
	page := pageJSON{
		ID:             p.ID.String(),
		CompetitorID:   p.CompetitorID.String(),
		URL:            p.URL,
		Title:          p.Title,
		CaptureProfile: p.CaptureProfile,
		DiffProfile:    p.DiffProfile,
		Status:         p.Status,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	if p.LastCheckedAt.Valid {
		page.LastCheckedAt = &p.LastCheckedAt.Time
	}

	return json.Marshal(page)
}

// UnmarshalJSON implements custom JSON unmarshaling for Page
func (p *Page) UnmarshalJSON(data []byte) error {
	var page pageJSON
	if err := json.Unmarshal(data, &page); err != nil {
		return fmt.Errorf("failed to unmarshal page: %w", err)
	}

	// Parse ID UUID
	id, err := uuid.Parse(page.ID)
	if err != nil {
		return fmt.Errorf("invalid page ID: %w", err)
	}
	p.ID = id

	// Parse CompetitorID UUID
	competitorID, err := uuid.Parse(page.CompetitorID)
	if err != nil {
		return fmt.Errorf("invalid competitor ID: %w", err)
	}
	p.CompetitorID = competitorID

	p.URL = page.URL
	p.CaptureProfile = page.CaptureProfile
	p.DiffProfile = page.DiffProfile
	p.Status = page.Status
	p.CreatedAt = page.CreatedAt
	p.UpdatedAt = page.UpdatedAt
	p.Title = page.Title

	// Handle nullable LastCheckedAt
	if page.LastCheckedAt != nil {
		p.LastCheckedAt = sql.NullTime{
			Time:  *page.LastCheckedAt,
			Valid: true,
		}
	} else {
		p.LastCheckedAt = sql.NullTime{Valid: false}
	}

	// Validate required fields and URL format
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// Set default values for DiffProfile if it's empty
func (p *Page) SetDefaultDiffProfile() {
	if len(p.DiffProfile) == 0 {
		p.DiffProfile = []string{
			"branding",
			"customers",
			"integration",
			"product",
			"pricing",
			"partnerships",
			"messaging",
		}
	}
}
