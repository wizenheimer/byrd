import axios from "axios";
import type {
	AddUsersToWorkspaceRequest,
	AddUsersToWorkspaceResponse,
	RemoveUserResponse,
	UpdateWorkspaceUserRoleRequest,
	UpdateWorkspaceUserRoleResponse,
	WorkspaceUsersQueryParams,
	WorkspaceUsersResponse,
} from "../types/api";
import type { WorkspaceRole } from "../types/workspace_user";

// Workspace User Management
// biome-ignore lint/complexity/noStaticOnlyClass:
export class WorkspaceUsers {
	static async list(
		workspaceId: string,
		params: WorkspaceUsersQueryParams,
		token: string,
		origin: string,
	): Promise<WorkspaceUsersResponse["data"]> {
		const { data } = await axios.get<WorkspaceUsersResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/users`,
			{
				params,
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	static async invite(
		workspaceId: string,
		emails: string[],
		token: string,
		origin: string,
	): Promise<AddUsersToWorkspaceResponse["data"]> {
		const request: AddUsersToWorkspaceRequest = { emails };
		const { data } = await axios.post<AddUsersToWorkspaceResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/users`,
			request,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async updateRole(
		workspaceId: string,
		userId: string,
		role: WorkspaceRole,
		token: string,
		origin: string,
	): Promise<UpdateWorkspaceUserRoleResponse["data"]> {
		const request: UpdateWorkspaceUserRoleRequest = { role };
		const { data } = await axios.put<UpdateWorkspaceUserRoleResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/users/${userId}`,
			request,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async remove(
		workspaceId: string,
		userId: string,
		token: string,
		origin: string,
	): Promise<RemoveUserResponse["data"]> {
		const { data } = await axios.delete<RemoveUserResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/users/${userId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}
}
