package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

// CreateCompetitorRequest is the request to create brand new competitors
// A new competitor is created for each page in the list
type CreateCompetitorRequest = []CreatePageRequest

type UpdateCompetitorRequest struct {
	CompetitorName string `json:"name" validate:"required,min=1,max=255"`
}

type CompetitorResponse struct {
	Competitor *models.Competitor `json:"competitor"`
	Pages      []models.Page      `json:"pages"`
}

func NewCompetitorResponse(competitor *models.Competitor, pages []models.Page) CompetitorResponse {
	return CompetitorResponse{
		Competitor: competitor,
		Pages:      pages,
	}
}
