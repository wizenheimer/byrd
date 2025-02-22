// app/(onboarding)/actions.ts
"use server";

import { headers } from "next/headers";

export async function handleSlackCallback(code: string, state: string) {
	try {
		const response = await fetch(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/integration/slack/oauth/callback?${new URLSearchParams(
				{
					code,
					state,
				},
			)}`,
			{
				method: "GET",
				headers: {
					Accept: "application/json",
					// Forward necessary headers
					Host: headers().get("host") || "",
				},
			},
		);

		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(errorText || "Failed to complete OAuth");
		}

		const data = await response.json();
		return { success: true, data };
	} catch (error) {
		return {
			success: false,
			error:
				error instanceof Error
					? error.message
					: "Failed to complete installation",
		};
	}
}
