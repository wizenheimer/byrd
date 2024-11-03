"use server";

import prisma from "@/lib/db";
import { createCompetitorProperties } from "@/services/property";
import { inviteTeamMember } from "@/services/team";
import { createUser, type UserCreateData } from "@/services/user";
import { createWorkspace } from "@/services/workspace";

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
	data: OnboardingData,
): Promise<{ success: boolean; workspaceId: string }> {
	try {
		console.log("Persisting onboarding data:", data);
		return await prisma.$transaction(async (tx) => {
			// Create primary user
			const userData: UserCreateData = {
				email: data.email,
				firstName: data.firstName,
				lastName: data.lastName,
				clerkId: data.clerkId,
			};

			const user = await createUser(userData);

			// Create workspace
			const workspace = await createWorkspace(
				`${data.firstName}'s Workspace`,
				user.id,
			);

			// Process competitors
			const competitorPromises = data.competitors.map((competitor) =>
				createCompetitorProperties(competitor.url, workspace.id),
			);
			await Promise.all(competitorPromises);

			// Invite team members
			const teamPromises = data.team
				.filter((member) => member.email !== data.email)
				.map((member) =>
					inviteTeamMember({
						email: member.email,
						workspaceId: workspace.id,
					}),
				);
			await Promise.all(teamPromises);

			return { success: true, workspaceId: workspace.id };
		});
	} catch (error) {
		console.error("Error persisting onboarding data:", error);
		throw new Error("Failed to persist onboarding data");
	}
}
