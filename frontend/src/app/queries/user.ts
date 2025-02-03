import axios from "axios";
import type {
	CreateOrUpdateUserResponse,
	DeleteUserResponse,
	GetUserResponse,
} from "../types/api";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Account {
	// get current user
	static async get(token: string): Promise<GetUserResponse["data"]> {
		const { data } = await axios.get<GetUserResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/users`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// delete a user account
	static async delete(token: string): Promise<DeleteUserResponse["data"]> {
		const { data } = await axios.delete<DeleteUserResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/users`,
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}

	// create a new user account
	static async create(
		token: string,
	): Promise<CreateOrUpdateUserResponse["data"]> {
		const { data } = await axios.post<CreateOrUpdateUserResponse>(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/users`,
			{},
			{ headers: { Authorization: `Bearer ${token}` } },
		);
		return data.data;
	}
}
