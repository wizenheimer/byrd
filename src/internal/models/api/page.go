package models

import models "github.com/wizenheimer/iris/src/internal/models/core"

// CreatePageRequest is the request to create a page
type CreatePageRequest struct {
	// URL is the page's URL
	URL string `json:"url" validate:"required,url"`

	// CaptureProfile is the profile used to capture the page
	// This is optional and defaults to an default capture profile
	CaptureProfile map[string]interface{} `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to an default diff profile
	DiffProfile map[string]interface{} `json:"diff_profile"`
}

// UpdatePageRequest is the request to update a page
type UpdatePageRequest struct {
	// URL is the page's URL
	URL string `json:"url" validate:"url"`

	// CaptureProfile is the profile used to capture the page
	// This is optional and defaults to last known capture profile
	CaptureProfile map[string]interface{} `json:"capture_profile"`

	// DiffProfile is the profile used to diff the page
	// This is optional and defaults to last known diff profile
	DiffProfile map[string]interface{} `json:"diff_profile"`
}

// GetPageResponse is the response to get a page
type GetPageResponse struct {
	// Page is the page that was retrieved
	Page models.Page `json:"page"`
	// History is the list of history for the page
	History []models.PageHistory `json:"history"`
}
