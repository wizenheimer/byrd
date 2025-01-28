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
	Time         time.Time        `json:"time"`
}

func NewReport(workspaceID, competitorID uuid.UUID, changes []CategoryChange) *Report {
	return &Report{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		CompetitorID: competitorID,
		Changes:      changes,
		Time:         time.Now(),
	}
}
