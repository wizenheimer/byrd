// ./src/internal/models/core/diff.go
package models

// -----------------------	Weekly Diff Report	------------------------------

// WeeklyReport represents a complete weekly analysis across all URLs
// type WeeklyReport struct {
// 	Data     WeeklyReportData     `json:"data"`
// 	Metadata WeeklyReportMetadata `json:"metadata"`
// }

// WeeklyReportData represents categorized changes for a weekly report
// type WeeklyReportData struct {
// 	Branding    WeeklyCategoryReport `json:"branding"`
// 	Integration WeeklyCategoryReport `json:"integration"`
// 	Pricing     WeeklyCategoryReport `json:"pricing"`
// 	Product     WeeklyCategoryReport `json:"product"`
// 	Positioning WeeklyCategoryReport `json:"positioning"`
// 	Partnership WeeklyCategoryReport `json:"partnership"`
// }

// WeeklyCategoryReport represents a weekly diff report for a single category
// type WeeklyCategoryReport struct {
// 	Changes []string            `json:"changes"`
// 	URLs    map[string][]string `json:"urls"`
// 	Summary string              `json:"summary"`
// }

// WeeklyReportMetadata represents metadata for a weekly report
// type WeeklyReportMetadata struct {
// 	GeneratedAt     string                `json:"generatedAt"`
// 	WeekNumber      string                `json:"weekNumber"`
// 	Competitor      string                `json:"competitor"`
// 	URLCount        int                   `json:"urlCount"`
// 	WeekDayRange    WeekDayRange          `json:"weekDayRange"`
// 	ProcessedURLs   WeeklyProcessedURLs   `json:"processedUrls"`
// 	ProcessingStats WeeklyProcessingStats `json:"processingStats"`
// 	Errors          map[string]string     `json:"errors"`
// 	Enriched        bool                  `json:"enriched"`
// }

// WeeklyProcessedURLs tracks the status of processed URLs
// type WeeklyProcessedURLs struct {
// 	Successful []string `json:"successful"`
// 	Failed     []string `json:"failed"`
// 	Skipped    []string `json:"skipped"`
// }

// WeeklyProcessingStats provides statistics about URL processing
// type WeeklyProcessingStats struct {
// 	TotalURLs    int `json:"totalUrls"`
// 	SuccessCount int `json:"successCount"`
// 	FailureCount int `json:"failureCount"`
// 	SkippedCount int `json:"skippedCount"`
// }

// WeekDayRange represents a time range for comparison
// type WeekDayRange struct {
// 	FromDay string `json:"fromDay"`
// 	ToDay   string `json:"toDay"`
// }
