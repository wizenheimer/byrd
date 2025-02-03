// api/competitor_report.ts
import axios from "axios";
import type {
	CreateReportResponse,
	DispatchReportRequest,
	DispatchReportResponse,
	ListReportsResponse,
	PaginationParams,
} from "../types/api";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class CompetitorReport {
	// list all reports for a competitor
	static async list(
		workspaceId: string,
		competitorId: string,
		params: PaginationParams,
		token: string,
	): Promise<ListReportsResponse["data"]> {
		const { _page = 1, _limit = 10 } = params;
		const { data } = await axios.get<ListReportsResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/reports`,
			{
				params: { _page, _limit },
				headers: { Authorization: `Bearer ${token}` },
			},
		);
		return data.data;
	}

	// create a new report for a competitor
	static async create(
		workspaceId: string,
		competitorId: string,
		token: string,
	): Promise<CreateReportResponse["data"]> {
		const { data } = await axios.post<CreateReportResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/reports`,
			{},
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// dispatch a report to a list of emails
	static async send(
		workspaceId: string,
		competitorId: string,
		emails: string[],
		token: string,
	): Promise<DispatchReportResponse["data"]> {
		const request: DispatchReportRequest = { emails };
		const { data } = await axios.post<DispatchReportResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/workspace/${workspaceId}/competitors/${competitorId}/reports/dispatch`,
			request,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}
}
