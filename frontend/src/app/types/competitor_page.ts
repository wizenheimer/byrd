export type PageStatus = "active" | "inactive";

// Screenshot related enums
export type WaitUntilOption =
	| "load"
	| "domcontentloaded"
	| "networkidle0"
	| "networkidle2";

export type WaitForSelectorAlgorithm = "at_least_one" | "at_least_by_count";

export type BlockResourceType =
	| "document"
	| "stylesheet"
	| "image"
	| "media"
	| "font"
	| "script"
	| "texttrack"
	| "xhr"
	| "fetch"
	| "eventsource"
	| "websocket"
	| "manifest"
	| "other";

export type FullPageAlgorithm = "by_sections" | "default";

// Timezone and IP Country types (subset shown for brevity)
export type Timezone =
	| "America/New_York"
	| "America/Los_Angeles"
	| "Europe/London"
	| "Europe/Berlin"
	| "Asia/Shanghai"
	| "Pacific/Auckland";

export type IpCountry =
	| "us"
	| "gb"
	| "de"
	| "fr"
	| "cn"
	| "ca"
	| "es"
	| "jp"
	| "kr"
	| "in"
	| "au";

// Clip options interface
export interface ClipOptions {
	x?: number;
	y?: number;
	width?: number;
	height?: number;
}

// Screenshot request options interface
export interface ScreenshotRequestOptions {
	// Target Options
	url: string;

	// Selector Options
	selector?: string;
	scroll_into_view?: string;
	adjust_top?: number;
	capture_beyond_viewport?: boolean;

	// Capture Options
	full_page?: boolean;
	full_page_scroll?: boolean;
	full_page_algorithm?: FullPageAlgorithm;
	scroll_delay?: number;
	scroll_by?: number;
	max_height?: number;
	format?: string;
	image_quality?: number;
	omit_background?: boolean;

	// Clip Options
	clip?: ClipOptions;

	// Resource Blocking Options
	block_ads?: boolean;
	block_cookie_banners?: boolean;
	block_banners_by_heuristics?: boolean;
	block_trackers?: boolean;
	block_chats?: boolean;
	block_request?: string[];
	block_resources?: BlockResourceType[];

	// Media Options
	dark_mode?: boolean;
	reduced_motion?: boolean;

	// Request Options
	user_agent?: string;
	authorization?: string;
	headers?: Record<string, string>;
	cookies?: string[];
	timezone?: Timezone;
	bypass_csp?: boolean;
	ip_country_code?: IpCountry;

	// Wait and Delay Options
	delay?: number;
	timeout?: number;
	navigation_timeout?: number;
	wait_for_selector?: string;
	wait_for_selector_algorithm?: WaitForSelectorAlgorithm;
	wait_until?: WaitUntilOption[];

	// Interaction Options
	click?: string;
	fail_if_click_not_found?: boolean;
	hide_selector?: string[];
	styles?: string;
	scripts?: string;
	scripts_wait_until?: WaitUntilOption[];

	// Metadata Options
	metadata_image_size?: boolean;
	metadata_page_title?: boolean;
	metadata_content?: boolean;
	metadata_http_response_status_code?: boolean;
	metadata_icon?: boolean;
}

export const DEFAULT_DIFF_PROFILE = [
	"branding",
	"customers",
	"integration",
	"product",
	"pricing",
	"partnerships",
	"messaging",
] as const;

export type DiffProfileType = (typeof DEFAULT_DIFF_PROFILE)[number];

// Page interface
export interface Page {
	id: string;
	competitor_id: string;
	url: string;
	capture_profile: ScreenshotRequestOptions;
	diff_profile: DiffProfileType[];
	last_checked_at?: string;
	status: PageStatus;
	created_at: string;
	updated_at: string;
}

export interface PageProps {
	url: string;
	capture_profile?: ScreenshotRequestOptions;
	diff_profile?: DiffProfileType[];
}
