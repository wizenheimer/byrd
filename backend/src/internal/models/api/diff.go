package models

// WeeklyReportRequest represents parameters for generating a weekly report for a list of URLs for a competitor
// type WeeklyReportRequest struct {
// 	// List of URLs to analyze
// 	URLs []string `json:"urls" validate:"required,min=1,max=100,dive,url"`

// 	// Range of days to compare
// 	WeekDay1 string `json:"weekDay1,omitempty"`
// 	WeekDay2 string `json:"weekDay2,omitempty"`

// 	// Week number for the report
// 	WeekNumber string `json:"weekNumber,omitempty"`

// 	// Competitor name
// 	Competitor string `json:"competitor" validate:"required"`

// 	// Whether to enrich the report with AI-generated summaries
// 	Enriched bool `json:"enriched,omitempty"`
// }

// -----------------------	Bi-weekly Diff per URL	------------------------------

// URLDiffRequest represents parameters for a diff comparison
// type URLDiffRequest struct {
// 	// URL to compare
// 	URL string `json:"url" validate:"required,url"`

// 	// From timestamp
// 	WeekDay1    int `json:"weekDay1" validate:"required"`
// 	WeekNumber1 int `json:"weekNumber1,omitempty"`

// 	// To timestamp
// 	WeekDay2    int `json:"weekDay2" validate:"required"`
// 	WeekNumber2 int `json:"weekNumber2,omitempty"`
// }
