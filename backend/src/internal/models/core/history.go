// ./src/internal/models/core/history.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// HistoryStatus is the status of a page history
type HistoryStatus string

const (
	// HistoryStatusActive is the status of an active page history
	HistoryStatusActive HistoryStatus = "active"
	// HistoryStatusInactive is the status of an inactive page history
	HistoryStatusInactive HistoryStatus = "inactive"
)

// PageHistory is a page history in the page_history table
type PageHistory struct {
	ID          uuid.UUID              `json:"id"`
	PageID      uuid.UUID              `json:"page_id"`
	DiffContent map[string]interface{} `json:"diff_content"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      HistoryStatus          `json:"history_status" default:"active"`
}
