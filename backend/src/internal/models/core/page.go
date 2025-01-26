// ./src/internal/models/core/page.go
package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// PageStatus is the status of a page
type PageStatus string

// DiffProfile is the profile used to diff the page
type DiffProfile []string

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
	CaptureProfile CaptureProfile `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	DiffProfile DiffProfile `json:"diff_profile" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`

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
	CaptureProfile *CaptureProfile `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to an default diff profile
	DiffProfile DiffProfile `json:"diff_profile" validate:"dive,oneof=branding customers integration product pricing partnerships messaging" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`
}

// Capture Profile defines the options for capturing a screenshot
// This is a sub set of the ScreenshotRequestOptions that are exposed to the user
type CaptureProfile struct {
	// Selector Options
	Selector              *string `json:"selector,omitempty"`
	ScrollIntoView        *string `json:"scroll_into_view,omitempty"`
	AdjustTop             *int    `json:"adjust_top,omitempty"`
	CaptureBeyondViewport *bool   `json:"capture_beyond_viewport,omitempty" default:"true"`

	// Capture Options
	FullPage          *bool              `json:"full_page,omitempty" default:"true"`
	FullPageScroll    *bool              `json:"full_page_scroll,omitempty"`
	FullPageAlgorithm *FullPageAlgorithm `json:"full_page_algorithm,omitempty" default:"default"`
	ScrollDelay       *int               `json:"scroll_delay,omitempty"`
	ScrollBy          *int               `json:"scroll_by,omitempty"`
	MaxHeight         *int               `json:"max_height,omitempty"`
	OmitBackground    *bool              `json:"omit_background,omitempty"`

	// Clip Options
	Clip *ClipOptions `json:"clip,omitempty"`

	// Resource Blocking Options
	BlockAds                 *bool               `json:"block_ads,omitempty" default:"true"`
	BlockCookieBanners       *bool               `json:"block_cookie_banners,omitempty" default:"true"`
	BlockBannersByHeuristics *bool               `json:"block_banners_by_heuristics,omitempty" default:"true"`
	BlockTrackers            *bool               `json:"block_trackers,omitempty" default:"true"`
	BlockChats               *bool               `json:"block_chats,omitempty" default:"true"`
	BlockRequests            []string            `json:"block_request,omitempty"`
	BlockResources           []BlockResourceType `json:"block_resources,omitempty"`

	// Media Options
	DarkMode      *bool `json:"dark_mode,omitempty" default:"false"`
	ReducedMotion *bool `json:"reduced_motion,omitempty" default:"true"`

	// Request Options
	UserAgent     *string           `json:"user_agent,omitempty"`
	Authorization *string           `json:"authorization,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Cookies       []string          `json:"cookies,omitempty"`
	Timezone      *Timezone         `json:"timezone,omitempty"`
	BypassCSP     *bool             `json:"bypass_csp,omitempty"`
	IpCountryCode *IpCountry        `json:"ip_country_code,omitempty"`

	// Wait and Delay Options
	Delay                    *int                      `json:"delay,omitempty" default:"0"`
	WaitForSelector          *string                   `json:"wait_for_selector,omitempty"`
	WaitForSelectorAlgorithm *WaitForSelectorAlgorithm `json:"wait_for_selector_algorithm,omitempty"`
	WaitUntil                []WaitUntilOption         `json:"wait_until,omitempty" default:"[\"networkidle2\",\"networkidle0\"]"`
}

func NewPageProps(pageURL string, diffProfile DiffProfile) (PageProps, error) {
	if _, err := url.Parse(pageURL); err != nil {
		return PageProps{}, fmt.Errorf("invalid URL: %w", err)
	}
	cp := GetDefaultCaptureProfile()
	if len(diffProfile) == 0 {
		diffProfile = GetDefaultDiffProfile()
	}
	title, err := utils.GetPageTitle(pageURL)
	if err != nil {
		title = pageURL
	}
	return PageProps{
		Title:          title,
		URL:            pageURL,
		CaptureProfile: &cp,
		DiffProfile:    diffProfile,
	}, nil
}

// pageJSON is an internal type for JSON marshaling/unmarshaling
type pageJSON struct {
	ID             string         `json:"id"`
	CompetitorID   string         `json:"competitor_id"`
	URL            string         `json:"url"`
	Title          string         `json:"title"`
	CaptureProfile CaptureProfile `json:"capture_profile"`
	DiffProfile    DiffProfile    `json:"diff_profile"`
	LastCheckedAt  *time.Time     `json:"last_checked_at,omitempty"`
	Status         PageStatus     `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
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

// GetDefaultDiffProfile returns the default diff profile
func GetDefaultDiffProfile() DiffProfile {
	return []string{
		"branding",
		"customers",
		"integration",
		"product",
		"pricing",
		"partnerships",
		"messaging",
	}
}

// Set default values for DiffProfile if it's empty
func (p *Page) SetDefaultDiffProfile() {
	if len(p.DiffProfile) == 0 {
		p.DiffProfile = GetDefaultDiffProfile()
	}
}

// SetDefaultCaptureProfile sets the default capture profile
func (p *Page) SetDefaultCaptureProfile() {
	p.CaptureProfile = GetDefaultCaptureProfile()
}

// GetDefaultCaptureProfile returns the default capture profile
// This is used when creating a new page
func GetDefaultCaptureProfile() CaptureProfile {
	return CaptureProfile{
		CaptureBeyondViewport: utils.ToPtr(true),
		FullPage:              utils.ToPtr(true),
		FullPageAlgorithm:     utils.ToPtr(FullPageAlgorithmDefault),

		// Resource blocking options
		BlockAds:                 utils.ToPtr(true),
		BlockCookieBanners:       utils.ToPtr(true),
		BlockBannersByHeuristics: utils.ToPtr(true),
		BlockTrackers:            utils.ToPtr(true),
		BlockChats:               utils.ToPtr(true),

		// Wait and delay options
		Delay:             utils.ToPtr(0),
		WaitUntil: []WaitUntilOption{
			WaitUntilNetworkIdle2,
			WaitUntilNetworkIdle0,
		},

		// Styling options
		DarkMode:      utils.ToPtr(false),
		ReducedMotion: utils.ToPtr(true),
	}
}

// MergeOptions merges the provided options with default options
func MergeScreenshotRequestOptions(defaults, override ScreenshotRequestOptions) ScreenshotRequestOptions {
	result := defaults

	// Use reflection to handle all fields
	rOverride := reflect.ValueOf(override)
	rResult := reflect.ValueOf(&result).Elem()

	for i := 0; i < rOverride.NumField(); i++ {
		field := rOverride.Field(i)
		resultField := rResult.Field(i)

		// Skip if the override field is nil or zero
		if field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				resultField.Set(field)
			}
		case reflect.String:
			if field.String() != "" {
				resultField.Set(field)
			}
		case reflect.Slice:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		case reflect.Map:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		}
	}

	return result
}

// MergeOptions merges the provided options with default options
func MergeScreenshotCaptureProfile(defaults, override CaptureProfile) CaptureProfile {
	result := defaults

	// Use reflection to handle all fields
	rOverride := reflect.ValueOf(override)
	rResult := reflect.ValueOf(&result).Elem()

	for i := 0; i < rOverride.NumField(); i++ {
		field := rOverride.Field(i)
		resultField := rResult.Field(i)

		// Skip if the override field is nil or zero
		if field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				resultField.Set(field)
			}
		case reflect.String:
			if field.String() != "" {
				resultField.Set(field)
			}
		case reflect.Slice:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		case reflect.Map:
			if field.Len() > 0 {
				resultField.Set(field)
			}
		}
	}

	return result
}
