export type WorkspaceStatus = "active" | "inactive";
export type WorkspacePlan = "trial" | "starter" | "scaler" | "enterprise";

export interface Workspace {
  /** Unique identifier of the workspace */
  id: string;
  /** Name of the workspace */
  name: string;
  /** Unique slug identifier */
  slug: string;
  /** Email address for billing information */
  billing_email: string;
  /** Current status of the workspace */
  workspace_status: WorkspaceStatus;
  /** Current plan of the workspace */
  workspace_plan: WorkspacePlan;
  /** Timestamp of creation */
  created_at: string;
  /** Timestamp of last update */
  updated_at: string;
}

export interface WorkspaceProps {
  /** Name of the workspace */
  name?: string;
  /** Email address for billing information */
  billing_email?: string;
}
