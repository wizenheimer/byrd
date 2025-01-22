import axios from "axios";
import type {
	AddPagesResponse,
	CreatePageRequest,
	DeletePageResponse,
	GetPageResponse,
	ListPageHistoryResponse,
	ListPagesResponse,
	PaginationParams,
	UpdatePageRequest,
	UpdatePageResponse,
} from "../types/api";

// Page Management
// biome-ignore lint/complexity/noStaticOnlyClass:
export class Pages {
	static async add(
		workspaceId: string,
		competitorId: string,
		pages: CreatePageRequest[],
		token: string,
		origin: string,
	): Promise<AddPagesResponse["data"]> {
		if (!Array.isArray(pages) || pages.length === 0) {
			throw new Error("Pages array must not be empty");
		}

		for (const page of pages) {
			if (!page.url) {
				throw new Error("URL is required for each page");
			}
			try {
				new URL(page.url);
			} catch {
				throw new Error(`Invalid URL format: ${page.url}`);
			}
		}

		const { data } = await axios.post<AddPagesResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages`,
			pages,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async list(
		workspaceId: string,
		competitorId: string,
		params: PaginationParams,
		token: string,
		origin: string,
	): Promise<ListPagesResponse["data"]> {
		const { _page = 1, _limit = 10 } = params;
		const { data } = await axios.get<ListPagesResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages`,
			{
				params: { _page, _limit },
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	static async get(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		token: string,
		origin: string,
	): Promise<GetPageResponse["data"]> {
		const { data } = await axios.get<GetPageResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async getHistory(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		params: PaginationParams,
		token: string,
		origin: string,
	): Promise<ListPageHistoryResponse["data"]> {
		const { _page = 1, _limit = 10 } = params;
		const { data } = await axios.get<ListPageHistoryResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}/history`,
			{
				params: { _page, _limit },
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	static async update(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		updateData: UpdatePageRequest,
		token: string,
		origin: string,
	): Promise<UpdatePageResponse["data"]> {
		const { data } = await axios.put<UpdatePageResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			updateData,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async delete(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		token: string,
		origin: string,
	): Promise<DeletePageResponse> {
		const { data } = await axios.delete<DeletePageResponse>(
			`${origin}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data;
	}
}
