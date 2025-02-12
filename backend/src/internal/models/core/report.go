package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID             uuid.UUID        `json:"id"`
	WorkspaceID    uuid.UUID        `json:"workspace_id"`
	CompetitorID   uuid.UUID        `json:"competitor_id"`
	CompetitorName string           `json:"competitor_name"`
	Changes        []CategoryChange `json:"changes"`
	URI            string           `json:"uri"`
	Time           time.Time        `json:"time"`
}

func NewReport(workspaceID, competitorID uuid.UUID, competitorName string, changes []CategoryChange, uri string) *Report {
	return &Report{
		ID:             uuid.New(),
		WorkspaceID:    workspaceID,
		CompetitorID:   competitorID,
		CompetitorName: competitorName,
		Changes:        changes,
		URI:            uri,
		Time:           time.Now(),
	}
}
