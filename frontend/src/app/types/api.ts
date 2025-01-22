import type { Competitor, CompetitorWithPages } from "./competitor";
import type {
	DiffProfileType,
	PageProps,
	PageStatus,
	ScreenshotRequestOptions,
} from "./competitor_page";
import type { PageHistory } from "./page_history";
import type { User } from "./user";
import type { Workspace } from "./workspace";
import type { WorkspaceRole, WorkspaceUser } from "./workspace_user";

// Request type for workspace creation
export interface WorkspaceCreationRequest {
	/** List of competitor URLs to track */
	competitors: string[];

	/** List of user emails for the team */
	team: string[];

	/** List of profile strings to track */
	profiles: string[];

	/** List of feature strings to track */
	features: string[];
}

// API response type for workspace creation
export interface WorkspaceCreationResponse {
	message: string;
	data: Workspace;
}

// API response type for listing workspaces
export interface WorkspaceListResponse {
	message: string;
	data: Workspace[];
}

export interface WorkspaceGetResponse {
	message: string;
	data: Workspace;
}

// Update request type matching WorkspaceProps from backend
export interface WorkspaceUpdateRequest {
	/** Name of the workspace */
	name?: string;

	/** Email address for billing information */
	billing_email?: string;
}

// Response data structure for workspace update
export interface WorkspaceUpdateData {
	/** Name of the workspace */
	name: string;

	/** Email address for billing information */
	billingEmail: string;

	/** ID of the updated workspace */
	workspaceId: string;
}

// API response type for workspace update
export interface WorkspaceUpdateResponse {
	message: string;
	data: WorkspaceUpdateData;
}

// Response data structure for joining a workspace
export interface WorkspaceJoinData {
	/** ID of the workspace that was joined */
	workspaceId: string;
}

// API response type for joining a workspace
export interface WorkspaceJoinResponse {
	message: string;
	data: WorkspaceJoinData;
}

// Response data structure for exiting a workspace
export interface WorkspaceExitData {
	/** ID of the workspace that was exited */
	workspaceId: string;
}

// API response type for exiting a workspace
export interface WorkspaceExitResponse {
	message: string;
	data: WorkspaceExitData;
}

// Response data structure for deleting a workspace
export interface WorkspaceDeleteData {
	/** Status of the workspace after deletion */
	status: string;

	/** ID of the workspace that was deleted */
	workspaceId: string;
}

// API response type for deleting a workspace
export interface WorkspaceDeleteResponse {
	message: string;
	data: WorkspaceDeleteData;
}

// API response interfaces
export interface ApiResponse<T> {
	message: string;
	data: T;
}

// Paginated response interface
export interface PaginatedResponse<T> {
	hasMore: boolean;
	users: T[];
}

// Combined type for workspace users list response
export type WorkspaceUsersResponse = ApiResponse<
	PaginatedResponse<WorkspaceUser>
>;

// Query parameters interface for listing users
export interface WorkspaceUsersQueryParams {
	_page?: number;
	_limit?: number;
	role?: WorkspaceRole;
}

// Request types
export interface AddUsersToWorkspaceRequest {
	emails: string[];
}

export interface UpdateWorkspaceUserRoleRequest {
	role: WorkspaceRole;
}

export type AddUsersToWorkspaceResponse = ApiResponse<WorkspaceUser[]>;

// Response data types
interface UpdateRoleResponseData {
	role: WorkspaceRole;
}

export type UpdateWorkspaceUserRoleResponse =
	ApiResponse<UpdateRoleResponseData>;

// Response type for delete operation
interface RemoveUserResponseData {
	userId: string;
	workspaceId: string;
}

export interface RemoveUserResponse {
	message: string;
	data: RemoveUserResponseData;
}

// Competitors List Response interface
export interface CompetitorsListResponse {
	competitors: CompetitorWithPages[];
	hasMore: boolean;
}

// Full API Response type for the competitors endpoint
export type CompetitorsApiResponse = ApiResponse<CompetitorsListResponse>;

// Helper type for pagination parameters
export interface PaginationParams {
	_page?: number;
	_limit?: number;
	includePages?: boolean;
}

export interface CreateCompetitorRequest {
	url: string;
}

export interface CreateCompetitorResponse {
	message: string;
	data: Competitor;
}

export type CreatePageRequest = PageProps;

export interface GetCompetitorResponse {
	message: string;
	data: CompetitorWithPages;
}

export interface UpdateCompetitorResponse {
	message: string;
	data: Competitor;
}

export interface AddPagesResponse {
	message: string;
	data: Array<{
		id: string;
		competitor_id: string;
		url: string;
		capture_profile: ScreenshotRequestOptions;
		diff_profile: DiffProfileType[];
		status: "active" | "inactive";
		created_at: string;
		updated_at: string;
	}>;
}

export interface ListPagesResponse {
	message: string;
	data: {
		pages: Array<{
			id: string;
			competitor_id: string;
			url: string;
			capture_profile: ScreenshotRequestOptions;
			diff_profile: DiffProfileType[];
			status: "active" | "inactive";
			created_at: string;
			updated_at: string;
		}>;
		hasMore: boolean;
	};
}

export interface GetPageResponse {
	message: string;
	data: {
		id: string;
		competitor_id: string;
		url: string;
		capture_profile: ScreenshotRequestOptions;
		diff_profile: DiffProfileType[];
		status: "active" | "inactive";
		created_at: string;
		updated_at: string;
	};
}

export interface ListPageHistoryResponse {
	message: string;
	data: {
		history: PageHistory[];
		hasMore: boolean;
	};
}

export interface UpdatePageResponse {
	message: string;
	data: {
		id: string;
		competitor_id: string;
		url: string;
		capture_profile: ScreenshotRequestOptions;
		diff_profile: DiffProfileType[];
		status: "active" | "inactive";
		created_at: string;
		updated_at: string;
	};
}

export interface GetUserResponse {
	message: string;
	data: User;
}

export interface DeleteUserResponse {
	message: string;
}

export interface DeleteCompetitorResponse {
	message: string;
}

export interface DeletePageResponse {
	message: string;
}

export interface UpdatePageRequest {
	url?: string;
	capture_profile?: Partial<ScreenshotRequestOptions>;
	diff_profile?: DiffProfileType[];
	status?: PageStatus;
}
