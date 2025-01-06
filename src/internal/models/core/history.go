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
	ID             uuid.UUID              `json:"id"`
	PageID         uuid.UUID              `json:"page_id"`
	WeekNumber1    int                    `json:"week_number_1"`
	WeekNumber2    int                    `json:"week_number_2"`
	YearNumber1    int                    `json:"year_number_1"`
	YearNumber2    int                    `json:"year_number_2"`
	BucketID1      string                 `json:"bucket_id_1"`
	BucketID2      string                 `json:"bucket_id_2"`
	DiffContent    map[string]interface{} `json:"diff_content"`
	ScreenshotURL1 string                 `json:"screenshot_url_1"`
	ScreenshotURL2 string                 `json:"screenshot_url_2"`
	CreatedAt      time.Time              `json:"created_at"`
	Status         HistoryStatus          `json:"history_status" default:"active"`
}
