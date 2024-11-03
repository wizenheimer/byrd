// user.ts
import prisma from "@/lib/db";

export interface UserCreateData {
	email: string;
	firstName: string;
	lastName: string;
	clerkId?: string;
}

export async function createUser(data: UserCreateData) {
	return prisma.user.create({
		data: {
			email: data.email,
			firstName: data.firstName,
			lastName: data.lastName,
		},
	});
}

export async function updateUser(
	userId: string,
	data: Partial<UserCreateData>,
) {
	return prisma.user.update({
		where: { id: userId },
		data,
	});
}

export async function getUserByEmail(email: string) {
	return prisma.user.findUnique({
		where: { email },
		include: {
			workspaces: {
				include: {
					workspace: true,
				},
			},
		},
	});
}

export async function getUserWorkspaces(userId: string) {
	return prisma.usersOnWorkspaces.findMany({
		where: { userId },
		include: {
			workspace: true,
		},
	});
}
