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

// biome-ignore lint/complexity/noStaticOnlyClass:
export class WorkspaceUsers {
	// list all users in a workspace
	static async list(
		workspaceId: string,
		params: WorkspaceUsersQueryParams,
		token: string,
	): Promise<WorkspaceUsersResponse["data"]> {
		const { data } = await axios.get<WorkspaceUsersResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/users`,
			{
				params,
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	// invite users to a workspace
	static async invite(
		workspaceId: string,
		emails: string[],
		token: string,
	): Promise<AddUsersToWorkspaceResponse["data"]> {
		const request: AddUsersToWorkspaceRequest = { emails };
		const { data } = await axios.post<AddUsersToWorkspaceResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/users`,
			request,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// update member role in a workspace
	static async updateRole(
		workspaceId: string,
		userId: string,
		role: WorkspaceRole,
		token: string,
	): Promise<UpdateWorkspaceUserRoleResponse["data"]> {
		const request: UpdateWorkspaceUserRoleRequest = { role };
		const { data } = await axios.put<UpdateWorkspaceUserRoleResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/users/${userId}`,
			request,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// remove a member from a workspace
	static async remove(
		workspaceId: string,
		userId: string,
		token: string,
	): Promise<RemoveUserResponse["data"]> {
		const { data } = await axios.delete<RemoveUserResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/users/${userId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}
}
