package models

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"time"
)

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
	URL string `json:"url"` // The URL of the website to take a screenshot of

	// Selector Options
	Selector              *string `json:"selector,omitempty"`                // A selector to take screenshot of
	ScrollIntoView        *string `json:"scroll_into_view,omitempty"`        // Selector to scroll into view
	AdjustTop             *int    `json:"adjust_top,omitempty"`              // Once reached the selector, scroll by this amount of pixels
	CaptureBeyondViewport *bool   `json:"capture_beyond_viewport,omitempty"` // Handle case where the page or element might not be visible on the viewport

	// Capture Options
	FullPage          *bool              `json:"full_page,omitempty"`           // Whether to capture the full page
	FullPageScroll    *bool              `json:"full_page_scroll,omitempty"`    // Whether to scroll the page before capturing
	FullPageAlgorithm *FullPageAlgorithm `json:"full_page_algorithm,omitempty"` // Algorithm to use for full page capture
	ScrollDelay       *int               `json:"scroll_delay,omitempty"`        // Milliseconds to wait between scrolls
	ScrollBy          *int               `json:"scroll_by,omitempty"`           // Scroll by how many pixels
	MaxHeight         *int               `json:"max_height,omitempty"`          // Maximum height of the screenshot
	Format            *string            `json:"format,omitempty"`              // Format of the image (jpg, png, webp)
	ImageQuality      *int               `json:"image_quality,omitempty"`       // Image quality from 0 to 100
	OmitBackground    *bool              `json:"omit_background,omitempty"`     // Whether to omit the background

	// Clip Options
	Clip *ClipOptions `json:"clip,omitempty"`

	// Resource Blocking Options
	BlockAds                 *bool               `json:"block_ads,omitempty"`
	BlockCookieBanners       *bool               `json:"block_cookie_banners,omitempty"`
	BlockBannersByHeuristics *bool               `json:"block_banners_by_heuristics,omitempty"`
	BlockTrackers            *bool               `json:"block_trackers,omitempty"`
	BlockChats               *bool               `json:"block_chats,omitempty"`
	BlockRequests            []string            `json:"block_request,omitempty"` // Changed from blockRequests
	BlockResources           []BlockResourceType `json:"block_resources,omitempty"`

	// Media Options
	DarkMode      *bool `json:"dark_mode,omitempty"`
	ReducedMotion *bool `json:"reduced_motion,omitempty"`

	// Request Options
	UserAgent     *string           `json:"user_agent,omitempty"`
	Authorization *string           `json:"authorization,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Cookies       []string          `json:"cookies,omitempty"`
	Timezone      *Timezone         `json:"timezone,omitempty"`
	BypassCSP     *bool             `json:"bypass_csp,omitempty"`
	IpCountryCode *IpCountry        `json:"ip_country_code,omitempty"`

	// Wait and Delay Options
	Delay                    *int                      `json:"delay,omitempty"`
	Timeout                  *int                      `json:"timeout,omitempty"`
	NavigationTimeout        *int                      `json:"navigation_timeout,omitempty"`
	WaitForSelector          *string                   `json:"wait_for_selector,omitempty"`
	WaitForSelectorAlgorithm *WaitForSelectorAlgorithm `json:"wait_for_selector_algorithm,omitempty"`
	WaitUntil                []WaitUntilOption         `json:"wait_until,omitempty"`

	// Interaction Options
	Click               *string           `json:"click,omitempty"`
	FailIfClickNotFound *bool             `json:"fail_if_click_not_found,omitempty"`
	HideSelectors       []string          `json:"hide_selector,omitempty"` // Changed from hideSelectors
	Styles              *string           `json:"styles,omitempty"`
	Scripts             *string           `json:"scripts,omitempty"`
	ScriptWaitUntil     []WaitUntilOption `json:"scripts_wait_until,omitempty"`

	// Metadata Options
	MetadataImageSize      *bool `json:"metadata_image_size,omitempty"`
	MetadataPageTitle      *bool `json:"metadata_page_title,omitempty"`
	MetadataContent        *bool `json:"metadata_content,omitempty"`
	MetadataHttpStatusCode *bool `json:"metadata_http_response_status_code,omitempty"` // Changed from metadataHttpStatusCode
	MetadataIcon           *bool `json:"metadata_icon,omitempty"`
}

// ScreenshotHTMLRequestOptions defines all possible options for capturing html
type ScreenshotHTMLRequestOptions struct {
	// Target Options
	URL string `json:"url"` // The URL of the website to take a screenshot of
}

type GetScreenshotOptions struct {
	// Target Options
	URL        string `json:"url"`                   // The URL of the website to take a screenshot of
	WeekNumber *int   `json:"week_number,omitempty"` // The week number of the screenshot
	WeekDay    *int   `json:"week_day,omitempty"`    // The day of the week of the screenshot
	Year       *int   `json:"year,omitempty"`        // The year of the screenshot
}

type ListScreenshotsOptions struct {
	// Target Options
	URL         string `json:"url"`          // The URL of the website to take a screenshot of
	ContentType string `json:"content_type"` // The type of content to list
	MaxItems    int    `json:"max_items"`    // The maximum number of items to list
}

// ScreenshotHTMLContentResponse defines the response structure for screenshot content requests
type ScreenshotHTMLContentResponse struct {
	Status      string              `json:"status"`
	HTMLContent string              `json:"content"`
	Metadata    *ScreenshotMetadata `json:"metadata,omitempty"`
}

// ScreenshotImageResponse defines the response structure for screenshot image requests
type ScreenshotImageResponse struct {
	Status      string              `json:"status"`
	Image       image.Image         `json:"image"`
	Metadata    *ScreenshotMetadata `json:"metadata,omitempty"`
	ImageHeight *int                `json:"image_height"`
	ImageWidth  *int                `json:"image_width"`
}

// ScreenshotListResponse defines the response structure for listing screenshots
type ScreenshotListResponse struct {
	Key          string
	LastModified time.Time
}

// ScreenshotPaths defines the paths for screenshot and content files
type ScreenshotPaths struct {
	Screenshot string `json:"screenshot"`
	Content    string `json:"content"`
}

// ScreenshotMeta defines metadata for screenshots
type ScreenshotMeta struct {
	ImageWidth  int     `json:"image_width"`
	ImageHeight int     `json:"image_height"`
	PageTitle   *string `json:"page_title"`
	ContentSize *int    `json:"content_size"`
}

// ScreenshotMetadata defines complete metadata for a screenshot
type ScreenshotMetadata struct {
	SourceURL   string `json:"source_url"`
	RenderedURL string `json:"rendered_url"`
	Year        int    `json:"year"`
	WeekNumber  int    `json:"week_number"`
	WeekDay     int    `json:"week_day"`
}

type ScreenshotServiceConfig struct {
	QPS       float64 `json:"qps"`
	Origin    string  `json:"origin"`
	Key       string  `json:"key"`
	Signature string  `json:"signature"`
}

func (s ScreenshotMetadata) ToMap() map[string]string {
	result := make(map[string]string)

	result["source_url"] = s.SourceURL
	result["rendered_url"] = s.RenderedURL
	result["year"] = strconv.Itoa(s.Year)
	result["week_day"] = strconv.Itoa(s.WeekDay)
	result["week_number"] = strconv.Itoa(s.WeekNumber)

	return result
}

// FromMap safely converts map[string]string to ScreenshotMetadata
func ScreenshotMetadataFromMap(m map[string]string) (ScreenshotMetadata, []error) {
	var result ScreenshotMetadata
	var errs []error

	// Required string fields
	if srcURL, exists := m["source_url"]; exists {
		result.SourceURL = srcURL
	} else {
		errs = append(errs, errors.New("missing required field: source_url"))
	}

	if rendURL, exists := m["rendered_url"]; exists {
		result.RenderedURL = rendURL
	} else {
		errs = append(errs, errors.New("missing required field: rendered_url"))
	}

	// Required integer fields
	if year, exists := m["year"]; exists {
		if y, err := strconv.Atoi(year); err == nil {
			result.Year = y
		} else {
			errs = append(errs, fmt.Errorf("invalid year: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: year"))
	}

	if weekday, exists := m["week_day"]; exists {
		if wd, err := strconv.Atoi(weekday); err == nil {
			result.WeekDay = wd
		} else {
			errs = append(errs, fmt.Errorf("invalid week_day: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: week_day"))
	}

	if weeknumber, exists := m["week_number"]; exists {
		if wn, err := strconv.Atoi(weeknumber); err == nil {
			result.WeekNumber = wn
		} else {
			errs = append(errs, fmt.Errorf("invalid week_number: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: week_number"))
	}

	// Return errs if any occurred
	if len(errs) > 0 {
		return result, errs
	}

	return result, nil
}
