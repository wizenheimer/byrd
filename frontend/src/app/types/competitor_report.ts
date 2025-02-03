export interface CategoryChange {
  /** The category of changes */
  category: string;
  /** Brief summary of the changes */
  summary: string;
  /** List of detailed changes */
  changes: string[];
}

export interface Report {
  id: string;
  workspace_id: string;
  competitor_id: string;
  changes: CategoryChange[];
  time: string;
}
