import axios from "axios";
import type { DeleteUserResponse, GetUserResponse } from "../types/api";

// User Management
// biome-ignore lint/complexity/noStaticOnlyClass:
export class Users {
	static async getCurrent(
		token: string,
		origin: string,
	): Promise<GetUserResponse["data"]> {
		const { data } = await axios.get<GetUserResponse>(
			`${origin}/api/public/v1/users`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	static async delete(
		token: string,
		origin: string,
	): Promise<DeleteUserResponse> {
		const { data } = await axios.delete<DeleteUserResponse>(
			`${origin}/api/public/v1/users`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data;
	}
}
