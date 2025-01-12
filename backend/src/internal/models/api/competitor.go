// ./src/internal/models/api/competitor.go
package models

// CreateCompetitorRequest is the request to create brand new competitors
// A new competitor is created for each page in the list
type CreateCompetitorRequest = []CreatePageRequest
