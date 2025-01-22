// Role types
export type WorkspaceRole = "admin" | "user";

// Membership status types
export type MembershipStatus = "pending" | "active" | "inactive";

// Base workspace user interface
export interface WorkspaceUser {
	user_id: string;
	workspace_id: string;
	name: string;
	email: string;
	workspace_role: WorkspaceRole;
	membership_status: MembershipStatus;
}

// Partial workspace user interface (for updates)
export interface PartialWorkspaceUser {
	user_id: string;
	workspace_role: WorkspaceRole;
	membership_status: MembershipStatus;
}

// Props for creating a new workspace user
export interface WorkspaceUserProps {
	email: string;
	workspace_role: WorkspaceRole;
}
