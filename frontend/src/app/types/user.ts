// User status type
export type UserStatus = "active" | "inactive";

// User interface
export interface User {
	id: string;
	clerk_id: string;
	email: string;
	name: string;
	status: UserStatus;
	created_at: string;
	updated_at: string;
}
