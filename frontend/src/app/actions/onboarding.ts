"use server";

import { randomUUID } from "node:crypto";

export interface OnboardingData {
	competitors: string[]; // urls
	profiles: string[]; // profile names
	features: string[]; // feature names
	team: string[]; // emails
}

export async function persistOnboardingData(
	data: OnboardingData,
): Promise<{ success: boolean; workspaceId: string }> {
	console.log("Persisting onboarding data", data);
	return { success: true, workspaceId: randomUUID() };
}
