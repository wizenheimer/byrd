// src/services/workspace.ts
import prisma from "@/lib/db";
import type { UserRole } from "@prisma/client";

export interface OnboardingData {
  clerkId: string;
  email: string;
  firstName: string;
  lastName: string;
  competitors: { url: string }[];
  features: { title: string }[];
  channels: { title: string }[];
  team: { email: string }[];
}

export async function createWorkspace(
  name: string,
  ownerId: string,
  ownerRole: UserRole = "ADMIN"
) {
  return prisma.workspace.create({
    data: {
      name,
      users: {
        create: {
          userId: ownerId,
          role: ownerRole,
        },
      },
    },
  });
}

export async function getWorkspaceData(workspaceId: string) {
  return prisma.workspace.findUnique({
    where: { id: workspaceId },
    include: {
      users: {
        include: {
          user: true,
        },
      },
      properties: {
        include: {
          property: true,
        },
      },
    },
  });
}

export async function updateWorkspace(
  workspaceId: string,
  data: {
    name?: string;
    subscriptionStatus?: string;
    subscriptionPlan?: string;
  }
) {
  return prisma.workspace.update({
    where: { id: workspaceId },
    data,
  });
}
