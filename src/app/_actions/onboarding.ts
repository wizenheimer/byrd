"use server";

import prisma from "@/lib/db";
import { createCompetitorProperties } from "@/services/property";
import { inviteTeamMember } from "@/services/team";
import { type UserCreateData, createUser } from "@/services/user";
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

// export async function persistOnboardingData(
// 	data: OnboardingData,
// ): Promise<{ success: boolean; workspaceId: string }> {
// 	try {
// 		// console.log("Persisting onboarding data:", data);
// 		return await prisma.$transaction(
// 			async () => {
// 				// Create primary user
// 				const userData: UserCreateData = {
// 					email: data.email,
// 					firstName: data.firstName,
// 					lastName: data.lastName,
// 					clerkId: data.clerkId,
// 				};

// 				const user = await createUser(userData);

// 				// Create workspace
// 				const workspace = await createWorkspace(
// 					`${data.firstName}'s Workspace`,
// 					user.id,
// 				);

// 				// Process competitors
// 				const competitorPromises = data.competitors.map((competitor) =>
// 					createCompetitorProperties(competitor.url, workspace.id),
// 				);
// 				await Promise.all(competitorPromises);

// 				// Invite team members
// 				const teamPromises = data.team
// 					.filter((member) => member.email !== data.email)
// 					.map((member) =>
// 						inviteTeamMember({
// 							email: member.email,
// 							workspaceId: workspace.id,
// 						}),
// 					);
// 				await Promise.all(teamPromises);

// 				return { success: true, workspaceId: workspace.id };
// 			},
// 			{
// 				// TODO: break down the transaction into smaller chunks
// 				maxWait: 6000, // default: 2000
// 				timeout: 15000, // default: 5000
// 			},
// 		);
// 	} catch (error) {
// 		console.error("Error persisting onboarding data:", error);
// 		throw new Error("Failed to persist onboarding data");
// 	}
// }

export async function persistOnboardingData(
	data: OnboardingData,
): Promise<{ success: boolean; workspaceId: string }> {
	try {
		// Step 1: Create user and workspace together (these need to be atomic)
		const { workspace } = await prisma.$transaction(
			async () => {
				const userData: UserCreateData = {
					email: data.email,
					firstName: data.firstName,
					lastName: data.lastName,
					clerkId: data.clerkId,
				};

				const user = await createUser(userData);
				const workspace = await createWorkspace(
					`${data.firstName}'s Workspace`,
					user.id,
				);

				return { user, workspace };
			},
			{
				maxWait: 2000,
				timeout: 5000,
			},
		);

		// Step 2: Process competitors in batches
		const BATCH_SIZE = 5;
		const competitors = data.competitors;

		for (let i = 0; i < competitors.length; i += BATCH_SIZE) {
			const batch = competitors.slice(
				i,
				Math.min(i + BATCH_SIZE, competitors.length),
			);
			await prisma.$transaction(
				async () => {
					const promises = batch.map((competitor) =>
						createCompetitorProperties(competitor.url, workspace.id),
					);
					await Promise.all(promises);
				},
				{
					maxWait: 2000,
					timeout: 5000,
				},
			);
		}

		// Step 3: Process team invites in batches
		const teamMembers = data.team.filter(
			(member) => member.email !== data.email,
		);

		for (let i = 0; i < teamMembers.length; i += BATCH_SIZE) {
			const batch = teamMembers.slice(
				i,
				Math.min(i + BATCH_SIZE, teamMembers.length),
			);
			await prisma.$transaction(
				async () => {
					const promises = batch.map((member) =>
						inviteTeamMember({
							email: member.email,
							workspaceId: workspace.id,
						}),
					);
					await Promise.all(promises);
				},
				{
					maxWait: 2000,
					timeout: 5000,
				},
			);
		}

		return { success: true, workspaceId: workspace.id };
	} catch (error) {
		console.error("Error persisting onboarding data:", error);
		throw new Error("Failed to persist onboarding data");
	}
}
