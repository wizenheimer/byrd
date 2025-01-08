// src/services/team.ts

import prisma from "@/lib/db";
import type { UserRole } from "@prisma/client";

export interface TeamMemberInvite {
  email: string;
  role?: UserRole;
  workspaceId: string;
}

export async function inviteTeamMember({
  email,
  role = "MEMBER",
  workspaceId,
}: TeamMemberInvite): Promise<void> {
  // Create user if doesn't exist and add to workspace
  await prisma.user.upsert({
    where: { email },
    create: {
      email,
      firstName: "",
      lastName: "",
      workspaces: {
        create: {
          workspaceId,
          role,
        },
      },
    },
    update: {
      workspaces: {
        create: {
          workspaceId,
          role,
        },
      },
    },
  });

  // TODO: Implement email invitation system
}

export async function removeTeamMember(userId: string, workspaceId: string) {
  return prisma.usersOnWorkspaces.delete({
    where: {
      userId_workspaceId: {
        userId,
        workspaceId,
      },
    },
  });
}

export async function updateTeamMemberRole(
  userId: string,
  workspaceId: string,
  role: UserRole
) {
  return prisma.usersOnWorkspaces.update({
    where: {
      userId_workspaceId: {
        userId,
        workspaceId,
      },
    },
    data: { role },
  });
}
