import type { Competitor, CompetitorWithPages } from "./competitor";
import type {
  CaptureProfile,
  ProfileType,
  Page,
  PageProps,
  PageStatus,
} from "./competitor_page";
import type { PageHistory } from "./page_history";
import type { User } from "./user";
import type { Workspace } from "./workspace";
import type { WorkspaceRole, WorkspaceUser } from "./workspace_user";

// Common Constants
export const DEFAULT_PAGE_NUMBER = 1;
export const DEFAULT_PAGE_SIZE = 10;

// Common Response Types
export interface ErrorResponse {
  /** Error message describing what went wrong */
  error: string;
  /** Optional error details (only included in development mode or with X-Debug header) */
  details?: unknown;
}

export interface DataResponse<T> {
  /** Success or informational message */
  message: string;
  /** Response data */
  data?: T;
}

export interface ApiResponse<T> {
  /** Response message */
  message: string;
  /** Response data */
  data: T;
}

// Pagination Types
export interface PaginationParams {
  // _page and _limit are used by reactquery to handle pagination
  _page?: number;
  _limit?: number;
  includePages?: boolean;
}

export interface PaginatedResponse<T> {
  /** Whether there are more items */
  has_more: boolean;
  /** Total number of items */
  total: number;
  /** Items for the current page */
  items: T[];
}

// API Request/Response Types
export interface WorkspaceCreationRequest {
  competitors: string[];
  profiles: ProfileType[];
  features: string[];
  team: string[];
}

export interface WorkspaceUpdateRequest {
  billing_email?: string;
  name?: string;
}

export interface WorkspaceCreationResponse extends ApiResponse<Workspace> {}
export interface WorkspaceListResponse extends ApiResponse<Workspace[]> {}
export interface WorkspaceGetResponse extends ApiResponse<Workspace> {}
export interface WorkspaceUpdateResponse
  extends ApiResponse<{
    name: string;
    billing_email: string;
    workspace_id: string;
  }> {}

export interface WorkspaceJoinResponse
  extends ApiResponse<{
    workspace_id: string;
  }> {}

export interface WorkspaceExitResponse
  extends ApiResponse<{
    workspace_id: string;
  }> {}

export interface WorkspaceDeleteResponse
  extends ApiResponse<{
    status: string;
    workspace_id: string;
  }> {}

export interface CompetitorsListResponse {
  competitors: CompetitorWithPages[];
  has_more: boolean;
}

export type CompetitorsApiResponse = ApiResponse<CompetitorsListResponse>;

export interface CreateCompetitorRequest {
  url: string;
}

export interface CreateCompetitorResponse extends ApiResponse<Competitor> {}
export interface GetCompetitorResponse
  extends ApiResponse<CompetitorWithPages> {}
export interface UpdateCompetitorResponse extends ApiResponse<Competitor> {}

export type CreatePageRequest = PageProps;

export interface AddPagesResponse
  extends ApiResponse<
    Array<{
      id: string;
      competitor_id: string;
      url: string;
      capture_profile: CaptureProfile;
      diff_profile: ProfileType[];
      status: PageStatus;
      created_at: string;
      updated_at: string;
    }>
  > {}

export interface ListPagesResponse
  extends ApiResponse<{
    pages: Array<{
      id: string;
      competitor_id: string;
      url: string;
      capture_profile: CaptureProfile;
      diff_profile: ProfileType[];
      status: PageStatus;
      created_at: string;
      updated_at: string;
    }>;
    has_more: boolean;
  }> {}

export type GetPageResponse = ApiResponse<Page>;
export type UpdatePageResponse = ApiResponse<Page>;

export interface ListPageHistoryResponse
  extends ApiResponse<{
    history: PageHistory[];
    has_more: boolean;
  }> {}

export interface GetUserResponse extends ApiResponse<User> {}
export interface DeleteUserResponse extends ApiResponse<void> {}
export interface DeleteCompetitorResponse extends ApiResponse<void> {}
export interface DeletePageResponse extends ApiResponse<void> {}

export interface UpdatePageRequest {
  url?: string;
  capture_profile?: Partial<CaptureProfile>;
  diff_profile?: ProfileType[];
  status?: PageStatus;
}

export interface WorkspaceUsersQueryParams extends PaginationParams {
  role?: WorkspaceRole;
}

export interface AddUsersToWorkspaceRequest {
  emails: string[];
}

export interface UpdateWorkspaceUserRoleRequest {
  role: WorkspaceRole;
}
export type WorkspaceUsersResponse = ApiResponse<
  PaginatedResponse<WorkspaceUser>
>;

export type AddUsersToWorkspaceResponse = ApiResponse<WorkspaceUser[]>;
export type UpdateWorkspaceUserRoleResponse = ApiResponse<{
  role: WorkspaceRole;
}>;

export interface RemoveUserResponse
  extends ApiResponse<{
    user_id: string;
    workspace_id: string;
  }> {}

export interface ListPagesQueryParams extends PaginationParams {
  include_pages?: boolean;
}
