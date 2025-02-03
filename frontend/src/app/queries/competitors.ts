import axios from "axios";
import type {
  CompetitorsApiResponse,
  CreateCompetitorRequest,
  CreateCompetitorResponse,
  DeleteCompetitorResponse,
  GetCompetitorResponse,
  PaginationParams,
  UpdateCompetitorResponse,
} from "../types/api";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Competitors {
  static async list(
    workspaceId: string,
    params: PaginationParams & { includePages?: boolean },
    token: string
  ): Promise<CompetitorsApiResponse["data"]> {
    const { _page = 1, _limit = 10, includePages = false } = params;
    const { data } = await axios.get<CompetitorsApiResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors`,
      {
        params: { _page, _limit, includePages },
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  static async get(
    workspaceId: string,
    competitorId: string,
    token: string,
    includePages: boolean = false
  ): Promise<GetCompetitorResponse["data"]> {
    const { data } = await axios.get<GetCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}`,
      {
        params: { includePages },
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  static async create(
    workspaceId: string,
    competitor: CreateCompetitorRequest,
    token: string
  ): Promise<CreateCompetitorResponse["data"]> {
    const { data } = await axios.post<CreateCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors`,
      [competitor],
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async update(
    workspaceId: string,
    competitorId: string,
    name: string,
    token: string
  ): Promise<UpdateCompetitorResponse["data"]> {
    if (!name || name.length === 0 || name.length > 255) {
      throw new Error("Name must be between 1 and 255 characters");
    }

    const { data } = await axios.put<UpdateCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}`,
      { name },
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async delete(
    workspaceId: string,
    competitorId: string,
    token: string
  ): Promise<DeleteCompetitorResponse> {
    const { data } = await axios.delete<DeleteCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data;
  }
}
