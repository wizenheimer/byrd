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

// Workspace Management
// biome-ignore lint/complexity/noStaticOnlyClass:
export class Workspace {
  static async create(
    request: WorkspaceCreationRequest,
    token: string,
    origin: string
  ): Promise<WorkspaceCreationResponse["data"]> {
    const { data } = await axios.post<WorkspaceCreationResponse>(
      `${origin}/api/public/v1/workspace`,
      request,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    console.log(token);
    return data.data;
  }

  static async list(
    token: string,
    origin: string
  ): Promise<WorkspaceListResponse["data"]> {
    const { data } = await axios.get<WorkspaceListResponse>(
      `${origin}/api/public/v1/workspace`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async get(
    id: string,
    token: string,
    origin: string
  ): Promise<WorkspaceGetResponse["data"]> {
    const { data } = await axios.get<WorkspaceGetResponse>(
      `${origin}/api/public/v1/workspace/${id}`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async update(
    workspaceId: string,
    request: WorkspaceUpdateRequest,
    token: string,
    origin: string
  ): Promise<WorkspaceUpdateResponse["data"]> {
    const { data } = await axios.put<WorkspaceUpdateResponse>(
      `${origin}/api/public/v1/workspace/${workspaceId}`,
      request,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async join(
    workspaceId: string,
    token: string,
    origin: string
  ): Promise<WorkspaceJoinResponse["data"]> {
    const { data } = await axios.post<WorkspaceJoinResponse>(
      `${origin}/api/public/v1/workspace/${workspaceId}/join`,
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async exit(
    workspaceId: string,
    token: string,
    origin: string
  ): Promise<WorkspaceExitResponse["data"]> {
    const { data } = await axios.post<WorkspaceExitResponse>(
      `${origin}/api/public/v1/workspace/${workspaceId}/exit`,
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async delete(
    workspaceId: string,
    token: string,
    origin: string
  ): Promise<WorkspaceDeleteResponse["data"]> {
    const { data } = await axios.delete<WorkspaceDeleteResponse>(
      `${origin}/api/public/v1/workspace/${workspaceId}`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }
}
