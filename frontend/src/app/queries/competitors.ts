import axios from "axios";
import type {
  CompetitorsApiResponse,
  CreateCompetitorRequest,
  CreateCompetitorResponse,
  DeleteCompetitorResponse,
  GetCompetitorResponse,
  PaginationParamsWithPageOptions,
  UpdateCompetitorResponse,
} from "../types/api";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Competitors {
  // list all competitors in a workspace
  static async list(
    workspaceId: string,
    params: PaginationParamsWithPageOptions,
    token: string
  ): Promise<CompetitorsApiResponse["data"]> {
    const { _page = 1, _limit = 10, include_pages = false } = params;
    const { data } = await axios.get<CompetitorsApiResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors`,
      {
        params: { _page, _limit, include_pages },
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // get a competitor by id
  static async get(
    workspaceId: string,
    competitorId: string,
    token: string,
    includePages: boolean = false
  ): Promise<GetCompetitorResponse["data"]> {
    const { data } = await axios.get<GetCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}`,
      {
        params: { include_pages: includePages },
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // create a new competitor for a workspace
  static async create(
    workspaceId: string,
    pages: CreateCompetitorRequest,
    token: string
  ): Promise<CreateCompetitorResponse["data"]> {
    const { data } = await axios.post<CreateCompetitorResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors`,
      pages,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // update a competitor by id
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
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    return data.data;
  }

  // delete a competitor by id
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
