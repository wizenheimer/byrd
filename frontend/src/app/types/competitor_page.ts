export type PageStatus = "active" | "inactive";

export type WaitUntilOption =
  | "load"
  | "domcontentloaded"
  | "networkidle0"
  | "networkidle2";

export type WaitForSelectorAlgorithm = "at_least_one" | "at_least_by_count";
export type FullPageAlgorithm = "by_sections" | "default";

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

export interface ClipOptions {
  x?: number;
  y?: number;
  width?: number;
  height?: number;
}

export interface CaptureProfile {
  /** Selector Options */
  selector?: string;
  scroll_into_view?: string;
  adjust_top?: number;
  capture_beyond_viewport?: boolean;

  /** Capture Options */
  full_page?: boolean;
  full_page_scroll?: boolean;
  full_page_algorithm?: FullPageAlgorithm;
  scroll_delay?: number;
  scroll_by?: number;
  max_height?: number;
  omit_background?: boolean;

  /** Clip Options */
  clip?: ClipOptions;

  /** Resource Blocking Options */
  block_ads?: boolean;
  block_cookie_banners?: boolean;
  block_banners_by_heuristics?: boolean;
  block_trackers?: boolean;
  block_chats?: boolean;
  block_requests?: string[];
  block_resources?: BlockResourceType[];

  /** Media Options */
  dark_mode?: boolean;
  reduced_motion?: boolean;

  /** Request Options */
  user_agent?: string;
  authorization?: string;
  headers?: Record<string, string>;
  cookies?: string[];
  timezone?: string;
  bypass_csp?: boolean;
  ip_country_code?: string;

  /** Wait and Delay Options */
  delay?: number;
  wait_for_selector?: string;
  wait_for_selector_algorithm?: WaitForSelectorAlgorithm;
  wait_until?: WaitUntilOption[];
}

export const DEFAULT_PROFILE = [
  "branding",
  "customers",
  "integration",
  "product",
  "pricing",
  "partnerships",
  "messaging",
] as const;

export type ProfileType = (typeof DEFAULT_PROFILE)[number];

export interface Page {
  /** Page's unique identifier */
  id: string;
  /** Competitor's unique identifier */
  competitor_id: string;
  /** Page's title */
  title: string;
  /** Page's URL */
  url: string;
  /** Profile used to capture the page */
  capture_profile: CaptureProfile;
  /** Profile used to diff the page */
  diff_profile: ProfileType[];
  /** Time the page was last checked */
  last_checked_at?: string;
  /** Page's status */
  status: PageStatus;
  /** Time the page was created */
  created_at: string;
  /** Time the page was last updated */
  updated_at: string;
}

export interface PageProps {
  url: string;
  capture_profile?: CaptureProfile;
  diff_profile?: ProfileType[];
}
