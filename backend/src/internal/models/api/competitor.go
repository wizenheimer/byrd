package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// CreateCompetitorRequest is the request to create brand new competitors
// A new competitor is created for each page in the list
type CreateCompetitorRequest struct {
	// PageURLs is the list of page URLs
	Pages []CreatePageRequest `json:"pages" validate:"required,dive"`
}

// CreateCompetitorResponse is the response to create a competitor
// This response contains the list of competitors that were created and the errors that occurred
type CreateWorkspaceCompetitorResponse struct {
	// Competitors is the list of competitors that were created
	Competitors []models.Competitor `json:"competitors"`
}

// GetWorkspaceCompetitorResponse is the response to get a competitor
// This response contains the competitor that was retrieved
// And pages that belong to the competitor
type GetWorkspaceCompetitorResponse struct {
	// Competitor is the competitor that was retrieved
	Competitor models.Competitor `json:"competitor"`
	// Pages is the list of pages that belong to the competitor
	Pages []models.Page `json:"pages"`
}

// UpdateCompetitorRequest is the request to update a competitor
type UpdateCompetitorRequest struct {
	// Name is the competitor's name
	Name string `json:"name" validate:"required"`
	// Status is the competitor's status
	Status models.CompetitorStatus `json:"status" validate:"required"`
}
