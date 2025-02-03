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

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Pages {
	// add pages to a competitor
	static async add(
		workspaceId: string,
		competitorId: string,
		pages: CreatePageRequest[],
		token: string,
	): Promise<AddPagesResponse["data"]> {
		if (!Array.isArray(pages) || pages.length === 0) {
			throw new Error("Pages array must not be empty");
		}

		// biome-ignore lint/complexity/noForEach: <explanation>
		pages.forEach((page) => {
			if (!page.url) throw new Error("URL is required for each page");
			try {
				new URL(page.url);
			} catch {
				throw new Error(`Invalid URL format: ${page.url}`);
			}
		});

		const { data } = await axios.post<AddPagesResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages`,
			pages,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// list all pages for a competitor
	static async list(
		workspaceId: string,
		competitorId: string,
		params: PaginationParams,
		token: string,
	): Promise<ListPagesResponse["data"]> {
		const { _page = 1, _limit = 10 } = params;
		const { data } = await axios.get<ListPagesResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages`,
			{
				params: { _page, _limit },
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	// get a page by id
	static async get(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		token: string,
	): Promise<GetPageResponse["data"]> {
		const { data } = await axios.get<GetPageResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// list page history
	static async listHistory(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		params: PaginationParams,
		token: string,
	): Promise<ListPageHistoryResponse["data"]> {
		const { _page = 1, _limit = 10 } = params;
		const { data } = await axios.get<ListPageHistoryResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}/history`,
			{
				params: { _page, _limit },
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	// update a page by id
	static async update(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		updateData: UpdatePageRequest,
		token: string,
	): Promise<UpdatePageResponse["data"]> {
		const { data } = await axios.put<UpdatePageResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			updateData,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// remove a page from a competitor
	static async delete(
		workspaceId: string,
		competitorId: string,
		pageId: string,
		token: string,
	): Promise<DeletePageResponse> {
		const { data } = await axios.delete<DeletePageResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/pages/${pageId}`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data;
	}
}
