package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// CreatePageRequest is the request to create a page
type CreatePageRequest = models.PageProps

// UpdatePageRequest is the request to update a page
type UpdatePageRequest = models.PageProps

// GetPageResponse is the response to get a page
type GetPageResponse struct {
	// Page is the page that was retrieved
	Page models.Page `json:"page"`
	// History is the list of history for the page
	History []models.PageHistory `json:"history"`
}
