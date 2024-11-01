"use client";
import { OnboardingHeader } from "@/components/OnboardingHeader";
import { OnboardingLayout } from "@/components/OnboardingLayout";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
	// ArrowLeft,
	// BarChart2,
	Inbox,
	Megaphone,
	Rss,
	Share2,
} from "lucide-react";
// import Link from "next/link";
import * as React from "react";

export default function ChannelsStep() {
	const [selected, setSelected] = React.useState<string[]>([
		"inbox",
		"mentions",
	]);

	const cards = [
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

	const toggleCard = (id: string) => {
		setSelected((prev) =>
			prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id],
		);
	};

	return (
		<OnboardingLayout
			previewImage="/onboarding/third.png"
			previewAlt="Dashboard Preview"
		>
			<OnboardingHeader
				title="Never Miss A Beat"
				description="Your competitors are everywhere. So are we."
			/>

			{/* Core */}
			<div className="space-y-6">
				{cards.map(({ id, icon: Icon, title, description }) => (
					<button
						key={id}
						type="button"
						onClick={() => toggleCard(id)}
						className={`relative flex w-full items-start gap-4 rounded-xl border-2 p-4 text-left transition-colors
				  ${selected.includes(id)
								? "border-primary bg-primary/5"
								: "border-border bg-background hover:border-primary/50"
							}`}
					>
						<Icon className="mt-1 size-5 shrink-0" />
						<div className="space-y-1">
							<div className="font-medium">{title}</div>
							<div className="text-sm text-muted-foreground">
								{description}
							</div>
						</div>
						{selected.includes(id) && (
							<div className="absolute right-3 top-1/2 size-4 -translate-y-1/2 rounded-full bg-primary" />
						)}
					</button>
				))}
			</div>

			{/* Footer */}
			<div className="space-y-6">
				<Button
					className={cn(
						"w-full h-12 text-base font-semibold",
						"bg-[#14171F] hover:bg-[#14171F]/90",
						"transition-colors duration-200",
						"disabled:opacity-50 disabled:cursor-not-allowed",
					)}
					size="lg"
				>
					Continue
				</Button>
				<p className="text-sm text-muted-foreground text-center">
					You can always customize them later
				</p>
			</div>

		</OnboardingLayout>
	);
}
