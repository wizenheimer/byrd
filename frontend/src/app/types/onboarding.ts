import type { LucideIcon } from "lucide-react";
import { z } from "zod";
import { DEFAULT_PROFILE } from "./competitor_page";

// URL validation schema
export const urlSchema = z
  .string()
  .trim()
  .toLowerCase()
  .refine((url) => {
    try {
      const parsedUrl = new URL(
        url.startsWith("http") ? url : `https://${url}`
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

export interface FeatureCard {
  id: string;
  icon: LucideIcon;
  title: string;
  description: string;
}

export const profileSchema = z.object({
  id: z.string(),
  title: z.string(),
  profile_key: z.enum(DEFAULT_PROFILE),
  description: z.string(),
  enabled: z.boolean(),
});

export const featureSchema = z.string();

const teamMemberSchema = z.object({
  email: z.string().email("Invalid email address").min(1, "Email is required"),
});

export const teamFormSchema = z.object({
  members: z
    .array(teamMemberSchema)
    .max(5, "Maximum 5 team members allowed")
    .refine((members) => {
      const emails = members.map((m) => m.email.toLowerCase());
      return new Set(emails).size === emails.length;
    }, "Duplicate email addresses are not allowed"),
});

export const onboardingFormSchema = z.object({
  competitors: competitorFormSchema.shape.competitors,
  profile: z.array(profileSchema),
  feature: z.array(featureSchema),
  team: z.array(teamFormSchema),
});

export type ProfileFormData = z.infer<typeof profileSchema>;
export type FeatureFormData = z.infer<typeof featureSchema>;
export type TeamFormData = z.infer<typeof teamFormSchema>;
export type OnboardingFormData = z.infer<typeof onboardingFormSchema>;
