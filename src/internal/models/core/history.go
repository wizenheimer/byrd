package models

import (
	"time"

	"github.com/google/uuid"
)

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
}
