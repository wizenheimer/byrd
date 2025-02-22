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

interface InstallationData {
	competitors: string[];
	features: string[];
	profiles: string[];
}

export async function handleSlackInit(data: InstallationData) {
	try {
		const response = await fetch(
			`${process.env.BACKEND_ORIGIN}/api/public/v1/integration/slack/oauth/init`,
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Host: headers().get("host") || "",
				},
				body: JSON.stringify(data),
			},
		);

		if (!response.ok) {
			throw new Error("Failed to initialize OAuth");
		}

		const { oauth_url } = await response.json();
		return { success: true, oauth_url };
	} catch (error) {
		console.error("Failed to initiate Slack OAuth:", error);
		return {
			success: false,
			error:
				error instanceof Error ? error.message : "Failed to initialize OAuth",
		};
	}
}
