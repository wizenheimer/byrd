// src/constants/onboarding.ts
import { Inbox, Megaphone, Rss, Share2 } from "lucide-react";
import type { ChannelCard } from "@/app/_types/onboarding";

export const STEPS = {
  COMPETITOR: 1,
  FEATURES: 2,
  CHANNELS: 3,
  TEAM: 4,
  AUTH: 5,
} as const;

export const INITIAL_FEATURES = [
  {
    id: "1",
    title: "Product",
    description: "Catch product evolution in real-time",
    enabled: true,
  },
  {
    id: "2",
    title: "Pricing",
    description: "Never be the last to know about a price war",
    enabled: false,
  },
  {
    id: "3",
    title: "Partnership",
    description: "Track who's teaming up with whom",
    enabled: false,
  },
  {
    id: "4",
    title: "Branding",
    description: "Monitor messaging shifts, and identity changes",
    enabled: false,
  },
  {
    id: "5",
    title: "Positioning",
    description: "Track narratives before they go mainstream",
    enabled: false,
  },
] as const;

export const CHANNEL_CARDS: ChannelCard[] = [
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
    description: "Follow their social playbook as it unfolds",
  },
  {
    id: "mentions",
    icon: Megaphone,
    title: "Mentions",
    description: "Beat them to their own announcement",
  },
  {
    id: "content",
    icon: Rss,
    title: "Content",
    description: "Watch their messaging evolve, post by post",
  },
];

export const STEP_INFO = {
  [STEPS.COMPETITOR]: {
    title: "Your Market, Your Rules",
    description: "Pick your targets. Add up to 5 competitors.",
    image: "/onboarding/first.png",
  },
  [STEPS.FEATURES]: {
    title: "Measure What Matters",
    description: "Choose your signals. Cut through the noise.",
    image: "/onboarding/second.png",
  },
  [STEPS.CHANNELS]: {
    title: "Never Miss A Beat",
    description: "Your competitors are everywhere. So are we.",
    image: "/onboarding/third.png",
  },
  [STEPS.TEAM]: {
    title: "Build Your War Room",
    description: "Business is a team sport. Bring in your heavy hitters.",
    image: "/onboarding/four.png",
  },
  [STEPS.AUTH]: {
    title: "You're almost there",
    description: "Quick auth, then let's get started.",
    image: "/onboarding/five.png",
  },
} as const;
