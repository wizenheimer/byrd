import axios from "axios";
import type { DeleteUserResponse, GetUserResponse } from "../types/api";

// biome-ignore lint/complexity/noStaticOnlyClass:
export class Users {
  static async getCurrent(token: string): Promise<GetUserResponse["data"]> {
    const { data } = await axios.get<GetUserResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/users`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data.data;
  }

  static async delete(token: string): Promise<DeleteUserResponse> {
    const { data } = await axios.delete<DeleteUserResponse>(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/users`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    return data;
  }
}
