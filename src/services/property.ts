// property.ts
import { PropertyCategory } from "@prisma/client";
import prisma from "@/lib/db";

interface PropertyCoordinates {
	points: { x: number; y: number }[];
}

export async function createCompetitorProperties(
	competitorUrl: string,
	workspaceId: string,
): Promise<void> {
	const defaultCoordinates = {
		points: [
			{ x: 0, y: 0 },
			{ x: 100, y: 0 },
			{ x: 100, y: 100 },
			{ x: 0, y: 100 },
		],
	};

	const propertyPromises = Object.values(PropertyCategory).map((category) => {
		return prisma.property.create({
			data: {
				origin: competitorUrl,
				route: "/",
				category,
				coordinates: defaultCoordinates,
				threshold: 0.5,
				workspaces: {
					create: {
						workspaceId,
					},
				},
			},
		});
	});

	await Promise.all(propertyPromises);
}

export async function getCompetitorProperties(
	workspaceId: string,
	competitorUrl: string,
) {
	return prisma.propertiesOnWorkspaces.findMany({
		where: {
			workspaceId,
			property: {
				origin: competitorUrl,
			},
		},
		include: {
			property: true,
		},
	});
}

export async function updatePropertyCoordinates(
	propertyId: string,
	coordinates: PropertyCoordinates,
) {
	return prisma.property.update({
		where: { id: propertyId },
		data: { coordinates: JSON.stringify(coordinates) },
	});
}

export async function updatePropertyThreshold(
	propertyId: string,
	threshold: number,
) {
	return prisma.property.update({
		where: { id: propertyId },
		data: { threshold },
	});
}
