package models

import "time"

type ScreenshotOptions struct {
	URL                      string   `json:"url"`
	RunID                    string   `json:"runId"`
	Format                   string   `json:"format,omitempty"`
	ImageQuality             int      `json:"imageQuality,omitempty"`
	CaptureBeyondViewport    bool     `json:"captureBeyondViewport,omitempty"`
	FullPage                 bool     `json:"fullPage,omitempty"`
	BlockAds                 bool     `json:"blockAds,omitempty"`
	BlockCookieBanners       bool     `json:"blockCookieBanners,omitempty"`
	BlockBannersByHeuristics bool     `json:"blockBannersByHeuristics,omitempty"`
	BlockTrackers            bool     `json:"blockTrackers,omitempty"`
	BlockChats               bool     `json:"blockChats,omitempty"`
	WaitUntil                []string `json:"waitUntil,omitempty"`
}

type ScreenshotResponse struct {
	Status   string           `json:"status"`
	Paths    *ScreenshotPaths `json:"paths,omitempty"`
	Metadata *ScreenshotMeta  `json:"metadata,omitempty"`
	Size     int64            `json:"size,omitempty"`
	Error    string           `json:"error,omitempty"`
}

type ScreenshotPaths struct {
	Screenshot string `json:"screenshot"`
	Content    string `json:"content"`
}

type ScreenshotMeta struct {
	ImageWidth  int    `json:"imageWidth"`
	ImageHeight int    `json:"imageHeight"`
	PageTitle   string `json:"pageTitle,omitempty"`
}

type ScreenshotMetadata struct {
	SourceURL         string    `json:"sourceUrl"`
	FetchedAt         time.Time `json:"fetchedAt"`
	ScreenshotService string    `json:"screenshotService"`
	ImageWidth        int       `json:"imageWidth"`
	ImageHeight       int       `json:"imageHeight"`
	ContentLength     int       `json:"contentLength"`
	PageTitle         string    `json:"pageTitle,omitempty"`
	ContentType       string    `json:"contentType,omitempty"`
}

type ScreenshotServiceConfig struct {
	APIKey     string            `json:"apiKey" validate:"required"`
	BaseURL    string            `json:"baseURL" validate:"required,url"`
	Timeout    int               `json:"timeout" validate:"required,min=1"`
	MaxRetries int               `json:"maxRetries" validate:"required,min=0"`
	RetryDelay int               `json:"retryDelay" validate:"required,min=0"`
	Headers    map[string]string `json:"headers,omitempty"`
	Options    struct {
		DefaultFormat     string   `json:"defaultFormat" validate:"oneof=jpg png webp"`
		DefaultQuality    int      `json:"defaultQuality" validate:"min=1,max=100"`
		DefaultTimeout    int      `json:"defaultTimeout" validate:"min=1"`
		MaxScreenshotSize int64    `json:"maxScreenshotSize" validate:"min=1"`
		AllowedDomains    []string `json:"allowedDomains,omitempty"`
	} `json:"options"`
}
