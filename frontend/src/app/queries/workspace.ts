import axios from "axios";
import type {
  WorkspaceCreationRequest,
  WorkspaceCreationResponse,
  WorkspaceDeleteResponse,
  WorkspaceExitResponse,
  WorkspaceGetResponse,
  WorkspaceJoinResponse,
  WorkspaceListResponse,
  WorkspaceUpdateRequest,
  WorkspaceUpdateResponse,
} from "../types/api";
import { MembershipStatus } from "../types/workspace_user";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Workspace {
  // create a new workspace for current user
  static async create(
    request: WorkspaceCreationRequest,
    token: string
  ): Promise<WorkspaceCreationResponse["data"]> {
    const { data } = await axios.post<WorkspaceCreationResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace`,
      request,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // list all workspaces for the current user
  static async list(
    token: string,
    membershipStatus: MembershipStatus = "active"
  ): Promise<WorkspaceListResponse["data"]> {
    const { data } = await axios.get<WorkspaceListResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace`,
      {
        headers: { Authorization: `Bearer ${token}` },
        params: { membership_status: membershipStatus },
      }
    );
    return data.data;
  }

  // get a workspace by id
  static async get(
    id: string,
    token: string
  ): Promise<WorkspaceGetResponse["data"]> {
    const { data } = await axios.get<WorkspaceGetResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${id}`,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // update a workspace by id
  static async update(
    workspaceId: string,
    request: WorkspaceUpdateRequest,
    token: string
  ): Promise<WorkspaceUpdateResponse["data"]> {
    const { data } = await axios.put<WorkspaceUpdateResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}`,
      request,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // join a workspace by id
  static async join(
    workspaceId: string,
    token: string
  ): Promise<WorkspaceJoinResponse["data"]> {
    const { data } = await axios.post<WorkspaceJoinResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/join`,
      {},
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // exit a workspace by id
  static async exit(
    workspaceId: string,
    token: string
  ): Promise<WorkspaceExitResponse["data"]> {
    const { data } = await axios.post<WorkspaceExitResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/exit`,
      {},
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // delete a workspace by id
  static async delete(
    workspaceId: string,
    token: string
  ): Promise<WorkspaceDeleteResponse["data"]> {
    const { data } = await axios.delete<WorkspaceDeleteResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }
}
