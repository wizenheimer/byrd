// src/app/_actions/onboarding.ts
"use server";

import { randomUUID } from "crypto";

interface Competitor {
  url: string;
}

interface Feature {
  title: string;
}

interface TeamMember {
  email: string;
}

interface Channel {
  title: string;
}

export interface OnboardingData {
  clerkId: string;
  email: string;
  firstName: string;
  lastName: string;

  competitors: Competitor[];
  features: Feature[];
  channels: Channel[];
  team: TeamMember[];
}

export async function persistOnboardingData(
  data: OnboardingData
): Promise<{ success: boolean; workspaceId: string }> {
  console.log("Persisting onboarding data", data);
  return { success: true, workspaceId: randomUUID() };
}
