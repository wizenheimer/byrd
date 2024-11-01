import type { LucideIcon } from "lucide-react";
import { z } from "zod";

// URL validation schema
export const urlSchema = z
	.string()
	.trim()
	.toLowerCase()
	.refine((url) => {
		try {
			const parsedUrl = new URL(
				url.startsWith("http") ? url : `https://${url}`,
			);
			return parsedUrl.protocol === "http:" || parsedUrl.protocol === "https:";
		} catch {
			return false;
		}
	}, "Please enter a valid website URL")
	.transform((url) => {
		if (!url.startsWith("http")) {
			return `https://${url}`;
		}
		return url;
	});

export const competitorSchema = z.object({
	url: urlSchema,
	favicon: z.string().optional(),
});

export const competitorFormSchema = z.object({
	competitors: z
		.array(competitorSchema)
		.min(1, "Add at least one competitor")
		.max(5, "Maximum 5 competitors allowed")
		.refine((competitors) => {
			const urls = competitors.map((c) => c.url);
			return new Set(urls).size === urls.length;
		}, "Duplicate websites are not allowed"),
});

export type CompetitorFormData = z.infer<typeof competitorFormSchema>;

export interface ChannelCard {
	id: string;
	icon: LucideIcon;
	title: string;
	description: string;
}
