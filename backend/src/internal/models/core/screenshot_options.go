// ./src/internal/models/core/screenshot_options.go
package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/wizenheimer/byrd/src/pkg/utils"
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

type ScreenshotRequestOptions struct {
	URL string `json:"url"`
	// Selector Options
	Selector              *string `json:"selector,omitempty"`
	ScrollIntoView        *string `json:"scroll_into_view,omitempty"`
	AdjustTop             *int    `json:"adjust_top,omitempty"`
	CaptureBeyondViewport *bool   `json:"capture_beyond_viewport,omitempty" default:"true"`

	// Capture Options
	FullPage          *bool              `json:"full_page,omitempty" default:"true"`
	FullPageScroll    *bool              `json:"full_page_scroll,omitempty"`
	FullPageAlgorithm *FullPageAlgorithm `json:"full_page_algorithm,omitempty" default:"default"`
	ScrollDelay       *int               `json:"scroll_delay,omitempty"`
	ScrollBy          *int               `json:"scroll_by,omitempty"`
	MaxHeight         *int               `json:"max_height,omitempty"`
	Format            *string            `json:"format" default:"png"`
	ImageQuality      *int               `json:"image_quality,omitempty" default:"80"`
	OmitBackground    *bool              `json:"omit_background,omitempty"`

	// Clip Options
	Clip *ClipOptions `json:"clip,omitempty"`

	// Resource Blocking Options
	BlockAds                 *bool               `json:"block_ads,omitempty" default:"true"`
	BlockCookieBanners       *bool               `json:"block_cookie_banners,omitempty" default:"true"`
	BlockBannersByHeuristics *bool               `json:"block_banners_by_heuristics,omitempty" default:"true"`
	BlockTrackers            *bool               `json:"block_trackers,omitempty" default:"true"`
	BlockChats               *bool               `json:"block_chats,omitempty" default:"true"`
	BlockRequests            []string            `json:"block_request,omitempty"`
	BlockResources           []BlockResourceType `json:"block_resources,omitempty"`

	// Media Options
	DarkMode      *bool `json:"dark_mode,omitempty" default:"false"`
	ReducedMotion *bool `json:"reduced_motion,omitempty" default:"true"`

	// Request Options
	UserAgent     *string           `json:"user_agent,omitempty"`
	Authorization *string           `json:"authorization,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Cookies       []string          `json:"cookies,omitempty"`
	Timezone      *Timezone         `json:"timezone,omitempty"`
	BypassCSP     *bool             `json:"bypass_csp,omitempty"`
	IpCountryCode *IpCountry        `json:"ip_country_code,omitempty"`

	// Wait and Delay Options
	Delay                    *int                      `json:"delay,omitempty" default:"0"`
	Timeout                  *int                      `json:"timeout,omitempty" default:"60"`
	NavigationTimeout        *int                      `json:"navigation_timeout,omitempty" default:"30"`
	WaitForSelector          *string                   `json:"wait_for_selector,omitempty"`
	WaitForSelectorAlgorithm *WaitForSelectorAlgorithm `json:"wait_for_selector_algorithm,omitempty"`
	WaitUntil                []WaitUntilOption         `json:"wait_until,omitempty" default:"[\"networkidle2\",\"networkidle0\"]"`

	// Interaction Options
	Click               *string           `json:"click,omitempty"`
	FailIfClickNotFound *bool             `json:"fail_if_click_not_found,omitempty"`
	HideSelectors       []string          `json:"hide_selector,omitempty"`
	Styles              *string           `json:"styles,omitempty"`
	Scripts             *string           `json:"scripts,omitempty"`
	ScriptWaitUntil     []WaitUntilOption `json:"scripts_wait_until,omitempty"`

	// Metadata Options
	MetadataImageSize      *bool `json:"metadata_image_size,omitempty" default:"true"`
	MetadataPageTitle      *bool `json:"metadata_page_title,omitempty" default:"true"`
	MetadataContent        *bool `json:"metadata_content,omitempty" default:"true"`
	MetadataHttpStatusCode *bool `json:"metadata_http_response_status_code,omitempty" default:"true"`
	MetadataIcon           *bool `json:"metadata_icon,omitempty"`
}

func GetScreenshotRequestOptions(url string, captureProfile CaptureProfile) ScreenshotRequestOptions {
	options := GetDefaultScreenshotRequestOptions(url)
	if captureProfile.Selector != nil {
		options.Selector = captureProfile.Selector
	}
	if captureProfile.ScrollIntoView != nil {
		options.ScrollIntoView = captureProfile.ScrollIntoView
	}
	if captureProfile.AdjustTop != nil {
		options.AdjustTop = captureProfile.AdjustTop
	}
	if captureProfile.CaptureBeyondViewport != nil {
		options.CaptureBeyondViewport = captureProfile.CaptureBeyondViewport
	}
	if captureProfile.FullPage != nil {
		options.FullPage = captureProfile.FullPage
	}
	if captureProfile.FullPageScroll != nil {
		options.FullPageScroll = captureProfile.FullPageScroll
	}
	if captureProfile.FullPageAlgorithm != nil {
		options.FullPageAlgorithm = captureProfile.FullPageAlgorithm
	}
	if captureProfile.ScrollDelay != nil {
		options.ScrollDelay = captureProfile.ScrollDelay
	}
	if captureProfile.ScrollBy != nil {
		options.ScrollBy = captureProfile.ScrollBy
	}
	if captureProfile.MaxHeight != nil {
		options.MaxHeight = captureProfile.MaxHeight
	}
	if captureProfile.OmitBackground != nil {
		options.OmitBackground = captureProfile.OmitBackground
	}
	if captureProfile.Clip != nil {
		options.Clip = captureProfile.Clip
	}
	if captureProfile.BlockAds != nil {
		options.BlockAds = captureProfile.BlockAds
	}
	if captureProfile.BlockCookieBanners != nil {
		options.BlockCookieBanners = captureProfile.BlockCookieBanners
	}
	if captureProfile.BlockBannersByHeuristics != nil {
		options.BlockBannersByHeuristics = captureProfile.BlockBannersByHeuristics
	}
	if captureProfile.BlockTrackers != nil {
		options.BlockTrackers = captureProfile.BlockTrackers
	}
	if captureProfile.BlockChats != nil {
		options.BlockChats = captureProfile.BlockChats
	}
	if len(captureProfile.BlockRequests) > 0 {
		options.BlockRequests = captureProfile.BlockRequests
	}
	if len(captureProfile.BlockResources) > 0 {
		options.BlockResources = captureProfile.BlockResources
	}
	if captureProfile.DarkMode != nil {
		options.DarkMode = captureProfile.DarkMode
	}
	if captureProfile.ReducedMotion != nil {
		options.ReducedMotion = captureProfile.ReducedMotion
	}
	if captureProfile.UserAgent != nil {
		options.UserAgent = captureProfile.UserAgent
	}
	if captureProfile.Authorization != nil {
		options.Authorization = captureProfile.Authorization
	}
	if len(captureProfile.Headers) > 0 {
		options.Headers = captureProfile.Headers
	}
	if len(captureProfile.Cookies) > 0 {
		options.Cookies = captureProfile.Cookies
	}
	if captureProfile.Timezone != nil {
		options.Timezone = captureProfile.Timezone
	}
	if captureProfile.BypassCSP != nil {
		options.BypassCSP = captureProfile.BypassCSP
	}
	if captureProfile.IpCountryCode != nil {
		options.IpCountryCode = captureProfile.IpCountryCode
	}
	if captureProfile.Delay != nil {
		options.Delay = captureProfile.Delay
	}
	if captureProfile.WaitForSelector != nil {
		options.WaitForSelector = captureProfile.WaitForSelector
	}
	if captureProfile.WaitForSelectorAlgorithm != nil {
		options.WaitForSelectorAlgorithm = captureProfile.WaitForSelectorAlgorithm
	}
	if len(captureProfile.WaitUntil) > 0 {
		options.WaitUntil = captureProfile.WaitUntil
	}

	return options
}

func GetDefaultScreenshotRequestOptions(url string) ScreenshotRequestOptions {
	defaultOpts := ScreenshotRequestOptions{
		URL: url,

		// Default capture options
		Format:                utils.ToPtr("png"),
		ImageQuality:          utils.ToPtr(80),
		CaptureBeyondViewport: utils.ToPtr(true),
		FullPage:              utils.ToPtr(true),
		FullPageAlgorithm:     utils.ToPtr(FullPageAlgorithmDefault),

		// Default resource blocking options
		BlockAds:                 utils.ToPtr(true),
		BlockCookieBanners:       utils.ToPtr(true),
		BlockBannersByHeuristics: utils.ToPtr(true),
		BlockTrackers:            utils.ToPtr(true),
		BlockChats:               utils.ToPtr(true),

		// Default wait and delay options
		Delay:             utils.ToPtr(0),
		Timeout:           utils.ToPtr(60),
		NavigationTimeout: utils.ToPtr(30),
		WaitUntil: []WaitUntilOption{
			WaitUntilNetworkIdle2,
			WaitUntilNetworkIdle0,
		},

		// Default styling options
		DarkMode:      utils.ToPtr(false),
		ReducedMotion: utils.ToPtr(true),

		// Default response options
		MetadataImageSize:      utils.ToPtr(true),
		MetadataPageTitle:      utils.ToPtr(true),
		MetadataContent:        utils.ToPtr(true),
		MetadataHttpStatusCode: utils.ToPtr(true),
	}

	return defaultOpts
}

// Hash generates a deterministic hash of the ScreenshotRequestOptions
func (s *ScreenshotRequestOptions) Hash() string {
	// Create a normalized version of the struct for consistent hashing
	normalized := normalizeOptions(s)

	// Marshal to JSON with sorted keys
	data, _ := json.Marshal(normalized)

	// Generate SHA-256 hash
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// normalizedOptions is a flat structure used for consistent hashing
type normalizedOptions struct {
	URL                      string            `json:"url"`
	Selector                 string            `json:"selector,omitempty"`
	ScrollIntoView           string            `json:"scroll_into_view,omitempty"`
	AdjustTop                int               `json:"adjust_top,omitempty"`
	CaptureBeyondViewport    bool              `json:"capture_beyond_viewport"`
	FullPage                 bool              `json:"full_page"`
	FullPageScroll           bool              `json:"full_page_scroll"`
	FullPageAlgorithm        string            `json:"full_page_algorithm,omitempty"`
	ScrollDelay              int               `json:"scroll_delay,omitempty"`
	ScrollBy                 int               `json:"scroll_by,omitempty"`
	MaxHeight                int               `json:"max_height,omitempty"`
	Format                   string            `json:"format"`
	ImageQuality             int               `json:"image_quality"`
	OmitBackground           bool              `json:"omit_background"`
	Clip                     json.RawMessage   `json:"clip,omitempty"`
	BlockAds                 bool              `json:"block_ads"`
	BlockCookieBanners       bool              `json:"block_cookie_banners"`
	BlockBannersByHeuristics bool              `json:"block_banners_by_heuristics"`
	BlockTrackers            bool              `json:"block_trackers"`
	BlockChats               bool              `json:"block_chats"`
	BlockRequests            []string          `json:"block_requests,omitempty"`
	BlockResources           []string          `json:"block_resources,omitempty"`
	DarkMode                 bool              `json:"dark_mode"`
	ReducedMotion            bool              `json:"reduced_motion"`
	UserAgent                string            `json:"user_agent,omitempty"`
	Authorization            string            `json:"authorization,omitempty"`
	Headers                  map[string]string `json:"headers,omitempty"`
	Cookies                  []string          `json:"cookies,omitempty"`
	Timezone                 string            `json:"timezone,omitempty"`
	BypassCSP                bool              `json:"bypass_csp"`
	IpCountryCode            string            `json:"ip_country_code,omitempty"`
	Delay                    int               `json:"delay"`
	Timeout                  int               `json:"timeout"`
	NavigationTimeout        int               `json:"navigation_timeout"`
	WaitForSelector          string            `json:"wait_for_selector,omitempty"`
	WaitForSelectorAlgorithm string            `json:"wait_for_selector_algorithm,omitempty"`
	WaitUntil                []string          `json:"wait_until,omitempty"`
	Click                    string            `json:"click,omitempty"`
	FailIfClickNotFound      bool              `json:"fail_if_click_not_found"`
	HideSelectors            []string          `json:"hide_selectors,omitempty"`
	Styles                   string            `json:"styles,omitempty"`
	Scripts                  string            `json:"scripts,omitempty"`
	ScriptWaitUntil          []string          `json:"script_wait_until,omitempty"`
	MetadataImageSize        bool              `json:"metadata_image_size"`
	MetadataPageTitle        bool              `json:"metadata_page_title"`
	MetadataContent          bool              `json:"metadata_content"`
	MetadataHttpStatusCode   bool              `json:"metadata_http_status_code"`
	MetadataIcon             bool              `json:"metadata_icon"`
}

func normalizeOptions(s *ScreenshotRequestOptions) normalizedOptions {
	normalized := normalizedOptions{
		URL: s.URL,
		// Handle pointer fields with their default values if nil
		CaptureBeyondViewport:    getPointerValue(s.CaptureBeyondViewport, true),
		FullPage:                 getPointerValue(s.FullPage, true),
		Format:                   getPointerValue(s.Format, "png"),
		ImageQuality:             getPointerValue(s.ImageQuality, 80),
		BlockAds:                 getPointerValue(s.BlockAds, true),
		BlockCookieBanners:       getPointerValue(s.BlockCookieBanners, true),
		BlockBannersByHeuristics: getPointerValue(s.BlockBannersByHeuristics, true),
		BlockTrackers:            getPointerValue(s.BlockTrackers, true),
		BlockChats:               getPointerValue(s.BlockChats, true),
		DarkMode:                 getPointerValue(s.DarkMode, false),
		ReducedMotion:            getPointerValue(s.ReducedMotion, true),
		Delay:                    getPointerValue(s.Delay, 0),
		Timeout:                  getPointerValue(s.Timeout, 60),
		NavigationTimeout:        getPointerValue(s.NavigationTimeout, 30),
	}

	// Handle optional string pointers
	if s.Selector != nil {
		normalized.Selector = *s.Selector
	}
	if s.ScrollIntoView != nil {
		normalized.ScrollIntoView = *s.ScrollIntoView
	}
	if s.UserAgent != nil {
		normalized.UserAgent = *s.UserAgent
	}

	// Handle slices by sorting them for consistency
	if s.BlockRequests != nil {
		normalized.BlockRequests = make([]string, len(s.BlockRequests))
		copy(normalized.BlockRequests, s.BlockRequests)
		sort.Strings(normalized.BlockRequests)
	}

	if s.BlockResources != nil {
		resources := make([]string, len(s.BlockResources))
		for i, r := range s.BlockResources {
			resources[i] = string(r)
		}
		sort.Strings(resources)
		normalized.BlockResources = resources
	}

	// Handle map by sorting keys
	if s.Headers != nil {
		normalized.Headers = make(map[string]string)
		for k, v := range s.Headers {
			normalized.Headers[k] = v
		}
	}

	// Handle complex objects by marshaling them to JSON
	if s.Clip != nil {
		clipData, _ := json.Marshal(s.Clip)
		normalized.Clip = clipData
	}

	return normalized
}

// Helper function to handle pointer values with defaults
func getPointerValue[T comparable](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
