// ./src/internal/models/api/page.go
package models

import (
	"net/url"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// CreatePageRequest is the request to create a page
type CreatePageRequest struct {
	// Title is the page's title
	Title string `json:"title,omitempty"`

	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	// This is optional and defaults to an default capture profile
	CaptureProfile *models.CaptureProfile `json:"capture_profile,omitempty"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to an default diff profile
	DiffProfile []string `json:"diff_profile,omitempty" validate:"omitempty,dive,oneof=branding customers integration product pricing partnerships messaging" default:"[\"branding\", \"customers\", \"integration\", \"product\", \"pricing\", \"partnerships\", \"messaging\"]"`
}

// ToProps converts the request to page properties.
// This is used to create a page.
// It ensures that the title is set and the URL is valid.
// If the title is not provided, it fetches the title from the URL.
// If the capture profile is not provided, it uses the default capture profile.
// If the diff profile is not provided, it uses the default diff profile.
func (r *CreatePageRequest) ToProps() (models.PageProps, error) {
	// Get the default capture profile
	defaults := models.GetDefaultCaptureProfile()

	if r.CaptureProfile == nil {
		// Incase the capture profile is not provided, use the default capture profile
		r.CaptureProfile = &defaults
	} else {
		// Incase the capture profile is provided, merge it with the default capture profile
		cp := models.MergeScreenshotCaptureProfile(defaults, *r.CaptureProfile)
		r.CaptureProfile = &cp
	}

	if r.Title == "" {
		var err error
		r.Title, err = utils.GetPageTitle(r.URL)
		if err != nil {
			r.Title = r.URL
		}
	}

	if len(r.DiffProfile) == 0 {
		r.DiffProfile = models.GetDefaultDiffProfile()
	}

	return models.PageProps{
		Title:          r.Title,
		URL:            r.URL,
		CaptureProfile: r.CaptureProfile,
		DiffProfile:    r.DiffProfile,
	}, nil
}

// UpdatePageRequest is the request to update a page
type UpdatePageRequest struct {
	Title          *string                `json:"title,omitempty"`
	URL            *string                `json:"url,omitempty" validate:"omitempty,url"`
	CaptureProfile *models.CaptureProfile `json:"capture_profile,omitempty"`
	DiffProfile    []string               `json:"diff_profile,omitempty" validate:"omitempty,dive,oneof=branding customers integration product pricing partnerships messaging"`
}

// ToProps converts the request to page properties.
// This is used to update a page.
// This ensures that the title is set and the URL is valid if provided.
// If the title is not provided, it fetches the title from the URL.
// If the capture profile is not provided, it does not change the capture profile.
// If the diff profile is not provided, it does not change the diff profile.
func (r *UpdatePageRequest) ToProps() (models.PageProps, error) {
	props := models.PageProps{}

	// If the title is provided, set it
	if r.Title != nil {
		props.Title = *r.Title
	}

	// If the URL is provided, set it
	if r.URL != nil {
		_, err := url.Parse(*r.URL)
		if err != nil {
			return props, err
		}
		props.URL = *r.URL

		// If the title is not provided, get the title from the URL
		if r.Title == nil {
			title, err := utils.GetPageTitle(*r.URL)
			if err != nil {
				// If we can't get the title, just use the URL
				title = *r.URL
			}
			props.Title = title
		}
	}

	// If the capture profile is provided, set it
	if r.CaptureProfile != nil {
		props.CaptureProfile = r.CaptureProfile
	}

	// If the diff profile is provided, set it
	if r.DiffProfile != nil {
		props.DiffProfile = r.DiffProfile
	}

	// Return the page properties
	return props, nil
}
