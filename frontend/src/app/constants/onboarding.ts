import type { FeatureCard } from "@/app/types/onboarding";
import { AppWindow, HandHeart, Inbox, Share2 } from "lucide-react";

export const STEPS = {
	COMPETITOR: 1,
	PROFILE: 2,
	FEATURES: 3,
	TEAM: 4,
	AUTH: 5,
} as const;

export const INITIAL_PROFILES = [
	{
		id: "1",
		title: "Product",
		profile_key: "product",
		description: "Catch product evolution in real-time",
		enabled: true,
	},
	{
		id: "2",
		title: "Pricing",
		profile_key: "pricing",
		description: "Never be the last to know about a price war",
		enabled: false,
	},
	{
		id: "3",
		title: "Partnerships",
		profile_key: "partnerships",
		description: "Track who's teaming up with whom",
		enabled: false,
	},
	{
		id: "4",
		title: "Branding",
		profile_key: "branding",
		description: "Monitor messaging shifts, and identity changes",
		enabled: false,
	},
	{
		id: "5",
		title: "Customers",
		profile_key: "customers",
		description: "They churn, you learn. Turn rage into raves",
		enabled: false,
	},
] as const;

export const FEATURE_CARD: FeatureCard[] = [
	{
		id: "page",
		icon: AppWindow,
		title: "Page",
		description: "Track every change that matters",
	},
	{
		id: "inbox",
		icon: Inbox,
		title: "Inbox",
		description: "Monitor the direct line to their customers",
	},
	{
		id: "social",
		icon: Share2,
		title: "Social",
		description: "Your front-row seat to their social strategy",
	},
	{
		id: "reviews",
		icon: HandHeart,
		title: "Reviews",
		description: "Turn their user feedback into your feature wins",
	},
];

export const STEP_INFO = {
	[STEPS.COMPETITOR]: {
		title: "Your Market, Your Rules",
		description: "Pick your targets. Add up to 5 competitors.",
		image: "/onboarding/first.png",
	},
	[STEPS.PROFILE]: {
		title: "Measure What Matters",
		description: "Choose your signals. Cut through the noise.",
		image: "/onboarding/second.png",
	},
	[STEPS.FEATURES]: {
		title: "Never Miss A Beat",
		description: "Your competitors are everywhere. So are we.",
		image: "/onboarding/third.png",
	},
	[STEPS.TEAM]: {
		title: "Build Your War Room",
		description: "Winning is a team sport. Let's bring in your heavy hitters.",
		image: "/onboarding/four.png",
	},
	[STEPS.AUTH]: {
		title: "You're almost there",
		description: "Quick auth, then let's get started.",
		image: "/onboarding/five.png",
	},
} as const;
