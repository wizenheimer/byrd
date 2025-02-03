export type HistoryStatus = "active" | "inactive";

export interface PageHistory {
	id: string;
	page_id: string;
	diff_content: Record<string, unknown>;
	history_status: HistoryStatus;
	prev: string;
	curr: string;
	created_at: string;
}
