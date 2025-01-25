export type HistoryStatus = "active" | "inactive";

export interface PageHistory {
  id: string;
  page_id: string;
  diff_content: Record<string, unknown>;
  created_at: string;
  history_status: HistoryStatus;
  prev: string;
  curr: string;
}
