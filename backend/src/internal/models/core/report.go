package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID           uuid.UUID        `json:"id"`
	WorkspaceID  uuid.UUID        `json:"workspace_id"`
	CompetitorID uuid.UUID        `json:"competitor_id"`
	Changes      []CategoryChange `json:"changes"`
	URI          string           `json:"uri"`
	Time         time.Time        `json:"time"`
}

func NewReport(workspaceID, competitorID uuid.UUID, changes []CategoryChange, uri string) *Report {
	return &Report{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		CompetitorID: competitorID,
		Changes:      changes,
		URI:          uri,
		Time:         time.Now(),
	}
}
