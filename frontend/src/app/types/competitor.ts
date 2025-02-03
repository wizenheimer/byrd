import type { Page } from "./competitor_page";

export type CompetitorStatus = "active" | "inactive";

export interface Competitor {
  /** Competitor's unique identifier */
  id: string;
  /** Workspace's unique identifier */
  workspace_id: string;
  /** Competitor's name */
  name: string;
  /** Competitor's status */
  status: CompetitorStatus;
  /** Time the competitor was created */
  created_at: string;
  /** Time the competitor was last updated */
  updated_at: string;
}

export interface CompetitorWithPages {
  competitor: Competitor;
  pages: Page[];
}
