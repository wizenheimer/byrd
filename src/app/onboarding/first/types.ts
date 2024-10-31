// types.ts
import { z } from "zod";

export const competitorFormSchema = z.object({
	competitors: z
		.array(
			z.object({
				name: z
					.string()
					.min(1, "Competitor name is required")
					.max(100, "Name must be less than 100 characters")
					.refine(
						(val) => !val.trim().match(/^\d+$/),
						"Name cannot be only numbers",
					),
			}),
		)
		.min(1, "At least one competitor is required"),
});

export type CompetitorFormData = z.infer<typeof competitorFormSchema>;
