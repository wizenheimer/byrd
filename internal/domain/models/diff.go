package models

// -----------------------	Weekly Diff Report	------------------------------

// WeeklyReportRequest represents parameters for generating a weekly report for a list of URLs for a competitor
type WeeklyReportRequest struct {
	// List of URLs to analyze
	URLs []string `json:"urls" validate:"required,min=1,max=100,dive,url"`

	// Range of days to compare
	WeekDay1 string `json:"weekDay1,omitempty"`
	WeekDay2 string `json:"weekDay2,omitempty"`

	// Week number for the report
	WeekNumber string `json:"weekNumber,omitempty"`

	// Competitor name
	Competitor string `json:"competitor" validate:"required"`

	// Whether to enrich the report with AI-generated summaries
	Enriched bool `json:"enriched,omitempty"`
}

// WeeklyReport represents a complete weekly analysis across all URLs
type WeeklyReport struct {
	Data     WeeklyReportData     `json:"data"`
	Metadata WeeklyReportMetadata `json:"metadata"`
}

// WeeklyReportData represents categorized changes for a weekly report
type WeeklyReportData struct {
	Branding    WeeklyCategoryReport `json:"branding"`
	Integration WeeklyCategoryReport `json:"integration"`
	Pricing     WeeklyCategoryReport `json:"pricing"`
	Product     WeeklyCategoryReport `json:"product"`
	Positioning WeeklyCategoryReport `json:"positioning"`
	Partnership WeeklyCategoryReport `json:"partnership"`
}

// WeeklyCategoryReport represents a weekly diff report for a single category
type WeeklyCategoryReport struct {
	Changes []string            `json:"changes"`
	URLs    map[string][]string `json:"urls"`
	Summary string              `json:"summary"`
}

// WeeklyReportMetadata represents metadata for a weekly report
type WeeklyReportMetadata struct {
	GeneratedAt     string                `json:"generatedAt"`
	WeekNumber      string                `json:"weekNumber"`
	Competitor      string                `json:"competitor"`
	URLCount        int                   `json:"urlCount"`
	WeekDayRange    WeekDayRange          `json:"weekDayRange"`
	ProcessedURLs   WeeklyProcessedURLs   `json:"processedUrls"`
	ProcessingStats WeeklyProcessingStats `json:"processingStats"`
	Errors          map[string]string     `json:"errors"`
	Enriched        bool                  `json:"enriched"`
}

// WeeklyProcessedURLs tracks the status of processed URLs
type WeeklyProcessedURLs struct {
	Successful []string `json:"successful"`
	Failed     []string `json:"failed"`
	Skipped    []string `json:"skipped"`
}

// WeeklyProcessingStats provides statistics about URL processing
type WeeklyProcessingStats struct {
	TotalURLs    int `json:"totalUrls"`
	SuccessCount int `json:"successCount"`
	FailureCount int `json:"failureCount"`
	SkippedCount int `json:"skippedCount"`
}

// WeekDayRange represents a time range for comparison
type WeekDayRange struct {
	FromDay string `json:"fromDay"`
	ToDay   string `json:"toDay"`
}

// -----------------------	Bi-weekly Diff per URL	------------------------------

// URLDiffRequest represents parameters for a diff comparison
type URLDiffRequest struct {
	// URL to compare
	URL string `json:"url" validate:"required,url"`

	// From timestamp
	WeekDay1    string `json:"weekDay1" validate:"required"`
	WeekNumber1 string `json:"weekNumber1,omitempty"`

	// To timestamp
	WeekDay2    string `json:"weekDay2" validate:"required"`
	WeekNumber2 string `json:"weekNumber2,omitempty"`
}

// URLDiffAnalysis represents differences found across all categories
type URLDiffAnalysis struct {
	Branding    []string `json:"branding"`
	Integration []string `json:"integration"`
	Pricing     []string `json:"pricing"`
	Product     []string `json:"product"`
	Positioning []string `json:"positioning"`
	Partnership []string `json:"partnership"`
}
