// Basic workspace properties
export interface WorkspaceProps {
	/** Name of the workspace */
	name: string;

	/** Email address for billing information */
	billingEmail: string;
}

// Workspace status enum matching backend constants
export enum WorkspaceStatus {
	ACTIVE = "active",
	INACTIVE = "inactive",
}

// Full workspace model matching the backend structure
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

	/** Timestamp of creation */
	created_at: string;

	/** Timestamp of last update */
	updated_at: string;
}
