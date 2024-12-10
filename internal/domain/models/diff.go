package models

import "time"

type DiffAnalysis struct {
	Branding    []string `json:"branding"`
	Integration []string `json:"integration"`
	Pricing     []string `json:"pricing"`
	Product     []string `json:"product"`
	Positioning []string `json:"positioning"`
	Partnership []string `json:"partnership"`
}

type DiffRequest struct {
	URL         string `json:"url" validate:"required,url"`
	RunID1      string `json:"runId1" validate:"required"`
	RunID2      string `json:"runId2" validate:"required"`
	WeekNumber1 string `json:"weekNumber1,omitempty"`
	WeekNumber2 string `json:"weekNumber2,omitempty"`
}

type CategoryBase struct {
	Changes []string            `json:"changes"`
	URLs    map[string][]string `json:"urls"`
}

type CategoryEnriched struct {
	CategoryBase
	Summary string `json:"summary"`
}

type DiffData struct {
	Branding    CategoryBase `json:"branding"`
	Integration CategoryBase `json:"integration"`
	Pricing     CategoryBase `json:"pricing"`
	Product     CategoryBase `json:"product"`
	Positioning CategoryBase `json:"positioning"`
	Partnership CategoryBase `json:"partnership"`
}

type AggregatedReport struct {
	Data     DiffData `json:"data"`
	Metadata struct {
		GeneratedAt string `json:"generatedAt"`
		WeekNumber  string `json:"weekNumber"`
		RunRange    struct {
			FromRun string `json:"fromRun"`
			ToRun   string `json:"toRun"`
		} `json:"runRange"`
		Competitor    string `json:"competitor"`
		URLCount      int    `json:"urlCount"`
		ProcessedURLs struct {
			Successful []string `json:"successful"`
			Failed     []string `json:"failed"`
			Skipped    []string `json:"skipped"`
		} `json:"processedUrls"`
		ProcessingStats struct {
			TotalURLs    int `json:"totalUrls"`
			SuccessCount int `json:"successCount"`
			FailureCount int `json:"failureCount"`
			SkippedCount int `json:"skippedCount"`
		} `json:"processingStats"`
		Errors   map[string]string `json:"errors"`
		Enriched bool              `json:"enriched"`
	} `json:"metadata"`
}

type ReportRequest struct {
	URLs       []string `json:"urls" validate:"required,min=1,max=100,dive,url"`
	RunID1     string   `json:"runId1,omitempty"`
	RunID2     string   `json:"runId2,omitempty"`
	WeekNumber string   `json:"weekNumber,omitempty"`
	Competitor string   `json:"competitor" validate:"required"`
	Enriched   bool     `json:"enriched,omitempty"`
}

// internal/domain/models/diff.go
// Add to existing file

type DiffHistoryParams struct {
	URL        string `json:"url" validate:"required,url"`
	FromRunID  string `json:"fromRunId,omitempty"`
	ToRunID    string `json:"toRunId,omitempty"`
	WeekNumber string `json:"weekNumber,omitempty" validate:"omitempty,len=2,numeric"`
	Limit      int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}

type DiffMetadata struct {
	URL         string    `json:"url"`
	RunID1      string    `json:"runId1"`
	RunID2      string    `json:"runId2"`
	WeekNumber  string    `json:"weekNumber"`
	CreatedAt   time.Time `json:"createdAt"`
	PageTitle   string    `json:"pageTitle,omitempty"`
	LastUpdated string    `json:"lastUpdated,omitempty"`
}

type DiffReport struct {
	URL         string       `json:"url"`
	Timestamp1  string       `json:"timestamp1"`
	Timestamp2  string       `json:"timestamp2"`
	Differences DiffAnalysis `json:"differences"`
	Metadata    struct {
		PageTitle   string `json:"pageTitle,omitempty"`
		LastUpdated string `json:"lastUpdated,omitempty"`
	} `json:"metadata,omitempty"`
}

type DiffHistoryResponse struct {
	Results  []DiffReport `json:"results"`
	Metadata struct {
		URL        string `json:"url"`
		WeekNumber string `json:"weekNumber"`
		DateRange  struct {
			FromRun string `json:"fromRun"`
			ToRun   string `json:"toRun"`
		} `json:"dateRange"`
		Count int `json:"count"`
		Limit int `json:"limit"`
	} `json:"metadata"`
}
