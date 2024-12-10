package models

// ClipOptions defines the coordinates and dimensions for screenshot clipping
type ClipOptions struct {
	X      *int `json:"x,omitempty"`
	Y      *int `json:"y,omitempty"`
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

// WaitUntilOption defines the possible wait conditions for screenshot capture
type WaitUntilOption string

const (
	WaitUntilLoad             WaitUntilOption = "load"
	WaitUntilDOMContentLoaded WaitUntilOption = "domcontentloaded"
	WaitUntilNetworkIdle0     WaitUntilOption = "networkidle0"
	WaitUntilNetworkIdle2     WaitUntilOption = "networkidle2"
)

// WaitForSelectorAlgorithm defines the algorithm for waiting for selectors
type WaitForSelectorAlgorithm string

const (
	WaitForSelectorAtLeastOne     WaitForSelectorAlgorithm = "at_least_one"
	WaitForSelectorAtLeastByCount WaitForSelectorAlgorithm = "at_least_by_count"
)

// BlockResourceType defines the types of resources that can be blocked
type BlockResourceType string

const (
	BlockResourceDocument    BlockResourceType = "document"
	BlockResourceStylesheet  BlockResourceType = "stylesheet"
	BlockResourceImage       BlockResourceType = "image"
	BlockResourceMedia       BlockResourceType = "media"
	BlockResourceFont        BlockResourceType = "font"
	BlockResourceScript      BlockResourceType = "script"
	BlockResourceTextTrack   BlockResourceType = "texttrack"
	BlockResourceXHR         BlockResourceType = "xhr"
	BlockResourceFetch       BlockResourceType = "fetch"
	BlockResourceEventSource BlockResourceType = "eventsource"
	BlockResourceWebSocket   BlockResourceType = "websocket"
	BlockResourceManifest    BlockResourceType = "manifest"
	BlockResourceOther       BlockResourceType = "other"
)

// Timezone defines the supported timezone values
type Timezone string

const (
	TimezoneAmericaBelize       Timezone = "America/Belize"
	TimezoneAmericaCayman       Timezone = "America/Cayman"
	TimezoneAmericaChicago      Timezone = "America/Chicago"
	TimezoneAmericaCostaRica    Timezone = "America/Costa_Rica"
	TimezoneAmericaDenver       Timezone = "America/Denver"
	TimezoneAmericaEdmonton     Timezone = "America/Edmonton"
	TimezoneAmericaElSalvador   Timezone = "America/El_Salvador"
	TimezoneAmericaGuatemala    Timezone = "America/Guatemala"
	TimezoneAmericaGuayaquil    Timezone = "America/Guayaquil"
	TimezoneAmericaHermosillo   Timezone = "America/Hermosillo"
	TimezoneAmericaJamaica      Timezone = "America/Jamaica"
	TimezoneAmericaLosAngeles   Timezone = "America/Los_Angeles"
	TimezoneAmericaMexicoCity   Timezone = "America/Mexico_City"
	TimezoneAmericaNassau       Timezone = "America/Nassau"
	TimezoneAmericaNewYork      Timezone = "America/New_York"
	TimezoneAmericaPanama       Timezone = "America/Panama"
	TimezoneAmericaPortAuPrince Timezone = "America/Port-au-Prince"
	TimezoneAmericaSantiago     Timezone = "America/Santiago"
	TimezoneAmericaTegucigalpa  Timezone = "America/Tegucigalpa"
	TimezoneAmericaTijuana      Timezone = "America/Tijuana"
	TimezoneAmericaToronto      Timezone = "America/Toronto"
	TimezoneAmericaVancouver    Timezone = "America/Vancouver"
	TimezoneAmericaWinnipeg     Timezone = "America/Winnipeg"
	TimezoneAsiaKualaLumpur     Timezone = "Asia/Kuala_Lumpur"
	TimezoneAsiaShanghai        Timezone = "Asia/Shanghai"
	TimezoneAsiaTashkent        Timezone = "Asia/Tashkent"
	TimezoneEuropeBerlin        Timezone = "Europe/Berlin"
	TimezoneEuropeKiev          Timezone = "Europe/Kiev"
	TimezoneEuropeLisbon        Timezone = "Europe/Lisbon"
	TimezoneEuropeLondon        Timezone = "Europe/London"
	TimezoneEuropeMadrid        Timezone = "Europe/Madrid"
	TimezonePacificAuckland     Timezone = "Pacific/Auckland"
	TimezonePacificMajuro       Timezone = "Pacific/Majuro"
)

// IpCountry defines the supported country codes for IP geolocation
type IpCountry string

const (
	IpCountryUS IpCountry = "us"
	IpCountryGB IpCountry = "gb"
	IpCountryDE IpCountry = "de"
	IpCountryIT IpCountry = "it"
	IpCountryFR IpCountry = "fr"
	IpCountryCN IpCountry = "cn"
	IpCountryCA IpCountry = "ca"
	IpCountryES IpCountry = "es"
	IpCountryJP IpCountry = "jp"
	IpCountryKR IpCountry = "kr"
	IpCountryIN IpCountry = "in"
	IpCountryAU IpCountry = "au"
	IpCountryBR IpCountry = "br"
	IpCountryMX IpCountry = "mx"
	IpCountryNZ IpCountry = "nz"
	IpCountryPE IpCountry = "pe"
	IpCountryIS IpCountry = "is"
	IpCountryIE IpCountry = "ie"
)

// FullPageAlgorithm defines the algorithm for full page screenshots
type FullPageAlgorithm string

const (
	FullPageAlgorithmBySections FullPageAlgorithm = "by_sections"
	FullPageAlgorithmDefault    FullPageAlgorithm = "default"
)

// ScreenshotRequestOptions defines all possible options for taking a screenshot
type ScreenshotRequestOptions struct {
	// Target Options
	URL   string `json:"url"`   // The URL of the website to take a screenshot of
	RunID string `json:"runId"` // Run ID for the screenshot

	// Selector Options
	Selector              *string `json:"selector,omitempty"`              // A selector to take screenshot of
	ScrollIntoView        *string `json:"scrollIntoView,omitempty"`        // Selector to scroll into view
	AdjustTop             *int    `json:"adjustTop,omitempty"`             // Once reached the selector, scroll by this amount of pixels
	CaptureBeyondViewport *bool   `json:"captureBeyondViewport,omitempty"` // Handle case where the page or element might not be visible on the viewport

	// Capture Options
	FullPage          *bool              `json:"fullPage,omitempty"`          // Whether to capture the full page
	FullPageScroll    *bool              `json:"fullPageScroll,omitempty"`    // Whether to scroll the page before capturing
	FullPageAlgorithm *FullPageAlgorithm `json:"fullPageAlgorithm,omitempty"` // Algorithm to use for full page capture
	ScrollDelay       *int               `json:"scrollDelay,omitempty"`       // Milliseconds to wait between scrolls
	ScrollBy          *int               `json:"scrollBy,omitempty"`          // Scroll by how many pixels
	MaxHeight         *int               `json:"maxHeight,omitempty"`         // Maximum height of the screenshot
	Format            *string            `json:"format,omitempty"`            // Format of the image (jpg, png, webp)
	ImageQuality      *int               `json:"imageQuality,omitempty"`      // Image quality from 0 to 100
	OmitBackground    *bool              `json:"omitBackground,omitempty"`    // Whether to omit the background

	// Clip Options
	Clip *ClipOptions `json:"clip,omitempty"`

	// Resource Blocking Options
	BlockAds                 *bool               `json:"blockAds,omitempty"`
	BlockCookieBanners       *bool               `json:"blockCookieBanners,omitempty"`
	BlockBannersByHeuristics *bool               `json:"blockBannersByHeuristics,omitempty"`
	BlockTrackers            *bool               `json:"blockTrackers,omitempty"`
	BlockChats               *bool               `json:"blockChats,omitempty"`
	BlockRequests            []string            `json:"blockRequests,omitempty"`
	BlockResources           []BlockResourceType `json:"blockResources,omitempty"`

	// Media Options
	DarkMode      *bool `json:"darkMode,omitempty"`
	ReducedMotion *bool `json:"reducedMotion,omitempty"`

	// Request Options
	UserAgent     *string           `json:"userAgent,omitempty"`
	Authorization *string           `json:"authorization,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Cookies       []string          `json:"cookies,omitempty"`
	Timezone      *Timezone         `json:"timezone,omitempty"`
	BypassCSP     *bool             `json:"bypassCSP,omitempty"`
	IpCountryCode *IpCountry        `json:"ipCountryCode,omitempty"`

	// Wait and Delay Options
	Delay                    *int                      `json:"delay,omitempty"`
	Timeout                  *int                      `json:"timeout,omitempty"`
	NavigationTimeout        *int                      `json:"navigationTimeout,omitempty"`
	WaitForSelector          *string                   `json:"waitForSelector,omitempty"`
	WaitForSelectorAlgorithm *WaitForSelectorAlgorithm `json:"waitForSelectorAlgorithm,omitempty"`
	WaitUntil                []WaitUntilOption         `json:"waitUntil,omitempty"`

	// Interaction Options
	Click               *string           `json:"click,omitempty"`
	FailIfClickNotFound *bool             `json:"failIfClickNotFound,omitempty"`
	HideSelectors       []string          `json:"hideSelectors,omitempty"`
	Styles              *string           `json:"styles,omitempty"`
	Scripts             *string           `json:"scripts,omitempty"`
	ScriptWaitUntil     []WaitUntilOption `json:"scriptWaitUntil,omitempty"`

	// Metadata Options
	MetadataImageSize      *bool `json:"metadataImageSize,omitempty"`
	MetadataPageTitle      *bool `json:"metadataPageTitle,omitempty"`
	MetadataContent        *bool `json:"metadataContent,omitempty"`
	MetadataHttpStatusCode *bool `json:"metadataHttpStatusCode,omitempty"`
	MetadataHttpHeaders    *bool `json:"metadataHttpHeaders,omitempty"`
}

// OutputFormat defines the possible output formats for screenshots
type OutputFormat string

const (
	OutputFormatBase64 OutputFormat = "base64"
	OutputFormatBinary OutputFormat = "binary"
	OutputFormatJSON   OutputFormat = "json"
)

// ScreenshotResponse defines the response structure for screenshot requests
type ScreenshotResponse struct {
	Status      string           `json:"status"`
	Paths       *ScreenshotPaths `json:"paths,omitempty"`
	Metadata    *ScreenshotMeta  `json:"metadata,omitempty"`
	Size        *int             `json:"size,omitempty"`
	URL         *string          `json:"url,omitempty"`
	ContentType *string          `json:"contentType,omitempty"`
	Error       *string          `json:"error,omitempty"`
	Details     *string          `json:"details,omitempty"`
}

// ScreenshotPaths defines the paths for screenshot and content files
type ScreenshotPaths struct {
	Screenshot string `json:"screenshot"`
	Content    string `json:"content"`
}

// ScreenshotMeta defines metadata for screenshots
type ScreenshotMeta struct {
	ImageWidth  int     `json:"imageWidth"`
	ImageHeight int     `json:"imageHeight"`
	PageTitle   *string `json:"pageTitle,omitempty"`
}

// ScreenshotMetadata defines complete metadata for a screenshot
type ScreenshotMetadata struct {
	SourceURL         string  `json:"sourceUrl"`
	FetchedAt         string  `json:"fetchedAt"`
	ScreenshotService string  `json:"screenshotService"`
	Options           string  `json:"options"`
	ImageWidth        int     `json:"imageWidth"`
	ImageHeight       int     `json:"imageHeight"`
	ContentLength     int     `json:"contentLength"`
	PageTitle         *string `json:"pageTitle,omitempty"`
	ContentType       *string `json:"contentType,omitempty"`
}

type ScreenshotServiceConfig struct {
	QPS       float64 `json:"qps"`
	Origin    string  `json:"origin"`
	Key       string  `json:"key"`
	Signature string  `json:"signature"`
}
