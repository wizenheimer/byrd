// src/app/_actions/onboarding.ts
"use server";

import { randomUUID } from "crypto";

export interface OnboardingData {
  competitors: string[]; // urls
  features: string[]; // feature names
  channels: string[]; // channel names
  team: string[]; // emails
}

export async function persistOnboardingData(
  data: OnboardingData
): Promise<{ success: boolean; workspaceId: string }> {
  console.log("Persisting onboarding data", data);
  return { success: true, workspaceId: randomUUID() };
}
