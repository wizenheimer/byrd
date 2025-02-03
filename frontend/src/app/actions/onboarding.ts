"use server";

import { Workspace } from "../queries/workspace";
import type { ProfileType } from "../types/competitor_page";

export interface OnboardingData {
	competitors: string[]; // urls
	profiles: ProfileType[]; // profile names
	features: string[]; // feature names
	team: string[]; // emails
}

export async function persistOnboardingData(
	data: OnboardingData,
	token: string,
): Promise<{ success: boolean; workspaceId: string }> {
	try {
		const origin = process.env.BACKEND_ORIGIN;
		if (!origin) {
			throw new Error("BACKEND_ORIGIN environment variable is not set");
		}

		// Create workspace in local database or state
		const workspaceData = await Workspace.create(
			{
				competitors: data.competitors,
				team: data.team,
				profiles: data.profiles,
				features: data.features,
			},
			token,
		);

		return {
			success: true,
			workspaceId: workspaceData.id,
		};
	} catch (error) {
		console.error("error persisting onboarding data:", error);
		throw error;
	}
}
