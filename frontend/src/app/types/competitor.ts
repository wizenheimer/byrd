import type { Page } from "./competitor_page";

// Competitor Status type
export type CompetitorStatus = "active" | "inactive";

// Competitor interface
export interface Competitor {
	id: string;
	workspace_id: string;
	name: string;
	status: CompetitorStatus;
	created_at: string;
	updated_at: string;
}

// CompetitorWithPages interface
export interface CompetitorWithPages {
	competitor: Competitor;
	pages: Page[];
}
