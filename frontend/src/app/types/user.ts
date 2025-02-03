export type UserStatus = "pending" | "active" | "inactive";

export interface User {
  /** User's unique identifier */
  id: string;
  /** References Clerk's user ID */
  clerk_id?: string;
  /** User's email address */
  email?: string;
  /** User's name */
  name?: string;
  /** User's account status */
  status: UserStatus;
  /** Time the user was created */
  created_at: string;
  /** Time the user was last updated */
  updated_at: string;
}
