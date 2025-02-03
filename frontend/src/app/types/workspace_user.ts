export type WorkspaceRole = "admin" | "user";
export type MembershipStatus = "pending" | "active" | "inactive";

export interface WorkspaceUser {
	/** User's unique identifier */
	user_id: string;
	/** Workspace's unique identifier */
	workspace_id: string;
	/** User's name */
	name: string;
	/** User's email */
	email: string;
	/** Role of the user in the workspace */
	workspace_role: WorkspaceRole;
	/** Status of the user's membership in the workspace */
	membership_status: MembershipStatus;
}

export interface WorkspaceUserProps {
	/** Email address of the user to create */
	email: string;
	/** Role of the user in the workspace */
	workspace_role: WorkspaceRole;
}
